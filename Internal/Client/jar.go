package client

import (
	"net/url"
	"sync"
	"time"

	"github.com/cicadaaio/httpclient/net/http"
)

type FJar struct {
	// mu locks the remaining fields.
	mu sync.Mutex

	// entries is a set of entries, keyed by their eTLD+1 and subkeyed by
	// their name/domain/path.
	entries map[string]*http.Cookie
}

func (F *FJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	now := time.Now()

	F.mu.Lock()
	defer F.mu.Unlock()
	for _, cookie := range cookies {
		if cookie.MaxAge > 0 {
			cookie.Expires = now.Add(time.Duration(cookie.MaxAge) * time.Second)
		} else if cookie.Expires.IsZero() {
			cookie.Expires = endOfTime
		}

		F.entries[cookie.Name] = cookie
	}
}

func New() *FJar {
	jar := &FJar{
		entries: make(map[string]*http.Cookie),
	}
	return jar
}

var endOfTime = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)

func (F *FJar) Cookies(_ *url.URL) []*http.Cookie {
	now := time.Now()
	F.mu.Lock()
	defer F.mu.Unlock()
	var cookies []*http.Cookie
	for _, cookie := range F.entries {
		if !cookie.Expires.After(now) {
			delete(F.entries, cookie.Name)
			continue
		}
		cookies = append(cookies, cookie)
	}
	return cookies
}
