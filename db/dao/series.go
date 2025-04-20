package dao

import (
	"context"
	"strings"

	"database/sql"
	"ignis/library/server/db/dbo"

	"github.com/rs/zerolog/log"
)

const querySeriesByBook = `SELECT s.seriesId, s.title, s.url, bs.sequence FROM series AS s
JOIN bookseries AS bs ON s.seriesId = bs.seriesId
WHERE bs.bookId = ?
ORDER BY s.title`

func (q *Queries) QuerySeriesByBook(ctx context.Context, bookid int64) ([]dbo.BookSeries, error) {
	rows, err := q.db.QueryContext(ctx, querySeriesByBook, bookid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dbo.BookSeries
	for rows.Next() {
		var i dbo.BookSeries
		if err := rows.Scan(&i.SeriesId, &i.Title, &i.Url, &i.Seqno); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createSeries = `INSERT INTO series (title, url) VALUES (?, ?)`

func (q *Queries) CreateSeries(ctx context.Context, arg dbo.Series) error {
	_, err := q.db.ExecContext(ctx, createSeries, arg.Title, arg.Url)
	return err
}

const updateSeries = `UPDATE series SET title=?, url=? WHERE seriesId =?`

func (q *Queries) UpdateSeries(ctx context.Context, series dbo.Series) error {
	_, err := q.db.ExecContext(ctx, updateSeries, series.Title, series.Url, series.SeriesId)
	return err
}

// bookSeries query func
const bindBookSeries = `INSERT INTO bookseries (bookId,seriesId,sequence) VALUES (?,?,?)`

func (q *Queries) BindBookSeries(ctx context.Context, bookId int64, seriesId int64, sequence sql.NullInt64) error {
	_, err := q.db.ExecContext(ctx, bindBookSeries, bookId, seriesId, sequence)
	return err
}

const updateBookSeries = `UPDATE bookseries SET sequence = ? WHERE bookId =? and seriesId=?`

func (q *Queries) UpdateBookSeries(ctx context.Context, bookId int64, seriesId int64, sequence sql.NullInt64) error {
	_, err := q.db.ExecContext(ctx, updateBookSeries, sequence, bookId, seriesId)
	return err
}

const divideAllSeriesFromBook = "DELETE FROM bookseries WHERE bookid=?"
const divideSeriesFromBookStart = "DELETE FROM bookseries WHERE bookid=? and seriesId not in (?"
const divideSeriesFromBookEnd = ")"

func (q *Queries) DivideSeriesFromBook(ctx context.Context, bookId int64, seriesids []int64) error {
	if len(seriesids) == 0 {
		_, err := q.db.ExecContext(ctx, divideAllSeriesFromBook, bookId)
		return err
	}
	args := make([]interface{}, len(seriesids)+1)
	args[0] = bookId
	for i, id := range seriesids {
		args[i+1] = id
	}
	query := divideSeriesFromBookStart + strings.Repeat(",?", len(seriesids)-1) + divideSeriesFromBookEnd
	_, err := q.db.ExecContext(ctx, query, args...)
	return err
}

const getSeries = `SELECT seriesId, title, url FROM series WHERE title = ?`

func (q *Queries) GetSeries(ctx context.Context, series string) (dbo.Series, error) {
	row := q.db.QueryRowContext(ctx, getSeries, series)
	var i dbo.Series
	err := row.Scan(&i.SeriesId, &i.Title, &i.Url)
	return i, err
}

const getSeriesUrl = `SELECT seriesId, title, url FROM series WHERE title = ? AND (url IS null OR url= ?)`

func (q *Queries) GetSeriesURL(ctx context.Context, series string, url string) (dbo.Series, error) {
	row := q.db.QueryRowContext(ctx, getSeriesUrl, series, url)
	var i dbo.Series
	err := row.Scan(&i.SeriesId, &i.Title, &i.Url)
	return i, err
}

const getSeriesById = `SELECT seriesId, title, url FROM series WHERE seriesId = ?`

func (q *Queries) GetSeriesById(ctx context.Context, seriesId int64) (dbo.Series, error) {
	row := q.db.QueryRowContext(ctx, getSeriesById, seriesId)
	var i dbo.Series
	err := row.Scan(&i.SeriesId, &i.Title, &i.Url)
	return i, err
}

const queryAllSeriesBegin = `SELECT s.seriesid, s.title, s.url, ifnull(bs.Cnt, 0) as books FROM series AS S
LEFT JOIN 
	(SELECT seriesId, count(1) Cnt FROM bookseries 
	GROUP BY seriesId) As bs
	ON bs.seriesId = s.seriesId
`
const queryAllSeriesEnd = `ORDER BY s.title
LIMIT ?,?`

const queryAllSeriesWhere = `WHERE s.title like concat('%',?,'%')
`
const queryAllSeries = queryAllSeriesBegin + queryAllSeriesEnd
const queryAllSeriesFilter = queryAllSeriesBegin + queryAllSeriesWhere + queryAllSeriesEnd

func (q *Queries) QueryAllSeries(ctx context.Context, title string, from int64, qty int64) ([]dbo.ListSeries, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if title == "" {
		rows, err = q.db.QueryContext(ctx, queryAllSeries, from, qty)
	} else {
		rows, err = q.db.QueryContext(ctx, queryAllSeriesFilter, title, from, qty)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dbo.ListSeries
	for rows.Next() {
		var i dbo.ListSeries
		if err := rows.Scan(&i.SeriesId, &i.Title, &i.Url, &i.BookQty); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSeriesQty = `SELECT count(*) FROM series AS s `
const getSeriesQtyFilter = getSeriesQty + queryAllSeriesWhere

func (q *Queries) GetSeriesQty(ctx context.Context, title string) (int64, error) {
	var row *sql.Row
	if title == "" {
		row = q.db.QueryRowContext(ctx, getSeriesQty)
	} else {
		row = q.db.QueryRowContext(ctx, getSeriesQtyFilter, title)

	}
	var count int64
	err := row.Scan(&count)
	return count, err
}

// ///////////////////////
func UpdateBookSeriesTBLK(q *Queries, ctx context.Context, bookId int64, series dbo.BookSeries) error {
	log.Logger.Debug().Int64("Series", series.SeriesId.Int64).Msg("Start UpdateSeriesTBLK")
	err := q.UpdateSeries(ctx, series.Series)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Series", series.SeriesId.Int64).Msg("UpdateSeriesTBLK failed")
		return err
	}
	err = q.UpdateBookSeries(ctx, bookId, series.SeriesId.Int64, series.Seqno)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Series", series.SeriesId.Int64).Msg("UpdateSeriesTBLK failed")
		return err
	}

	log.Logger.Debug().Int64("Series", series.SeriesId.Int64).Msg("End UpdateSeriesTBLK")
	return nil
}

func DivideBookAllOtherSeries(q *Queries, ctx context.Context, bookId int64, seriesId []int64) error {
	log.Logger.Debug().Int64("Book", bookId).Ints64("SeriesIds", seriesId).Msg("Start Divide Book All Other Author")
	err := q.DivideSeriesFromBook(ctx, bookId, seriesId)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Book", bookId).Ints64("SeriesIds", seriesId).Msg("Divide Book All Other Author Failed")
		return err
	}

	log.Logger.Debug().Int64("Book", bookId).Ints64("SeriesIds", seriesId).Msg("End Divide Book All Other Author")
	return nil
}

func AddSeriesToBookTBLK(q *Queries, ctx context.Context, bookId int64, series dbo.Series, seqno sql.NullInt64) error {
	log.Logger.Debug().Int64("Book", bookId).Str("Series", series.Title).Int64("SeriesId", series.SeriesId.Int64).Bool("SeriesId is Valid", series.SeriesId.Valid).Msg("Start Add Series")
	var err error
	if !series.SeriesId.Valid {
		seriesTitle := series.Title
		seriesUrl := series.Url
		if series.Url.Valid && series.Url.String != "" {
			series, err = q.GetSeriesURL(ctx, seriesTitle, seriesUrl.String)
		} else {
			series, err = q.GetSeries(ctx, seriesTitle)
		}
		if err == sql.ErrNoRows {
			log.Logger.Trace().Str("Series", seriesTitle).Msg("Series not found, create")
			series = dbo.Series{Title: seriesTitle, Url: seriesUrl}
			err = q.CreateSeries(ctx, series)
			if err != nil {
				log.Logger.Error().Err(err).Str("Series", series.Title).Msg("Create Series failed")
				return err
			}
			seriesId, err := q.GetLastId(ctx)
			if err != nil {
				log.Logger.Error().Err(err).Str("Series", series.Title).Msg("Create Series failed")
				return err
			}
			series.SeriesId = sql.NullInt64{Int64: seriesId, Valid: true}
		} else if err != nil {
			log.Logger.Error().Err(err).Str("Series", series.Title).Msg("Get Series failed")
			return err
		}
	}
	err = q.BindBookSeries(ctx, bookId, series.SeriesId.Int64, seqno)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Book", bookId).Str("Series", series.Title).Int64("SeriesId", series.SeriesId.Int64).Msg("Binding Series to Book failed")
		return err
	}
	log.Logger.Debug().Int64("Book", bookId).Str("Series", series.Title).Int64("SeriesId", series.SeriesId.Int64).Msg("End Add Series")
	return nil
}
func QuerySeriesByBookTBLK(q *Queries, ctx context.Context, bookId int64) ([]dbo.BookSeries, error) {
	log.Logger.Debug().Int64("Book", bookId).Msg("Start Query Series by Book")
	series, err := q.QuerySeriesByBook(ctx, bookId)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Book", bookId).Msg("Query Series by Book failed")
		return nil, err
	}
	log.Logger.Debug().Int64("Book", bookId).Msg("End Query Series by Book ")
	return series, nil
}
func GetSeriesById(db *sql.DB, ctx context.Context, seriesId int64) (dbo.Series, error) {
	log.Logger.Debug().Int64("Series", seriesId).Msg("Start Get Series By Id")
	queries := NewQueries(db)
	series, err := queries.GetSeriesById(ctx, seriesId)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Series", seriesId).Msg("Get Series By Id failed")
		return dbo.Series{}, err
	}
	log.Logger.Debug().Int64("Series", seriesId).Msg("End Get Series By Id")
	return series, nil
}

func QueryAllSeries(db *sql.DB, ctx context.Context, query string, from int64, qty int64) ([]dbo.ListSeries, error) {
	log.Logger.Debug().Str("Query", query).Int64("From", from).Int64("Qty", qty).Msg("Start Query All Series")
	queries := NewQueries(db)
	series, err := queries.QueryAllSeries(ctx, query, from, qty)
	if err != nil {
		log.Logger.Error().Str("Query", query).Int64("From", from).Int64("Qty", qty).Err(err).Msg("Query All Series Failed")
		return make([]dbo.ListSeries, 0), err
	}
	log.Logger.Debug().Str("Query", query).Int64("From", from).Int64("Qty", qty).Msg("End Query All Series")
	return series, nil
}

func GetSeriesQty(db *sql.DB, query string, ctx context.Context) (int64, error) {
	log.Logger.Debug().Str("Query", query).Msg("Start Get Series Qty")
	queries := NewQueries(db)
	qty, err := queries.GetSeriesQty(ctx, query)
	if err != nil {
		log.Logger.Error().Err(err).Str("Query", query).Msg("Get Series Qty Failed")
		return 0, err
	}
	log.Logger.Debug().Str("Query", query).Msg("End Get Series Qty")
	return qty, nil
}
