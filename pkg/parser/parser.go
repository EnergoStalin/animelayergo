package parser

import (
	"io"
	"slices"
	"strconv"
	"strings"
	"time"

	lib "github.com/EnergoStalin/torrent-feed/pkg"
	"github.com/PuerkitoBio/goquery"
)

func monthFromVerb(s string) int {
	return slices.Index([]string{
		"января",
		"февраля",
		"марта",
		"апреля",
		"мая",
		"июня",
		"июля",
		"августа",
		"сентября",
		"октября",
		"ноября",
		"декабря",
	}, s) + 1
}

func parseDate(s string) time.Time {
	s = s[strings.Index(s, "н:")+4:]

	// formats := []string{"DD MMMM YYYY в HH:mm", "DD MMMM в HH:mm"}
	var day, month, year, hour, minute int
	values := strings.Split(s, " ")

	if len(values) == 4 {
		day, _ = strconv.Atoi(values[0])
		month = monthFromVerb(values[1])
		year = time.Now().Year()
		v2 := strings.Split(values[3], ":")
		hour, _ = strconv.Atoi(v2[0])
		minute, _ = strconv.Atoi(v2[1])
	} else if len(values) == 5 {
		day, _ = strconv.Atoi(values[0])
		month = monthFromVerb(values[1])
		year, _ = strconv.Atoi(values[2])
		v2 := strings.Split(values[4], ":")
		hour, _ = strconv.Atoi(v2[0])
		minute, _ = strconv.Atoi(v2[1])
	}

	loc, _ := time.LoadLocation("Europe/Moscow")
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, loc)
}

func parseSize(s string) int {
	a := strings.Split(s, " ")
	if len(a) != 2 {
		return -1
	}

	mult := map[string]int{
		"GB": 1e9,
		"MB": 1e6,
	}[a[1]]

	size, err := strconv.ParseFloat(a[0], 32)
	if err != nil {
		return -1
	}

	return int(size * float64(mult))
}

func trim(s string) string {
	return strings.Trim(strings.ReplaceAll(s, "\u00a0", " "), " \t\n")
}

func parseEntry(s *goquery.Selection) (info AnimeInfo) {
	linkElement := s.Find("h3 > a")
	link, _ := linkElement.Attr("href")
	link = lib.BaseUrl + link + "download"

	safeString := trim(s.Find("div.info").Text())
	infodiv := strings.Split(safeString, "|")

	seed, _ := strconv.Atoi(trim(infodiv[0]))
	leech, _ := strconv.Atoi(trim(infodiv[1]))
	size := parseSize(trim(infodiv[2]))
	uploader := trim(infodiv[3])
	updated := parseDate(trim(infodiv[4]))

	return AnimeInfo{
		Title:        trim(linkElement.Text()),
		Seed:         seed,
		Leech:        leech,
		Size:         size,
		Uploader:     uploader,
		Date:         updated,
		DownloadLink: link,
	}
}

func ParseEntries(body io.Reader) (info []AnimeInfo, err error) {
	doc, err := goquery.NewDocumentFromReader(body)

	torrents := doc.Find("li.torrent-item")

	info = []AnimeInfo{}
	torrents.Map(func(i int, s *goquery.Selection) (empty string) {
		info = append(info, parseEntry(s))
		return
	})

	return
}
