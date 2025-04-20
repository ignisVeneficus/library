package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Pagination struct {
	Qty          int    `json:"qty"`
	Pages        int    `json:"pages"`
	PerPage      int    `json:"perPage"`
	SelectedPage int    `json:"selectedPage"`
	BaseUrl      string `json:"base"`
}
type Filter struct {
	FilterType  string `json:"type"`
	FilterValue string `json:"value"`
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
