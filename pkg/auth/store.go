package auth

import (
	"net/http"
	"time"
)

type CookieStore interface {
	GetCookies() http.CookieJar
	SetCookies(jar http.CookieJar)
}

type DefaultCookieStore struct {
	cookie      http.CookieJar
	lastUpdated time.Time
}

func (s *DefaultCookieStore) GetCookies() http.CookieJar {
	return s.cookie
}

func (s *DefaultCookieStore) SetCookies(jar http.CookieJar) {
	s.cookie = jar
	s.lastUpdated = time.Now()
}
