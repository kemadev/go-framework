// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package headutil

import (
	"mime"
	"net/http"
	"time"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
)

// Date returns the Date of the request, based on "Date" header, as a [time.Date]. If an error occurs,
// it returns [time.Time]{}.
func Date(h http.Header) time.Time {
	date, err := time.Parse(time.RFC1123, h.Get(headkey.Date))
	if err != nil {
		return time.Time{}
	}

	return date
}

// IsMIME returns whether the request satisfies given MIME.
func IsMIME(h http.Header, mim string) bool {
	typ, _, _ := mime.ParseMediaType(h.Get(headkey.ContentType))
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
