package kctx

import (
	"strconv"
	"strings"
	"time"

	"github.com/kemadev/go-framework/pkg/header"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cache-Control#cache_directives
type CacheHeader struct {
	MaxAge               time.Duration
	MaxStale             time.Duration
	MinFresh             time.Duration
	SMaxAge              time.Duration
	NoCache              bool
	NoStore              bool
	NoTransform          bool
	OnlyIfCached         bool
	MustRevalidate       bool
	ProxyRevalidate      bool
	MustUnderstand       bool
	Private              bool
	Public               bool
	Immutable            bool
	StaleWhileRevalidate bool
	StaleIfError         bool
}

// CacheDecision represents a decision on how to handle caching
type CacheDecision int

const (
	// Allow caching
	CacheAllow CacheDecision = iota
	// Revalidate cache
	CacheRevalidate
	// Refresh cache entry
	CacheRefresh
	// Skip cache
	CacheBypass
)

// Build constructs the Cache-Control header value for given CacheHeader
func (head *CacheHeader) Build() string {
	if head == nil {
		return ""
	}

	reasonableHeadValueSize := 32

	var b strings.Builder
	b.Grow(reasonableHeadValueSize)

	head.addDuration(&b, "max-age", head.MaxAge)
	head.addDuration(&b, "max-stale", head.MaxStale)
	head.addDuration(&b, "min-fresh", head.MinFresh)
	head.addDuration(&b, "s-maxage", head.SMaxAge)

	head.addDirective(&b, "no-cache", head.NoCache)
	head.addDirective(&b, "no-store", head.NoStore)
	head.addDirective(&b, "no-transform", head.NoTransform)
	head.addDirective(&b, "only-if-cached", head.OnlyIfCached)
	head.addDirective(&b, "must-revalidate", head.MustRevalidate)
	head.addDirective(&b, "proxy-revalidate", head.ProxyRevalidate)
	head.addDirective(&b, "must-understand", head.MustUnderstand)
	head.addDirective(&b, "private", head.Private)
	head.addDirective(&b, "public", head.Public)
	head.addDirective(&b, "immutable", head.Immutable)
	head.addDirective(&b, "stale-while-revalidate", head.StaleWhileRevalidate)
	head.addDirective(&b, "stale-if-error", head.StaleIfError)

	return b.String()
}

func (head *CacheHeader) addDirective(b *strings.Builder, directive string, enabled bool) {
	if !enabled {
		return
	}

	if b.Len() > 0 {
		b.WriteString(", ")
	}

	b.WriteString(directive)
}

func (head *CacheHeader) addDuration(b *strings.Builder, directive string, duration time.Duration) {
	if duration <= 0 {
		return
	}

	head.addDirective(b, directive, true)
	b.WriteByte('=')
	b.WriteString(strconv.Itoa(int(duration.Seconds())))
}

// CachePolicySet sets cache control header with given cache header
func (c *Kctx) CachePolicySet(head CacheHeader) {
	c.w.Header().Set(header.CacheControl, head.Build())
}

// CacheDecisionFromRequest return a caching decision from based on the request headers
func (c *Kctx) CacheDecisionFromRequest() CacheDecision {
	cacheControl := c.r.Header.Get(header.CacheControl)
	if cacheControl != "" {
		if strings.Contains(cacheControl, "no-cache") {
			return CacheBypass
		}
		if strings.Contains(cacheControl, "max-age=0") {
			return CacheRefresh
		}
	}

	if c.r.Header.Get(header.IfNoneMatch) != "" ||
		c.r.Header.Get(header.IfModifiedSince) != "" ||
		c.r.Header.Get(header.IfUnmodifiedSince) != "" ||
		c.r.Header.Get(header.IfMatch) != "" ||
		c.r.Header.Get(header.IfRange) != "" {
		return CacheRevalidate
	}

	return CacheAllow
}
