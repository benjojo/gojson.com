package main

import (
	"github.com/codegangsta/martini"
	"net/http"
)

func main() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Post("/cnv", Convertjson)
	m.Run()
}

func Convertjson(rw http.ResponseWriter, req *http.Request) string {
	if output, err := generate(req.Body, "Foo", "main"); err != nil {
		return "Error!, Is that json?"
	} else {
		return string(output)
	}
}
