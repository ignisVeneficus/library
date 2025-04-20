package api

import (
	"context"
	"database/sql"
	"github.com/ignisVeneficus/library/db"
	"github.com/ignisVeneficus/library/db/dao"
	"github.com/ignisVeneficus/library/db/dbo"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const seriesPageQty = 100

type Series struct {
	SeriesId NullNumber `json:"id"`
	Name     string     `json:"name"`
	Seqno    NullNumber `json:"seqno"`
	Url      NullString `json:"url"`
}

type ListSeries struct {
	Series
	BookQty int `json:"books"`
}
type SeriesResponse struct {
	Pagination Pagination   `json:"pagination"`
	Filters    []Filter     `json:"filter"`
	Series     []ListSeries `json:"result"`
}

func convertDBOSeriesToApiSeries(dbo dbo.Series) Series {
	return Series{
		SeriesId: NullNumber{dbo.SeriesId},
		Name:     dbo.Title,
		Url:      NullString{dbo.Url},
		Seqno:    NullNumber{sql.NullInt64{Valid: false}},
	}
}
func convertDBOBookSeriesToApiSeries(dbo dbo.BookSeries) Series {
	return Series{
		SeriesId: NullNumber{dbo.SeriesId},
		Name:     dbo.Title,
		Url:      NullString{dbo.Url},
		Seqno:    NullNumber{dbo.Seqno},
	}
}
func convertApiSeriesToDBOBookSeries(api Series) dbo.BookSeries {
	return dbo.BookSeries{
		Series: dbo.Series{
			SeriesId: api.SeriesId.NullInt64,
			Title:    api.Name,
			Url:      api.Url.NullString,
		},
		Seqno: api.Seqno.NullInt64,
	}
}
func convertDBOListSeriesToApiListSeries(dbo dbo.ListSeries) ListSeries {
	return ListSeries{
		Series:  convertDBOSeriesToApiSeries(dbo.Series),
		BookQty: dbo.BookQty,
	}
}

func GetAllSeries(c *gin.Context) {
	baseUrl := c.FullPath() + "?"
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	page = page - 1
	if page < 0 {
		page = 0
	}
	query := c.DefaultQuery("q", "")
	if query != "" {
		baseUrl += "q=" + query
	}

	log.Logger.Debug().Str("Query", query).Int("Page", page).Msg("Start Api.GetAllSeries")

	ctx := context.Background()
	database := db.GetDatabase()

	// query

	qtySeries, err := dao.GetSeriesQty(database, query, ctx)
	if err != nil {
		log.Logger.Debug().Str("Query", query).Int("Page", page).Err(err).Msg("Api.GetAllSeries Failed")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	dboSeries, err := dao.QueryAllSeries(database, ctx, query, int64(page*seriesPageQty), seriesPageQty)
	if err != nil {
		log.Logger.Debug().Str("Query", query).Int("Page", page).Err(err).Msg("Api.GetAllSeries Failed")
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	pagination := getPagination(baseUrl, qtySeries, page, authorPageQty)
	log.Logger.Trace().Str("Query", query).Int("Series", len(dboSeries)).Msg("Got the Series")

	series := make([]ListSeries, len(dboSeries))
	for i, dboSeries := range dboSeries {
		series[i] = convertDBOListSeriesToApiListSeries(dboSeries)
	}
	ret := SeriesResponse{
		Pagination: pagination,
		Series:     series,
	}
	if query != "" {
		ret.Filters = make([]Filter, 1)
		ret.Filters[0] = Filter{
			FilterType:  "Title",
			FilterValue: `"*` + query + `*"`,
		}
	}

	c.IndentedJSON(http.StatusOK, ret)
	log.Logger.Debug().Str("Query", query).Msg("End Api.GetAllSeries")
}
