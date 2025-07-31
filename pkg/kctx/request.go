package kctx

import (
	"strconv"
	"strings"
	"time"

	"github.com/kemadev/go-framework/pkg/header"
)

// CacheDirective represents cache handling instructions
type CacheDirective int

const (
	CacheAllow CacheDirective = iota
	CacheRevalidate
	CacheBypass
	CacheForceRefresh
)

// Date returns the Date of the request, based on "Date" header, as a [time.Date]. If an error occurs,
// it returns [time.Time]{}.
func (c *Kctx) Date() time.Time {
	date, err := time.Parse(time.RFC1123, c.r.Header.Get(header.Date))
	if err != nil {
		return time.Time{}
	}
	return date
}

// CachePolicy determines how to handle caching based on request headers
func (c *Kctx) CachePolicy() CacheDirective {
	cacheControl := c.r.Header.Get(header.CacheControl)

	if cacheControl != "" {
		switch {
		case strings.Contains(cacheControl, "no-cache"):
			return CacheForceRefresh
		case strings.Contains(cacheControl, "no-store"):
			return CacheBypass
		case strings.Contains(cacheControl, "max-age=0"):
			return CacheRevalidate
		case strings.Contains(cacheControl, "must-revalidate"):
			return CacheRevalidate
		}
	}

	return CacheAllow
}

// IsStale checks if cached content should be considered stale
func (c *Kctx) IsStale(cacheTime time.Time) bool {
	return time.Since(cacheTime) > c.getMaxAge()
}

func (c *Kctx) getMaxAge() time.Duration {
	cacheControl := c.r.Header.Get(header.CacheControl)
	if cacheControl == "" {
		return 0
	}

	maxAgePrefix := "max-age="
	if idx := strings.Index(cacheControl, maxAgePrefix); idx != -1 {
		start := idx + len(maxAgePrefix)
		end := start
		for end < len(cacheControl) && cacheControl[end] >= '0' && cacheControl[end] <= '9' {
			end++
		}
		if end > start {
			seconds, err := strconv.Atoi(cacheControl[start:end])
			if err != nil {
				return 0
			}

			return time.Duration(seconds) * time.Second
		}
	}

	return 0
}
