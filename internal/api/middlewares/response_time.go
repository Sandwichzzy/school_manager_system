package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Response Time Middleware being returned...")
		start := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Calculate the duration
		duration := time.Since(start)
		w.Header().Set("X-Response-Time", duration.String())
		//Create a custom response writer to capture status code
		next.ServeHTTP(wrappedWriter, r)
		// Log the request details
		duration = time.Since(start)
		fmt.Printf("Method: %s, URL: %s, Status: %d, Duration: %v\n", r.Method, r.URL, wrappedWriter.status, duration.String())
		fmt.Println("Sent Response from Response Time Middleware")
	})
}

// responseWriter
// 结构体内部嵌入了 http.ResponseWriter 接口，
// 这意味着它自动实现了 http.ResponseWriter 接口的所有方法（除非需要重写）
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
