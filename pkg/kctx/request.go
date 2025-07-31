package kctx

import (
	"mime"
	"time"

	"github.com/kemadev/go-framework/pkg/header"
)

// Date returns the Date of the request, based on "Date" header, as a [time.Date]. If an error occurs,
// it returns [time.Time]{}.
func (c *Kctx) Date() time.Time {
	date, err := time.Parse(time.RFC1123, c.r.Header.Get(header.Date))
	if err != nil {
		return time.Time{}
	}
	return date
}

// IsMIME returns whether the request satisfies given MIME
func (c *Kctx) IsMIME(mim string) bool {
	typ, _, _ := mime.ParseMediaType(c.r.Header.Get(header.ContentType))
	if typ == "" {
		// Invalid mime
		return false
	}

	mim, _, _ = mime.ParseMediaType(mim)
	if mim == "" {
		// Invalid mime
		return false
	}

	return typ == mim
}
