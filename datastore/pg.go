package datastore

import (
	"github.com/semirm-dev/findhotel/geo"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	db.Logger = logger.Default.LogMode(logger.Silent)

	return &pgStore{
		db: db,
	}
}

func (storer *pgStore) Store(geoData []*geo.Geo) (int, error) {
	var bulk []*Geo

	for _, g := range geoData {
		existing, _ := storer.ByIp(g.Ip)
		if existing != nil {
			continue
		}

		bulk = append(bulk, geoToEntity(g))
	}

	if len(bulk) == 0 {
		return 0, nil
	}

	return len(bulk), storer.db.Create(bulk).Error
}

func (storer *pgStore) ByIp(ip string) (*geo.Geo, error) {
	var geoData *Geo
	if result := storer.db.Where("ip", ip).Find(&geoData); result.Error != nil {
		return nil, result.Error
	}
	if geoData.Id == 0 {
		return nil, nil
	}
	return entityToGeo(geoData), nil
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

func entityToGeo(entity *Geo) *geo.Geo {
	return &geo.Geo{
		Ip:           entity.Ip,
		CountryCode:  entity.CountryCode,
		Country:      entity.Country,
		City:         entity.City,
		Latitude:     entity.Latitude,
		Longitude:    entity.Longitude,
		MysteryValue: entity.MysteryValue,
	}
}
