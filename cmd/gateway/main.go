package main

import (
	"flag"
	"github.com/semirm-dev/findhotel/gateway"
	"github.com/semirm-dev/findhotel/geo"
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

	var srch geo.Search

	router.GET("geo", gateway.GetProducts(srch))

	web.ServeHttp(*httpAddr, "gateway", router)
}
