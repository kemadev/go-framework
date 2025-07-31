package kctx

import (
	"net/http"
	"net/url"

	"github.com/kemadev/go-framework/pkg/header"
)

func RedirectPermanent(w http.ResponseWriter, url url.URL) {
	redirect(w, url, http.StatusPermanentRedirect)
}

func RedirectTemporary(w http.ResponseWriter, url url.URL) {
	redirect(w, url, http.StatusTemporaryRedirect)
}

func redirect(w http.ResponseWriter, url url.URL, code int) {
	w.Header().Set(header.Location, url.String())
	w.WriteHeader(code)
}
