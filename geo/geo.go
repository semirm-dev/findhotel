package geo

import (
	"context"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Geo struct {
	Ip           string
	CountryCode  string
	Country      string
	City         string
	Latitude     float64
	Longitude    float64
	MysteryValue int
}

// Importer will import *geo data from its source
type Importer interface {
	Import(context.Context) *Imported
}

// Storer will store *geo data in data store
type Storer interface {
	Store([]*Geo) (int, error)
}

// Search will get *geo data from its source
type Search interface {
	ByIp(ip string) (*Geo, error)
}

// Cache is used to store all previously *geo data.
// It's mainly used for validation to check if there are duplicate entries,
// which is to make less database calls on *geo data insert
type Cache interface {
	Store(string, string) error
	Get(string) (string, error)
}

// Imported presents each imported *geo data record/row
type Imported struct {
	GeoDataBatch chan []*Geo
	OnError      chan error
}

type loader struct {
	importer Importer
	storer   Storer
	cache    Cache
}

// NewLoader will initialize *loader.
// Loader will load *geo data from Importer and store it in data store using Storer
func NewLoader(importer Importer, storer Storer, cache Cache) *loader {
	return &loader{
		importer: importer,
		storer:   storer,
		cache:    cache,
	}
}

// Load will start loading *geo data from Importer to Storer
func (ldr *loader) Load(ctx context.Context) {
	t := time.Now()

	imported := ldr.importer.Import(ctx)
	filtered := ldr.filterValidGeoData(ctx, imported)
	cached := ldr.cacheGeoData(ctx, filtered)
	ldr.storeGeoData(ctx, cached)

	logrus.Info("---")
	logrus.Infof("total time finished in %v", time.Now().Sub(t))
}

func (g *Geo) valid() bool {
	if strings.TrimSpace(g.Ip) == "" {
		return false
	}

	return true
}

func (ldr *loader) filterValidGeoData(ctx context.Context, imported *Imported) chan []*Geo {
	filtered := make(chan []*Geo)

	go func() {
		b := 0
		i := 0
		e := 0
		c := 0
		t := time.Now()

		defer func() {
			logrus.Info("--- import csv")
			logrus.Infof("total records for import = %d", b+e)
			logrus.Infof("successfully imported = %d", i)
			logrus.Infof("previously cached (duplicate entries) = %d", c)
			logrus.Infof("invalid imports = %d", e)
			logrus.Infof("skipped records = %d", c+e)
			logrus.Infof("import finished in %v", time.Now().Sub(t))

			close(filtered)
		}()

		for {
			select {
			case batch, ok := <-imported.GeoDataBatch:
				if !ok {
					return
				}
				b += len(batch)

				buf := make([]*Geo, 0)
				for _, g := range batch {
					if g.valid() {
						cached, _ := ldr.cache.Get(g.Ip)
						if strings.TrimSpace(cached) != "" {
							c++
							continue
						}

						i++
						buf = append(buf, g)
					}
				}
				filtered <- buf
			case <-imported.OnError:
				e++
			case <-ctx.Done():
				return
			}
		}
	}()

	return filtered
}

func (ldr *loader) cacheGeoData(ctx context.Context, geoData <-chan []*Geo) <-chan []*Geo {
	cached := make(chan []*Geo)

	go func() {
		b := 0
		i := 0
		e := 0
		t := time.Now()

		defer func() {
			logrus.Info("--- cache")
			logrus.Infof("total records to cache = %d", b)
			logrus.Infof("successfully cached = %d", i)
			logrus.Infof("failed to cache = %d", e)
			logrus.Infof("cache finished in %v", time.Now().Sub(t))

			close(cached)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case batch, ok := <-geoData:
				if !ok {
					return
				}
				b += len(batch)

				buf := make([]*Geo, 0)
				for _, g := range batch {
					if err := ldr.cache.Store(g.Ip, g.Ip); err != nil {
						e++
						continue
					}
					i++
					buf = append(buf, g)
				}

				cached <- buf
			}
		}
	}()

	return cached
}

// storeGeoData must be last in the line, all data should be already checked and validated
func (ldr *loader) storeGeoData(ctx context.Context, geoData <-chan []*Geo) {
	b := 0
	i := 0
	e := 0
	t := time.Now()

	defer func() {
		logrus.Info("--- store in db")
		logrus.Infof("total records to store = %d", b)
		logrus.Infof("successfully stored = %d", i)
		logrus.Infof("failed to store = %d", e)
		logrus.Infof("store finished in %v", time.Now().Sub(t))
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case batch, ok := <-geoData:
			if !ok {
				return
			}
			b += len(batch)

			stored, err := ldr.storer.Store(batch)
			i += stored
			if err != nil {
				e++
				continue
			}
		}
	}
}
