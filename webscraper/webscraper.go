package webscraper

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

var ErrScraperNotFound = errors.New("scraper not found")

type Link struct {
	Value string
	Url   string
}
type SeriesLink struct {
	Value    string
	Url      string
	Seqno    int
	HasSeqno bool
}

type Metadata struct {
	Blurb   string
	Authors []Link
	Series  []SeriesLink
	Title   string
	Tags    []string
}

type Scraper interface {
	Name() string
	Scrape(url string) (Metadata, error)
	CheckUrl(url string) bool
}

var Registry = make(map[string]Scraper)

func RegisterScraper(p Scraper) {
	Registry[p.Name()] = p
	log.Logger.Info().Str("Scraper", p.Name()).Msg("Register a scraper")
}

func Scrape(url string) (Metadata, error) {
	for _, scraper := range Registry {
		if scraper.CheckUrl(url) {
			return scraper.Scrape(url)
		}
	}
	return Metadata{}, fmt.Errorf("%w, for url: %s", ErrScraperNotFound, url)
}
