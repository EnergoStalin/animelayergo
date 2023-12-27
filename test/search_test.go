package test

import (
	"testing"

	"github.com/EnergoStalin/animelayergo/pkg/auth"
	"github.com/EnergoStalin/animelayergo/pkg/client"
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
