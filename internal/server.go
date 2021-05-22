package internal

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type Server struct {
	routeFile *RouteFile
	port      int

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
	}

	// TODO move this onto routeFile struct as a func when we implement hot reload option
	for dnsPrefix, fileloc := range w.routeFile.Routes {
		w.routes[dnsPrefix] = http.FileServer(http.Dir(fileloc)).ServeHTTP
	}

	return w, nil

}

func (s *Server) Serve() error {
	http.HandleFunc("/", s.handleRoot)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
	// TODO close everything
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
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

// func infoHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "%s", requestInfo(r))
// }
//
// func requestInfo(r *http.Request) string {
//
// 	w := &bytes.Buffer{}
// 	scheme := "http://"
// 	if r.TLS != nil {
// 		scheme = "https://"
// 	}
// 	full := fmt.Sprintf("%s%s%s", scheme, r.Host, r.URL.Path)
//
// 	fmt.Fprintf(w, "Full: %s\n\n", full)
// 	fmt.Fprintf(w, "r.ContentLength: %d\n", r.ContentLength)
// 	fmt.Fprintf(w, "r.Form.Encode: %s\n", r.Form.Encode())
// 	fmt.Fprintf(w, "r.Host: %s\n", r.Host)
// 	fmt.Fprintf(w, "r.Method: %s\n", r.Method)
// 	fmt.Fprintf(w, "r.Proto: %s\n", r.Proto)
// 	fmt.Fprintf(w, "r.Referer: %s\n", r.Referer())
// 	fmt.Fprintf(w, "r.RemoteAddr: %s\n", r.RemoteAddr)
// 	fmt.Fprintf(w, "r.RequestURI: %s\n", r.RequestURI)
// 	fmt.Fprintf(w, "r.TransferEncoding: %s\n", strings.Join(r.TransferEncoding, ","))
// 	fmt.Fprintf(w, "r.UserAgent: %s\n", r.UserAgent())
//
// 	u, err := url.Parse(full)
// 	if err != nil {
// 		fmt.Fprintf(w, "r.URL: %s\n", err.Error())
// 		return w.String()
// 	}
// 	fmt.Fprintf(w, "----------------------------------------------------------------------\n")
// 	fmt.Fprintf(w, "urlInfo\n")
// 	fmt.Fprintf(w, "%s", urlInfo(u))
// 	return w.String()
// }
//
// func urlInfo(u *url.URL) string {
// 	w := &bytes.Buffer{}
// 	fmt.Fprintf(w, "URL.EscapedFragment: %s\n", u.EscapedFragment())
// 	fmt.Fprintf(w, "URL.EscapedPath: %s\n", u.EscapedPath())
// 	fmt.Fprintf(w, "URL.ForceQuery: %t\n", u.ForceQuery)
// 	fmt.Fprintf(w, "URL.Fragment: %s\n", u.Fragment)
// 	fmt.Fprintf(w, "URL.Host: %s\n", u.Host)
// 	fmt.Fprintf(w, "URL.Hostname: %s\n", u.Hostname())
// 	fmt.Fprintf(w, "URL.IsAbs: %t\n", u.IsAbs())
// 	fmt.Fprintf(w, "URL.Opaque: %s\n", u.Opaque)
// 	fmt.Fprintf(w, "URL.Path: %s\n", u.Path)
// 	fmt.Fprintf(w, "URL.Port: %s\n", u.Port())
// 	fmt.Fprintf(w, "URL.Query.Encode: %s\n", u.Query().Encode())
// 	fmt.Fprintf(w, "URL.RawFragment: %s\n", u.RawFragment)
// 	fmt.Fprintf(w, "URL.RawPath: %s\n", u.RawPath)
// 	fmt.Fprintf(w, "URL.RawQuery: %s\n", u.RawQuery)
// 	fmt.Fprintf(w, "URL.Redacted: %s\n", u.Redacted())
// 	fmt.Fprintf(w, "URL.RemoteAddr: %s\n", u.RequestURI())
// 	fmt.Fprintf(w, "URL.Scheme: %s\n", u.Scheme)
// 	fmt.Fprintf(w, "URL.String: %s\n", u.String())
// 	return w.String()
// }
