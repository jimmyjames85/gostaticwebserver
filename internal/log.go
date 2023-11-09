package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	l.Printf("%s\n", string(b))
}

func (*log) Printf(format string, a ...interface{}) {
	// fmt.Printf("[%s] ", time.Now().UTC().Format(time.RFC822))
	fmt.Printf(format, a...)
	fmt.Printf("\n")
}

func requestInfo(r *http.Request) map[string]interface{} {
	b, err := ioutil.ReadAll(r.Body)
	body := string(b)
	if err != nil {
		body = fmt.Sprintf("body read err: %s\n%s", err.Error(), body)
	}

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
		"body":              body,
	}
}

func urlInfo(u *url.URL) map[string]interface{} {
	return map[string]interface{}{
		"EscapedFragment": u.EscapedFragment(),
		"EscapedPath":     u.EscapedPath(),
		"ForceQuery":      u.ForceQuery,
		"Fragment":        u.Fragment,
		"Host":            u.Host,
		"Hostname":        u.Hostname(),
		"IsAbs":           u.IsAbs(),
		"Opaque":          u.Opaque,
		"Path":            u.Path,
		"Port":            u.Port(),
		"Query.Encode":    u.Query().Encode(),
		"RawFragment":     u.RawFragment,
		"RawPath":         u.RawPath,
		"RawQuery":        u.RawQuery,
		"Redacted":        u.Redacted(),
		"RemoteAddr":      u.RequestURI(),
		"Scheme":          u.Scheme,
		"String":          u.String(),
	}
}
