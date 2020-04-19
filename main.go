package main

import (
	"net/http"

	"github.com/codegangsta/martini"
	"google.golang.org/appengine"
)

func main() {
	m := martini.Classic()
	m.Post("/cnv", Convertjson)
	http.Handle("/", m)
	appengine.Main()
}

func Convertjson(rw http.ResponseWriter, req *http.Request) string {
	if output, err := generate(req.Body, "Foo", "main"); err != nil {
		return "Error!, Is that json?"
	} else {
		return string(output)
	}
}
