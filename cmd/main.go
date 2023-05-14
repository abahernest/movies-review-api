package main

import (
	"flag"
	"fmt"
	"log"
	httpDelivery "movies-review-api/delivery/http"
	port "movies-review-api/delivery/http"
	"movies-review-api/pkg/logger"
	"os"
)

func init() {

}

func main() {
	l, err := logger.InitLogger()

	if err != nil {
		panic(err)
	}

	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "dev"
	}

	l.Info(fmt.Sprintf("Loading %s env", env))

	l.Info("Starting App")

	httpConfig := httpDelivery.Config{}

	app := port.RunHttpServer(httpConfig)

	port := os.Getenv("PORT")

	if port == "" {
		port = "6001"
	}

	addr := flag.String("addr", fmt.Sprintf(":%s", port), "http service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))
}
