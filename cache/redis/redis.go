package redis

import (
	"github.com/go-redis/redis"
	"github.com/semirm-dev/findhotel/geo"
)

// pipeLength defines limit whether to use pipeline or not
const pipeLength = 1

type cache struct {
	*redis.Client
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
	client := redis.NewClient(&redis.Options{
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

func (c *cache) Get(key string) (string, error) {
	cacheValue, err := c.Client.Get(key).Result()

	switch {
	// key does not exist
	case err == redis.Nil:
		// errors.New(fmt.Sprintf("key %v does not exist", key))
		return "", nil
	// some other error
	case err != nil:
		return "", err
	}

	return cacheValue, nil
}
