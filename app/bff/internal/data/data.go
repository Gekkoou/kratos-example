package data

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kratos-example/api/user/v1"
	"kratos-example/app/bff/internal/conf"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewUserRepo,
	NewRegistry,
	NewDiscovery,
	NewGormClient,
	NewUserServiceClient,
	NewRedisClient,
)

// Data .
type Data struct {
	db    *gorm.DB
	redis *redis.Client
	uc    v1.UserClient
	conf  *conf.Bootstrap
}

// NewData .
func NewData(conf *conf.Bootstrap, logger log.Logger, db *gorm.DB, redis *redis.Client, uc v1.UserClient) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db, redis: redis, uc: uc, conf: conf}, cleanup, nil
}

func NewRegistry(conf *conf.Bootstrap, logger log.Logger) (registry.Registrar, func(), error) {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Registry.Consul.Addr
	c.Scheme = conf.Registry.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		log.NewHelper(logger).Fatalf("consul registry failed: %v", err)
		panic(err)
	}
	r := consul.New(cli)
	cleanup := func() {
		log.NewHelper(logger).Info("closing the consul registry")
	}
	return r, cleanup, nil
}

func NewDiscovery(conf *conf.Bootstrap, logger log.Logger) (registry.Discovery, func(), error) {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Registry.Consul.Addr
	c.Scheme = conf.Registry.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		log.NewHelper(logger).Fatalf("consul discovery failed: %v", err)
		panic(err)
	}
	dis := consul.New(cli)
	cleanup := func() {
		log.NewHelper(logger).Info("closing the consul discovery")
	}
	return dis, cleanup, nil
}

func NewGormClient(conf *conf.Bootstrap, logger log.Logger) *gorm.DB {
	client, err := gorm.Open(mysql.Open(conf.Data.Database.Source), &gorm.Config{})
	if conf.Data.Database.Debug {
		client = client.Debug()
	}
	if err != nil {
		log.NewHelper(logger).Fatalf("failed opening connection to db: %v", err)
		panic(err)
	}
	return client
}

func NewRedisClient(conf *conf.Bootstrap, logger log.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Data.Redis.Addr,
		ReadTimeout:  conf.Data.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Data.Redis.WriteTimeout.AsDuration(),
		Password:     conf.Data.Redis.Password,
		DialTimeout:  time.Second * 2,
		DB:           int(conf.Data.Redis.Db),
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.NewHelper(logger).Fatalf("redis connect error: %v", err)
		panic(err)
	}
	return client
}

func NewUserServiceClient(r registry.Discovery, logger log.Logger, conf *conf.Bootstrap) v1.UserClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///kratos.user.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
			jwt.Client(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(conf.Jwt.Client.SigningKey), nil
			}),
			logging.Client(logger),
			circuitbreaker.Client(),
		),
	)
	if err != nil {
		log.NewHelper(logger).Fatalf("grpc user service failed: %v", err)
		panic(err)
	}
	c := v1.NewUserClient(conn)
	return c
}
