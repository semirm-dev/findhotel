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

func (c *cache) Get(keys []string) ([]string, error) {
	values := make([]string, 0)
	for _, k := range keys {
		values = append(values, c.items[k])
	}

	return values, nil
}
