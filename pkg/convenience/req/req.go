package req

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"

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
	if ct != headval.AcceptJSON {
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

func extractFromForwardedHeaderFirst(r *http.Request, extractKey string) (string, error) {
	forwarded := r.Header.Get(headkey.Forwarded)
	if forwarded == "" {
		return "", ErrNoForwardedHeader
	}

	forwardedDirective := func() string {
		parts := strings.SplitSeq(forwarded, ";")
		for part := range parts {
			extractor := extractKey + "="
			found := strings.Index(part, extractor)
			if found >= 0 {
				return part[len(extractor):]
			}
		}
		return ""
	}()

	if forwardedDirective == "" {
		return "", ErrForwardedDirectiveNotFound
	}

	return forwardedDirective, nil
}

func IP(r *http.Request) (net.IP, error) {
	ipString, err := extractFromForwardedHeaderFirst(r, "for")
	if err != nil {
		return net.IP{}, fmt.Errorf("error extracting IP from request: %w", err)
	}

	ip := net.ParseIP(ipString)
	if ip == nil {
		return net.IP{}, ErrForwardedDirectiveInvalid
	}

	return ip, nil
}

func IPs(r *http.Request) ([]net.IP, error) {
	forwarded := r.Header.Get(headkey.Forwarded)
	if forwarded == "" {
		return []net.IP{}, ErrNoForwardedHeader
	}

	parts := strings.SplitSeq(forwarded, ";")
	ips := []net.IP{}
	for part := range parts {
		extractor := "for="
		ipsChain := strings.SplitSeq(part, ", ")
		for chainPart := range ipsChain {
			ipString := chainPart[len(extractor):]
			slog.Debug(ipString)
			ip := net.ParseIP(ipString)
			if ip == nil {
				return []net.IP{}, ErrForwardedDirectiveInvalid
			}
			ips = append(ips, ip)
		}
	}

	return ips, nil
}

func Host(r *http.Request) (*url.URL, error) {
	hostString, err := extractFromForwardedHeaderFirst(r, "host")
	if err != nil {
		return nil, fmt.Errorf("error extracting host from request: %w", err)
	}

	host, err := url.Parse(hostString)
	if err != nil {
		return nil, fmt.Errorf("error extracting host from request: %w", err)
	}

	if host == nil {
		return nil, ErrForwardedDirectiveInvalid
	}

	return host, nil
}
