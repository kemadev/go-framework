package resp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
)

var ErrTemplateNotFound = errors.New("template not found")

// Redirect sends an HTTP redirect with given code and URL
func Redirect(w http.ResponseWriter, code int, url url.URL) {
	w.Header().Set(headkey.Location, url.String())
	w.WriteHeader(code)
}

// JSON sends payload after marshalling it to JSON, returning an error if marshalling fails.
// It also sets correct content type header.
func JSON(w http.ResponseWriter, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling json: %w", err)
	}

	w.Header().Set(headkey.ContentType, headval.AcceptJSON)
	w.Write(body)

	return nil
}
