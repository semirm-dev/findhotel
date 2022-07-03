package main

import (
	"context"
	"flag"
	"github.com/semirm-dev/findhotel/gateway"
	"github.com/semirm-dev/findhotel/geo"
	"github.com/semirm-dev/findhotel/internal/web"
	"github.com/sirupsen/logrus"
)

const defaultConnStr = "host=localhost port=5432 dbname=findhotel_geo user=postgres password=postgres sslmode=disable"

var (
	httpAddr   = flag.String("http", ":8000", "Http address")
	connString = flag.String("connStr", defaultConnStr, "Condition Service connection string")
	path       = flag.String("path ", "cmd/importer/data_dump_part.csv", "")
)

func main() {
	flag.Parse()

	impCtx, impCancel := context.WithCancel(context.Background())
	defer impCancel()

	importer := geo.NewCsvImporter(*path)
	ldr := geo.NewLoader(importer)
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

	var srch geo.Search

	router.GET("geo", gateway.GetProducts(srch))

	web.ServeHttp(*httpAddr, "gateway", router)
}
