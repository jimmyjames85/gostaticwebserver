package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port   int    `envconfig:"PORT" required:"false" default:"8080"`
	WebDir string `envconfig:"WEBDIR" required:"true"`
}

// LoadConfig loads environment variables
func LoadConfig() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s\n", err.Error())
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), http.FileServer(http.Dir(cfg.WebDir)))
	if err != nil {
		log.Fatalf("unable to server: %v\n", err)
	}
}
