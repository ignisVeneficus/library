package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"math"
	"net/http"

	"github.com/ignisVeneficus/library/db/dbo"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// pagination structure
type Pagination struct {
	//number of the all results
	Qty int `json:"qty"`
	//number of tha pages
	Pages int `json:"pages"`
	//results per page
	PerPage int `json:"perPage"`
	//selected, actual page number
	SelectedPage int `json:"selectedPage"`
	//bas url for reloding
	BaseUrl string `json:"base"`
}

// filter response
type Filter struct {
	// type of the filter
	FilterType string `json:"type"`
	// value of the filter
	FilterValue string `json:"value"`
}

// autocomplete data for replace item in UI
type SuggestionData struct {
	//id of the data
	Id int `json:"id"`
	//Text of the data
	Name string `json:"name"`
	//url of the data
	Url string `json:"url"`
}

// packed structure for autocomplete
type SuggestionItem struct {
	//display text
	Value string `json:"value"`
	//data for processing
	Data SuggestionData `json:"data"`
}

// response of the autocomplete apis
type Suggestions struct {
	//query string
	Query string `json:"query"`
	//result list
	List []SuggestionItem `json:"suggestions"`
}

func getPagination(base string, qty int64, page int, qtyPage int) Pagination {
	pageQty := int(math.Ceil(float64(qty) / float64(qtyPage)))
	pagination := Pagination{
		Qty:          int(qty),
		PerPage:      qtyPage,
		Pages:        pageQty,
		SelectedPage: page + 1,
		BaseUrl:      base,
	}
	return pagination
}

type NullNumber struct{ sql.NullInt64 }

func (nn *NullNumber) UnmarshalJSON(data []byte) error {
	var x *json.Number
	if bytes.Equal(data, []byte(`""`)) {
		nn.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	this, err := x.Int64()
	if err != nil {
		return err
	}
	nn.Valid = true
	nn.Int64 = this
	return nil
}

func (nn NullNumber) MarshalJSON() ([]byte, error) {
	if nn.Valid {
		return json.Marshal(nn.Int64)
	}
	return json.Marshal(nil)
}

type NullString struct{ sql.NullString }

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var x *string
	if bytes.Equal(data, []byte(`""`)) {
		ns.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	ns.Valid = true
	ns.String = *x
	return nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	ret := ""
	if ns.Valid {
		ret = ns.String
	}
	return json.Marshal(ret)
}

// DownloadAllBook endpoint
// @Summary create a json file for all books in the system
// @Description The endpont give back file for every ebook in the database
// @Tags book
// @Produce application/json
// @Success 200 {file} file "JSON file containing book list"
// @Router /export [get]
func DownloadAllBook(c *gin.Context) {
	log.Logger.Debug().Msg("Start Api.DownloadAllBook")
	data, err := GetAllBookAsJSON()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Api.DownloadAllBook Error")
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	c.Header("Content-Disposition", "attachment; filename=books.json")
	c.Data(http.StatusOK, "application/octet-stream", data)
	log.Logger.Debug().Msg("End Api.DownloadAllBook")
}

func CreateSuggestionItem(data SuggestionData) SuggestionItem {
	return SuggestionItem{
		Value: data.Name,
		Data:  data,
	}
}
func CreateSuggestion(query string, result []dbo.AutoComplete) Suggestions {
	list := make([]SuggestionItem, len(result))
	for i, dboAutocomplete := range result {
		list[i] = CreateSuggestionItem(SuggestionData{
			Name: dboAutocomplete.Name,
			Id:   dboAutocomplete.Id,
			Url:  dboAutocomplete.StrUrl(),
		})
	}
	return Suggestions{
		Query: query,
		List:  list,
	}

}
