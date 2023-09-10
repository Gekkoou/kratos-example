package data

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kratos-example/app/user/internal/biz"
	"kratos-example/app/user/internal/conf"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewUserRepo,
	NewRegistry,
	NewDiscovery,
	NewGormClient,
	NewRedisClient,
)

// Data .
type Data struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB, redis *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db, redis: redis}, cleanup, nil
}

func NewRegistry(conf *conf.Registry, logger log.Logger) (registry.Registrar, func(), error) {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Addr
	c.Scheme = conf.Consul.Scheme
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

func NewDiscovery(conf *conf.Registry, logger log.Logger) (registry.Discovery, func(), error) {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Addr
	c.Scheme = conf.Consul.Scheme
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

func NewGormClient(conf *conf.Data, logger log.Logger) *gorm.DB {
	client, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{})
	if conf.Database.Debug {
		client = client.Debug()
	}
	if err != nil {
		log.NewHelper(logger).Fatalf("failed opening connection to db: %v", err)
		panic(err)
	}
	if err = client.AutoMigrate(&biz.User{}); err != nil {
		panic(err)
	}
	return client
}

func NewRedisClient(conf *conf.Data, logger log.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		Password:     conf.Redis.Password,
		DialTimeout:  time.Second * 2,
		DB:           int(conf.Redis.Db),
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.NewHelper(logger).Fatalf("redis connect error: %v", err)
		panic(err)
	}
	return client
}
