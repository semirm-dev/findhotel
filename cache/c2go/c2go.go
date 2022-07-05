package c2go

import (
	"encoding/json"
	"github.com/muesli/cache2go"
	"github.com/semirm-dev/findhotel/geo"
	"time"
)

type cache struct {
	engine *cache2go.CacheTable
}

func NewC2Go(tableName string) *cache {
	return &cache{
		engine: cache2go.Cache(tableName),
	}
}

func (c *cache) Store(items geo.CacheBucket) error {
	for k, v := range items {
		c.engine.Add(k, 24*time.Hour*7, v)
	}

	return nil
}

func (c *cache) Get(key string) (string, error) {
	cacheValue, err := c.engine.Value(key)
	if err != nil {
		return "", err
	}

	value, err := json.Marshal(cacheValue.Data())
	if err != nil {
		return "", err
	}

	return string(value), nil
}
