package swagger

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func WithSwagger(swaggerFile []byte, next http.Handler) http.Handler {
	static := http.FileServer(http.FS(data))
	cacheControl := fmt.Sprintf("max-age=%v", (time.Hour * 24 * 30).Seconds())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, basePath) {
			if r.URL.Path == basePath+"/swagger.json" {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("content-type", "application/json")
				w.Header().Set("cache-control", "no-cache")
				_, _ = w.Write(swaggerFile)
				return
			}
			w.Header().Set("cache-control", cacheControl)
			static.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

const basePath = "/swagger"
