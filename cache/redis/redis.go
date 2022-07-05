package redis

import (
	redisLib "github.com/go-redis/redis"
	"github.com/semirm-dev/findhotel/geo"
)

// pipeLength defines limit whether to use pipeline or not
const pipeLength = 1

type cache struct {
	*redisLib.Client
	*Config
}

type Config struct {
	Host       string
	Port       string
	Password   string
	DB         int
	PipeLength int
}

func NewConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}
}

func NewCache(conf *Config) *cache {
	return &cache{
		Config: conf,
	}
}

func (c *cache) Initialize() error {
	client := redisLib.NewClient(&redisLib.Options{
		Addr:     c.Config.Host + ":" + c.Config.Port,
		Password: c.Config.Password, // no password set
		DB:       c.Config.DB,       // use default DB
	})

	if c.Config.PipeLength == 0 {
		c.Config.PipeLength = pipeLength
	}

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	c.Client = client

	return nil
}

func (c *cache) Store(items geo.CacheBucket) error {
	pipe := c.Pipeline()

	for k, v := range items {
		pipe.Set(k, v, -1)
	}

	_, err := pipe.Exec()
	return err
}

func (c *cache) Get(keys []string) ([]string, error) {
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
