package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ignisVeneficus/library/db"
	"github.com/ignisVeneficus/library/db/dao"
	"github.com/ignisVeneficus/library/db/dbo"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const authorPageQty = 100

// author of a book
type Author struct {
	//Id of the author
	AuthorId NullNumber `json:"id" swaggertype:"integer"`
	//Name of the author
	Name string `json:"name"`
	//external url of the author
	Url NullString `json:"url" swaggertype:"string"`
}

// Author with book count
type ListAuthor struct {
	Author
	//book count of the author
	BookQty int `json:"books"`
}

// Response of the author query
type AuthorResponse struct {
	//list pagination
	Pagination Pagination `json:"pagination"`
	//filter/query information
	Filters []Filter `json:"filter"`
	//authors
	Authors []ListAuthor `json:"result"`
}

func convertDBOAuthorToApiAuthor(dbo dbo.Author) Author {
	return Author{
		AuthorId: NullNumber{dbo.Authorid},
		Name:     dbo.Name,
		Url:      NullString{dbo.Url},
	}
}
func convertDBOListAuthorToApiListAuthor(dbo dbo.ListAuthor) ListAuthor {
	return ListAuthor{
		Author:  convertDBOAuthorToApiAuthor(dbo.Author),
		BookQty: dbo.BookQty,
	}
}

func convertApiAuthorToDBOAuthor(api Author) dbo.Author {
	return dbo.Author{
		Authorid: api.AuthorId.NullInt64,
		Name:     api.Name,
		Url:      api.Url.NullString,
	}
}

// GetAllAuthor endpoint
// @Summary Query for Authors metadata
// @Description The endpont give back a list of authors with pagging, filtering information
// @Param page query int false "page number, start with 1, default is 1"
// @Param q query string false "A query string used to filter authors based on their names."
// @Tags author
// @Produce json
// @Success 200 {object} AuthorResponse
// @Router /author [get]
func GetAllAuthor(c *gin.Context) {
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

	log.Logger.Debug().Int("Page", page).Str("Query", query).Msg("Start Api.GetAllAuthor")

	ctx := context.Background()
	database := db.GetDatabase()

	// query

	qtyAuthor, err := dao.GetAuthorQty(database, ctx, query)
	if err != nil {
		log.Logger.Debug().Int("Page", page).Str("Query", query).Err(err).Msg("Api.GetAllAuthor Failed")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	dboAuthors, err := dao.QueryAllAuthor(database, ctx, query, int64(page*authorPageQty), authorPageQty)
	if err != nil {
		log.Logger.Debug().Int("Page", page).Str("Query", query).Err(err).Msg("Api.GetAllAuthor Failed")
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	pagination := getPagination(baseUrl, qtyAuthor, page, authorPageQty)
	log.Logger.Trace().Int("Authors", len(dboAuthors)).Msg("Got the Authors")

	authors := make([]ListAuthor, len(dboAuthors))
	for i, dboAuthor := range dboAuthors {
		authors[i] = convertDBOListAuthorToApiListAuthor(dboAuthor)
	}
	ret := AuthorResponse{
		Pagination: pagination,
		Authors:    authors,
	}
	if query != "" {
		ret.Filters = make([]Filter, 1)
		ret.Filters[0] = Filter{
			FilterType:  "Name",
			FilterValue: `"*` + query + `*"`,
		}

	}

	c.IndentedJSON(http.StatusOK, ret)
	log.Logger.Debug().Msg("End Api.GetAllAuthor")
}

// QueryAuthors endpoint
// @Summary Query Authors for autocomplete
// @Description The endpont give back a list of authors for the given query string
// @Param query query string false "A query string used to filter authors based on their names."
// @Tags author
// @Produce json
// @Success 200 {object} Suggestions
// @Router /query/author [get]
func QueryAuthors(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	ctx := context.Background()
	database := db.GetDatabase()
	log.Logger.Debug().Str("Query", query).Msg("Start Api.QueryAuthors")

	dboAutocomplete, err := dao.QueryAuthorAutocomplete(database, ctx, query)
	if err != nil {
		log.Logger.Debug().Str("Query", query).Err(err).Msg("Api.QueryAuthors Failed")
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	result := CreateSuggestion(query, dboAutocomplete)
	c.IndentedJSON(http.StatusOK, result)
	log.Logger.Debug().Str("Query", query).Msg("End Api.QueryAuthors")
}
