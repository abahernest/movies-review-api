package main

import (
	"flag"
	"fmt"
	"log"
	httpDelivery "movies-review-api/delivery/http"
	port "movies-review-api/delivery/http"
	"movies-review-api/domain"
	"movies-review-api/pkg/logger"
	"movies-review-api/repository/mongodb"
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

	domain.GetSecrets(l)

	repo := mongodb.New(l)

	httpConfig := httpDelivery.Config{
		UserRepo:    repo.UserRepo,
		FilmRepo:    repo.FilmRepo,
		CommentRepo: repo.CommentRepo,
	}

	app := port.RunHttpServer(httpConfig)

	port := os.Getenv("PORT")

	if port == "" {
		port = "6001"
	}

	addr := flag.String("addr", fmt.Sprintf(":%s", port), "http service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))
}
