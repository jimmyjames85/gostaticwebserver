package internal

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"text/template"
)

// TODO: hot reload route file
type RouteFile struct {
	fileloc string

	// json API
	Basedir string            `json:"basedir"`
	Routes  map[string]string `json:"routes"`
}

func (r *RouteFile) String() string {
	b, _ := json.Marshal(r)
	return string(b)

}
func (r *RouteFile) Fileloc() string { return r.fileloc }

func parseRouteFile(fileloc string) (*RouteFile, error) {
	f, err := os.OpenFile(fileloc, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	cfg := &RouteFile{fileloc: fileloc}
	err = json.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}

	for key, route := range cfg.Routes {
		b, err := parseTemplate(route, map[string]interface{}{"basedir": cfg.Basedir})
		if err != nil {
			panic(err)
		}
		cfg.Routes[key] = string(b)
	}

	return cfg, nil
}

func parseTemplate(text string, data map[string]interface{}) ([]byte, error) {
	// var helpers = map[string]interface{}{
	// 	"base64": func(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) },
	// 	"date":   func() string { return fmt.Sprintf("%d", time.Now().Unix()) },
	// }
	// t, err := template.New("gotemplate").Funcs(helpers).Parse(text)

	t, err := template.New("gotemplate").Parse(text)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
