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
	csvPath    = flag.String("p", "cmd/loader/data_dump.csv", "path to csv file")
	connString = flag.String("c", defaultConnStr, "Database connection string")
	redisHost  = flag.String("r", "localhost", "Redis host")
	batch      = flag.Int("b", 400, "Batch size")
)

func main() {
	flag.Parse()

	impCtx, impCancel := context.WithCancel(context.Background())
	defer impCancel()

	ds := datastore.NewPg(db.PostgresDb(*connString))
	conf := redis.NewConfig()
	conf.Host = *redisHost
	cacheStore := redis.NewCache(conf)
	if err := cacheStore.Initialize(); err != nil {
		logrus.Fatal(err)
	}
	ldr := geo.NewLoader(importer.NewCsvImporter(*csvPath, *batch), ds, cacheStore)
	ldr.Load(impCtx)
}
