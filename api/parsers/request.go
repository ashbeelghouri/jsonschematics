package parsers

import (
	"encoding/json"
	"github.com/ashbeelghouri/jsonschematics"
	"io"
	"net/http"
	"strings"
)

func ParseRequest(r *http.Request) (map[string]interface{}, error) {
	var headers map[string]string
	for key, values := range r.Header {
		headers[key] = values[0]
	}
	var body map[string]interface{}
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bodyBytes, &body)
		if err != nil {
			return nil, err
		}
	}
	body = jsonschematics.DeflateMap(body, ".")
	splitPath := strings.Split(r.RequestURI, "?")
	// get query parameters
	query := map[string]interface{}{}
	if splitPath[1] != "" {
		for _, param := range strings.Split(splitPath[1], "&") {
			kv := strings.Split(param, "=")
			if len(kv) == 2 {
				query[kv[0]] = kv[1]
			}
		}
	}
	// already in the FLAT mode
	return map[string]interface{}{
		"headers": headers,
		"body":    body,
		"path":    splitPath[0],
		"method":  r.Method,
		"query":   query,
	}, nil
}
