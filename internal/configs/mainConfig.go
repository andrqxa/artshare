package configs

import "os"

type Config struct {
	Port string
}

func Read() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port: port,
	}, nil
}
