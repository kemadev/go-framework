// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package log

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

// Colors definitions using ANSI escape codes
// These colors are used to format the log output in the terminal.
const (
	// Gray color as ANSI escape code.
	gray = "\033[90m"
	// Cyan color as ANSI escape code.
	cyan = "\033[36m"
	// Green color as ANSI escape code.
	green = "\033[32m"
	// Yellow color as ANSI escape code.
	yellow = "\033[33m"
	// Red color as ANSI escape code.
	red = "\033[31m"
	// Magenta color as ANSI escape code.
	magenta = "\033[35m"
	// Reset color as ANSI escape code, reverts to default terminal color.
	reset = "\033[0m"
)

// Exporter is a custom log exporter that formats and outputs log records to stdout.
// It uses the OpenTelemetry SDK's stdoutlog.Exporter as a base, but formats the output
// to make it human-readable and color-coded for different severity levels.
type Exporter struct {
	*stdoutlog.Exporter
}

// formatAttr returns a formatted string for a log attribute, using ANSI escape codes
// to colorize the key and value. The key is colored gray, and the value is colored cyan.
func formatAttr(key string, value any) string {
	return fmt.Sprintf("%s%s%s=%s%v%s", gray, key, reset, cyan, value, reset)
}

func appendSourceAttrs(
	attrs []string,
	filepath string,
	lineno int64,
	function string,
) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return []string{}, fmt.Errorf("failed to get working directory: %w", err)
	}

	filepath = strings.TrimPrefix(filepath, wd+"/")

	if filepath != "" && function != "" && lineno != 0 {
		f := formatAttr(string(semconv.CodeFunctionNameKey), function)
		attrs = append(attrs, f)
		// Add a leading space to enable hyperlink detection in terminals, remove leading mount
		// path to make it relative to application root.
		reg, _ := regexp.Compile(`^/` + wd + `\w+/(.+)`)
		path := reg.ReplaceAllString(filepath, "$1")
		source := formatAttr(string(semconv.CodeFilePathKey), fmt.Sprintf(" %s:%v", path, lineno))
		attrs = append(attrs, source)
	}

	return attrs, nil
}

// Export exports log records to writer.
func (e *Exporter) Export(ctx context.Context, records []sdklog.Record) error {
	for _, record := range records {
		// Honor context cancellation.
		err := ctx.Err()
		if err != nil {
			return fmt.Errorf("error present in context: %w", err)
		}

		var attrs []string

		var filepath, function string

		var lineno int64

		severity := record.Severity().String()

		// severityColor maps log severity levels to their corresponding colors.
		severityColor := map[string]string{
			"DEBUG": cyan,
			"INFO":  green,
			"WARN":  yellow,
			"ERROR": red,
			"FATAL": magenta,
		}

		color, ok := severityColor[severity]
		if !ok {
			color = ""
		}

		record.WalkAttributes(func(logKV log.KeyValue) bool {
			switch logKV.Key {
			case string(semconv.CodeFilePathKey):
				filepath = logKV.Value.AsString()
			case string(semconv.CodeFunctionNameKey):
				function = logKV.Value.AsString()
			case string(semconv.CodeLineNumberKey):
				lineno = logKV.Value.AsInt64()
			default:
				attrs = append(attrs, formatAttr(logKV.Key, logKV.Value.String()))
			}

			return true
		})

		attrs, err = appendSourceAttrs(attrs, filepath, lineno, function)
		if err != nil {
			return fmt.Errorf("error appending source attributes: %w", err)
		}

		fmt.Printf(
			"%s%s%s %s%s%s\t%s\t%s\n",
			gray, record.Timestamp().Format(time.TimeOnly), reset,
			color, severity, reset,
			record.Body().AsString(),
			strings.Join(attrs, " "),
		)
	}

	return nil
}

func NewExporter(opts ...stdoutlog.Option) (*Exporter, error) {
	exp, err := stdoutlog.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating stdout exporter: %w", err)
	}

	return &Exporter{Exporter: exp}, nil
}
