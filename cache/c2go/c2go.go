package c2go

import (
	"encoding/json"
	"github.com/muesli/cache2go"
)

type cache struct {
	engine *cache2go.CacheTable
}

func NewC2Go(tableName string) *cache {
	return &cache{
		engine: cache2go.Cache(tableName),
	}
}

func (c *cache) Store(key, value string) error {
	c.engine.Add(key, -1, value)

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
