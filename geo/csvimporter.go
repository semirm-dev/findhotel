package geo

import (
	"context"
	"encoding/csv"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

type csvImporter struct {
	path string
}

func NewCsvImporter(path string) Importer {
	return &csvImporter{
		path: path,
	}
}

func (imp *csvImporter) Import(ctx context.Context) *Imported {
	imported := &Imported{
		GeoData:  make(chan *Geo),
		OnError:  make(chan error),
		Finished: make(chan bool),
	}

	go func(ctx context.Context, imported *Imported) {
		defer func() {
			close(imported.Finished)
			logrus.Warn("csv importer finished")
		}()

		csvFile, csvErr := os.Open(imp.path)
		if csvErr != nil {
			logrus.Fatal(csvErr)
		}
		defer func() {
			if err := csvFile.Close(); err != nil {
				logrus.Error(err)
				return
			}
			logrus.Warn("csv file closed")
		}()

		csvr := csv.NewReader(csvFile)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				row, err := csvr.Read()
				if err != nil {
					if err == io.EOF {
						return
					}
					imported.OnError <- err
					continue
				}

				geoData, err := encodeToGeo(row)
				if err != nil {
					imported.OnError <- err
					continue
				}
				imported.GeoData <- geoData
			}
		}
	}(ctx, imported)

	return imported
}

func encodeToGeo(row []string) (*Geo, error) {
	ip, ccode, country, city, lat, long, myst := row[0], row[1], row[2], row[3], row[4], row[5], row[6]

	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return nil, err
	}
	longitude, err := strconv.ParseFloat(long, 64)
	if err != nil {
		return nil, err
	}
	mystVal, err := strconv.Atoi(myst)
	if err != nil {
		return nil, err
	}

	return &Geo{
		Ip:           ip,
		CountryCode:  ccode,
		Country:      country,
		City:         city,
		Latitude:     latitude,
		Longitude:    longitude,
		MysteryValue: mystVal,
	}, nil
}
