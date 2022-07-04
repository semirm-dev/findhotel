package datastore

import (
	"github.com/semirm-dev/findhotel/geo"
	"gorm.io/gorm"
)

type pgStore struct {
	db *gorm.DB
}

func NewPgStore(db *gorm.DB) *pgStore {
	return &pgStore{
		db: db,
	}
}

func (store *pgStore) Store(geoData *geo.Geo) error {
	return nil
}

func (store *pgStore) StoreMultiple(geoData []*geo.Geo) error {
	return nil
}

func (store *pgStore) ByIp(ip string) (*geo.Geo, error) {
	return nil, nil
}
