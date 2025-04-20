package api

import (
	"context"
	"ignis/library/server/db"
	"ignis/library/server/db/dao"
	"ignis/library/server/db/dbo"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const authorPageQty = 100

type Author struct {
	AuthorId NullNumber `json:"id"`
	Name     string     `json:"name"`
	Url      NullString `json:"url"`
}

type ListAuthor struct {
	Author
	BookQty int `json:"books"`
}
type AuthorResponse struct {
	Pagination Pagination   `json:"pagination"`
	Filters    []Filter     `json:"filter"`
	Authors    []ListAuthor `json:"result"`
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
