package geo_test

import (
	"context"
	"github.com/semirm-dev/findhotel/cache"
	"github.com/semirm-dev/findhotel/datastore"
	"github.com/semirm-dev/findhotel/geo"
	"github.com/semirm-dev/findhotel/importer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoader_Load(t *testing.T) {
	testTable := map[string]struct {
		given               []*geo.Geo
		expectedCachedTotal int
		expectedStoredTotal int
	}{
		"valid data should store in cache and datastore": {
			given: []*geo.Geo{
				{
					Ip:           "1.1.1.1",
					CountryCode:  "cc1",
					Country:      "c1",
					City:         "city1",
					Latitude:     123.123,
					Longitude:    123.123,
					MysteryValue: 123,
				},
			},
			expectedCachedTotal: 1,
			expectedStoredTotal: 1,
		},
		"ignore duplicate data in cache and datastore": {
			given: []*geo.Geo{
				{
					Ip:           "1.1.1.1",
					CountryCode:  "cc1",
					Country:      "c1",
					City:         "city1",
					Latitude:     123.123,
					Longitude:    123.123,
					MysteryValue: 123,
				},
				{
					Ip:           "1.1.1.1",
					CountryCode:  "cc1",
					Country:      "c1",
					City:         "city1",
					Latitude:     123.123,
					Longitude:    123.123,
					MysteryValue: 123,
				},
				{
					Ip:           "2.2.2.2",
					CountryCode:  "cc2",
					Country:      "c2",
					City:         "city2",
					Latitude:     123.123,
					Longitude:    123.123,
					MysteryValue: 123,
				},
			},
			expectedCachedTotal: 2,
			expectedStoredTotal: 2,
		},
	}

	impCtx, impCancel := context.WithCancel(context.Background())
	defer impCancel()

	for name, suite := range testTable {
		t.Run(name, func(t *testing.T) {
			mockImporter := importer.NewInMemory(suite.given, 2)
			mockStorer := datastore.NewInMemory()
			mockCache := cache.NewInMemory()

			ldr := geo.NewLoader(mockImporter, mockStorer, mockCache)
			assert.NotNil(t, ldr)

			ldr.Load(impCtx, 1)

			cached := mockCache.All()
			assert.Equal(t, suite.expectedCachedTotal, len(cached))
			stored := mockStorer.All()
			assert.Equal(t, suite.expectedStoredTotal, len(stored))
		})
	}
}
