// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package http

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
