package api

import (
	"fmt"
	"github.com/ignisVeneficus/library/webscraper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ScraperResult struct {
	Title   string          `json:"title"`
	Authors []ScraperAuthor `json:"authors"`
	Series  []ScraperSeries `json:"series"`
	Tags    []string        `json:"tags"`
	Blurb   string          `json:"blurb"`
}
type ScraperAuthor struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}
type ScraperSeries struct {
	Url   string `json:"url"`
	Name  string `json:"name"`
	SeqNo string `json:"seqno"`
}

func convertScraperAuthors(authors []webscraper.Link) []ScraperAuthor {
	ret := make([]ScraperAuthor, len(authors))
	for i, au := range authors {
		ret[i] = ScraperAuthor{
			Url:  au.Url,
			Name: au.Value,
		}
	}
	return ret
}
func convertScraperSeries(series []webscraper.SeriesLink) []ScraperSeries {
	ret := make([]ScraperSeries, len(series))
	for i, se := range series {
		rse := ScraperSeries{
			Url:  se.Url,
			Name: se.Value,
		}
		if se.HasSeqno {
			rse.SeqNo = fmt.Sprintf("%d", se.Seqno)
		}
		ret[i] = rse
	}
	return ret
}
func convertScraperMetadata(metadata webscraper.Metadata) ScraperResult {
	ret := ScraperResult{
		Title:   metadata.Title,
		Blurb:   metadata.Blurb,
		Authors: convertScraperAuthors(metadata.Authors),
		Series:  convertScraperSeries(metadata.Series),
		Tags:    metadata.Tags,
	}
	return ret
}

func Scrape(c *gin.Context) {
	url := c.Query("url")
	log.Logger.Debug().Str("Url", url).Msg("Start Api.Scrape")
	metadata, err := webscraper.Scrape(url)
	if err != nil {
		log.Logger.Error().Err(err).Str("Url", url).Msg("Api.Scrape failed")
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	ret := convertScraperMetadata(metadata)
	c.IndentedJSON(http.StatusOK, ret)
	log.Logger.Debug().Str("Url", url).Msg("End Api.Scrape")

}
