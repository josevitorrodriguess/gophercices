package urlshortner

import (
	"encoding/json"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathToUrls := buildYAMLMap(pathUrls)
	return MapHandler(pathToUrls, fallback), nil
}

func buildYAMLMap(pathUrls []pathUrlYAML) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.URL
	}
	return pathToUrls
}

func parseYaml(data []byte) ([]pathUrlYAML, error) {
	var pathUrls []pathUrlYAML
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

type pathUrlYAML struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func JSONHandler(jsonByes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseJson(jsonByes)
	if err != nil {
		return nil, err
	}
	pathToUrls := buildJSONMap(pathUrls)
	return MapHandler(pathToUrls, fallback), nil
}

func buildJSONMap(pathUrls []pathUrlJSON) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.URL
	}
	return pathToUrls
} 	

func parseJson(data []byte) ([]pathUrlJSON, error) {
	var pathUrls []pathUrlJSON
	err := json.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

type pathUrlJSON struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}
