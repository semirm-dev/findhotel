package inmemory

type cache struct {
	items map[string]string
}

func NewInMemory() *cache {
	return &cache{}
}

func (c *cache) Store(key, value string) error {
	c.items[key] = value
	return nil
}

func (c *cache) Get(key string) (string, error) {
	return c.items[key], nil
}
