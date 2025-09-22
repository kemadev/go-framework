// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package resp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
)

var ErrTemplateNotFound = errors.New("template not found")

// JSON sends payload after marshalling it to JSON, returning an error if marshalling fails.
// It also sets correct content type header.
func JSON(w http.ResponseWriter, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling json: %w", err)
	}

	w.Header().Set(headkey.ContentType, headval.MIMEApplicationJSONCharsetUTF8)
	w.Write(body)

	return nil
}
