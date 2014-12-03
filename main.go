package gojson

import (
	"github.com/codegangsta/martini"
	"net/http"
)

func init() {
	m := martini.Classic()
	m.Post("/cnv", Convertjson)
	http.Handle("/", m)
}

func Convertjson(rw http.ResponseWriter, req *http.Request) string {
	if output, err := generate(req.Body, "Foo", "main"); err != nil {
		return "Error!, Is that json?"
	} else {
		return string(output)
	}
}
