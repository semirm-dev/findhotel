package geo

import (
	"context"
	"github.com/sirupsen/logrus"
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

// Search will get *geo data from its source
type Search interface {
	ByIp(ip string) (*Geo, error)
}

// Importer will import *geo data from its source
type Importer interface {
	Import(context.Context) *Imported
}

// Storer will store *geo data in data store
type Storer interface {
	Store([]*Geo) error
}

type Imported struct {
	GeoData  chan *Geo
	OnError  chan error
	Finished chan bool
}

type loader struct {
	importer Importer
}

// NewLoader will initialize *loader.
// Loader will load *geo data from Importer and store it in data store using Storer
func NewLoader(importer Importer) *loader {
	return &loader{
		importer: importer,
	}
}

// Load will start loading *geo data from Importer to Storer
func (ldr *loader) Load(ctx context.Context) chan bool {
	finished := make(chan bool)
	imported := ldr.importer.Import(ctx)

	go func(ctx context.Context) {
		c := 0
		e := 0
		t := time.Now()

		defer func() {
			close(finished)
			logrus.Infof("successfully imported %d records", c)
			logrus.Infof("failed to import %d records", e)
			logrus.Infof("finished in %v", time.Now().Sub(t))
		}()

		for {
			select {
			case <-imported.GeoData:
				c++
			case <-imported.OnError:
				e++
			case <-imported.Finished:
				return
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return finished
}
