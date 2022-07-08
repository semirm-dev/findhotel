package cache

import "github.com/semirm-dev/findhotel/geo"

type inmemory struct {
	items map[string]string
}

func NewInMemory() *inmemory {
	return &inmemory{
		items: make(map[string]string),
	}
}

func (c *inmemory) Store(items geo.CacheBucket) error {
	for k, v := range items {
		c.items[k] = v
	}

	return nil
}

func (c *inmemory) Get(keys []string) ([]string, error) {
	values := make([]string, 0)
	for _, k := range keys {
		values = append(values, c.items[k])
	}

	return values, nil
}

func (c *inmemory) All() map[string]string {
	return c.items
}
