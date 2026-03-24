package rest

import (
	"net/http"
	"strings"
)

func AliasPrefixHandler(fromPrefix, toPrefix string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, fromPrefix) {
			rewritten := r.Clone(r.Context())
			newURL := *r.URL
			newPath := toPrefix + strings.TrimPrefix(r.URL.Path, fromPrefix)
			if newPath == "" {
				newPath = "/"
			}
			newURL.Path = newPath
			rewritten.URL = &newURL
			next.ServeHTTP(w, rewritten)
			return
		}

		next.ServeHTTP(w, r)
	})
}
