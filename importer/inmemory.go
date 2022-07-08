package importer

import (
	"context"
	"github.com/semirm-dev/findhotel/geo"
)

type inmemory struct {
	geoDataBatch []*geo.Geo
	batchSize    int
}

func NewInMemory(geoDataBatch []*geo.Geo, batchSize int) geo.Importer {
	return &inmemory{
		geoDataBatch: geoDataBatch,
		batchSize:    batchSize,
	}
}

func (imp inmemory) Import(ctx context.Context) *geo.Imported {
	imported := &geo.Imported{
		GeoDataBatch: make(chan []*geo.Geo),
		OnError:      make(chan error),
	}

	go func() {
		buf := make([]*geo.Geo, 0, imp.batchSize)

		defer func() {
			if len(buf) > 0 {
				imported.GeoDataBatch <- buf
			}

			close(imported.GeoDataBatch)
			close(imported.OnError)
		}()

		for _, b := range imp.geoDataBatch {
			select {
			case <-ctx.Done():
				return
			default:
				buf = append(buf, b)

				if len(buf) >= imp.batchSize {
					imported.GeoDataBatch <- buf
					buf = nil // reset buf
				}
			}
		}
	}()

	return imported
}
