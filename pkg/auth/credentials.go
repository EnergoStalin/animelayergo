package auth

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/EnergoStalin/torrent-feed/pkg/utils"
)

func findLayerIdCookie(cookies []*http.Cookie) *http.Cookie {
	for _, cookie := range cookies {
		name := strings.ToLower(cookie.Name)
		if strings.Contains(name, "layer_id") {
			return cookie
		}
	}

	return nil
}

func filterCookies(cookies []*http.Cookie) []*http.Cookie {
	var jar []*http.Cookie
	for _, cookie := range cookies {
		if slices.Contains(authCookies, cookie.Name) {
			jar = append(jar, cookie)
		}
	}

	return jar
}

type Credentials interface {
	GetCookies() (http.CookieJar, error)
}

type LoginCredentials struct {
	login    string
	password string

	store CookieStore
}

func NewLoginCredentials(login string, password string, store CookieStore) *LoginCredentials {
	if store == nil {
		store = &DefaultCookieStore{}
	}
	return &LoginCredentials{
		login:    login,
		password: password,
		store:    store,
	}
}

func (c *LoginCredentials) expired() bool {
	cookies := c.store.GetCookies()
	if cookies == nil {
		return true
	}

	cookie := findLayerIdCookie(cookies.Cookies(utils.ShouldParseURL(authUrl)))
	if cookies == nil {
		return true
	}

	return time.Now().After(cookie.Expires)
}

func (c *LoginCredentials) refreshCookies() error {
	h := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := h.PostForm(authUrl, url.Values{
		"login":    []string{c.login},
		"password": []string{c.password},
	})
	if err != nil {
		return errors.Join(errors.New("login request failed"), err)
	}

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(
		utils.ShouldParseURL(authUrl),
		filterCookies(res.Cookies()),
	)
	c.store.SetCookies(jar)

	return nil
}

func (c *LoginCredentials) GetCookies() (jar http.CookieJar, err error) {
	if c.expired() {
		err = c.refreshCookies()
		if err != nil {
			return
		}
	}

	jar = c.store.GetCookies()

	return
}
