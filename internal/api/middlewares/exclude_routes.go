package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

func MiddlewareExcludePaths(middleware func(http.Handler) http.Handler, excludedPaths ...string) func(http.Handler) http.Handler {
	fmt.Println("Middleware Exclude Paths Initialized")
	return func(next http.Handler) http.Handler {
		fmt.Println("================= Middleware Exclude Paths ================")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range excludedPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}
			middleware(next).ServeHTTP(w, r)
			fmt.Println("Response sent from Middleware Exclude Paths")
		})
	}
}
