package test

import (
	"testing"

	"github.com/EnergoStalin/torrent-feed/pkg/auth"
	"github.com/EnergoStalin/torrent-feed/pkg/client"
)

func TestSearch(t *testing.T) {
	login, password, err := getCreds()
	if err != nil {
		t.Error(err)
	}

	creds := client.NewAnimeLayerClient(auth.NewLoginCredentials(login, password, nil))

	_, err = creds.Search("hik")
	if err != nil {
		t.Error(err)
	}
}
