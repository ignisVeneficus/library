package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ignisVeneficus/library/db"
	"github.com/ignisVeneficus/library/db/dao"
	"github.com/ignisVeneficus/library/db/dbo"
	"github.com/rs/zerolog/log"
)

type Tag struct {
	TagId NullNumber `json:"id"`
	Name  string     `json:"name"`
	Color NullString `json:"color"`
}

func convertDBOTagToApiTag(dbo dbo.Tag) Tag {
	return Tag{
		TagId: NullNumber{dbo.TagId},
		Name:  dbo.Name,
		Color: NullString{dbo.Color},
	}
}
func convertApiTagToDBOTag(api Tag) dbo.Tag {
	return dbo.Tag{
		TagId: api.TagId.NullInt64,
		Name:  api.Name,
		Color: api.Color.NullString,
	}
}
func QueryTags(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	ctx := context.Background()
	database := db.GetDatabase()
	log.Logger.Debug().Str("Query", query).Msg("Start Api.QueryTags")

	dboAutocompletes, err := dao.QueryTagsAutocomplete(database, ctx, query)
	if err != nil {
		log.Logger.Debug().Str("Query", query).Err(err).Msg("Api.QueryTags Failed")
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	result := CreateSuggestion(query, dboAutocompletes)
	c.IndentedJSON(http.StatusOK, result)
	log.Logger.Debug().Str("Query", query).Msg("End Api.QueryTags")
}
