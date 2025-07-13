package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/ignisVeneficus/library/db"
	"github.com/ignisVeneficus/library/db/dao"
	"github.com/ignisVeneficus/library/db/dbo"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const seriesPageQty = 100

// Book series item
type Series struct {
	//Id of the series
	SeriesId NullNumber `json:"id" swaggertype:"integer"`
	//Name (title) of the series
	Name string `json:"name"`
	//Book order in the series
	Seqno NullNumber `json:"seqno" swaggertype:"integer"`
	//external Url of the series
	Url NullString `json:"url" swaggertype:"string"`
}

// Series with book count
type ListSeries struct {
	Series
	//book count in this series
	BookQty int `json:"books"`
}

// Response of the series query
type SeriesResponse struct {
	//list pagination
	Pagination Pagination `json:"pagination"`
	//filter/query information
	Filters []Filter `json:"filter"`
	//Series
	Series []ListSeries `json:"result"`
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

// GetAllSeries endpoint
// @Summary Query for Series metadata
// @Description The endpont give back a list of series with pagging, filtering information
// @Param page query int false "page number, start with 1, default is 1"
// @Param q query string false "A query string used to filter series based on their title/name."
// @Tags series
// @Produce json
// @Success 200 {object} SeriesResponse
// @Router /series [get]
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

// QuerySeries endpoint
// @Summary Query Series for autocomplete
// @Description The endpont give back a list of series for the given query string
// @Param query query string false "A query string used to filter series based on their title."
// @Tags series
// @Produce json
// @Success 200 {object} Suggestions
// @Router /query/series [get]

func QuerySeries(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	ctx := context.Background()
	database := db.GetDatabase()
	log.Logger.Debug().Str("Query", query).Msg("Start Api.QuerySeries")

	dboAutocomplete, err := dao.QuerySeriesAutocomplete(database, ctx, query)
	if err != nil {
		log.Logger.Debug().Str("Query", query).Err(err).Msg("Api.QuerySeries Failed")
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	result := CreateSuggestion(query, dboAutocomplete)
	c.IndentedJSON(http.StatusOK, result)
	log.Logger.Debug().Str("Query", query).Msg("End Api.QuerySeries")
}
