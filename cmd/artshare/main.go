package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andrqxa/artshare/internal/configs"
	artistcontroller "github.com/andrqxa/artshare/internal/controller/artist"
	artworkcontroller "github.com/andrqxa/artshare/internal/controller/artwork"
	exchangecontroller "github.com/andrqxa/artshare/internal/controller/exchange"
	likecontroller "github.com/andrqxa/artshare/internal/controller/like"
	usercontroller "github.com/andrqxa/artshare/internal/controller/user"
	artworkrepository "github.com/andrqxa/artshare/internal/repository/memory/artwork"
	exchangerepository "github.com/andrqxa/artshare/internal/repository/memory/exchange"
	likerepository "github.com/andrqxa/artshare/internal/repository/memory/like"
	userrepository "github.com/andrqxa/artshare/internal/repository/memory/user"
	"github.com/andrqxa/artshare/internal/repository/postgres"
	artistrepository "github.com/andrqxa/artshare/internal/repository/postgres/artist"
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
	if len(os.Args) > 1 {
		cfg.Port = os.Args[1]
	}

	postgresConnection, err := postgres.NewConnection(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer postgresConnection.Close()

	userRepository := userrepository.NewRepository()
	artistRepository := artistrepository.NewRepository(postgresConnection)
	artworkRepository := artworkrepository.NewRepository()
	likeRepository := likerepository.NewRepository()
	exchangeRepository := exchangerepository.NewRepository()

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
