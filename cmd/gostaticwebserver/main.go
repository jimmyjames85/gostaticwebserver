package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/net/publicsuffix"
)

type Config struct {
	Port   int    `envconfig:"PORT" required:"false" default:"8080"`
	WebDir string `envconfig:"WEBDIR" required:"false"` // TODO remove
}

func main() {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)

	http.HandleFunc("/", rootHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
	if err != nil {
		log.Fatalf("unable to serve: %v\n", err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	// TODO serve custom 404 and LOG ip addr
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 yo\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	handler, err := lookupRouteHandler(r)
	if err != nil {
		notFoundHandler(w, r)
		return
	}
	handler.ServeHTTP(w, r)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", requestInfo(r))
}

var basedir = "/tmp"

// dnsPrefix -> http.HandlerFunc
//
// if request is for mail.foo.com then dnsPrefix is mail
//
var routeHandlers = map[string]http.HandlerFunc{
	"":     http.FileServer(http.Dir(basedir)).ServeHTTP,
	"dir1": http.FileServer(http.Dir(basedir + "/dir1")).ServeHTTP,
	"dir2": http.FileServer(http.Dir(basedir + "/dir2")).ServeHTTP,
	"info": infoHandler,
	// "mail": http.RedirectHandler("http://gmail.com", http.StatusTemporaryRedirect).ServeHTTP,
}

func parseDNSPrefix(u *url.URL) (string, error) {
	hostname := strings.ToLower(u.Hostname())

	if net.ParseIP(hostname) != nil {
		// hostname is an ip address
		return "", fmt.Errorf("no prefix for IP address")
	}

	tldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return "", err // probably localhost
	}

	prefix := strings.TrimSuffix(hostname, tldPlusOne)
	prefix = strings.TrimSuffix(prefix, ".")
	return prefix, nil
}

func lookupRouteHandler(r *http.Request) (http.Handler, error) {
	u, err := parseRequestURL(r)
	if err != nil {
		return nil, err
	}

	dnsPrefix, err := parseDNSPrefix(u)
	if err != nil {
		return nil, err
	}

	handler, ok := routeHandlers[dnsPrefix]
	if !ok {
		return nil, fmt.Errorf("no fileserver found for %s", dnsPrefix)
	}

	return handler, nil
}

func parseRequestURL(r *http.Request) (*url.URL, error) {
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
