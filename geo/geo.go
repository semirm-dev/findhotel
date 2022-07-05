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

// Imported presents each imported *geo data record/row
type Imported struct {
	GeoDataBatch chan []*Geo
	OnError      chan error
}

type loader struct {
	importer Importer
	storer   Storer
}

// NewLoader will initialize *loader.
// Loader will load *geo data from Importer and store it in data store using Storer
func NewLoader(importer Importer, storer Storer) *loader {
	return &loader{
		importer: importer,
		storer:   storer,
	}
}

// Load will start loading *geo data from Importer to Storer
func (ldr *loader) Load(ctx context.Context) {
	t := time.Now()

	imported := ldr.importer.Import(ctx)
	filtered := ldr.filterValidGeoData(ctx, imported)
	filteredUnique := ldr.filterUnique(ctx, filtered)
	ldr.storeGeoData(ctx, filteredUnique)

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
		c := 0
		b := 0
		e := 0
		t := time.Now()

		defer func() {
			logrus.Info("---")
			logrus.Infof("total records = %d", b+e)
			logrus.Infof("successfully imported %d records", c)
			logrus.Infof("failed to import %d records", (b-c)+e)
			logrus.Infof("finished in %v", time.Now().Sub(t))

			close(filtered)
		}()

		for {
			select {
			case geoData, ok := <-imported.GeoDataBatch:
				if !ok {
					return
				}
				b += len(geoData)

				batch := make([]*Geo, 0)
				for _, g := range geoData {
					if g.valid() {
						batch = append(batch, g)
					}
				}
				c += len(batch)
				filtered <- batch
			case <-imported.OnError:
				e++
			case <-ctx.Done():
				return
			}
		}
	}()

	return filtered
}

func (ldr *loader) filterUnique(ctx context.Context, filtered <-chan []*Geo) <-chan []*Geo {
	unique := make(chan []*Geo)

	go func() {
		defer close(unique)

		buf := make([]*Geo, 0)

		for {
			select {
			case <-ctx.Done():
				return
			case geoData, ok := <-filtered:
				if !ok {
					return
				}

				batch := make([]*Geo, 0)
				for _, g := range geoData {
					if !exists(g, buf) {
						batch = append(batch, g)
						buf = append(buf, g)
					}
				}
				if len(batch) > 0 {
					unique <- batch
				}
			}
		}
	}()

	return unique
}

func (ldr *loader) storeGeoData(ctx context.Context, filtered <-chan []*Geo) {
	c := 0
	b := 0
	t := time.Now()

	defer func() {
		logrus.Info("---")
		logrus.Infof("total records to store = %d", b)
		logrus.Infof("successfully stored %d records", c)
		logrus.Infof("failed to store %d records", b-c)
		logrus.Infof("store finished in %v", time.Now().Sub(t))
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case batch, ok := <-filtered:
			if !ok {
				return
			}
			b += len(batch)

			//stored, err := ldr.storer.Store(batch)
			//c += stored
			//if err != nil {
			//	continue
			//}
		}
	}
}

func exists(g *Geo, buf []*Geo) bool {
	for _, b := range buf {
		if g.Ip == b.Ip {
			return true
		}
	}

	return false
}
