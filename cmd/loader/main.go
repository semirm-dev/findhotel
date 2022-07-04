package main

import (
	"context"
	"flag"
	"github.com/semirm-dev/findhotel/datastore"
	"github.com/semirm-dev/findhotel/geo"
	"github.com/semirm-dev/findhotel/importer"
)

const defaultConnStr = "host=localhost port=5432 dbname=findhotel_geo user=postgres password=postgres sslmode=disable"

var (
	path       = flag.String("path ", "cmd/loader/data_dump.csv", "path to csv file")
	connString = flag.String("connStr", defaultConnStr, "Condition Service connection string")
)

func main() {
	flag.Parse()

	impCtx, impCancel := context.WithCancel(context.Background())
	defer impCancel()

	//ds := datastore.NewPgStore(db.PostgresDb(*connString))
	ds := datastore.NewInMemoryStore()
	ldr := geo.NewLoader(importer.NewCsvImporter(*path, 100), ds)
	ldr.Load(impCtx)
}
