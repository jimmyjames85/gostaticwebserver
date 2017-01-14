package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"fmt"
)

var (
	host string
	webdir string
	port int
)

func main() {
	host = os.Getenv("HOST")
	if len(host) == 0 {
		log.Fatalf("could not read environment variable HOST\n")
	}

	var err error
	port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("could not read environment variable PORT: %v\n", err)
	}

	webdir = os.Getenv("WEBDIR")
	if len(webdir) == 0 {
		log.Fatalf("could not read environment variable WEBDIR\n")
	}
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(webdir)))
	if err != nil {
		log.Fatalf("unable to server: %v\n", err)
	}
}
