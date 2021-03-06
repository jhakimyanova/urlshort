package handlers

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// pathURL is representing a member of the list in provided YAML file
type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsUrls, err := parseYAML(yamlData)
	if err != nil {
		return nil, err
	}
	pathURLMap := buildMap(pathsUrls)
	return MapHandler(pathURLMap, fallback), nil
}

func parseYAML(yamlData []byte) ([]pathURL, error) {
	pathsUrls := []pathURL{}
	err := yaml.Unmarshal(yamlData, &pathsUrls)
	// Make sure that the first part returned is nil in case of error:
	if err != nil {
		return nil, err
	}
	return pathsUrls, nil
}

func buildMap(pathsUrls []pathURL) map[string]string {
	pathURLMap := make(map[string]string)
	for _, pu := range pathsUrls {
		pathURLMap[pu.Path] = pu.URL
	}
	return pathURLMap
}
