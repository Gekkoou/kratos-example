package server

import (
	"context"
	"errors"
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/status"
	v1 "kratos-example/api/bff/v1"
	"kratos-example/app/bff/internal/conf"
	"kratos-example/app/bff/internal/pkg/jwt"
	"kratos-example/app/bff/internal/service"
	stdhttp "net/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Bootstrap, user *service.UserService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			selector.Server(JwtMiddleware(c.Jwt.Server.SigningKey, c.Jwt.Server.TokenKey)).Match(NewWhiteListMatcher()).Build(),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricRequests)),
			),
			ratelimit.Server(),
		),
		http.ErrorEncoder(EncoderError()),
	}
	if c.Server.Http.Network != "" {
		opts = append(opts, http.Network(c.Server.Http.Network))
	}
	if c.Server.Http.Addr != "" {
		opts = append(opts, http.Address(c.Server.Http.Addr))
	}
	if c.Server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Server.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	// prometheus 数据采集
	prometheus.MustRegister(_metricSeconds, _metricRequests)
	srv.Handle("/metrics", promhttp.Handler())

	v1.RegisterUserHTTPServer(srv, user)
	return srv
}

func JwtMiddleware(SigningKey string, tokenKey string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				jwtToken := header.RequestHeader().Get(tokenKey)
				if jwtToken == "" {
					return nil, errors.New("JWT token is missing")
				}
				j := jwt.NewJWT(SigningKey)
				claims, err := j.ParseToken(jwtToken)
				if err != nil {
					return nil, err
				}
				ctx = context.WithValue(ctx, "claims", claims)
				return handler(ctx, req)
			}
			return nil, errors.New("wrong context for middleware")
		}
	}
}

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/api.bff.v1.User/Login"] = struct{}{}
	whiteList["/api.bff.v1.User/Register"] = struct{}{}
	whiteList["/api.bff.v1.User/GetUserInfo"] = struct{}{}
	whiteList["/api.bff.v1.User/GetUserList"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

var (
	_metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_sec",
		Help:      "server requests duratio(sec).",
		Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1},
	}, []string{"kind", "operation"})

	_metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})
)

type httpResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func EncoderError() http.EncodeErrorFunc {
	return func(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
		if err == nil {
			return
		}
		se := &httpResponse{}
		gs, ok := status.FromError(err)
		if !ok {
			se = &httpResponse{Code: stdhttp.StatusInternalServerError}
		}
		se = &httpResponse{
			Code:    httpstatus.FromGRPCCode(gs.Code()),
			Message: gs.Message(),
			Data:    nil,
		}
		codec, _ := http.CodecForRequest(r, "Accept")
		body, err := codec.Marshal(se)
		if err != nil {
			w.WriteHeader(stdhttp.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/"+codec.Name())
		w.WriteHeader(se.Code)
		_, _ = w.Write(body)
	}
}
