package datastore

import "github.com/semirm-dev/findhotel/geo"

type inmemory struct {
	data []*geo.Geo
}

func NewInMemoryStore() *inmemory {
	return &inmemory{}
}

func (store *inmemory) Store(geoData *geo.Geo) error {
	store.data = append(store.data, geoData)

	return nil
}

func (store *inmemory) StoreMultiple(geoData []*geo.Geo) error {
	store.data = append(store.data, geoData...)

	return nil
}

func (store *inmemory) ByIp(ip string) (*geo.Geo, error) {
	for _, g := range store.data {
		if g.Ip == ip {
			return g, nil
		}
	}

	return nil, nil
}
