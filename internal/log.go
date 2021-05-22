package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type log struct{}

func (l *log) Event(m map[string]interface{}) {
	_, ok := m["processed"]
	if !ok {
		m["processed"] = time.Now().Unix()
	}

	b, err := json.Marshal(m)
	if err != nil {
		l.Printf("failed to log: %#v", m)
		return
	}
	l.Printf(string(b))
}

func (*log) Printf(format string, a ...interface{}) {
	// fmt.Printf("[%s] ", time.Now().UTC().Format(time.RFC822))
	fmt.Printf(format, a...)
	fmt.Printf("\n")
}

func requestInfo(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"content_length":    r.ContentLength,
		"form.encode":       r.Form.Encode(),
		"host":              r.Host,
		"method":            r.Method,
		"proto":             r.Proto,
		"referer":           r.Referer(),
		"remote_addr":       r.RemoteAddr,
		"request_uri":       r.RequestURI,
		"transfer_encoding": r.TransferEncoding,
		"user_agent":        r.UserAgent(),
	}
}

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
