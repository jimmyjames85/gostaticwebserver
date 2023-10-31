package internal

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/publicsuffix"
)

type logger interface {
	Printf(format string, a ...interface{})
	Event(map[string]interface{})
}

type Server struct {
	log

	RouteFileloc string
	Port         int
	SSLPort      int    // if 0 no ssl is served
	CertDir      string // location to store ssl cert cache, if empty no ssl is served

	// dnsPrefix -> http.HandlerFunc
	//
	// if request is for mail.foo.com then dnsPrefix is mail
	//
	routes map[string]http.HandlerFunc
}

// func requestInfoHandler(w http.ResponseWriter, r *http.Request) {
// 	m := requestInfo(r)
// 	b, _ := json.MarshalIndent(m, "", " ")
// 	fmt.Fprintf(w, "%s\n", string(b))
// }
//
// func urlInfoHandler(w http.ResponseWriter, r *http.Request) {
// 	m := urlInfo(r.URL)
// 	b, _ := json.MarshalIndent(m, "", " ")
// 	fmt.Fprintf(w, "%s\n", string(b))
// }

func (s *Server) initRoutes() error {
	routeFile, err := parseRouteFile(s.RouteFileloc)
	if err != nil {
		return err
	}
	s.routes = map[string]http.HandlerFunc{}
	// TODO move this onto routeFile struct as a func when we implement hot reload option
	for dnsPrefix, fileloc := range routeFile.Routes {
		s.routes[dnsPrefix] = http.FileServer(http.Dir(fileloc)).ServeHTTP
	}
	s.log.Printf("using routes: %s", routeFile.String())
	return nil
}

func (s *Server) Serve() error {
	err := s.initRoutes()
	if err != nil {
		return err
	}

	addr := fmt.Sprintf(":%d", s.Port)
	mux := http.NewServeMux() // TODO add middleware here
	// http.HandleFunc("/request/info", requestInfoHandler)
	// http.HandleFunc("/url/info", urlInfoHandler)
	mux.HandleFunc("/", s.rootHandler)

	s.log.Printf("webserver started on port: %d", s.Port)

	if len(s.CertDir) == 0 || s.SSLPort == 0 { // no SSL
		return http.ListenAndServe(addr, mux)
	}

	s.log.Printf("SSL webserver started on port: %d", s.SSLPort)
	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(s.CertDir),
	}

	tlsAddr := fmt.Sprintf(":%d", s.SSLPort)
	server := &http.Server{
		Addr:    tlsAddr,
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
				// refuse requests for which we don't have routes for
				dnsPrefix, err := parseDNSPrefix(&url.URL{Host: hello.ServerName})
				if err != nil {
					return nil, err
				}
				if _, ok := s.routes[dnsPrefix]; !ok {
					return nil, fmt.Errorf("invalid domain: %s", hello.ServerName)
				}

				ret, err := certManager.GetCertificate(hello)
				if err != nil {
					s.log.Event(map[string]interface{}{
						"error":      err.Error(),
						"line":       "getCertCall(hello)",
						"serverName": hello.ServerName,
					})
					return nil, err
				}

				return ret, nil
			},
		},
	}

	go http.ListenAndServe(addr, certManager.HTTPHandler(nil))
	return server.ListenAndServeTLS("", "")
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

	// routes foo.website.com to website.com/foo
	dnsPrefix, err := parseDNSPrefix(u)
	if err != nil {
		// uncomment useful for local development
		// dnsPrefix = "" // this causes it to route to root handler: useful for localhost

		// don't serve
		s.log.Event(map[string]interface{}{
			"event":   "error",
			"message": fmt.Sprintf("parseDNSPrefix: %s: routing to root handler", err.Error()),
		})
		return nil, err
	}

	handler, ok := s.routes[dnsPrefix]
	if !ok {
		return nil, fmt.Errorf("no route found for %s", dnsPrefix)
	}

	return handler, nil
}

// returns an error if url is an ip address or simply localhost
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
