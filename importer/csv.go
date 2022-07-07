package importer

import (
	"context"
	"encoding/csv"
	"github.com/semirm-dev/findhotel/geo"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

type csvImporter struct {
	path      string
	batchSize int
}

func NewCsvImporter(path string, batchSize int) geo.Importer {
	return &csvImporter{
		path:      path,
		batchSize: batchSize,
	}
}

func (imp *csvImporter) Import(ctx context.Context) *geo.Imported {
	imported := &geo.Imported{
		GeoDataBatch: make(chan []*geo.Geo),
		OnError:      make(chan error),
	}

	go func(ctx context.Context, imported *geo.Imported) {
		defer func() {
			close(imported.GeoDataBatch)
			close(imported.OnError)
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
		buf := make([]*geo.Geo, 0, imp.batchSize)
		first := true

		for {
			select {
			case <-ctx.Done():
				return
			default:
				row, err := csvr.Read()
				if err != nil {
					if err == io.EOF {
						// check for leftover, incomplete buf
						if len(buf) > 0 {
							imported.GeoDataBatch <- buf
						}
						return
					}
					imported.OnError <- err
					continue
				}

				// skip first record, it's csv header
				if first {
					first = false
					continue
				}

				geoData, err := encodeToGeo(row)
				if err != nil {
					imported.OnError <- err
					continue
				}

				buf = append(buf, geoData)

				if len(buf) >= imp.batchSize {
					imported.GeoDataBatch <- buf
					buf = nil // reset buf
				}
			}
		}
	}(ctx, imported)

	return imported
}

func encodeToGeo(row []string) (*geo.Geo, error) {
	ip, ccode, country, city, lat, long, myst := row[0], row[1], row[2], row[3], row[4], row[5], row[6]

	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil && lat != "" {
		return nil, err
	}
	longitude, err := strconv.ParseFloat(long, 64)
	if err != nil && long != "" {
		return nil, err
	}
	mystVal, err := strconv.Atoi(myst)
	if err != nil && myst != "" {
		return nil, err
	}

	return &geo.Geo{
		Ip:           ip,
		CountryCode:  ccode,
		Country:      country,
		City:         city,
		Latitude:     latitude,
		Longitude:    longitude,
		MysteryValue: mystVal,
	}, nil
}
