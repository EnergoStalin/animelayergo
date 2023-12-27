package auth

import lib "github.com/EnergoStalin/animelayergo/pkg"

const authUrl = lib.BaseUrl + "/auth/login/"

var authCookies = []string{
	"layer_id",
	"layer_hash",
	"PHPSESSID",
}
