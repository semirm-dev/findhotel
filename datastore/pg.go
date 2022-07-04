package datastore

import (
	"github.com/semirm-dev/findhotel/geo"
	"gorm.io/gorm"
	"time"
)

type pgStore struct {
	db *gorm.DB
}

type Geo struct {
	Id           int    `gorm:"primarykey"`
	Ip           string `gorm:"uniqueIndex"`
	CountryCode  string
	Country      string
	City         string
	Latitude     float64
	Longitude    float64
	MysteryValue int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func NewPgStore(db *gorm.DB) *pgStore {
	db.AutoMigrate(&Geo{})

	return &pgStore{
		db: db,
	}
}

func (store *pgStore) Store(geoData *geo.Geo) error {
	return store.db.Create(geoToEntity(geoData)).Error
}

func (store *pgStore) StoreMultiple(geoData []*geo.Geo) error {
	var bulk []*Geo

	for _, g := range geoData {
		bulk = append(bulk, geoToEntity(g))
	}

	return store.db.Create(bulk).Error
}

func (store *pgStore) ByIp(ip string) (*geo.Geo, error) {
	return nil, nil
}

func geoToEntity(geoData *geo.Geo) *Geo {
	return &Geo{
		Ip:           geoData.Ip,
		CountryCode:  geoData.CountryCode,
		Country:      geoData.Country,
		City:         geoData.City,
		Latitude:     geoData.Latitude,
		Longitude:    geoData.Longitude,
		MysteryValue: geoData.MysteryValue,
	}
}
