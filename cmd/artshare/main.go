package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andrqxa/artshare/internal/configs"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <port>")
		os.Exit(1)
	}
	port := os.Args[1]

	if err := run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello, artshare!, your port is ", port)

	os.Exit(0)
}

func run() error {
	// read config from env
	_ = configs.Read()

	return nil
}
