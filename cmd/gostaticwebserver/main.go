package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

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
		log.Fatalf("unable to serve: %v\n", err)
	}
}

func newMain() {
	// cfg, err := LoadConfig()
	// if err != nil {
	// 	log.Fatalf("unable to load config: %s\n", err.Error())
	// }

	http.HandleFunc("/", rootHandler)

	fs := http.FileServer(http.Dir("/tmp"))
	fs = nil
	err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), fs)
	if err != nil {
		log.Fatalf("unable to serve: %v\n", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", requestInfo(r))

	rewrite, err := rewrite(r)
	if err != nil {
		fmt.Fprintf(w, "rewrite err: %s\n", err.Error())
		return
	}
	fmt.Fprintf(w, "rewrite: %s\n", rewrite)

}

func rewrite(r *http.Request) (string, error) {

	u, err := requestURL(r)
	if err != nil {
		return "", err
	}
	path := strings.Split(u.Host, ".")
	if len(path) <= 1 {
		return u.String(), nil
	}
	newpath := path[:len(path)-1] // TODO what if instead of localhost its jimmyjames.tld ... need len -2
	newhost := path[len(path)-1:] // TODO what if instead of localhost its jimmyjames.tld ... need len -2

	return fmt.Sprintf("%s://%s%s%s", u.Scheme, strings.Join(newhost, "."), u.Path, strings.Join(newpath, "/")), nil

	// better yet for FOO.jimmyjames.tld
	// 1. check if FOO dir exists
	// 2. yes? serve static fileserver location at that dir
	// 3. no? serve 404
}

func requestURL(r *http.Request) (*url.URL, error) {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	full := fmt.Sprintf("%s%s%s", scheme, r.Host, r.URL.Path)
	return url.Parse(full)
}

func requestInfo(r *http.Request) string {
	w := &bytes.Buffer{}
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	full := fmt.Sprintf("%s%s%s", scheme, r.Host, r.URL.Path)

	fmt.Fprintf(w, "Full: %s\n\n", full)
	fmt.Fprintf(w, "r.ContentLength: %d\n", r.ContentLength)
	fmt.Fprintf(w, "r.Form.Encode: %s\n", r.Form.Encode())
	fmt.Fprintf(w, "r.Host: %s\n", r.Host)
	fmt.Fprintf(w, "r.Method: %s\n", r.Method)
	fmt.Fprintf(w, "r.Proto: %s\n", r.Proto)
	fmt.Fprintf(w, "r.Referer: %s\n", r.Referer())
	fmt.Fprintf(w, "r.RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "r.RequestURI: %s\n", r.RequestURI)
	fmt.Fprintf(w, "r.TransferEncoding: %s\n", strings.Join(r.TransferEncoding, ","))
	fmt.Fprintf(w, "r.UserAgent: %s\n", r.UserAgent())

	u, err := url.Parse(full)
	if err != nil {
		fmt.Fprintf(w, "r.URL: %s\n", err.Error())
		return w.String()
	}
	fmt.Fprintf(w, "----------------------------------------------------------------------\n")
	fmt.Fprintf(w, "urlInfo\n")
	fmt.Fprintf(w, "%s", urlInfo(u))
	return w.String()
}

func urlInfo(u *url.URL) string {
	w := &bytes.Buffer{}
	fmt.Fprintf(w, "URL.EscapedFragment: %s\n", u.EscapedFragment())
	fmt.Fprintf(w, "URL.EscapedPath: %s\n", u.EscapedPath())
	fmt.Fprintf(w, "URL.ForceQuery: %t\n", u.ForceQuery)
	fmt.Fprintf(w, "URL.Fragment: %s\n", u.Fragment)
	fmt.Fprintf(w, "URL.Host: %s\n", u.Host)
	fmt.Fprintf(w, "URL.Hostname: %s\n", u.Hostname())
	fmt.Fprintf(w, "URL.IsAbs: %t\n", u.IsAbs())
	fmt.Fprintf(w, "URL.Opaque: %s\n", u.Opaque)
	fmt.Fprintf(w, "URL.Path: %s\n", u.Path)
	fmt.Fprintf(w, "URL.Port: %s\n", u.Port())
	fmt.Fprintf(w, "URL.Query.Encode: %s\n", u.Query().Encode())
	fmt.Fprintf(w, "URL.RawFragment: %s\n", u.RawFragment)
	fmt.Fprintf(w, "URL.RawPath: %s\n", u.RawPath)
	fmt.Fprintf(w, "URL.RawQuery: %s\n", u.RawQuery)
	fmt.Fprintf(w, "URL.Redacted: %s\n", u.Redacted())
	fmt.Fprintf(w, "URL.RemoteAddr: %s\n", u.RequestURI())
	fmt.Fprintf(w, "URL.Scheme: %s\n", u.Scheme)
	fmt.Fprintf(w, "URL.String: %s\n", u.String())
	return w.String()
}
