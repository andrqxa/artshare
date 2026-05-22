package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/andrqxa/artshare/internal/configs"
	artistcontroller "github.com/andrqxa/artshare/internal/controller/artist"
	artworkcontroller "github.com/andrqxa/artshare/internal/controller/artwork"
	exchangecontroller "github.com/andrqxa/artshare/internal/controller/exchange"
	likecontroller "github.com/andrqxa/artshare/internal/controller/like"
	usercontroller "github.com/andrqxa/artshare/internal/controller/user"
	"github.com/andrqxa/artshare/internal/repository/postgres"
	artistrepository "github.com/andrqxa/artshare/internal/repository/postgres/artist"
	artworkrepository "github.com/andrqxa/artshare/internal/repository/postgres/artwork"
	exchangerepository "github.com/andrqxa/artshare/internal/repository/postgres/exchange"
	likerepository "github.com/andrqxa/artshare/internal/repository/postgres/like"
	userrepository "github.com/andrqxa/artshare/internal/repository/postgres/user"
	"github.com/andrqxa/artshare/internal/router"
	artistrouter "github.com/andrqxa/artshare/internal/router/artist"
	artworkrouter "github.com/andrqxa/artshare/internal/router/artwork"
	exchangerouter "github.com/andrqxa/artshare/internal/router/exchange"
	likerouter "github.com/andrqxa/artshare/internal/router/like"
	userrouter "github.com/andrqxa/artshare/internal/router/user"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := configs.Read()
	if err != nil {
		return err
	}

	postgresConnection, err := postgres.NewConnection(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer postgresConnection.Close()

	userRepository := userrepository.NewRepository(postgresConnection)
	artistRepository := artistrepository.NewRepository(postgresConnection)
	artworkRepository := artworkrepository.NewRepository(postgresConnection)
	likeRepository := likerepository.NewRepository(postgresConnection)
	exchangeRepository := exchangerepository.NewRepository(postgresConnection)

	userController := usercontroller.NewController(userRepository)
	artistController := artistcontroller.NewController(artistRepository)
	artworkController := artworkcontroller.NewController(artworkRepository)
	likeController := likecontroller.NewController(likeRepository)
	exchangeController := exchangecontroller.NewController(exchangeRepository)

	r := router.NewRouter(
		userrouter.NewHandler(userController),
		artistrouter.NewHandler(artistController),
		artworkrouter.NewHandler(artworkController),
		likerouter.NewHandler(likeController),
		exchangerouter.NewHandler(exchangeController),
	)
	r.Mount()

	addr := ":" + cfg.Port
	fmt.Println("ArtShare API listening on", addr)
	return http.ListenAndServe(addr, r.Router())
}
