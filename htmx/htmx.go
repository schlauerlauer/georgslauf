package htmx

import "net/http"

// HTTP request headers
const (
	HeaderRequest = "Hx-Request"
)

// HTTP response headers
const (
	HeaderRedirect = "HX-Redirect"
)

func IsHTMX(r *http.Request) bool {
	return r.Header.Get(HeaderRequest) == "true"
}
