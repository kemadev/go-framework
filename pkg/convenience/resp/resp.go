package resp

import (
	"net/http"
	"net/url"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
)

// Redirect sends an HTTP redirect with given code and URL
func Redirect(w http.ResponseWriter, code int, url url.URL) {
	w.Header().Set(headkey.Location, url.String())
	w.WriteHeader(code)
}
