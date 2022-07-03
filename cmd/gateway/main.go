package main

import (
	"context"
	"flag"
	"github.com/semirm-dev/findhotel/datastore"
	"github.com/semirm-dev/findhotel/gateway"
	"github.com/semirm-dev/findhotel/geo"
	"github.com/semirm-dev/findhotel/importer"
	"github.com/semirm-dev/findhotel/internal/web"
	"github.com/sirupsen/logrus"
)

const defaultConnStr = "host=localhost port=5432 dbname=findhotel_geo user=postgres password=postgres sslmode=disable"

var (
	httpAddr   = flag.String("http", ":8000", "Http address")
	connString = flag.String("connStr", defaultConnStr, "Condition Service connection string")
	path       = flag.String("path ", "cmd/loader/data_dump.csv", "path to csv file")
)

func main() {
	flag.Parse()

	impCtx, impCancel := context.WithCancel(context.Background())
	defer impCancel()

	ds := datastore.NewInMemoryStore()
	ldr := geo.NewLoader(importer.NewCsvImporter(*path), ds)
	ldrFinished := ldr.Load(impCtx)

	go func() {
		defer logrus.Warn("geo loader finished")

		for {
			select {
			case <-ldrFinished:
				return
			}
		}
	}()

	router := web.NewRouter()

	router.GET("geo", gateway.GetGeoLocation(ds))

	web.ServeHttp(*httpAddr, "gateway", router)
}
