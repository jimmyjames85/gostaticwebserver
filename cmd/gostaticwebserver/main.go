package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type userPayload struct {
	// NOTE: The tags `json` specify what to parse. This means the
	// incoming JSON payload needs to specify "userid" and not
	// ID. However, when the `json` tag is omitted the variable
	// name is what is expected, and is case insensitive.
	ID   int `json:"userid"`
	Info string
	Name string
}

// handleParseJSON expects a payload akin to the following curl calls
//
// curl localhost:8080/parse/json -d '{"userid":123,"NAME":"jimbo", "Info":"extra info"}'
func handleParseJSON(w http.ResponseWriter, r *http.Request) {
	byts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("unable to read request body: %v\n", err) // write to our internal logs
		w.WriteHeader(http.StatusInternalServerError)        // let client know we have an unkown error
		w.Write([]byte("Internal Server Error"))             // let client know we have an unkown error
		return                                               // stop proccessing
	}

	log.Printf("received: %s\n", string(byts))
	log.Printf("attempting to unmarshal request into struct")

	var user userPayload
	err = json.Unmarshal(byts, &user)
	if err != nil {
		log.Printf("unable to unmarshal incoming payload: %v: %s", err, string(byts)) // save to our internal logs
		w.WriteHeader(http.StatusBadRequest)                                          // let client know they sent us a bad payload
		w.Write([]byte("bad request: expecting a JSON payload\n"))                    // let client know they sent us a bad payload
		return                                                                        // stop proccessing
	}

	log.Printf("JSON unmarshalled successfully: %v", user)

	bw := bufio.NewWriter(w)
	bw.WriteString(fmt.Sprintf("  user.ID: %d\n", user.ID))
	bw.WriteString(fmt.Sprintf("user.Name: %s\n", user.Name))
	bw.WriteString(fmt.Sprintf("user.Info: %s\n", user.Info))
	bw.Flush()
}

// handleParseParameters expects a payload akin to the following curl call
//
// curl localhost:8080/parse/parameters -d foo="foo1" -d foo="foo2" -d foo=3 -d bar="hello world"
func handleParseParameters(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// this should never happen
		log.Printf("unable to parse form: %v\n", err) // write to our internal logs
		w.WriteHeader(http.StatusInternalServerError) // let client know we have an unkown error
		w.Write([]byte("Internal Server Error"))      // let client know we have an unkown error
		return                                        // stop proccessing
	}

	bw := bufio.NewWriter(w)
	for name, vals := range r.Form {
		bw.WriteString(fmt.Sprintf("%s\n", name))
		for _, val := range vals {
			bw.WriteString(fmt.Sprintf("\t%s\n", val))
		}
		bw.WriteString("\n")
	}
	bw.Flush()
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s\n", err.Error())
	}

	http.HandleFunc("/parse/json", handleParseJSON)
	http.HandleFunc("/parse/parameters", handleParseParameters)

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
	if err != nil {
		log.Fatalf("unable to serve: %v\n", err)
	}
}
