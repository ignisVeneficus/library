package api

import (
	"fmt"
	"net/http"

	"github.com/ignisVeneficus/library/webscraper"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// result of the web-scraper
type ScraperResult struct {
	//title of the book
	Title string `json:"title"`
	//list of authors
	Authors []ScraperAuthor `json:"authors"`
	//list of the series
	Series []ScraperSeries `json:"series"`
	//list of the tags
	Tags []string `json:"tags"`
	//blurb of the book
	Blurb string `json:"blurb"`
}

// Author of the scraped book
type ScraperAuthor struct {
	//external Url
	Url string `json:"url"`
	//name of the author
	Name string `json:"name"`
}

// Series of the scraped book
type ScraperSeries struct {
	//external Url
	Url string `json:"url"`
	//name of the series
	Name string `json:"name"`
	//book's position in the series
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

// Scraper endpoint
// @Summary Scrape and process an url
// @Description The endpont give a processed url andpoint as a book
// @Param url query string true "The url need to scrape"
// @Produce json
// @Success 200 {object} ScraperResult
// @Router /api/scraper/ [get]
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
