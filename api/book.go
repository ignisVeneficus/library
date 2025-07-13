package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ignisVeneficus/library/db"
	"github.com/ignisVeneficus/library/db/dao"
	"github.com/ignisVeneficus/library/db/dbo"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const bookPageQty = 50 //100

// Book metadata
type Book struct {
	//id of the book
	BookId NullNumber `json:"id" swaggertype:"integer"`
	//Title of the book (if present)
	Title NullString `json:"title"  swaggertype:"string"`
	//Book has cover image
	HasCover bool `json:"hasCover"`
	//If has cover image, itas the avarage cover color
	CoverColor NullString `json:"coverColor" swaggertype:"string"`
	//external Url of the book
	Url NullString `json:"url" swaggertype:"string"`
	//list of authors
	Authors []Author `json:"authors"`
	//list of series
	Series []Series `json:"series"`
	//list of tags
	Tags []Tag `json:"tags"`
	//cover filetype (if present cover image)
	CoverType NullString `json:"coverType" swaggertype:"string"`
	//file of the book
	File string `json:"file"`
	//blurb if the book
	Blurb NullString `json:"blurb" swaggertype:"string"`
	//True if the metadata at least once has be saved (it will not overwriten by any automatic process)
	Edited int `json:"edited"`
	//type of the ebook: epub or mobipocket (prc, azw)
	FileType string `json:"fileType"`
}

// Extra information about the filtering, example fiulter by author: terry
type ExtraDisplay struct {
	//type, description of extra information
	Type string `json:"type"`
	//data for the filtering
	Data int `json:"data"`
}

// Response of the book query
type BookResponse struct {
	//list pagination
	Pagination Pagination `json:"pagination"`
	//filter/query information
	Filters []Filter `json:"filter"`
	//list opf the books
	Books []Book `json:"result"`
	//Extra display information
	Display ExtraDisplay `json:"display"`
}

func convertDBOBookToApiBook(dbo dbo.Book) Book {
	book := Book{
		BookId:     NullNumber{dbo.Bookid},
		Title:      NullString{dbo.Title},
		HasCover:   dbo.Hascover > 0,
		CoverColor: NullString{dbo.CoverColor},
		Url:        NullString{dbo.Url},
		CoverType:  NullString{dbo.CoverType},
		File:       dbo.File,
		Edited:     int(dbo.Edited),
		Blurb:      NullString{dbo.Blurb},
		FileType:   dbo.Format,
	}
	log.Logger.Trace().Str("coverType", dbo.CoverType.String).Bool("valid", dbo.CoverType.Valid).Send()
	authors := make([]Author, len(dbo.Authors))
	for i, dboAuthor := range dbo.Authors {
		authors[i] = convertDBOAuthorToApiAuthor(dboAuthor)
	}
	book.Authors = authors
	series := make([]Series, len(dbo.Series))
	for i, dboSeries := range dbo.Series {
		series[i] = convertDBOBookSeriesToApiSeries(dboSeries)
	}
	book.Series = series
	tags := make([]Tag, len(dbo.Tags))
	for i, dboTag := range dbo.Tags {
		tags[i] = convertDBOTagToApiTag(dboTag)
	}
	book.Tags = tags
	return book
}

func convertApiBookToDBOBook(api Book) dbo.Book {
	ret := dbo.Book{
		Bookid:     api.BookId.NullInt64,
		Title:      api.Title.NullString,
		CoverColor: api.CoverColor.NullString,
		Url:        api.Url.NullString,
		CoverType:  api.CoverType.NullString,
		File:       api.File,
		Edited:     int32(api.Edited),
		Hascover:   0,
		Blurb:      api.Blurb.NullString,
		Format:     api.FileType,
	}
	if api.HasCover {
		ret.Hascover = 1
	}
	ret.Authors = make([]dbo.Author, len(api.Authors))
	for i, apiAuthors := range api.Authors {
		ret.Authors[i] = convertApiAuthorToDBOAuthor(apiAuthors)
	}
	ret.Series = make([]dbo.BookSeries, len(api.Series))
	for i, apiSeries := range api.Series {
		ret.Series[i] = convertApiSeriesToDBOBookSeries(apiSeries)
	}
	ret.Tags = make([]dbo.Tag, len(api.Tags))
	for i, apiTag := range api.Tags {
		ret.Tags[i] = convertApiTagToDBOTag(apiTag)
	}

	return ret
}
func convertAllDBOBookToAPIBook(dboBooks []dbo.Book) []Book {
	books := make([]Book, len(dboBooks))
	for i, dboBook := range dboBooks {
		books[i] = convertDBOBookToApiBook(dboBook)
	}
	return books
}

func GetAllBookAsJSON() ([]byte, error) {
	database := db.GetDatabase()
	ctx := context.Background()
	books, err := dao.QueryAllBooks(database, ctx)
	if err != nil {
		return nil, err
	}
	apiBooks := convertAllDBOBookToAPIBook(books)
	ret, err := json.Marshal(apiBooks)
	return ret, err
}

func getBookByAuthor(baseUrl string, page int, authorId int) (BookResponse, error) {
	log.Logger.Debug().Int("Author", authorId).Int("Page", page).Msg("Start Get Book by AuthorId")
	ctx := context.Background()
	database := db.GetDatabase()
	qtyBooks, err := dao.GetBookByAuthorIdQty(database, ctx, int64(authorId))
	if err != nil {
		log.Logger.Error().Int("Author", authorId).Int("Page", page).Err(err).Msg("Get Book by AuthorId Failed")
		return BookResponse{}, err
	}
	author, err := dao.GetAuthorsById(database, ctx, int64(authorId))
	if err != nil {
		log.Logger.Error().Int("Author", authorId).Int("Page", page).Err(err).Msg("Get Book by AuthorId Failed")
		return BookResponse{}, err
	}
	dboBooks, err := dao.QueryAllBookByAuthorId(database, ctx, int64(authorId), int64(page*bookPageQty), bookPageQty)
	if err != nil {
		log.Logger.Error().Int("Author", authorId).Int("Page", page).Err(err).Msg("Get Book by AuthorId Failed")
		return BookResponse{}, err
	}
	pagination := getPagination(baseUrl, qtyBooks, page, bookPageQty)
	log.Logger.Trace().Int("Books", len(dboBooks)).Msg("Got the Books")
	books := convertAllDBOBookToAPIBook(dboBooks)
	filters := make([]Filter, 1)
	filters[0] = Filter{
		FilterType:  "Author",
		FilterValue: `"` + author.Name + `"`,
	}

	ret := BookResponse{
		Pagination: pagination,
		Filters:    filters,
		Books:      books,
	}
	log.Logger.Debug().Int("Author", authorId).Int("Page", page).Msg("End Get Book by AuthorId")
	return ret, nil

}
func getBookBySeries(baseUrl string, page int, seriesId int) (BookResponse, error) {
	log.Logger.Debug().Int("Series", seriesId).Int("Page", page).Msg("Start Get Book by SeriesId")
	ctx := context.Background()
	database := db.GetDatabase()
	qtyBooks, err := dao.GetBookBySeriesIdQty(database, ctx, int64(seriesId))
	if err != nil {
		log.Logger.Error().Int("Series", seriesId).Int("Page", page).Err(err).Msg("Get Book by SeriesId Failed")
		return BookResponse{}, err
	}
	series, err := dao.GetSeriesById(database, ctx, int64(seriesId))
	if err != nil {
		log.Logger.Error().Int("Series", seriesId).Int("Page", page).Err(err).Msg("Get Book by SeriesId Failed")
		return BookResponse{}, err
	}
	dboBooks, err := dao.QueryAllBookBySeriesId(database, ctx, int64(seriesId), int64(page*bookPageQty), bookPageQty)
	if err != nil {
		log.Logger.Error().Int("Series", seriesId).Int("Page", page).Err(err).Msg("Get Book by SeriesId Failed")
		return BookResponse{}, err
	}
	pagination := getPagination(baseUrl, qtyBooks, page, bookPageQty)
	log.Logger.Trace().Int("Books", len(dboBooks)).Msg("Got the Books")
	books := convertAllDBOBookToAPIBook(dboBooks)
	filters := make([]Filter, 1)
	filters[0] = Filter{
		FilterType:  "Series",
		FilterValue: `"` + series.Title + `"`,
	}

	ret := BookResponse{
		Pagination: pagination,
		Filters:    filters,
		Books:      books,
		Display:    ExtraDisplay{Type: "series", Data: seriesId},
	}
	log.Logger.Debug().Int("Series", seriesId).Int("Page", page).Msg("End Get Book by SeriesId")
	return ret, nil
}

func getBookByTag(baseUrl string, page int, tagId int) (BookResponse, error) {
	log.Logger.Debug().Int("Tag", tagId).Int("Page", page).Msg("Start Get Book by TagId")
	ctx := context.Background()
	database := db.GetDatabase()
	qtyBooks, err := dao.GetBookByTagIdQty(database, ctx, int64(tagId))
	if err != nil {
		log.Logger.Error().Int("Tag", tagId).Int("Page", page).Err(err).Msg("Get Book by TagId Failed")
		return BookResponse{}, err
	}
	tag, err := dao.GetTagById(database, ctx, int64(tagId))
	if err != nil {
		log.Logger.Error().Int("Tag", tagId).Int("Page", page).Err(err).Msg("Get Book by TagId Failed")
		return BookResponse{}, err
	}
	dboBooks, err := dao.QueryAllBookByTagId(database, ctx, int64(tagId), int64(page*bookPageQty), bookPageQty)
	if err != nil {
		log.Logger.Error().Int("Tag", tagId).Int("Page", page).Err(err).Msg("Get Book by TagId Failed")
		return BookResponse{}, err
	}
	pagination := getPagination(baseUrl, qtyBooks, page, bookPageQty)
	log.Logger.Trace().Int("Books", len(dboBooks)).Msg("Got the Books")
	books := convertAllDBOBookToAPIBook(dboBooks)
	filters := make([]Filter, 1)
	filters[0] = Filter{
		FilterType:  "Tag",
		FilterValue: `"` + tag.Name + `"`,
	}

	ret := BookResponse{
		Pagination: pagination,
		Filters:    filters,
		Books:      books,
		Display:    ExtraDisplay{Type: "series", Data: tagId},
	}
	log.Logger.Debug().Int("Tag", tagId).Int("Page", page).Msg("End Get Book by TagId")
	return ret, nil

}

func getBookByQuery(baseUrl string, query string, page int) (BookResponse, error) {
	var (
		qtyBooks int64
		err      error
		dboBooks []dbo.Book
	)
	filters := make([]Filter, 0)
	log.Logger.Debug().Str("Query", query).Int("Page", page).Msg("Start Get Book by query")

	ctx := context.Background()
	database := db.GetDatabase()

	qtyBooks, err = dao.GetBookQty(database, ctx, query)
	if err == nil {
		dboBooks, err = dao.QueryBook(database, ctx, query, int64(page*bookPageQty), bookPageQty)
	}

	if err != nil {
		log.Logger.Error().Err(err).Msg("Get Book by query Error")
		return BookResponse{}, err
	}
	pagination := getPagination(baseUrl, qtyBooks, page, bookPageQty)
	log.Logger.Trace().Int("Books", len(dboBooks)).Msg("Got the Books")
	books := convertAllDBOBookToAPIBook(dboBooks)
	ret := BookResponse{
		Pagination: pagination,
		Filters:    filters,
		Books:      books,
	}
	log.Logger.Debug().Str("Query", query).Int("Page", page).Msg("End Get Book by query")
	return ret, nil
}

// GetBooks endpoint
// @Summary Query for books metadata
// @Description The endpont give back a list of book metadata with pagging, filtering information
// @Param page query int false "page number, start with 1, default is 1"
// @Param ai query int false "author ID used to filter books based on authors"
// @Param si query int false "series ID used to filter books based on series"
// @Param ti query int false "tag ID used to filter books based on tags"
// @Param q query string false "A query string used to filter books based on their description, title, author, series, and tags."
// @Tags book
// @Produce json
// @Success 200 {object} BookResponse
// @Router /book [get]
func GetAllBook(c *gin.Context) {
	baseUrl := c.FullPath() + "?"
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	page = page - 1
	if page < 0 {
		page = 0
	}
	// author by id
	author := c.Query("ai")
	if author != "" {
		baseUrl += "ai=" + author
		ai, err := strconv.Atoi(author)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Api.GetAllBook Error")
			c.JSON(http.StatusUnprocessableEntity, "")
			return
		}
		ret, err := getBookByAuthor(baseUrl, page, ai)
		if err != nil {
			log.Logger.Debug().Int("Page", page).Int("Suthor", ai).Err(err).Msg("Api.GetAllBook Failed")
			c.JSON(http.StatusInternalServerError, "")
			return
		}
		c.IndentedJSON(http.StatusOK, ret)
		log.Logger.Debug().Msg("End Api.GetAllBook")
		return

	}
	series := c.Query("si")
	if series != "" {
		baseUrl += "si=" + series
		si, err := strconv.Atoi(series)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Api.GetAllBook Error")
			c.JSON(http.StatusUnprocessableEntity, "")
			return
		}
		ret, err := getBookBySeries(baseUrl, page, si)
		if err != nil {
			log.Logger.Debug().Int("Page", page).Int("Series", si).Err(err).Msg("Api.GetAllBook Failed")
			c.JSON(http.StatusInternalServerError, "")
			return
		}
		c.IndentedJSON(http.StatusOK, ret)
		log.Logger.Debug().Msg("End Api.GetAllBook")
		return
	}
	tags := c.Query("ti")
	if tags != "" {
		baseUrl += "ti=" + tags
		ti, err := strconv.Atoi(tags)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Api.GetAllBook Error")
			c.JSON(http.StatusUnprocessableEntity, "")
			return
		}
		ret, err := getBookByTag(baseUrl, page, ti)
		if err != nil {
			log.Logger.Debug().Int("Page", page).Int("Tag", ti).Err(err).Msg("Api.GetAllBook Failed")
			c.JSON(http.StatusInternalServerError, "")
			return
		}
		c.IndentedJSON(http.StatusOK, ret)
		log.Logger.Debug().Msg("End Api.GetAllBook")
		return

	}

	query := c.Query("q")
	if query != "" {
		baseUrl += "q=" + query
	}

	// query
	log.Logger.Debug().Int("Page", page).Str("Query", query).Msg("Start Api.GetAllBook")
	ret, err := getBookByQuery(baseUrl, query, page)
	if err != nil {
		log.Logger.Debug().Int("Page", page).Str("Query", query).Err(err).Msg("Api.GetAllBook Failed")
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	c.IndentedJSON(http.StatusOK, ret)
	log.Logger.Debug().Msg("End Api.GetAllBook")
}

// GetBook endpoint
// @Summary Get a book metadata
// @Description The endpont give back the given book netadata
// @Param bookId path int true "ID of the book"
// @Tags book
// @Produce json
// @Success 200 {object} Book
// @Failure 404 {string} string "Book not found"
// @Router /book/{bookId} [get]
func GetBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Logger.Error().Err(err).Msg("Api.GetBook Error")
		c.JSON(http.StatusUnprocessableEntity, "")
		return
	}
	log.Logger.Debug().Int("Id", id).Msg("Start Api.GetBook")
	ctx := context.Background()
	database := db.GetDatabase()

	dboBook, err := dao.GetBookById(database, ctx, int64(id))
	if errors.Is(err, dao.ErrDataNotFound) {
		log.Logger.Debug().Int("Id", id).Msg("Stop Api.GetBook: not found")
		c.JSON(http.StatusNotFound, "")
		return
	}
	if err != nil {
		log.Logger.Error().Int("Id", id).Err(err).Msg("Api.GetBook Error")
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	book := convertDBOBookToApiBook(dboBook)
	c.IndentedJSON(http.StatusOK, book)
	log.Logger.Debug().Int("Id", id).Msg("Stop Api.GetBook")
}

// @Summary Update a book
// @Description Uopdate a book's metadata
// @Tags book
// @Accept json
// @Produce json
// @Param book body Book true "Book to update"
// @Success 200
// @Router /book [post]
func PostBook(c *gin.Context) {
	var book Book
	c.BindJSON(&book)
	log.Logger.Debug().Int64("Id", book.BookId.Int64).Msg("Start Api.PostBook")
	dboBook := convertApiBookToDBOBook(book)

	ctx := context.Background()
	database := db.GetDatabase()

	err := dao.UpdateBook(database, ctx, dboBook)

	if err != nil {
		log.Logger.Error().Int64("Id", book.BookId.Int64).Err(err).Msg("Api.PostBook Error")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	fmt.Println(book)
	c.IndentedJSON(http.StatusOK, "")
	log.Logger.Debug().Int64("Id", book.BookId.Int64).Msg("End Api.PostBook")
}
