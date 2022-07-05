package inmemory

import "github.com/semirm-dev/findhotel/geo"

type cache struct {
	items map[string]string
}

func NewInMemory() *cache {
	return &cache{}
}

func (c *cache) Store(items geo.CacheBucket) error {
	for k, v := range items {
		c.items[k] = v
	}

	return nil
}

func (c *cache) Get(key string) (string, error) {
	return c.items[key], nil
}
