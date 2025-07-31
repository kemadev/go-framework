package kctx

import (
	"net/url"

	"github.com/kemadev/go-framework/pkg/header"
)

// Redirect sends an HTTP redirect with given code and URL
func (c *Kctx) Redirect(code int, url url.URL) {
	c.w.Header().Set(header.Location, url.String())
	c.w.WriteHeader(code)
}
