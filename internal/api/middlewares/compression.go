package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

func Comporession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check of the client can accept gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
		}
		// set the respone header
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		//Wrap the responseWritter
		w = &gzipResponseWriter{ResponseWriter: w, Writer: gz}

		// call the next handler
		next.ServeHTTP(w, r)
		fmt.Println("Send response from Compression Middleware")
	})
}

// gizpResponseWriter wraps http.ResponseWriter to write gzipped responses
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}
