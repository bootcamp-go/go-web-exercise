package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// NewLogger creates a new logger.
func NewLogger() *Logger {
	return &Logger{}
}

// Logger handles logging.
type Logger struct {
	// ...
}

// Log logs requests.
func (l *Logger) Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before
		// - start timer
		start := time.Now()
		
		// call
		handler.ServeHTTP(w, r)

		// after
		// - stop timer
		duration := time.Since(start)
		// - message
		messageRequest := fmt.Sprintf("[%s] %s -size: %d", r.Method, r.URL, r.ContentLength)
		messageResponse:= fmt.Sprintf("| -time: %v", duration)
		// - log
		fmt.Println(messageRequest, messageResponse)
	})
}