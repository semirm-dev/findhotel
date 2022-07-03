package main

import (
	"context"
	"flag"
	"github.com/semirm-dev/findhotel/geo"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var (
	path = flag.String("path ", "cmd/importer/data_dump_part.csv", "")
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

	logrus.Info("listening for messages...")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
