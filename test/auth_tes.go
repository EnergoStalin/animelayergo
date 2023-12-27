package test

import (
	"fmt"
	"testing"

	"github.com/EnergoStalin/torrent-feed/pkg/auth"
)

func TestLogin(t *testing.T) {
	login, password, err := getCreds()
	if err != nil {
		t.Error(err)
	}

	creds := auth.NewLoginCredentials(login, password, nil)

	cookies, _ := creds.GetCookies()
	fmt.Printf("%+v\n", cookies)
}
