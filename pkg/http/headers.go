package http

import (
	"net/http"
	"strings"

	"github.com/kemadev/go-framework/pkg/config"
)

const (
	AcceptEncodingHeaderKey                = "Accept-Encoding"
	AcceptHeaderKey                        = "Accept"
	AcceptLanguageHeaderKey                = "Accept-Language"
	AcceptRangesHeaderKey                  = "Accept-Ranges"
	AccessControlAllowCredentialsHeaderKey = "Access-Control-Allow-Credentials"
	AccessControlAllowHeadersHeaderKey     = "Access-Control-Allow-Headers"
	AccessControlAllowMethodsHeaderKey     = "Access-Control-Allow-Methods"
	AccessControlAllowOriginHeaderKey      = "Access-Control-Allow-Origin"
	AccessControlExposeHeadersHeaderKey    = "Access-Control-Expose-Headers"
	AccessControlMaxAgeHeaderKey           = "Access-Control-Max-Age"
	AccessControlRequestHeadersHeaderKey   = "Access-Control-Request-Headers"
	AccessControlRequestMethodHeaderKey    = "Access-Control-Request-Method"
	AgeHeaderKey                           = "Age"
	AllowHeaderKey                         = "Allow"
	AuthorizationHeaderKey                 = "Authorization"
	CacheControlHeaderKey                  = "Cache-Control"
	ContentDigestHeaderKey                 = "Content-Digest"
	ContentEncodingHeaderKey               = "Content-Encoding"
	ContentLanguageHeaderKey               = "Content-Language"
	ContentLocationHeaderKey               = "Content-Location"
	ContentSecurityPolicyHeaderKey         = "Content-Security-Policy"
	ContentTypeHeaderKey                   = "Content-Type"
	CrossOriginEmbedderPolicyHeaderKey     = "Cross-Origin-Embedder-Policy"
	CrossOriginOpenerPolicyHeaderKey       = "Cross-Origin-Opener-Policy"
	CrossOriginResourcePolicyHeaderKey     = "Cross-Origin-Resource-Policy"
	ETagHeaderKey                          = "ETag"
	ExpectHeaderKey                        = "Expect"
	ExpiresHeaderKey                       = "Expires"
	ForwardedHeaderKey                     = "Forwarded"
	IfMatchHeaderKey                       = "If-Match"
	IfModifiedSinceHeaderKey               = "If-Modified-Since"
	IfNoneMatchHeaderKey                   = "If-None-Match"
	IfRangeHeaderKey                       = "If-Range"
	IfUnmodifiedSinceHeaderKey             = "If-Unmodified-Since"
	IntegrityPolicyHeaderKey               = "Integrity-Policy"
	KeepAliveHeaderKey                     = "Keep-Alive"
	LastModifiedHeaderKey                  = "Last-Modified"
	PreferHeaderKey                        = "Prefer"
	ReferrerPolicyHeaderKey                = "Referrer-Policy"
	SecFetchDestHeaderKey                  = "Sec-Fetch-Dest"
	SecFetchModeHeaderKey                  = "Sec-Fetch-Mode"
	SecFetchSiteHeaderKey                  = "Sec-Fetch-Site"
	SecFetchUserHeaderKey                  = "Sec-Fetch-User"
	SecWebSocketAcceptHeaderKey            = "Sec-WebSocket-Accept"
	SecWebSocketExtensionsHeaderKey        = "Sec-WebSocket-Extensions"
	SecWebSocketKeyHeaderKey               = "Sec-WebSocket-Key"
	SecWebSocketProtocolHeaderKey          = "Sec-WebSocket-Protocol"
	SecWebSocketVersionHeaderKey           = "Sec-WebSocket-Version"
	StrictTransportSecurityHeaderKey       = "Strict-Transport-Security"
	TEHeaderKey                            = "TE"
	TransferEncodingHeaderKey              = "Transfer-Encoding"
	UpgradeInsecureRequestsHeaderKey       = "Upgrade-Insecure-Requests"
	UserAgentHeaderKey                     = "User-Agent"
	VaryHeaderKey                          = "Vary"
	ViaHeaderKey                           = "Via"
	WantContentDigestHeaderKey             = "Want-Content-Digest"
	WantReprDigestHeaderKey                = "Want-Repr-Digest"
	WWWAuthenticateHeaderKey               = "WWW-Authenticate"
	XContentTypeOptionsHeaderKey           = "X-Content-Type-Options"
	XCSRFTokenHeaderKey                    = "X-CSRF-Token"
	XFrameOptionsHeaderKey                 = "X-Frame-Options"
)

type HeadersConfig struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Credentials
	AccessControlAllowCredentials string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Headers
	AccessControlAllowHeaders string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Methods
	AccessControlAllowMethods string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Origin
	AccessControlAllowOrigin string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Expose-Headers
	AccessControlExposeHeaders string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Max-Age
	AccessControlMaxAge string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Request-Headers
	AccessControlRequestHeaders string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Request-Method
	AccessControlRequestMethod string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy
	ContentSecurityPolicy string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Embedder-Policy
	CrossOriginEmbedderPolicy string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Opener-Policy
	CrossOriginOpenerPolicy string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Resource-Policy
	CrossOriginResourcePolicy string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Integrity-Policy
	IntegrityPolicy string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy
	ReferrerPolicy string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest
	SecFetchDest string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Mode
	SecFetchMode string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Site
	SecFetchSite string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-User
	SecFetchUser string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Accept
	SecWebSocketAccept string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Extensions
	SecWebSocketExtensions string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Key
	SecWebSocketKey string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Protocol
	SecWebSocketProtocol string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Version
	SecWebSocketVersion string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Strict-Transport-Security
	StrictTransportSecurity string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-Frame-Options
	XFrameOptions string
	// Custom header, used as CSRF token (2nd submit)
	XCSRFToken string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Authorization
	Authorization string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cache-Control
	CacheControl string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Age
	Age string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/ETag
	ETag string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Expires
	Expires string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Encoding
	AcceptEncoding string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept
	Accept string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Encoding
	ContentEncoding string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Language
	AcceptLanguage string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Ranges
	AcceptRanges string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Language
	ContentLanguage string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Type
	ContentType string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Last-Modified
	LastModified string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Prefer
	Prefer string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/TE
	TE string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Transfer-Encoding
	TransferEncoding string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Match
	IfMatch string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Modified-Since
	IfModifiedSince string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-None-Match
	IfNoneMatch string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Range
	IfRange string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Unmodified-Since
	IfUnmodifiedSince string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Expect
	Expect string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Allow
	Allow string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Digest
	ContentDigest string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Location
	ContentLocation string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Forwarded
	Forwarded string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Keep-Alive
	KeepAlive string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Upgrade-Insecure-Requests
	UpgradeInsecureRequests string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/User-Agent
	UserAgent string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Vary
	Vary string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Via
	Via string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Want-Content-Digest
	WantContentDigest string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Want-Repr-Digest
	WantReprDigest string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/WWW-Authenticate
	WWWAuthenticate string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-Content-Type-Options
	XContentTypeOptions string
}

func defaultHeadersConfig(
	cspConfig CSPConfig,
) HeadersConfig {
	return HeadersConfig{
		AccessControlAllowHeaders: strings.Join([]string{
			AuthorizationHeaderKey,
			AcceptEncodingHeaderKey,
			AcceptHeaderKey,
			PreferHeaderKey,
			IfMatchHeaderKey,
			IfModifiedSinceHeaderKey,
			IfNoneMatchHeaderKey,
			IfUnmodifiedSinceHeaderKey,
			IfRangeHeaderKey,
			IfUnmodifiedSinceHeaderKey,
			ExpectHeaderKey,
			UserAgentHeaderKey,
			WantContentDigestHeaderKey,
			WantReprDigestHeaderKey,
		}, ", "),
		// TODO integrate with [net/http.Handler] methods
		AccessControlAllowMethods: "PUT, DELETE, OPTIONS, PATCH",
		// TODO integrate with [net/http.Handler] methods
		Allow: "PUT, DELETE, OPTIONS, PATCH",
		// TODO set it to self origin dynamically
		AccessControlAllowOrigin: "",
		AccessControlExposeHeaders: strings.Join([]string{
			ETagHeaderKey,
			ContentEncodingHeaderKey,
			TEHeaderKey,
			TransferEncodingHeaderKey,
			ContentDigestHeaderKey,
			VaryHeaderKey,
		}, ", "),
		AccessControlMaxAge:       "300",
		ContentSecurityPolicy:     NewCSP(cspConfig),
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
		IntegrityPolicy:           "blocked-destinations=(script)",
		ReferrerPolicy:            "strict-origin",
		StrictTransportSecurity:   "max-age=63072000; includeSubDomains; preload",
		XFrameOptions:             "DENY",
		// TODO integrate csrf double commit handling
		XCSRFToken:          "",
		CacheControl:        "no-cache",
		AcceptEncoding:      "gzip",
		Accept:              "application/json",
		ContentEncoding:     "gzip",
		AcceptRanges:        "bytes",
		WWWAuthenticate:     "Bearer",
		XContentTypeOptions: "nosniff",
	}
}

func setHeaders(w http.ResponseWriter, config HeadersConfig) {
	if AcceptEncodingHeaderKey != "" {
		w.Header().Set(AcceptEncodingHeaderKey, config.AcceptEncoding)
	}
	if AcceptHeaderKey != "" {
		w.Header().Set(AcceptHeaderKey, config.Accept)
	}
	if AcceptLanguageHeaderKey != "" {
		w.Header().Set(AcceptLanguageHeaderKey, config.AcceptLanguage)
	}
	if AcceptRangesHeaderKey != "" {
		w.Header().Set(AcceptRangesHeaderKey, config.AcceptRanges)
	}
	if AccessControlAllowCredentialsHeaderKey != "" {
		w.Header().Set(AccessControlAllowCredentialsHeaderKey, config.AccessControlAllowCredentials)
	}
	if AccessControlAllowHeadersHeaderKey != "" {
		w.Header().Set(AccessControlAllowHeadersHeaderKey, config.AccessControlAllowHeaders)
	}
	if AccessControlAllowMethodsHeaderKey != "" {
		w.Header().Set(AccessControlAllowMethodsHeaderKey, config.AccessControlAllowMethods)
	}
	if AccessControlAllowOriginHeaderKey != "" {
		w.Header().Set(AccessControlAllowOriginHeaderKey, config.AccessControlAllowOrigin)
	}
	if AccessControlExposeHeadersHeaderKey != "" {
		w.Header().Set(AccessControlExposeHeadersHeaderKey, config.AccessControlExposeHeaders)
	}
	if AccessControlMaxAgeHeaderKey != "" {
		w.Header().Set(AccessControlMaxAgeHeaderKey, config.AccessControlMaxAge)
	}
	if AccessControlRequestHeadersHeaderKey != "" {
		w.Header().Set(AccessControlRequestHeadersHeaderKey, config.AccessControlRequestHeaders)
	}
	if AccessControlRequestMethodHeaderKey != "" {
		w.Header().Set(AccessControlRequestMethodHeaderKey, config.AccessControlRequestMethod)
	}
	if AgeHeaderKey != "" {
		w.Header().Set(AgeHeaderKey, config.Age)
	}
	if AllowHeaderKey != "" {
		w.Header().Set(AllowHeaderKey, config.Allow)
	}
	if AuthorizationHeaderKey != "" {
		w.Header().Set(AuthorizationHeaderKey, config.Authorization)
	}
	if CacheControlHeaderKey != "" {
		w.Header().Set(CacheControlHeaderKey, config.CacheControl)
	}
	if ContentDigestHeaderKey != "" {
		w.Header().Set(ContentDigestHeaderKey, config.ContentDigest)
	}
	if ContentEncodingHeaderKey != "" {
		w.Header().Set(ContentEncodingHeaderKey, config.ContentEncoding)
	}
	if ContentLanguageHeaderKey != "" {
		w.Header().Set(ContentLanguageHeaderKey, config.ContentLanguage)
	}
	if ContentLocationHeaderKey != "" {
		w.Header().Set(ContentLocationHeaderKey, config.ContentLocation)
	}
	if ContentSecurityPolicyHeaderKey != "" {
		w.Header().Set(ContentSecurityPolicyHeaderKey, config.ContentSecurityPolicy)
	}
	if ContentTypeHeaderKey != "" {
		w.Header().Set(ContentTypeHeaderKey, config.ContentType)
	}
	if CrossOriginEmbedderPolicyHeaderKey != "" {
		w.Header().Set(CrossOriginEmbedderPolicyHeaderKey, config.CrossOriginEmbedderPolicy)
	}
	if CrossOriginOpenerPolicyHeaderKey != "" {
		w.Header().Set(CrossOriginOpenerPolicyHeaderKey, config.CrossOriginOpenerPolicy)
	}
	if CrossOriginResourcePolicyHeaderKey != "" {
		w.Header().Set(CrossOriginResourcePolicyHeaderKey, config.CrossOriginResourcePolicy)
	}
	if ETagHeaderKey != "" {
		w.Header().Set(ETagHeaderKey, config.ETag)
	}
	if ExpectHeaderKey != "" {
		w.Header().Set(ExpectHeaderKey, config.Expect)
	}
	if ExpiresHeaderKey != "" {
		w.Header().Set(ExpiresHeaderKey, config.Expires)
	}
	if ForwardedHeaderKey != "" {
		w.Header().Set(ForwardedHeaderKey, config.Forwarded)
	}
	if IfMatchHeaderKey != "" {
		w.Header().Set(IfMatchHeaderKey, config.IfMatch)
	}
	if IfModifiedSinceHeaderKey != "" {
		w.Header().Set(IfModifiedSinceHeaderKey, config.IfModifiedSince)
	}
	if IfNoneMatchHeaderKey != "" {
		w.Header().Set(IfNoneMatchHeaderKey, config.IfNoneMatch)
	}
	if IfRangeHeaderKey != "" {
		w.Header().Set(IfRangeHeaderKey, config.IfRange)
	}
	if IfUnmodifiedSinceHeaderKey != "" {
		w.Header().Set(IfUnmodifiedSinceHeaderKey, config.IfUnmodifiedSince)
	}
	if IntegrityPolicyHeaderKey != "" {
		w.Header().Set(IntegrityPolicyHeaderKey, config.IntegrityPolicy)
	}
	if KeepAliveHeaderKey != "" {
		w.Header().Set(KeepAliveHeaderKey, config.KeepAlive)
	}
	if LastModifiedHeaderKey != "" {
		w.Header().Set(LastModifiedHeaderKey, config.LastModified)
	}
	if PreferHeaderKey != "" {
		w.Header().Set(PreferHeaderKey, config.Prefer)
	}
	if ReferrerPolicyHeaderKey != "" {
		w.Header().Set(ReferrerPolicyHeaderKey, config.ReferrerPolicy)
	}
	if SecFetchDestHeaderKey != "" {
		w.Header().Set(SecFetchDestHeaderKey, config.SecFetchDest)
	}
	if SecFetchModeHeaderKey != "" {
		w.Header().Set(SecFetchModeHeaderKey, config.SecFetchMode)
	}
	if SecFetchSiteHeaderKey != "" {
		w.Header().Set(SecFetchSiteHeaderKey, config.SecFetchSite)
	}
	if SecFetchUserHeaderKey != "" {
		w.Header().Set(SecFetchUserHeaderKey, config.SecFetchUser)
	}
	if SecWebSocketAcceptHeaderKey != "" {
		w.Header().Set(SecWebSocketAcceptHeaderKey, config.SecWebSocketAccept)
	}
	if SecWebSocketExtensionsHeaderKey != "" {
		w.Header().Set(SecWebSocketExtensionsHeaderKey, config.SecWebSocketExtensions)
	}
	if SecWebSocketKeyHeaderKey != "" {
		w.Header().Set(SecWebSocketKeyHeaderKey, config.SecWebSocketKey)
	}
	if SecWebSocketProtocolHeaderKey != "" {
		w.Header().Set(SecWebSocketProtocolHeaderKey, config.SecWebSocketProtocol)
	}
	if SecWebSocketVersionHeaderKey != "" {
		w.Header().Set(SecWebSocketVersionHeaderKey, config.SecWebSocketVersion)
	}
	if StrictTransportSecurityHeaderKey != "" {
		w.Header().Set(StrictTransportSecurityHeaderKey, config.StrictTransportSecurity)
	}
	if TEHeaderKey != "" {
		w.Header().Set(TEHeaderKey, config.TE)
	}
	if TransferEncodingHeaderKey != "" {
		w.Header().Set(TransferEncodingHeaderKey, config.TransferEncoding)
	}
	if UpgradeInsecureRequestsHeaderKey != "" {
		w.Header().Set(UpgradeInsecureRequestsHeaderKey, config.UpgradeInsecureRequests)
	}
	if UserAgentHeaderKey != "" {
		w.Header().Set(UserAgentHeaderKey, config.UserAgent)
	}
	if VaryHeaderKey != "" {
		w.Header().Set(VaryHeaderKey, config.Vary)
	}
	if ViaHeaderKey != "" {
		w.Header().Set(ViaHeaderKey, config.Via)
	}
	if WantContentDigestHeaderKey != "" {
		w.Header().Set(WantContentDigestHeaderKey, config.WantContentDigest)
	}
	if WantReprDigestHeaderKey != "" {
		w.Header().Set(WantReprDigestHeaderKey, config.WantReprDigest)
	}
	if WWWAuthenticateHeaderKey != "" {
		w.Header().Set(WWWAuthenticateHeaderKey, config.WWWAuthenticate)
	}
	if XContentTypeOptionsHeaderKey != "" {
		w.Header().Set(XContentTypeOptionsHeaderKey, config.XContentTypeOptions)
	}
	if XCSRFTokenHeaderKey != "" {
		w.Header().Set(XCSRFTokenHeaderKey, config.XCSRFToken)
	}
	if XFrameOptionsHeaderKey != "" {
		w.Header().Set(XFrameOptionsHeaderKey, config.XFrameOptions)
	}
}

func SetSecurityHeaders(
	appConfig config.Config,
	clientInfo ClientInfo,
	cspConfig CSPConfig,
) error {
	setHeaders(clientInfo.Writer, defaultHeadersConfig(cspConfig))

	return nil
}
