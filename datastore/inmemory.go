package datastore

import "github.com/semirm-dev/findhotel/geo"

type inmemory struct {
	data []*geo.Geo
}

func NewInMemory() *inmemory {
	return &inmemory{}
}

func (storer *inmemory) Store(geoData []*geo.Geo) (int, error) {
	storer.data = append(storer.data, geoData...)

	return len(geoData), nil
}

func (storer *inmemory) ByIp(ip string) (*geo.Geo, error) {
	for _, g := range storer.data {
		if g.Ip == ip {
			return g, nil
		}
	}

	return nil, nil
}

func (storer *inmemory) All() []*geo.Geo {
	return storer.data
}
