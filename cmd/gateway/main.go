package main

import (
	"flag"
	"github.com/semirm-dev/findhotel/datastore"
	"github.com/semirm-dev/findhotel/gateway"
	"github.com/semirm-dev/findhotel/internal/db"
	"github.com/semirm-dev/findhotel/internal/web"
)

const defaultConnStr = "host=localhost port=5432 dbname=findhotel_geo user=postgres password=postgres sslmode=disable"

var (
	httpAddr   = flag.String("http", ":8000", "Http address")
	connString = flag.String("connStr", defaultConnStr, "Condition Service connection string")
)

func main() {
	flag.Parse()

	router := web.NewRouter()

	router.GET("geo", gateway.GetGeoLocation(datastore.NewPgStore(db.PostgresDb(*connString))))

	web.ServeHttp(*httpAddr, "gateway", router)
}
