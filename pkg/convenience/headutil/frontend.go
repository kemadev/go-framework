package headutil

import (
	"net/url"
	"time"
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
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Request-Headers
	AccessControlRequestHeaders []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Request-Method
	AccessControlRequestMethod string
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

type ContentSecurityPolicyBaseURI struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/base#href
	HRef url.URL
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/base#target
	Target string
}

const (
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
	BaseURI ContentSecurityPolicyBaseURI
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
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#audio
	SecFetchDestAudio = "audio"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#audioworklet
	SecFetchDestAudioWorklet = "audioworklet"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#document
	SecFetchDestDocument = "document"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#embed
	SecFetchDestEmbed = "embed"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#empty
	SecFetchDestEmpty = "empty"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#font
	SecFetchDestFont = "font"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#frame
	SecFetchDestFrame = "frame"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#iframe
	SecFetchDestIFrame = "iframe"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#image
	SecFetchDestImage = "image"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#manifest
	SecFetchDestManifest = "manifest"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#object
	SecFetchDestObject = "object"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#paintworklet
	SecFetchDestPaintWorklet = "paintworklet"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#report
	SecFetchDestReport = "report"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#script
	SecFetchDestScript = "script"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#serviceworker
	SecFetchDestServiceWorker = "serviceworker"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#sharedworker
	SecFetchDestSharedWorker = "sharedworker"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#style
	SecFetchDestStyle = "style"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#track
	SecFetchDestTrack = "track"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#video
	SecFetchDestVideo = "video"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#webidentity
	SecFetchDestWebIdentity = "webidentity"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#worker
	SecFetchDestWorker = "worker"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest#xslt
	SecFetchDestXSLT = "xslt"
)

const (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Mode#cors
	SecFetchModeCORS = "cors"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Mode#navigate
	SecFetchModeNavigate = "navigate"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Mode#no-cors
	SecFetchModeNoCORS = "no-cors"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Mode#same-origin
	SecFetchModeSameOrigin = "same-origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Mode#websocket
	SecFetchModeWebSocket = "websocket"
)

const (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Site#cross-site
	SecFetchSiteCrossSite = "cross-site"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Site#same-origin
	SecFetchSiteSameOrigin = "same-origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Site#same-site
	SecFetchSiteSameSite = "same-site"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Site#none
	SecFetchSiteNone = "none"
)

type SecFetch struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest
	SecFetchDest string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Mode
	SecFetchMode string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Site
	SecFetchSite string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-User
	SecFetchUser bool
}

type SecWebSocket struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Accept
	SecWebSocketAccept string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Extensions
	SecWebSocketExtensions []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Key
	SecWebSocketKey string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Protocol
	SecWebSocketProtocol string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Version
	SecWebSocketVersion int
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

type CSRFToken struct {
	Token string
}

type SecurityHeadersConfig struct {
	AccessControl         AccessControl
	ContentSecurityPolicy ContentSecurityPolicy
	CrossOriginPolicy     CrossOriginPolicy
	ReferrerPolicy        ReferrerPolicy
	SecFetch              SecFetch
	SecWebSocket          SecWebSocket
	OtherOptions          OtherOptions
}

var DefaultStrict = SecurityHeadersConfig{
	AccessControl: AccessControl{},
}
