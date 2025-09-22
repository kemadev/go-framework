// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package headutil

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
)

// AcceptedValue represents a parsed Accept header value with its quality factor.
type AcceptedValue struct {
	Value   string
	Quality float64
}

// Accepts returns whether client signals accepting given media type (based on Accept header).
func Accepts(h http.Header, encoding string) bool {
	return accepts(h.Get(headkey.Accept), encoding)
}

// AcceptsEncoding returns whether client signals accepting given encoding.
func AcceptsEncoding(h http.Header, encoding string) bool {
	return accepts(h.Get(headkey.AcceptEncoding), encoding)
}

// AcceptsLanguage returns whether client signals accepting given language.
func AcceptsLanguage(h http.Header, language string) bool {
	return accepts(h.Get(headkey.AcceptLanguage), language)
}

func accepts(head, val string) bool {
	if val == "" {
		return false
	}

	if head == "" {
		return true
	}

	accepted := parseAcceptHeader(head)

	for _, value := range accepted {
		if value.Quality == 0 {
			continue
		}

		valueLower := strings.ToLower(value.Value)
		if valueLower == val || valueLower == "*" {
			return true
		}
	}

	return false
}

func parseAcceptHeader(header string) []AcceptedValue {
	if header == "" {
		return nil
	}

	values := strings.Split(header, ",")
	accepted := make([]AcceptedValue, 0, len(values))

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		acceptedValue := AcceptedValue{
			Quality: 1.0,
		}

		parts := strings.Split(value, ";")
		acceptedValue.Value = strings.TrimSpace(parts[0])

		if len(parts) < 2 {
			accepted = append(accepted, acceptedValue)

			continue
		}

		if len(parts) != 2 {
			// Malformed quality
			continue
		}

		for _, part := range parts[1:] {
			param := strings.TrimSpace(part)

			qualPrefix := "q="
			lenQualPrefix := len(qualPrefix)

			if !strings.HasPrefix(param, qualPrefix) || len(param) < (lenQualPrefix+1) {
				continue
			}

			qual, err := strconv.ParseFloat(param[lenQualPrefix:], 32)
			if err != nil {
				// Fallback
				continue
			}

			acceptedValue.Quality = qual
		}

		accepted = append(accepted, acceptedValue)
	}

	return accepted
}
