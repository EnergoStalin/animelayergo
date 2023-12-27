package client

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/EnergoStalin/animelayergo/pkg/auth"
	"github.com/EnergoStalin/animelayergo/pkg/parser"
	"github.com/anacrolix/torrent/metainfo"
)

type AnimeLayer struct {
	credentials auth.Credentials
	c           http.Client
}

func (a *AnimeLayer) expandTorrent(info []parser.AnimeInfo) (err error) {
	for i := 0; i < len(info); i++ {
		elem := &info[i]
		req, _ := http.NewRequest(http.MethodGet, elem.DownloadLink, nil)
		a.applyCookie(req)

		res, err := a.c.Do(req)
		if err != nil {
			return err
		}

		meta, err := metainfo.Load(res.Body)
		if err != nil {
			return err
		}

		magnet := meta.Magnet(nil, nil)
		elem.MagnetUri = magnet.String()
		elem.Hash = magnet.InfoHash.String()
	}

	return nil
}

func (a *AnimeLayer) applyCookie(req *http.Request) (err error) {
	jar, err := a.credentials.GetCookies()
	if err != nil {
		return
	}

	for _, cookie := range jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	return
}

func (a *AnimeLayer) Search(searchTerm string) (info []parser.AnimeInfo, err error) {
	query := url.Values{
		"q":      []string{searchTerm},
		"sort":   []string{"seeders"},
		"action": []string{"list"},
		"page":   []string{"1"},
	}

	req, _ := http.NewRequest(http.MethodPost, searchUrl, strings.NewReader(query.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	err = a.applyCookie(req)
	if err != nil {
		return
	}

	res, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}

	info, err = parser.ParseEntries(res.Body)
	if err != nil {
		return nil, err
	}

	a.expandTorrent(info)

	return
}

func NewAnimeLayerClient(credentials auth.Credentials) *AnimeLayer {
	return &AnimeLayer{
		credentials: credentials,
		c:           http.Client{},
	}
}
