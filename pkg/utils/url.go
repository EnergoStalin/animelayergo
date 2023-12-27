package utils

import "net/url"

func ShouldParseURL(urL string) *url.URL {
	url, _ := url.Parse(urL)

	return url
}
