package http

import (
	"github.com/kemadev/go-framework/pkg/config"
)

const (
	// Security
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Credentials
	AccessControlAllowCredentialsHeaderKey = "Access-Control-Allow-Credentials"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Headers
	AccessControlAllowHeadersHeaderKey = "Access-Control-Allow-Headers"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Methods
	AccessControlAllowMethodsHeaderKey = "Access-Control-Allow-Methods"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Allow-Origin
	AccessControlAllowOriginHeaderKey = "Access-Control-Allow-Origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Expose-Headers
	AccessControlExposeHeadersHeaderKey = "Access-Control-Expose-Headers"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Max-Age
	AccessControlMaxAgeHeaderKey = "Access-Control-Max-Age"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Request-Headers
	AccessControlRequestHeadersHeaderKey = "Access-Control-Request-Headers"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Access-Control-Request-Method
	AccessControlRequestMethodHeaderKey = "Access-Control-Request-Method"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy
	ContentSecurityPoliHeaderKey = "Content-Security-Policy"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Embedder-Policy
	CrossOriginEmbedderPolicyHeaderKey = "Cross-Origin-Embedder-Policy"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Opener-Policy
	CrossOriginOpenerPolicyHeaderKey = "Cross-Origin-Opener-Policy"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cross-Origin-Resource-Policy
	CrossOriginResourcePolicyHeaderKey = "Cross-Origin-Resource-Policy"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Integrity-Policy
	IntegrityPolicyHeaderKey = "Integrity-Policy"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy
	ReferrerPolicyHeaderKey = "Referrer-Policy"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Dest
	SecFetchDestHeaderKey = "Sec-Fetch-Dest"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Mode
	SecFetchModeHeaderKey = "Sec-Fetch-Mode"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Site
	SecFetchSiteHeaderKey = "Sec-Fetch-Site"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-User
	SecFetchUserHeaderKey = "Sec-Fetch-User"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Accept
	SecWebSocketAcceptHeaderKey = "Sec-WebSocket-Accept"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Extensions
	SecWebSocketExtensionsHeaderKey = "Sec-WebSocket-Extensions"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Key
	SecWebSocketKeyHeaderKey = "Sec-WebSocket-Key"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Protocol
	SecWebSocketProtocolHeaderKey = "Sec-WebSocket-Protocol"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-WebSocket-Version
	SecWebSocketVersionHeaderKey = "Sec-WebSocket-Version"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Strict-Transport-Security
	StrictTransportSecurityHeaderKey = "Strict-Transport-Security"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-Frame-Options
	XFrameOptionsHeaderKey = "X-Frame-Options"
	// Custom header, used as CSRF token (2nd submit)
	XCSRFTokenHeaderKey = "X-CSRF-Token"

	// Auth
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Authorization
	AuthorizationHeaderKey = "Authorization"

	// Cache
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cache-Control
	CacheControlHeaderKey = "Cache-Control"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Age
	AgeHeaderKey = "Age"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/ETag
	ETagHeaderKey = "ETag"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Expires
	ExpiresHeaderKey = "Expires"

	// Content negotiation
	AcceptEncodingHeaderKey = "Accept-Encoding"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept
	AcceptHeaderKey = "Accept"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Encoding
	ContentEncodingHeaderKey = "Content-Encoding"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Language
	AcceptLanguageHeaderKey = "Accept-Language"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Ranges
	AcceptRangesHeaderKey = "Accept-Ranges"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Language
	ContentLanguageHeaderKey = "Content-Language"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Type
	ContentTypeHeaderKey = "Content-Type"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Last-Modified
	LastModifiedHeaderKey = "Last-Modified"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Prefer
	PreferHeaderKey = "Prefer"

	// Encoding
	// Encoding for HTTP/2 and HTTP/3, https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/TE
	TEHeaderKey = "TE"
	// Encoding for HTTP/1.1, https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Transfer-Encoding
	TransferEncodingHeaderKey = "Transfer-Encoding"

	// Condition
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Match
	IfMatchHeaderKey = "If-Match"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Modified-Since
	IfModifiedSinceHeaderKey = "If-Modified-Since"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-None-Match
	IfNoneMatchHeaderKey = "If-None-Match"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Range
	IfRangeHeaderKey = "If-Range"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Unmodified-Since
	IfUnmodifiedSinceHeaderKey = "If-Unmodified-Since"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Expect
	ExpectHeaderKey = "Expect"

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Allow
	AllowHeaderKey = "Allow"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Digest
	ContentDigestHeaderKey = "Content-Digest"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Location
	ContentLocationHeaderKey = "Content-Location"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Forwarded
	ForwardedHeaderKey = "Forwarded"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Keep-Alive
	KeepAliveHeaderKey = "Keep-Alive"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Upgrade-Insecure-Requests
	UpgradeInsecureRequestsHeaderKey = "Upgrade-Insecure-Requests"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/User-Agent
	UserAgentHeaderKey = "User-Agent"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Vary
	VaryHeaderKey = "Vary"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Via
	ViaHeaderKey = "Via"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Want-Content-Digest
	WantContentDigestHeaderKey = "Want-Content-Digest"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Want-Repr-Digest
	WantReprDigestHeaderKey = "Want-Repr-Digest"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/WWW-Authenticate
	WWWAuthenticateHeaderKey = "WWW-Authenticate"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-Content-Type-Options
	XContentTypeOptionsHeacyderKey = "X-Content-Type-Options"
)

func SetSecurityHeaders(config config.Config, clientInfo ClientInfo) error {
	writer := clientInfo.Writer

	// if config.IsBrowserFacing {
	// 	writer.Header().Set(http.)
	// }

	return nil
}
