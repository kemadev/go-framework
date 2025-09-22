// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package sechead

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
)

type AccessControl struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Credentials
	AccessControlAllowCredentials bool
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Headers
	AccessControlAllowHeaders []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Methods
	AccessControlAllowMethods []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Origin
	AccessControlAllowOrigin url.URL
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Expose-Headers
	AccessControlExposeHeaders []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Max-Age
	AccessControlMaxAge time.Duration
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#fetch_directives
type ContentSecurityPolicyFetchDirectives struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/child-src
	ChildSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/connect-src
	ConnectSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/default-src
	DefaultSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/font-src
	FontSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/frame-src
	FrameSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/img-src
	ImageSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/manifest-src
	ManifestSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/media-src
	MediaSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/object-src
	ObjectSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/script-src
	ScriptSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/style-src
	StyleSource string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/worker-src
	WorkerSource string
}

const (
	// Custom, special value to set strict sandbox (renders as `sandbox;`).
	ContentSecurityPolicySandboxStrict = "true"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-downloads
	ContentSecurityPolicySandboxAllowDownloads = "allow-downloads"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-forms
	ContentSecurityPolicySandboxAllowForms = "allow-forms"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-modals
	ContentSecurityPolicySandboxAllowModals = "allow-modals"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-orientation-lock
	ContentSecurityPolicySandboxAllowOrientationLock = "allow-orientation-lock"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-pointer-lock
	ContentSecurityPolicySandboxAllowPointerLock = "allow-pointer-lock"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-popups
	ContentSecurityPolicySandboxAllowAllowPopups = "allow-popups"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-popups-to-escape-sandbox
	ContentSecurityPolicySandboxAllowAllowPopupsToEscapeSandbox = "allow-popups-to-escape-sandbox"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-presentation
	ContentSecurityPolicySandboxAllowPresentation = "allow-presentation"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-same-origin
	ContentSecurityPolicySandboxAllowSameOrigin = "allow-same-origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-scripts
	ContentSecurityPolicySandboxAllowScripts = "allow-scripts"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-top-navigation
	ContentSecurityPolicySandboxAllowTopNavigation = "allow-top-navigation"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-top-navigation-by-user-activation
	ContentSecurityPolicySandboxAllowTopNavigationByUserActivation = "allow-top-navigation-by-user-activation"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox#allow-top-navigation-to-custom-protocols
	ContentSecurityPolicySandboxAllowTopNavigationToCustomProtocols = "allow-top-navigation-to-custom-protocols"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#document_directives
type ContentSecurityPolicyDocumentDirectives struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/base-uri
	BaseURI string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/sandbox
	Sandbox string
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#navigation_directives
type ContentSecurityPolicyNavigationDirectives struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/form-action
	FormAction string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/frame-ancestors
	FrameAncestors string
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy
type ContentSecurityPolicy struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/child-src
	FetchDirectives ContentSecurityPolicyFetchDirectives
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#document_directives
	DocumentDirectives ContentSecurityPolicyDocumentDirectives
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#navigation_directives
	NavigationDirectives ContentSecurityPolicyNavigationDirectives
}

const (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Embedder-Policy#unsafe-none
	CrossOriginEmbedderPolicyUnsafeNone = "unsafe-none"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Embedder-Policy#require-corp
	CrossOriginEmbedderPolicyRequireCORP = "require-corp"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Embedder-Policy#credentialless
	CrossOriginEmbedderPolicyCredentialLess = "credentialless"
)

const (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Opener-Policy#unsafe-none
	CrossOriginOpenerPolicyUnsafeNone = "unsafe-none"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Opener-Policy#same-origin
	CrossOriginOpenerPolicySameOrigin = "same-origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Opener-Policy#same-origin-allow-popups
	CrossOriginOpenerPolicySameOriginAllowPopups = "same-origin-allow-popups"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Opener-Policy#noopener-allow-popups
	CrossOriginOpenerPolicyNoOpenerAllowPopups = "noopener-allow-popups"
)

const (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Resource-Policy#same-site
	CrossOriginResourcePolicySameSite = "same-site"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Resource-Policy#same-origin
	CrossOriginResourcePolicySameOrigin = "same-origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Resource-Policy#cross-origin
	CrossOriginResourcePolicyCrossOrigin = "cross-origin"
)

type CrossOriginPolicy struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Embedder-Policy
	CrossOriginEmbedderPolicy string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Opener-Policy
	CrossOriginOpenerPolicy string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Resource-Policy
	CrossOriginResourcePolicy string
}

const (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy#no-referrer_2
	ReferrerPolicyNoReferrer = "no-referrer"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy#no-referrer-when-downgrade_2
	ReferrerPolicyNoReferrerWhenDowngrade = "no-referrer-when-downgrade"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy#origin_2
	ReferrerPolicyOrigin = "origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy#origin-when-cross-origin_2
	ReferrerPolicyOriginWhenCrossOrigin = "origin-when-cross-origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy#same-origin_2
	ReferrerPolicySameOrigin = "same-origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy#strict-origin_2
	ReferrerPolicyStrictOrigin = "strict-origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy#strict-origin-when-cross-origin_2
	ReferrerPolicyStrictOriginWhenCrossOrigin = "strict-origin-when-cross-origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy#unsafe-url_2
	ReferrerPolicyUnsafeURL = "unsafe-url"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy
type ReferrerPolicy struct {
	ReferrerPolicy string
}

const (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-Frame-Options#deny
	FrameOptionsDeny = "DENY"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-Frame-Options#sameorigin
	FrameOptionsSameOrigin = "SAMEORIGIN"
)

type OtherOptions struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-Content-Type-Options
	ContentTypeOptionsNoSniff bool
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-Frame-Options
	FrameOptions string
}

type SecurityHeadersConfig struct {
	AccessControl         AccessControl
	ContentSecurityPolicy ContentSecurityPolicy
	CrossOriginPolicy     CrossOriginPolicy
	ReferrerPolicy        ReferrerPolicy
	OtherOptions          OtherOptions
}

func (conf *SecurityHeadersConfig) Headers() http.Header {
	head := http.Header{}

	if conf.AccessControl.AccessControlAllowCredentials {
		head.Set(
			headkey.AccessControlAllowCredentials,
			strconv.FormatBool(conf.AccessControl.AccessControlAllowCredentials),
		)
	}

	if len(conf.AccessControl.AccessControlAllowHeaders) > 0 {
		head.Set(
			headkey.AccessControlAllowHeaders,
			strings.Join(conf.AccessControl.AccessControlAllowHeaders, ", "),
		)
	}

	if len(conf.AccessControl.AccessControlAllowMethods) > 0 {
		head.Set(
			headkey.AccessControlAllowMethods,
			strings.Join(conf.AccessControl.AccessControlAllowMethods, ", "),
		)
	}

	if conf.AccessControl.AccessControlAllowOrigin.String() != "" {
		head.Set(
			headkey.AccessControlAllowOrigin,
			conf.AccessControl.AccessControlAllowOrigin.String(),
		)
	}

	if len(conf.AccessControl.AccessControlExposeHeaders) > 0 {
		head.Set(
			headkey.AccessControlExposeHeaders,
			strings.Join(conf.AccessControl.AccessControlExposeHeaders, ", "),
		)
	}

	if conf.AccessControl.AccessControlMaxAge > 0 {
		head.Set(
			headkey.AccessControlMaxAge,
			strconv.FormatInt(int64(conf.AccessControl.AccessControlMaxAge.Seconds()), 10),
		)
	}

	csp := conf.ContentSecurityPolicy.buildCSPHeader()
	if csp != "" {
		head.Set(headkey.ContentSecurityPolicy, csp)
	}

	if conf.CrossOriginPolicy.CrossOriginEmbedderPolicy != "" {
		head.Set(
			headkey.CrossOriginEmbedderPolicy,
			conf.CrossOriginPolicy.CrossOriginEmbedderPolicy,
		)
	}

	if conf.CrossOriginPolicy.CrossOriginOpenerPolicy != "" {
		head.Set(headkey.CrossOriginOpenerPolicy, conf.CrossOriginPolicy.CrossOriginOpenerPolicy)
	}

	if conf.CrossOriginPolicy.CrossOriginResourcePolicy != "" {
		head.Set(
			headkey.CrossOriginResourcePolicy,
			conf.CrossOriginPolicy.CrossOriginResourcePolicy,
		)
	}

	if conf.ReferrerPolicy.ReferrerPolicy != "" {
		head.Set(headkey.ReferrerPolicy, conf.ReferrerPolicy.ReferrerPolicy)
	}

	if conf.OtherOptions.ContentTypeOptionsNoSniff {
		head.Set(headkey.XContentTypeOptions, "nosniff")
	}

	if conf.OtherOptions.FrameOptions != "" {
		head.Set(headkey.XFrameOptions, conf.OtherOptions.FrameOptions)
	}

	return head
}

func (conf *ContentSecurityPolicy) buildCSPHeader() string {
	var directives []string

	if conf.FetchDirectives.DefaultSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("default-src '%s'", conf.FetchDirectives.DefaultSource),
		)
	}

	if conf.FetchDirectives.ChildSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("child-src '%s'", conf.FetchDirectives.ChildSource),
		)
	}

	if conf.FetchDirectives.ConnectSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("connect-src '%s'", conf.FetchDirectives.ConnectSource),
		)
	}

	if conf.FetchDirectives.FontSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("font-src '%s'", conf.FetchDirectives.FontSource),
		)
	}

	if conf.FetchDirectives.FrameSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("frame-src '%s'", conf.FetchDirectives.FrameSource),
		)
	}

	if conf.FetchDirectives.ImageSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("img-src '%s'", conf.FetchDirectives.ImageSource),
		)
	}

	if conf.FetchDirectives.ManifestSource != "" {
		directives = append(
			directives,
			fmt.Sprintf(
				"manifest-src '%s'",
				conf.FetchDirectives.ManifestSource,
			),
		)
	}

	if conf.FetchDirectives.MediaSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("media-src '%s'", conf.FetchDirectives.MediaSource),
		)
	}

	if conf.FetchDirectives.ObjectSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("object-src '%s'", conf.FetchDirectives.ObjectSource),
		)
	}

	if conf.FetchDirectives.ScriptSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("script-src '%s'", conf.FetchDirectives.ScriptSource),
		)
	}

	if conf.FetchDirectives.StyleSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("style-src '%s'", conf.FetchDirectives.StyleSource),
		)
	}

	if conf.FetchDirectives.WorkerSource != "" {
		directives = append(
			directives,
			fmt.Sprintf("worker-src '%s'", conf.FetchDirectives.WorkerSource),
		)
	}

	if conf.DocumentDirectives.BaseURI != "" {
		directives = append(
			directives,
			fmt.Sprintf("base-uri '%s'", conf.DocumentDirectives.BaseURI),
		)
	}

	if conf.DocumentDirectives.Sandbox != "" {
		if conf.DocumentDirectives.Sandbox == "true" {
			directives = append(directives, "sandbox")
		} else {
			directives = append(directives, fmt.Sprintf("sandbox '%s'", conf.DocumentDirectives.Sandbox))
		}
	}

	if conf.NavigationDirectives.FormAction != "" {
		directives = append(
			directives,
			fmt.Sprintf(
				"form-action '%s'",
				conf.NavigationDirectives.FormAction,
			),
		)
	}

	if conf.NavigationDirectives.FrameAncestors != "" {
		directives = append(
			directives,
			fmt.Sprintf(
				"frame-ancestors '%s'",
				conf.NavigationDirectives.FrameAncestors,
			),
		)
	}

	return strings.Join(directives, "; ")
}

var SecHeadersDefaultStrict = SecurityHeadersConfig{
	AccessControl: AccessControl{
		AccessControlAllowCredentials: false,
		AccessControlAllowHeaders:     []string{},
		AccessControlAllowMethods:     []string{},
		AccessControlAllowOrigin:      url.URL{},
		AccessControlExposeHeaders:    []string{},
		AccessControlMaxAge:           5 * time.Second,
	},
	ContentSecurityPolicy: ContentSecurityPolicy{
		FetchDirectives: ContentSecurityPolicyFetchDirectives{
			ChildSource:    "none",
			ConnectSource:  "none",
			DefaultSource:  "none",
			FontSource:     "none",
			FrameSource:    "none",
			ImageSource:    "none",
			ManifestSource: "none",
			MediaSource:    "none",
			ObjectSource:   "none",
			ScriptSource:   "none",
			StyleSource:    "none",
			WorkerSource:   "none",
		},
		DocumentDirectives: ContentSecurityPolicyDocumentDirectives{
			BaseURI: "none",
			Sandbox: "",
		},
		NavigationDirectives: ContentSecurityPolicyNavigationDirectives{
			FormAction:     "none",
			FrameAncestors: "none",
		},
	},
	CrossOriginPolicy: CrossOriginPolicy{
		CrossOriginEmbedderPolicy: CrossOriginEmbedderPolicyRequireCORP,
		CrossOriginOpenerPolicy:   CrossOriginOpenerPolicySameOrigin,
		CrossOriginResourcePolicy: CrossOriginResourcePolicySameOrigin,
	},
	ReferrerPolicy: ReferrerPolicy{
		ReferrerPolicy: ReferrerPolicySameOrigin,
	},
	OtherOptions: OtherOptions{
		ContentTypeOptionsNoSniff: true,
		FrameOptions:              FrameOptionsDeny,
	},
}
