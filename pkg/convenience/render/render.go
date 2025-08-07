package render

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
	"github.com/kemadev/go-framework/pkg/convenience/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

const packageName = "github.com/kemadev/go-framework/pkg/convenience/render"

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

		// Add leading `/` to permit matching request URLs
		tr.templates["/"+path] = t

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

// HandlerFuncWithData creates an HTTP handler function with dynamic data, using [templatePathResolver] to form
// template name from request URL path, or a default resolver if none is provided.
func (tr *TemplateRenderer) HandlerFuncWithData(
	dataFunc func(*http.Request) (any, error),
	templatePathResolver func(originalPath string) string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := dataFunc(r)
		if err != nil {
			log.Logger(packageName).Error(
				"can't load template data",
				slog.String(string(semconv.ErrorMessageKey), err.Error()),
				slog.String(string(semconv.URLPathKey), r.URL.Path),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		path := r.URL.Path

		if templatePathResolver != nil {
			path = templatePathResolver(path)
		} else {
			path = tr.resolvePath(path)
		}

		_, exists := tr.templates[path]
		if !exists {
			http.NotFound(w, r)
			return
		}

		err = tr.Execute(w, path, data)
		if err != nil {
			log.Logger(packageName).Error(
				"can't execute template",
				slog.String(string(semconv.ErrorMessageKey), err.Error()),
				slog.String(string(semconv.URLPathKey), r.URL.Path),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (tr *TemplateRenderer) resolvePath(urlPath string) string {
	cleanPath := path.Clean(urlPath)

	if cleanPath == "." || cleanPath == "/" {
		return "/index.html"
	}

	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}

	if strings.HasSuffix(cleanPath, "/") {
		return cleanPath + "index.html"
	}

	if path.Ext(cleanPath) != "" {
		return cleanPath
	}

	return cleanPath + ".html"
}
