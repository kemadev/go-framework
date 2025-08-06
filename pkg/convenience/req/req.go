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
