package req

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
)

var ErrNotJSON = errors.New("content type is not JSON")

// JSONFromBody parses [r]'s body as JSON into [dest] (which must be an allocated object) and returns appropriate
// HTTP status code for any error encountered during processing (or ok status if there is no error).
// Please note that [net/http.MaxBytesReader] is not called, it is the responsability of the caller to set it accordingly.
func JSONFromBody(w http.ResponseWriter, r *http.Request, dest any) (int, error) {
	ct := r.Header.Get(headkey.ContentType)
	if ct != headval.AcceptJSON {
		return http.StatusBadRequest, fmt.Errorf(
			"error processing JSON with content type %q: %w",
			ct,
			ErrNotJSON,
		)
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dest)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return http.StatusRequestEntityTooLarge, fmt.Errorf("request body too large: %w", err)
		}

		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			return http.StatusBadRequest, fmt.Errorf(
				"invalid JSON syntax at position %d: %w",
				syntaxErr.Offset,
				err,
			)
		}

		var unmarshalTypeErr *json.UnmarshalTypeError
		if errors.As(err, &unmarshalTypeErr) {
			return http.StatusBadRequest, fmt.Errorf(
				"invalid value for field %q: %w",
				unmarshalTypeErr.Field,
				err,
			)
		}

		return http.StatusBadRequest, fmt.Errorf("error processing JSON: %w", err)
	}

	err = dec.Decode(&struct{}{})
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return http.StatusBadRequest, fmt.Errorf(
				"JSON stream not supported: %w",
				err,
			)
		}
	}

	return http.StatusOK, nil
}
