package main

import (
	"fmt"
	"os"

	"github.com/jimmyjames85/gostaticwebserver/internal"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port      int    `envconfig:"PORT" required:"false" default:"8080"`
	SSLPort   int    `envconfig:"SSL_PORT" required:"false" default:"0"`
	CertDir   string `envconfig:"CERT_DIR" required:"false" default:""`
	RouteFile string `envconfig:"ROUTE_FILE" required:"true"`
}

func main() {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)

	s := &internal.Server{
		RouteFileloc: cfg.RouteFile,
		Port:         cfg.Port,
		SSLPort:      cfg.SSLPort,
		CertDir:      cfg.CertDir,
	}

	err := s.Serve()
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
