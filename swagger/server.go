package swagger

import (
	"net/http"
	"strings"
)

func WithSwagger(swaggerFile []byte, next http.Handler) http.Handler {
	static := http.FileServer(http.FS(data))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, basePath) {
			if r.URL.Path == basePath+"/swagger.json" {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("content-type", "application/json")
				// todo: set no-cache header
				_, _ = w.Write(swaggerFile)
				return
			}
			// todo: set cache header
			static.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

const basePath = "/swagger"
