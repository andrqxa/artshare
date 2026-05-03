package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andrqxa/artshare/internal/configs"
	"github.com/andrqxa/artshare/internal/router"
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

	r := router.NewRouter()
	r.Mount()

	addr := ":" + cfg.Port
	fmt.Println("ArtShare API listening on", addr)
	return http.ListenAndServe(addr, r.Router())
}
