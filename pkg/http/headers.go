package http

import (
	"github.com/kemadev/go-framework/pkg/config"
)

const (
	AcceptEncodingHeaderKey                = "Accept-Encoding"
	AcceptHeaderKey                        = "Accept"
	AcceptLanguageHeaderKey                = "Accept-Language"
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
	ContentLengthHeaderKey                 = "Content-Length"
	ContentLocationHeaderKey               = "Content-Location"
	ContentSecurityPoliHeaderKey           = "Content-Security-Policy"
	ContentTypeHeaderKey                   = "Content-Type"
	CrossOriginEmbedderPolicyHeaderKey     = "Cross-Origin-Embedder-Policy"
	CrossOriginOpenerPolicyHeaderKey       = "Cross-Origin-Opener-Policy"
	CrossOriginResourcePolicyHeaderKey     = "Cross-Origin-Resource-Policy"
	ETagHeaderKey                          = "ETag"
	ExpiresHeaderKey                       = "Expires"
	ForwardedHeaderKey                     = "Forwarded"
	IfMatchHeaderKey                       = "If-Match"
	IfModifiedSinceHeaderKey               = "If-Modified-Since"
	IfNoneMatchHeaderKey                   = "If-None-Match"
	IfRangeHeaderKey                       = "If-Range"
	IfUnmodifiedHeaderKey                  = "If-Unmodified"
	IntegrityPolicyHeaderKey               = "Integrity-Policy"
	KeepAliveHeaderKey                     = "Keep-Alive"
	LastModifiedHeaderKey                  = "Last-Modified"
	PreferHeaderKey                        = "Prefer"
	ReferrerPolicyHeaderKey                = "Referrer-Policy"
	SecFetchModeHeaderKey                  = "Sec-Fetch-Mode"
	SecFetchSiteHeaderKey                  = "Sec-Fetch-Site"
	SecFetchUserHeaderKey                  = "Sec-Fetch-User"
	SourceMapHeaderKey                     = "SourceMap"
	StrictTransportSecurityHeaderKey       = "Strict-Transport-Security"
	TEHeaderKey                            = "TE"
	TimingAllowOriginHeaderKey             = "Timing-Allow-Origin"
	TransferEncodingHeaderKey              = "Transfer-Encoding"
	UpgradeInsecureRequestsHeaderKey       = "Upgrade-Insecure-Requests"
	UserAgentHeaderKey                     = "User-Agent"
	VaryHeaderKey                          = "Vary"
	ViaHeaderKey                           = "Via"
	WantContentDigestHeaderKey             = "Want-Content-Digest"
	WantReprDigestHeaderKey                = "Want-Repr-Digest"
	WWWAuthenticateHeaderKey               = "WWW-Authenticate"
	XContentTypeOptionsHeacyderKey         = "X-Content-Type-Options"
	XCSRFTokenHeaderKey                    = "X-CSRF-Token"
	XFrameOptionsHeaderKey                 = "X-Frame-Options"
)

func SetSecurityHeaders(config config.Config, clientInfo ClientInfo) error {
	writer := clientInfo.Writer

	if config.IsBrowserFacing {
		writer.Header().Set(http.)
	}

	return nil
}

// https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
// Cookie policy
