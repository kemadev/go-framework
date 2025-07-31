package kctx

import (
	"strconv"
	"strings"

	"github.com/kemadev/go-framework/pkg/header"
)

// AcceptedValue represents a parsed Accept header value with its quality factor
type AcceptedValue struct {
	Value   string
	Quality float64
}

// AcceptsEncoding returns whether client signals accepting given encoding
func (c *Kctx) AcceptsEncoding(encoding string) bool {
	return accepts(c.r.Header.Get(header.AcceptEncoding), encoding)
}

// AcceptsLanguage returns whether client signals accepting given language
func (c *Kctx) AcceptsLanguage(language string) bool {
	return accepts(c.r.Header.Get(header.AcceptLanguage), language)
}

func accepts(head string, val string) bool {
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

		for i := range parts[1:] {
			param := strings.TrimSpace(parts[i])

			qualPrefix := "q="
			lenQualPrefix := len(qualPrefix)

			if !strings.HasPrefix(param, qualPrefix) || len(param) < (lenQualPrefix+1) {
				continue
			}

			qual, err := strconv.ParseFloat(param[lenQualPrefix:], 10)
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
