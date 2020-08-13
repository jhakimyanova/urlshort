package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jhakimyanova/urlshort/pkg/handlers"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := handlers.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	// TODO: accept yaml file as a flag:
	yamlFile := "pathsToUrls.yaml"
	yamlData, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		log.Fatal(err)
	}
	yamlHandler, err := handlers.YAMLHandler(yamlData, mapHandler)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
