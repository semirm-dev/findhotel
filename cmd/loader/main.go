package main

import (
	"context"
	"flag"
	"github.com/semirm-dev/findhotel/geo"
	"github.com/semirm-dev/findhotel/importer"
	"github.com/sirupsen/logrus"
)

var (
	path = flag.String("path ", "cmd/loader/data_dump_part.csv", "")
)

func main() {
	flag.Parse()

	impCtx, impCancel := context.WithCancel(context.Background())
	defer impCancel()

	ldr := geo.NewLoader(importer.NewCsvImporter(*path))
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

	logrus.Info("listening for messages...")

	<-ldrFinished
}
