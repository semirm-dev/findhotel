package main

import (
	"context"
	"flag"
	"github.com/semirm-dev/findhotel/cache/redis"
	"github.com/semirm-dev/findhotel/datastore"
	"github.com/semirm-dev/findhotel/geo"
	"github.com/semirm-dev/findhotel/importer"
	"github.com/semirm-dev/findhotel/internal/db"
	"github.com/sirupsen/logrus"
)

const defaultConnStr = "host=localhost port=5432 dbname=findhotel_geo user=postgres password=postgres sslmode=disable"

var (
	path       = flag.String("path ", "cmd/loader/data_dump_med.csv", "path to csv file")
	connString = flag.String("connStr", defaultConnStr, "Condition Service connection string")
)

func main() {
	flag.Parse()

	impCtx, impCancel := context.WithCancel(context.Background())
	defer impCancel()

	ds := datastore.NewPg(db.PostgresDb(*connString))
	conf := redis.NewConfig()
	cacheStore := redis.NewCache(conf)
	if err := cacheStore.Initialize(); err != nil {
		logrus.Fatal(err)
	}
	//cacheStore := c2go.NewC2Go("geo_data")
	ldr := geo.NewLoader(importer.NewCsvImporter(*path, 100), ds, cacheStore)
	ldr.Load(impCtx)
}
