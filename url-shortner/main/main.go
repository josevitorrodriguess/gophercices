package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	urlshortner "github.com/josevitorrodriguess/url-shortner"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortner.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := flag.String("yaml", "urls.yaml", " a yaml file with the urls")
	json := flag.String("json", "urls.json", "a json file with the urls")
	flag.Parse()
	contentYAML, _ := readFile(yaml)

	yamlHandler, err := urlshortner.YAMLHandler([]byte(contentYAML), mapHandler)
	if err != nil {
		panic(err)
	}

	//JSONHandler
	contentJSON, _ := readFile(json)

	jsonHandler, err := urlshortner.JSONHandler([]byte(contentJSON), yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultUndef)
	return mux
}

func defaultUndef(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "default or invalid url")
}

func readFile(yaml *string) ([]byte, error) {
	file, err := os.Open(*yaml)
	if err != nil {
		panic(err)
		return nil, err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
		return nil, err
	}
	return content, nil
}
