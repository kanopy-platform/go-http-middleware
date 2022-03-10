package middleware

import "net/http"

type Provider interface {
	Middleware(http.Handler) http.Handler
}
