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

type CacheBucket map[string]string

// Cache is used to keep track of all previously saved *geo data.
// It's mainly used for validation to check if there are duplicate entries,
// that is to make less database calls on *geo data insert
type Cache interface {
	Store(CacheBucket) error
	Get([]string) ([]string, error)
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
	ldr.storeGeoData(ctx, filtered)

	logrus.Infof("=== total time finished in %v ===", time.Now().Sub(t))
}

func (g *Geo) valid() bool {
	if strings.TrimSpace(g.Ip) == "" {
		return false
	}

	return true
}

// filterValidGeoData will sanitize *geo data. Duplicate and corrupted entries will be removed/skipped.
func (ldr *loader) filterValidGeoData(ctx context.Context, imported *Imported) <-chan []*Geo {
	filtered := make(chan []*Geo)

	go func() {
		b := 0
		i := 0
		e := 0
		t := time.Now()

		defer func() {
			close(filtered)

			logrus.Infof("=== import csv ===\n"+
				"- total records for import = %d\n"+
				"- successfully imported = %d\n"+
				"- skipped records = %d\n"+
				"- import finished in %v", b+e, i, (b+e)-i, time.Now().Sub(t))
		}()

		for {
			select {
			case batch, ok := <-imported.GeoDataBatch:
				if !ok {
					return
				}
				b += len(batch)

				keysToCheck := make([]string, 0)
				validBatch := make([]*Geo, 0)
				for _, g := range batch {
					if g.valid() && !exists(g.Ip, keysToCheck) { // check duplicate ips per batch
						keysToCheck = append(keysToCheck, g.Ip)
						validBatch = append(validBatch, g)
					}
				}

				keysExist, err := ldr.cache.Get(keysToCheck)
				if err != nil {
					e += len(keysToCheck)
					break
				}

				cacheBucket := make(CacheBucket)
				buf := make([]*Geo, 0)
				for _, vb := range validBatch {
					if exists(vb.Ip, keysExist) { // check duplicate ips for entire storage
						continue
					}
					cacheBucket[vb.Ip] = vb.Ip
					buf = append(buf, vb)
					i++
				}

				if err = ldr.cache.Store(cacheBucket); err != nil {
					e += len(cacheBucket)
					break
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

// storeGeoData will store *geo data in database.
// It must be last in the line, all data should already be checked and validated.
func (ldr *loader) storeGeoData(ctx context.Context, geoData <-chan []*Geo) {
	b := 0
	i := 0
	e := 0
	t := time.Now()

	defer func() {
		logrus.Infof("=== store in db ===\n"+
			"- total records to store = %d\n"+
			"- successfully stored = %d\n"+
			"- failed to store = %d\n"+
			"- store finished in %v", b, i, e, time.Now().Sub(t))
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
				e += len(batch)
			}
		}
	}
}

func exists(key string, keys []string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}

	return false
}
