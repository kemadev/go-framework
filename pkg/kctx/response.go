package kctx

import (
	"net/http"
	"net/url"

	"github.com/kemadev/go-framework/pkg/header"
)

func Redirect(w http.ResponseWriter, url url.URL, code int) {
	w.Header().Set(header.Location, url.String())
	w.WriteHeader(code)
}
