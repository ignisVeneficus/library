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

// book tag item
type Tag struct {
	//id fo the tag
	TagId NullNumber `json:"id" swaggertype:"integer"`
	//name of the tag
	Name string `json:"name"`
	//Color of the tag
	Color NullString `json:"color" swaggertype:"string"`
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

// QueryTags endpoint
// @Summary Query Tags for autocomplete
// @Description The endpont give back a list of tags for the given query string
// @Param query query string false "A query string used to filter tags based on their title."
// @Tags tags
// @Produce json
// @Success 200 {object} Suggestions
// @Router /query/tags [get]
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
