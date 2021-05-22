package main

import (
	"fmt"
	"os"

	"github.com/jimmyjames85/gostaticwebserver/internal"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port      int    `envconfig:"PORT" required:"false" default:"8080"`
	RouteFile string `envconfig:"ROUTE_FILE" required:"true"`
}

func main() {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)

	s, err := internal.NewServer(cfg.Port, cfg.RouteFile)
	if err != nil {
		exitf(-1, "%s\n", err.Error())
	}

	err = s.Serve()
	if err != nil {
		exitf(-1, "%s\n", err.Error())
	}
}

func exitf(code int, format string, a ...interface{}) {
	w := os.Stderr
	if code == 0 {
		w = os.Stdout
	}

	fmt.Fprintf(w, format, a...)
	os.Exit(code)
}
