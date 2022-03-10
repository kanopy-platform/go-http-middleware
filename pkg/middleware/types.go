package middleware

import "net/http"

type Provider interface {
	Middeleware(http.Handler) http.Handler
}
