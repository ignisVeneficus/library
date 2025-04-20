package webscraper

import (
	"ignis/library/server/utils"
	"regexp"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/rs/zerolog/log"
)

const urlRoot = "https://moly.hu"

type Moly struct {
}

func (m *Moly) Name() string {
	return "Moly.hu"
}
func (m *Moly) CheckUrl(url string) bool {
	return strings.HasPrefix(url, urlRoot)
}

func (m *Moly) Scrape(url string) (Metadata, error) {
	log.Logger.Debug().Str("Url", url).Msg("Start Scraping")
	ret := Metadata{
		Authors: make([]Link, 0),
		Series:  make([]SeriesLink, 0),
	}
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		log.Logger.Error().Err(err).Str("Url", url).Msg("Scraping faild")
		return ret, err
	}
	// content
	contentNodes, err := htmlquery.Query(doc, "//div[@id='content']")
	if err != nil {
		log.Logger.Error().Err(err).Str("Url", url).Msg("Scraping faild")
		return ret, err
	}

	// blurb
	blurbNodes, err := htmlquery.QueryAll(doc, "//div[@id='full_description']")
	if err != nil {
		log.Logger.Error().Err(err).Str("Url", url).Msg("Scraping faild")
		return ret, err
	}
	blurb := ""
	for _, b := range blurbNodes {
		blurb += htmlquery.InnerText(b)
	}

	log.Logger.Trace().Str("Blurb", blurb).Msg("Blurb in the description node")
	//diferent type of blurb
	if blurb == "" {
		blurbNodes, err := htmlquery.QueryAll(contentNodes, "/div[@class='text']")
		if err != nil {
			log.Logger.Error().Err(err).Str("Url", url).Msg("Scraping faild")
			return ret, err
		}
		for _, b := range blurbNodes {
			blurb += htmlquery.InnerText(b) + "\n"
		}
		tblurb := ""
		if len(blurb) > 50 {
			tblurb = blurb[0:50]
		} else {
			tblurb = blurb
		}
		log.Logger.Trace().Str("Url", url).Str("Blurb", tblurb).Msg("Blurb in the text nodes")
	}
	blurb = utils.CleanString(blurb)
	ret.Blurb = strings.Replace(blurb, "Vigyázat! Cselekményleírást tartalmaz.\n", "", 1)
	// title
	titleNodes, err := htmlquery.QueryAll(contentNodes, "//h1[@class='book']//span[@class='item']/text()")
	if err != nil {
		log.Logger.Error().Err(err).Str("Url", url).Msg("Scraping faild")
		return ret, err
	}
	title := ""
	for _, tn := range titleNodes {
		title += htmlquery.InnerText(tn)
	}
	log.Logger.Trace().Str("Url", url).Str("Title", title).Msg("Found Title")
	ret.Title = strings.TrimSpace(utils.CleanString(title))
	// authors
	authorNodes1, err := htmlquery.QueryAll(contentNodes, "//div[@class='authors']//a[not(following-sibling::span[@class='data'])]")
	if err != nil {
		log.Logger.Error().Err(err).Str("Url", url).Msg("Scraping faild")
		return ret, err
	}
	authors := make([]Link, 0)
	for _, an := range authorNodes1 {
		nodeUrl := urlRoot + htmlquery.SelectAttr(an, "href")
		value := utils.CleanString(htmlquery.InnerText(an))
		log.Logger.Trace().Str("Url", url).Str("Author", value).Msg("Found author")
		link := Link{
			Url:   nodeUrl,
			Value: value,
		}
		authors = append(authors, link)
	}
	ret.Authors = authors

	series := make([]SeriesLink, 0)
	seriesNode1, err := htmlquery.QueryAll(contentNodes, "//h1[@class='book']//span[@class='item']//a")
	if err != nil {
		log.Logger.Error().Err(err).Str("Url", url).Msg("Scraping faild")
		return ret, err
	}
	regEx := regexp.MustCompile(`(.*)( (\d+)\.)`)
	for _, an := range seriesNode1 {
		nodeUrl := urlRoot + htmlquery.SelectAttr(an, "href")
		text := strings.Trim(utils.CleanString(htmlquery.InnerText(an)), "()")
		match := regEx.FindStringSubmatch(text)
		name := ""
		if len(match) > 2 {
			name = match[1]
		} else {
			name = text
		}
		seriesItem := SeriesLink{
			Url:      nodeUrl,
			Value:    name,
			HasSeqno: false,
		}
		if len(match) > 3 {
			seqno, err := strconv.Atoi(match[3])
			if err == nil {
				seriesItem.Seqno = seqno
				seriesItem.HasSeqno = true
			}
		}
		log.Logger.Trace().Str("Url", url).Str("Series", name).Msg("Found series")
		series = append(series, seriesItem)

	}
	ret.Series = series
	tags := make([]string, 0)
	// authors + series from the description
	otherNodes, err := htmlquery.QueryAll(contentNodes, "//p/a[not(@onclick)]")
	if err != nil {
		log.Logger.Error().Err(err).Str("Url", url).Msg("Scraping faild")
		return ret, err
	}
	for _, on := range otherNodes {
		nodeUrl := htmlquery.SelectAttr(on, "href")
		fullUrl := urlRoot + nodeUrl
		value := utils.CleanString(htmlquery.InnerText(on))
		if strings.HasPrefix(nodeUrl, "/alkotok/") {
			link := Link{
				Url:   fullUrl,
				Value: value,
			}
			log.Logger.Trace().Str("Url", url).Str("Author", value).Msg("Found author in heap")
			authors = append(authors, link)
		}
		if strings.HasPrefix(nodeUrl, "/sorozatok/") {
			seriesLink := SeriesLink{
				Url:      fullUrl,
				Value:    value,
				HasSeqno: false,
			}
			series = append(series, seriesLink)
			log.Logger.Trace().Str("Url", url).Str("Series", value).Msg("Found series in heap")
		}
		if strings.HasPrefix(nodeUrl, "/cimkek/") {
			tags = append(tags, value)
			log.Logger.Trace().Str("Url", url).Str("Tag", value).Msg("Found tag in heap")
		}

	}
	ret.Authors = authors
	ret.Series = series
	ret.Tags = tags
	log.Logger.Debug().Str("Url", url).Msg("End Scraping")

	return ret, nil
}

// Automatikus regisztráció init-ben
func init() {
	log.Logger.Info().Msg("Moly.hu scraper loaded")
	RegisterScraper(&Moly{})
}
