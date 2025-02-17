package handler

import "net/http"

const (
	HeaderRequest = "Hx-Request"
)

func IsHTMX(r *http.Request) bool {
	return r.Header.Get(HeaderRequest) == "true"
}
