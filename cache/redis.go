package cache

import (
	redisLib "github.com/go-redis/redis"
	"github.com/semirm-dev/findhotel/geo"
)

// pipeLength defines limit whether to use pipeline or not
const pipeLength = 1

type redis struct {
	*redisLib.Client
	*redisConfig
}

type redisConfig struct {
	Host       string
	Port       string
	Password   string
	DB         int
	PipeLength int
}

func NewRedisConfig() *redisConfig {
	return &redisConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}
}

func NewRedis(conf *redisConfig) *redis {
	return &redis{
		redisConfig: conf,
	}
}

func (c *redis) Initialize() error {
	client := redisLib.NewClient(&redisLib.Options{
		Addr:     c.redisConfig.Host + ":" + c.redisConfig.Port,
		Password: c.redisConfig.Password, // no password set
		DB:       c.redisConfig.DB,       // use default DB
	})

	if c.redisConfig.PipeLength == 0 {
		c.redisConfig.PipeLength = pipeLength
	}

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	c.Client = client

	return nil
}

func (c *redis) Store(items geo.CacheBucket) error {
	pipe := c.Pipeline()

	for k, v := range items {
		pipe.Set(k, v, -1)
	}

	_, err := pipe.Exec()
	return err
}

func (c *redis) Get(keys []string) ([]string, error) {
	pipe := c.Pipeline()

	for _, k := range keys {
		pipe.Get(k)
	}

	res, err := pipe.Exec()
	if err != nil && err != redisLib.Nil {
		return nil, err
	}

	values := make([]string, 0)
	for _, item := range res {
		values = append(values, item.(*redisLib.StringCmd).Val())
	}

	return values, nil
}
