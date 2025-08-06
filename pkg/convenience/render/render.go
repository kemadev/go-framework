package render

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
)

var ErrTemplateNotFound = errors.New("template not found")

// TemplateRenderer handles template parsing and rendering
type TemplateRenderer struct {
	templates map[string]*template.Template
}

// New creates a new template renderer with all templates parsed
func New(tmpl embed.FS) (*TemplateRenderer, error) {
	tr := &TemplateRenderer{
		templates: make(map[string]*template.Template),
	}

	err := tr.loadTemplates(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return tr, nil
}

func (tr *TemplateRenderer) loadTemplates(tmpl embed.FS) error {
	return fs.WalkDir(tmpl, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		content, err := tmpl.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		t := template.New(filepath.Base(path))

		_, err = t.Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		templateName := strings.TrimPrefix(path, "tmpl/")
		tr.templates[templateName] = t

		return nil
	})
}

// Execute executes a template with the given data
func (tr *TemplateRenderer) Execute(
	w http.ResponseWriter,
	templateName string,
	data any,
) error {
	t, exists := tr.templates[templateName]

	if !exists {
		return fmt.Errorf("%s: %w", templateName, ErrTemplateNotFound)
	}

	w.Header().Set(headkey.ContentType, headval.ContentTypeHTMLUTF8)

	err := t.Execute(w, data)
	if err != nil {
		return fmt.Errorf("error executing template %s: %w", templateName, err)
	}

	return nil
}

// HandlerFuncWithData creates an HTTP handler function with dynamic data
func (tr *TemplateRenderer) HandlerFuncWithData(
	dataFunc func(*http.Request) (any, error),
	pathValueKey string,
	tmplExt string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := dataFunc(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tr.Execute(w, r.PathValue(pathValueKey)+tmplExt, data)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}
}
