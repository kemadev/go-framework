package http

import (
	"net/url"
	"strings"
)

type IntegrityConfig struct {
	Algorithm string
	Value     string
}

type Source struct {
	None       bool
	Expression string
	Host       url.URL
	Scheme     string
}

type FetchSrc struct {
	Sources     []Source
	NonceValues []string
	Integrities []IntegrityConfig
}

type CSPFetchConfig struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/default-src
	DefaultSrc FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/script-src
	ScriptSrc FetchSrc
	// // https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/script-src-elem
	// ScriptSrcElem FetchSrc
	// // https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/script-src-attr
	// ScriptSrcAttr FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/style-src
	StyleSrc FetchSrc
	// // https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/style-src-elem
	// StyleSrcElem FetchSrc
	// // https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/style-src-attr
	// StyleSrcAttr FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/img-src
	ImgSrc FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/font-src
	FontSrc FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/connect-src
	ConnectSrc FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/manifest-src
	MediaSrc FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/object-src
	ObjectSrc FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/child-src
	ChildSrc FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/frame-src
	FrameSrc FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/worker-src
	WorkerSrc FetchSrc
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/manifest-src
	ManifestSrc FetchSrc
}

type CSPDocumentConfig struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/base-uri
	BaseURI Source
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox
	Sandbox string
}

type CSPNavigationConfig struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/form-action
	FormAction Source
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/frame-ancestors
	FrameAncestors Source
}

type CSPReportingConfig struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/report-to
	ReportTo string
}

type CSPOtherConfig struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/upgrade-insecure-requests
	UpgradeInsecureRequests bool
}

type CSPConfig struct {
	Fetch      CSPFetchConfig
	Document   CSPDocumentConfig
	Navigation CSPNavigationConfig
	Other      CSPOtherConfig
}

func sourceToString(src Source) string {
	if src.None {
		return "'none'"
	}
	if src.Expression != "" {
		return src.Expression
	}
	if src.Host.String() != "" {
		return src.Host.String()
	}
	if src.Scheme != "" {
		return src.Scheme + ":"
	}
	return "'none'"
}

func fetchSrcToString(fetchSrc FetchSrc) string {
	var parts []string

	for _, source := range fetchSrc.Sources {
		parts = append(parts, sourceToString(source))
	}

	for _, nonce := range fetchSrc.NonceValues {
		parts = append(parts, "'nonce-"+nonce+"'")
	}

	for _, integrity := range fetchSrc.Integrities {
		parts = append(parts, "'"+integrity.Algorithm+"-"+integrity.Value+"'")
	}

	if len(parts) == 0 {
		return "'none'"
	}

	return strings.Join(parts, " ")
}

func sandboxToString(sandbox string) string {
	if sandbox != "" {
		return sandbox
	}
	return ""
}

func newCSPConfigWithIntegrity(
	scriptIntegrities []IntegrityConfig,
	styleIntegrities []IntegrityConfig,
) CSPConfig {
	noneSource := Source{None: true}
	selfSource := Source{Expression: "'self'"}

	return CSPConfig{
		Fetch: CSPFetchConfig{
			DefaultSrc: FetchSrc{Sources: []Source{noneSource}},
			ScriptSrc: FetchSrc{
				Sources:     []Source{selfSource},
				Integrities: scriptIntegrities,
			},
			StyleSrc: FetchSrc{
				Sources:     []Source{selfSource},
				Integrities: styleIntegrities,
			},
			ImgSrc:      FetchSrc{Sources: []Source{selfSource}},
			FontSrc:     FetchSrc{Sources: []Source{selfSource}},
			ConnectSrc:  FetchSrc{Sources: []Source{selfSource}},
			MediaSrc:    FetchSrc{Sources: []Source{noneSource}},
			ObjectSrc:   FetchSrc{Sources: []Source{noneSource}},
			ChildSrc:    FetchSrc{Sources: []Source{noneSource}},
			FrameSrc:    FetchSrc{Sources: []Source{noneSource}},
			WorkerSrc:   FetchSrc{Sources: []Source{noneSource}},
			ManifestSrc: FetchSrc{Sources: []Source{selfSource}},
		},
		Document: CSPDocumentConfig{
			BaseURI: selfSource,
			Sandbox: "",
		},
		Navigation: CSPNavigationConfig{
			FormAction:     selfSource,
			FrameAncestors: noneSource,
		},
		Other: CSPOtherConfig{
			UpgradeInsecureRequests: false,
		},
	}
}

func newCSPFromConfig(config CSPConfig) string {
	var csp []string

	csp = append(csp, "default-src "+fetchSrcToString(config.Fetch.DefaultSrc))
	csp = append(csp, "script-src "+fetchSrcToString(config.Fetch.ScriptSrc))
	// csp = append(csp, "script-src-elem "+fetchSrcToString(config.Fetch.ScriptSrcElem))
	// csp = append(csp, "script-src-attr "+fetchSrcToString(config.Fetch.ScriptSrcAttr))
	csp = append(csp, "style-src "+fetchSrcToString(config.Fetch.StyleSrc))
	// csp = append(csp, "style-src-elem "+fetchSrcToString(config.Fetch.StyleSrcElem))
	// csp = append(csp, "style-src-attr "+fetchSrcToString(config.Fetch.StyleSrcAttr))
	csp = append(csp, "img-src "+fetchSrcToString(config.Fetch.ImgSrc))
	csp = append(csp, "font-src "+fetchSrcToString(config.Fetch.FontSrc))
	csp = append(csp, "connect-src "+fetchSrcToString(config.Fetch.ConnectSrc))
	csp = append(csp, "media-src "+fetchSrcToString(config.Fetch.MediaSrc))
	csp = append(csp, "object-src "+fetchSrcToString(config.Fetch.ObjectSrc))
	csp = append(csp, "child-src "+fetchSrcToString(config.Fetch.ChildSrc))
	csp = append(csp, "frame-src "+fetchSrcToString(config.Fetch.FrameSrc))
	csp = append(csp, "worker-src "+fetchSrcToString(config.Fetch.WorkerSrc))
	csp = append(csp, "manifest-src "+fetchSrcToString(config.Fetch.ManifestSrc))

	csp = append(csp, "base-uri "+sourceToString(config.Document.BaseURI))
	csp = append(csp, "sandbox "+sandboxToString(config.Document.Sandbox))

	csp = append(csp, "form-action "+sourceToString(config.Navigation.FormAction))
	csp = append(csp, "frame-ancestors "+sourceToString(config.Navigation.FrameAncestors))

	if config.Other.UpgradeInsecureRequests {
		csp = append(csp, "upgrade-insecure-requests")
	}

	return strings.Join(csp, "; ")
}

func NewCSP(
	scriptIntegrities []IntegrityConfig,
	styleIntegrities []IntegrityConfig,
	config CSPConfig,
) string {
	base := newCSPConfigWithIntegrity(
		scriptIntegrities,
		styleIntegrities,
	)

	merged := mergeCSPConfigs(base, config)

	final := newCSPFromConfig(merged)

	return final
}

func mergeCSPConfigs(base CSPConfig, override CSPConfig) CSPConfig {
	merged := CSPConfig{
		Fetch: CSPFetchConfig{
			DefaultSrc:  mergeFetchSrc(base.Fetch.DefaultSrc, override.Fetch.DefaultSrc),
			ScriptSrc:   mergeFetchSrc(base.Fetch.ScriptSrc, override.Fetch.ScriptSrc),
			StyleSrc:    mergeFetchSrc(base.Fetch.StyleSrc, override.Fetch.StyleSrc),
			ImgSrc:      mergeFetchSrc(base.Fetch.ImgSrc, override.Fetch.ImgSrc),
			FontSrc:     mergeFetchSrc(base.Fetch.FontSrc, override.Fetch.FontSrc),
			ConnectSrc:  mergeFetchSrc(base.Fetch.ConnectSrc, override.Fetch.ConnectSrc),
			MediaSrc:    mergeFetchSrc(base.Fetch.MediaSrc, override.Fetch.MediaSrc),
			ObjectSrc:   mergeFetchSrc(base.Fetch.ObjectSrc, override.Fetch.ObjectSrc),
			ChildSrc:    mergeFetchSrc(base.Fetch.ChildSrc, override.Fetch.ChildSrc),
			FrameSrc:    mergeFetchSrc(base.Fetch.FrameSrc, override.Fetch.FrameSrc),
			WorkerSrc:   mergeFetchSrc(base.Fetch.WorkerSrc, override.Fetch.WorkerSrc),
			ManifestSrc: mergeFetchSrc(base.Fetch.ManifestSrc, override.Fetch.ManifestSrc),
		},
		Document: CSPDocumentConfig{
			BaseURI: mergeSource(base.Document.BaseURI, override.Document.BaseURI),
			Sandbox: mergeString(base.Document.Sandbox, override.Document.Sandbox),
		},
		Navigation: CSPNavigationConfig{
			FormAction: mergeSource(base.Navigation.FormAction, override.Navigation.FormAction),
			FrameAncestors: mergeSource(
				base.Navigation.FrameAncestors,
				override.Navigation.FrameAncestors,
			),
		},
		Other: CSPOtherConfig{
			UpgradeInsecureRequests: mergeBool(
				base.Other.UpgradeInsecureRequests,
				override.Other.UpgradeInsecureRequests,
			),
		},
	}

	return merged
}

func mergeFetchSrc(base FetchSrc, override FetchSrc) FetchSrc {
	if hasFetchSrcContent(override) {
		return FetchSrc{
			Sources:     mergeSourceSlices(base.Sources, override.Sources),
			NonceValues: mergeStringSlices(base.NonceValues, override.NonceValues),
			Integrities: mergeIntegritySlices(base.Integrities, override.Integrities),
		}
	}
	return base
}

func mergeSource(base Source, override Source) Source {
	if hasSourceContent(override) {
		return override
	}
	return base
}

func mergeString(base string, override string) string {
	if override != "" {
		return override
	}
	return base
}

func mergeBool(base bool, override bool) bool {
	if override {
		return override
	}
	return base
}

func mergeSourceSlices(base []Source, override []Source) []Source {
	if len(override) > 0 {
		merged := make([]Source, len(base))
		copy(merged, base)

		for _, overrideSrc := range override {
			if !containsSource(merged, overrideSrc) {
				merged = append(merged, overrideSrc)
			}
		}
		return merged
	}
	return base
}

func mergeStringSlices(base []string, override []string) []string {
	if len(override) > 0 {
		merged := make([]string, len(base))
		copy(merged, base)

		for _, overrideStr := range override {
			if !containsString(merged, overrideStr) {
				merged = append(merged, overrideStr)
			}
		}
		return merged
	}
	return base
}

func mergeIntegritySlices(base []IntegrityConfig, override []IntegrityConfig) []IntegrityConfig {
	if len(override) > 0 {
		merged := make([]IntegrityConfig, len(base))
		copy(merged, base)

		for _, overrideIntegrity := range override {
			if !containsIntegrity(merged, overrideIntegrity) {
				merged = append(merged, overrideIntegrity)
			}
		}
		return merged
	}
	return base
}

func hasFetchSrcContent(fetchSrc FetchSrc) bool {
	return len(fetchSrc.Sources) > 0 || len(fetchSrc.NonceValues) > 0 ||
		len(fetchSrc.Integrities) > 0
}

func hasSourceContent(src Source) bool {
	return src.None || src.Expression != "" || src.Host.String() != "" || src.Scheme != ""
}

func containsSource(sources []Source, target Source) bool {
	for _, src := range sources {
		if src.None == target.None &&
			src.Expression == target.Expression &&
			src.Host.String() == target.Host.String() &&
			src.Scheme == target.Scheme {
			return true
		}
	}
	return false
}

func containsString(strings []string, target string) bool {
	for _, str := range strings {
		if str == target {
			return true
		}
	}
	return false
}

func containsIntegrity(integrities []IntegrityConfig, target IntegrityConfig) bool {
	for _, integrity := range integrities {
		if integrity.Algorithm == target.Algorithm && integrity.Value == target.Value {
			return true
		}
	}
	return false
}
