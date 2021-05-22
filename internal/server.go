package internal

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type logger interface {
	Printf(format string, a ...interface{})
	Event(map[string]interface{})
}

type Server struct {
	routeFile *RouteFile
	port      int

	log logger
	// dnsPrefix -> http.HandlerFunc
	//
	// if request is for mail.foo.com then dnsPrefix is mail
	//
	routes map[string]http.HandlerFunc
}

func NewServer(port int, routeFile string) (*Server, error) {
	cfg, err := parseRouteFile(routeFile)
	if err != nil {
		return nil, err
	}

	w := &Server{
		port:      port,
		routes:    map[string]http.HandlerFunc{},
		routeFile: cfg,
		log:       &log{},
	}

	// TODO move this onto routeFile struct as a func when we implement hot reload option
	for dnsPrefix, fileloc := range w.routeFile.Routes {
		w.routes[dnsPrefix] = http.FileServer(http.Dir(fileloc)).ServeHTTP
	}

	return w, nil

}

// func infoHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "%s", requestInfo(r))
// }

func (s *Server) Serve() error {
	// TODO middle ware here

	// http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/", s.rootHandler)

	s.log.Printf("webserver start on port: %d", s.port)
	s.log.Printf("%s", s.routeFile.String())

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
	// TODO close everything
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Event(requestInfo(r))

	h, err := s.routeRequest(r)
	if err != nil {
		handleNotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	// TODO serve custom 404 and LOG ip addr
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 yo\n")
}

func (s *Server) routeRequest(r *http.Request) (http.Handler, error) {
	u, err := parseRequestURL(r)
	if err != nil {
		return nil, err
	}

	dnsPrefix, err := parseDNSPrefix(u)
	if err != nil {
		return nil, err
	}

	handler, ok := s.routes[dnsPrefix]
	if !ok {
		return nil, fmt.Errorf("no route found for %s", dnsPrefix)
	}

	return handler, nil
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

func parseRequestURL(r *http.Request) (*url.URL, error) {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	full := fmt.Sprintf("%s%s%s", scheme, r.Host, r.URL.Path)
	return url.Parse(full)
}
