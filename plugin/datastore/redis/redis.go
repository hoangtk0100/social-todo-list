package redisdb

import (
	"context"
	"flag"

	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultPoolSize      = 0 // 0 is unlimited number of socket connections
	DefaultMintIdleConns = 10
)

// <user>:<password>@<host>:<port>/<db_number>
type RedisDBOpt struct {
	Prefix       string
	URL          string
	PoolSize     int
	MinIdleConns int
}

type redisDB struct {
	name   string
	client *redis.Client
	logger logger.Logger
	*RedisDBOpt
}

func NewRedisDB(name, prefix string) *redisDB {
	return &redisDB{
		name: name,
		RedisDBOpt: &RedisDBOpt{
			Prefix:       prefix,
			PoolSize:     DefaultPoolSize,
			MinIdleConns: DefaultMintIdleConns,
		},
	}
}

func (r *redisDB) GetPrefix() string {
	return r.Prefix
}

func (r *redisDB) Name() string {
	return r.name
}

func (r *redisDB) Get() interface{} {
	return r.client
}

func (r *redisDB) InitFlags() {
	prefix := r.Prefix
	flag.StringVar(&r.URL,
		prefix+"-url",
		"redis://localhost:6379",
		"Redis connection-string. Ex: redis:<user>:<password>@<host>:<port>/<db_name>",
	)

	flag.IntVar(&r.PoolSize,
		prefix+"-pool-size",
		DefaultPoolSize,
		"Redis pool size",
	)

	flag.IntVar(&r.MinIdleConns,
		prefix+"-pool-min-idle",
		DefaultMintIdleConns,
		"Redis min idle connections",
	)
}

func (r *redisDB) isDisabled() bool {
	return r.URL == ""
}

func (r *redisDB) Configure() error {
	if r.isDisabled() {
		return nil
	}

	r.logger = logger.GetCurrent().GetLogger(r.name)

	opt, err := redis.ParseURL(r.URL)
	if err != nil {
		r.logger.Error("Cannot parse Redis URL", err.Error())
		return err
	}

	opt.PoolSize = r.PoolSize
	opt.MinIdleConns = r.MinIdleConns

	client := redis.NewClient(opt)

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		r.logger.Error("Cannot connect Redis", err.Error())
		return err
	}

	r.client = client

	r.logger.Print("Connecting Redis on %s", opt.Addr)

	return nil
}

func (r *redisDB) Run() error {
	return r.Configure()
}

func (r *redisDB) Stop() <-chan bool {
	if r.client != nil {
		if err := r.client.Close(); err != nil {
			r.logger.Info("Cannot close ", r.name)
		}
	}

	c := make(chan bool)
	go func() {
		c <- true
	}()

	return c
}
