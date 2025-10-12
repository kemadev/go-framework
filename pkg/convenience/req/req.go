// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package req

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
)

var (
	ErrNotJSON                    = errors.New("content type is not JSON")
	ErrNoForwardedHeader          = errors.New("no Forwarded header found in request")
	ErrForwardedDirectiveNotFound = errors.New("Forwarded directive not found in request header")
	ErrForwardedDirectiveInvalid  = errors.New("Forwarded directive is invalid")
)

// JSONFromBody parses [r]'s body as JSON into a new instance of type T and returns
// the parsed object along with an appropriate HTTP status code for any error encountered
// during processing (or ok status if there is no error).
// Please note that [net/http.MaxBytesReader] is not called, it is the responsibility of the caller to set it accordingly.
func JSONFromBody[T any](w http.ResponseWriter, r *http.Request) (T, int, error) {
	var zero T

	ct := r.Header.Get(headkey.ContentType)
	if !strings.HasPrefix(ct, headval.MIMEApplicationJSON) {
		return zero, http.StatusBadRequest, fmt.Errorf(
			"error processing JSON with content type %q: %w",
			ct,
			ErrNotJSON,
		)
	}

	var result T

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&result)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return zero, http.StatusRequestEntityTooLarge, fmt.Errorf(
				"request body too large: %w",
				err,
			)
		}

		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			return zero, http.StatusBadRequest, fmt.Errorf(
				"invalid JSON syntax at position %d: %w",
				syntaxErr.Offset,
				err,
			)
		}

		var unmarshalTypeErr *json.UnmarshalTypeError
		if errors.As(err, &unmarshalTypeErr) {
			return zero, http.StatusBadRequest, fmt.Errorf(
				"invalid value for field %q: %w",
				unmarshalTypeErr.Field,
				err,
			)
		}

		return zero, http.StatusBadRequest, fmt.Errorf("error processing JSON: %w", err)
	}

	// Check for additional JSON content (streaming not supported)
	err = dec.Decode(&struct{}{})
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return zero, http.StatusBadRequest, fmt.Errorf(
				"JSON stream not supported: %w",
				err,
			)
		}
	}

	return result, http.StatusOK, nil
}

// IP parses the Forwarded headers and returns the forwarded IP (the first to appear in the headers), and an error if any occurs.
// Please note that this functions does not check proxy trust, it is the caller's responsibility to ensure the header is trusted.
func IP(r *http.Request) (net.IP, error) {
	conf, err := config.NewManager().Get()
	if err != nil {
		return nil, fmt.Errorf("error retrieving app configuration: %w", err)
	}

	for _, head := range r.Header[conf.Server.ProxyHeader] {
		if head == "" {
			return nil, ErrNoForwardedHeader
		}

		entries := strings.SplitSeq(head, ",")
		for entry := range entries {
			entry = strings.TrimSpace(entry)
			if entry == "" {
				return nil, ErrForwardedDirectiveInvalid
			}

			directives := strings.SplitSeq(entry, ";")
			for directive := range directives {
				directive = strings.TrimSpace(directive)

				extractor := "for="
				if !strings.HasPrefix(directive, extractor) {
					continue
				}

				ipString := directive[len(extractor):]
				ipString = strings.TrimPrefix(ipString, `"`)
				ipString = strings.TrimSuffix(ipString, `"`)
				ipString = strings.TrimPrefix(ipString, `[`)
				ipString = strings.TrimSuffix(ipString, `]`)

				ip := net.ParseIP(ipString)
				if ip == nil {
					return nil, ErrForwardedDirectiveInvalid
				}

				return ip, nil
			}
		}
	}

	return nil, ErrForwardedDirectiveInvalid
}

// IPs parses the Forwarded headers and returns the forwarded IPs, and an error if any occurs.
// Please note that this functions does not check proxy trust, it is the caller's responsibility to ensure the header is trusted.
func IPs(r *http.Request) ([]*net.IP, error) {
	conf, err := config.NewManager().Get()
	if err != nil {
		return nil, fmt.Errorf("error retrieving app configuration: %w", err)
	}

	var ips []*net.IP

	for _, head := range r.Header[conf.Server.ProxyHeader] {
		if head == "" {
			return nil, ErrNoForwardedHeader
		}

		entries := strings.SplitSeq(head, ",")
		for entry := range entries {
			entry = strings.TrimSpace(entry)
			if entry == "" {
				return nil, ErrForwardedDirectiveInvalid
			}

			directives := strings.SplitSeq(entry, ";")
			for directive := range directives {
				directive = strings.TrimSpace(directive)

				extractor := "for="
				if !strings.HasPrefix(directive, extractor) {
					continue
				}

				ipString := directive[len(extractor):]
				ipString = strings.TrimPrefix(ipString, `"`)
				ipString = strings.TrimSuffix(ipString, `"`)
				ipString = strings.TrimPrefix(ipString, `[`)
				ipString = strings.TrimSuffix(ipString, `]`)

				ip := net.ParseIP(ipString)
				if ip == nil {
					return nil, ErrForwardedDirectiveInvalid
				}

				ips = append(ips, &ip)
			}
		}

		if len(ips) == 0 {
			return nil, ErrForwardedDirectiveInvalid
		}
	}

	return ips, nil
}

// Host parses the Forwarded headers and returns the forwarded host (the first to appear in the headers), and an error if any occurs.
// Please note that this functions does not check proxy trust, it is the caller's responsibility to ensure the header is trusted.
func Host(r *http.Request) (*url.URL, error) {
	conf, err := config.NewManager().Get()
	if err != nil {
		return nil, fmt.Errorf("error retrieving app configuration: %w", err)
	}

	for _, head := range r.Header[conf.Server.ProxyHeader] {
		if head == "" {
			return nil, ErrNoForwardedHeader
		}

		entries := strings.SplitSeq(head, ",")
		for entry := range entries {
			entry = strings.TrimSpace(entry)
			if entry == "" {
				return nil, ErrForwardedDirectiveInvalid
			}

			directives := strings.SplitSeq(entry, ";")
			for directive := range directives {
				directive = strings.TrimSpace(directive)

				extractor := "host="
				if !strings.HasPrefix(directive, extractor) {
					continue
				}

				hostString := directive[len(extractor):]

				host, err := url.Parse(hostString)
				if err != nil || host == nil {
					return nil, ErrForwardedDirectiveInvalid
				}

				return host, nil
			}
		}
	}

	return nil, ErrForwardedDirectiveInvalid
}

// Hosts parses the Forwarded headers and returns the forwarded hosts, and an error if any occurs.
// Please note that this functions does not check proxy trust, it is the caller's responsibility to ensure the header is trusted.
func Hosts(r *http.Request) ([]*url.URL, error) {
	var hosts []*url.URL

	for _, head := range r.Header[headkey.Forwarded] {
		if head == "" {
			return nil, ErrNoForwardedHeader
		}

		entries := strings.SplitSeq(head, ",")
		for entry := range entries {
			entry = strings.TrimSpace(entry)
			if entry == "" {
				return nil, ErrForwardedDirectiveInvalid
			}

			directives := strings.SplitSeq(entry, ";")
			for directive := range directives {
				directive = strings.TrimSpace(directive)

				extractor := "host="
				if !strings.HasPrefix(directive, extractor) {
					continue
				}

				hostString := directive[len(extractor):]

				host, err := url.Parse(hostString)
				if err != nil || host == nil {
					return nil, ErrForwardedDirectiveInvalid
				}

				hosts = append(hosts, host)
			}
		}

		if len(hosts) == 0 {
			return nil, ErrForwardedDirectiveInvalid
		}
	}

	return hosts, nil
}
