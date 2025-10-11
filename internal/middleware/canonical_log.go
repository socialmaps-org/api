package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"codeberg.org/socialmaps/auth/internal/session"
)

// CanonicalLog is middleware that logs canonical HTTP request lines
func CanonicalLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		// Wrap response writer to capture status code
		rw := NewResponseWriter(w)

		// Extract session ID if present
		var sessionID string
		cookie, err := r.Cookie(session.COOKIE_NAME)
		if err == nil && cookie != nil {
			// Try to extract session ID, but don't panic if invalid
			func() {
				defer func() {
					if recover() != nil {
						sessionID = ""
					}
				}()
				// We need the cookie secret to decode, but we don't have access to it here
				// For now, just note that a cookie is present
				sessionID = "present"
			}()
		}

		// Recover from panics in handlers
		var panicErr error
		defer func() {
			if rec := recover(); rec != nil {
				panicErr = fmt.Errorf("panic: %v", rec)
				// Re-panic after logging
				defer panic(rec)
			}

			duration := time.Since(start)

			// Build log attributes
			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.StatusCode,
				"duration_ms", duration.Milliseconds(),
				"remote_addr", r.RemoteAddr,
			}

			if r.URL.RawQuery != "" {
				attrs = append(attrs, "query", r.URL.RawQuery)
			}

			if sessionID != "" {
				attrs = append(attrs, "session_present", true)
			}

			if userAgent := r.UserAgent(); userAgent != "" {
				attrs = append(attrs, "user_agent", userAgent)
			}

			if referer := r.Referer(); referer != "" {
				attrs = append(attrs, "referer", referer)
			}

			if panicErr != nil {
				attrs = append(attrs, "error", panicErr.Error())
			}

			slog.InfoContext(ctx, "CANONICAL-HTTP-REQUEST", attrs...)
		}()

		// Call the next handler
		next.ServeHTTP(rw, r)
	})
}
