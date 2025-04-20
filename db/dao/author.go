package dao

import (
	"context"
	"database/sql"
	"errors"
	"ignis/library/server/db/dbo"
	"strings"

	"github.com/rs/zerolog/log"
)

// delete not used authors
const deleteOrphanAuthors = `delete a from library.author a left join library.bookauthors as ba on ba.authorId = a.authorId where ba.bookId is null`

// Author query func
const getAuthor = `SELECT authorId, name, url FROM author WHERE name = ?`

func (q *Queries) GetAuthor(ctx context.Context, author string) (dbo.Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthor, author)
	var i dbo.Author
	err := row.Scan(&i.Authorid, &i.Name, &i.Url)
	return i, err
}

const getAuthorUrl = `SELECT authorId, name, url FROM author WHERE name = ? AND (url IS null OR url= ?)`

func (q *Queries) GetAuthorURL(ctx context.Context, author string, url string) (dbo.Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthorUrl, author, url)
	var i dbo.Author
	err := row.Scan(&i.Authorid, &i.Name, &i.Url)
	return i, err
}

const getAuthorById = `SELECT authorId, name, url FROM author WHERE authorId = ?`

func (q *Queries) GetAuthorById(ctx context.Context, authorId int64) (dbo.Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthorById, authorId)
	var i dbo.Author
	err := row.Scan(&i.Authorid, &i.Name, &i.Url)
	return i, err
}

const createAuthor = `INSERT INTO author (name, url) VALUES (?, ?)`

func (q *Queries) CreateAuthor(ctx context.Context, arg dbo.Author) error {
	_, err := q.db.ExecContext(ctx, createAuthor, arg.Name, arg.Url)
	return err
}

const updateAuthor = `UPDATE author SET name=?, url=? WHERE authorId =?`

func (q *Queries) UpdateAuthor(ctx context.Context, author dbo.Author) error {
	_, err := q.db.ExecContext(ctx, updateAuthor, author.Name, author.Url, author.Authorid)
	return err
}

// bookAuthor query func
const bindBookAuthor = `INSERT INTO bookauthors (bookId,authorId) VALUES (?,?)`

func (q *Queries) BindBookAuthor(ctx context.Context, bookId int64, authorId int64) error {
	_, err := q.db.ExecContext(ctx, bindBookAuthor, bookId, authorId)
	return err
}

const rebindBookAuthor = `UPDATE bookAuthors SET authorId=? where authorId=?`

func (q *Queries) RebindBookAuthor(ctx context.Context, oldAuthor int64, newAuthor int64) error {
	_, err := q.db.ExecContext(ctx, rebindBookAuthor, newAuthor, oldAuthor)
	return err
}

const queryAuthorsByBook = `SELECT a.authorid, a.name, a.url FROM author AS a
JOIN bookauthors AS ba ON a.authorId = ba.authorId
WHERE ba.bookId = ?
ORDER BY a.name`

func (q *Queries) QueryAuthorsByBook(ctx context.Context, bookid int64) ([]dbo.Author, error) {
	rows, err := q.db.QueryContext(ctx, queryAuthorsByBook, bookid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dbo.Author
	for rows.Next() {
		var i dbo.Author
		if err := rows.Scan(&i.Authorid, &i.Name, &i.Url); err != nil {
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

const queryAllAuthorBegin = `SELECT a.authorid, a.name, a.url, ifnull(ba.Cnt, 0) as books FROM author AS a
LEFT JOIN 
	(SELECT authorId, count(1) Cnt FROM bookAuthors 
	GROUP BY authorId) As ba
	ON ba.authorId = a.authorId
`
const queryAllAuthorEnd = `ORDER BY a.name
LIMIT ?,?`

const queryAllAuthorWhere = `WHERE a.name like concat('%',?,'%')
`
const queryAllAuthor = queryAllAuthorBegin + queryAllAuthorEnd
const queryAllAuthorFilter = queryAllAuthorBegin + queryAllAuthorWhere + queryAllAuthorEnd

func (q *Queries) QueryAllAuthor(ctx context.Context, name string, from int64, qty int64) ([]dbo.ListAuthor, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if name == "" {
		rows, err = q.db.QueryContext(ctx, queryAllAuthor, from, qty)
	} else {
		rows, err = q.db.QueryContext(ctx, queryAllAuthorFilter, name, from, qty)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dbo.ListAuthor
	for rows.Next() {
		var i dbo.ListAuthor
		if err := rows.Scan(&i.Authorid, &i.Name, &i.Url, &i.BookQty); err != nil {
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

const divideAllAuthorFromBook = "DELETE FROM bookauthors WHERE bookid=?"
const divideAuthosFromBookStart = "DELETE FROM bookauthors WHERE bookid=? and authorId not in (?"
const divideAuthosFromBookEnd = ")"

func (q *Queries) DivideAuthorsFromBook(ctx context.Context, bookId int64, authorids []int64) error {
	if len(authorids) == 0 {
		_, err := q.db.ExecContext(ctx, divideAllAuthorFromBook, bookId)
		return err
	}
	args := make([]interface{}, len(authorids)+1)
	args[0] = bookId
	for i, id := range authorids {
		args[i+1] = id
	}
	query := divideAuthosFromBookStart + strings.Repeat(",?", len(authorids)-1) + divideAuthosFromBookEnd
	_, err := q.db.ExecContext(ctx, query, args...)
	return err
}

const getAuthorQty = `SELECT count(*) FROM author AS a `
const getAuthorQtyFilter = getAuthorQty + queryAllAuthorWhere

func (q *Queries) GetAuthorQty(ctx context.Context, name string) (int64, error) {
	var row *sql.Row
	if name == "" {
		row = q.db.QueryRowContext(ctx, getAuthorQty)
	} else {
		row = q.db.QueryRowContext(ctx, getAuthorQtyFilter, name)
	}
	var count int64
	err := row.Scan(&count)
	return count, err
}

// <======================[ Public functions ]==============================>

func GetAuthorByName(db *sql.DB, ctx context.Context, authorName string) (dbo.Author, error) {
	log.Logger.Debug().Str("Author", authorName).Msg("Start Get Author")
	queries := NewQueries(db)
	author, err := queries.GetAuthor(ctx, authorName)
	if err == sql.ErrNoRows {
		log.Logger.Debug().Str("Author", authorName).Msg("Get Author Not found")
		return dbo.Author{}, err
	} else if err != nil {
		log.Logger.Error().Err(err).Str("Author", authorName).Msg("Get Author failed")
		return dbo.Author{}, err
	}
	log.Logger.Debug().Str("Author", authorName).Int("Author", int(author.Authorid.Int64)).Msg("End Get Author")
	return author, err
}

func AddAuthorNameToBookTBLK(q *Queries, ctx context.Context, bookId int64, authorName string, url string) error {
	log.Logger.Debug().Int64("BookId", bookId).Str("Author", authorName).Msg("Add Author")
	var author dbo.Author
	var err error
	if url == "" {
		author, err = q.GetAuthor(ctx, authorName)
	} else {
		author, err = q.GetAuthorURL(ctx, authorName, url)
	}
	if errors.Is(err, sql.ErrNoRows) {
		author = dbo.Author{Name: authorName}
		err = q.CreateAuthor(ctx, author)
		if err != nil {
			log.Logger.Error().Err(err).Str("Author", authorName).Msg("Create Author failed")
			return err
		}
		authorId, err := q.GetLastId(ctx)
		if err != nil {
			log.Logger.Error().Err(err).Str("Author", authorName).Msg("Create Author failed")
			return err
		}
		log.Logger.Trace().Str("Author", authorName).Int64("CreatedId", authorId).Msg("Create Author failed")
		author.Authorid = sql.NullInt64{Int64: authorId, Valid: true}
	} else if err != nil {
		log.Logger.Error().Err(err).Str("Author", authorName).Msg("Get Author failed")
		return err
	}
	err = q.BindBookAuthor(ctx, bookId, author.Authorid.Int64)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Book", bookId).Str("Author", authorName).Int64("AuthorId", author.Authorid.Int64).Msg("Binding Author to Book failed")
		return err
	}
	log.Logger.Debug().Int64("Book", bookId).Str("Author", author.Name).Int64("AuthorId", author.Authorid.Int64).Msg("End Add Author")
	return nil
}
func AddAuthorToBookTBLK(q *Queries, ctx context.Context, bookId int64, author dbo.Author) error {
	log.Logger.Debug().Int64("Book", bookId).Str("Author", author.Name).Int64("AuthorId", author.Authorid.Int64).Bool("authorId is Valid", author.Authorid.Valid).Msg("Start Add Author")
	var err error
	if !author.Authorid.Valid {
		authorName := author.Name
		authorUrl := author.Url
		if author.Url.Valid && author.Url.String != "" {
			author, err = q.GetAuthorURL(ctx, authorName, authorUrl.String)
		} else {
			author, err = q.GetAuthor(ctx, authorName)
		}
		if err == sql.ErrNoRows {
			log.Logger.Trace().Str("Author", authorName).Msg("Author not found, create")
			author = dbo.Author{Name: authorName, Url: authorUrl}
			err = q.CreateAuthor(ctx, author)
			if err != nil {
				log.Logger.Error().Err(err).Str("Author", authorName).Msg("Create Author failed")
				return err
			}
			authorId, err := q.GetLastId(ctx)
			if err != nil {
				log.Logger.Error().Err(err).Str("Author", authorName).Msg("Create Author failed")
				return err
			}
			author.Authorid = sql.NullInt64{Int64: authorId, Valid: true}
		} else if err != nil {
			log.Logger.Error().Err(err).Str("Author", authorName).Msg("Get Author failed")
			return err
		}
		if authorUrl.Valid && !author.Url.Valid {
			author.Url = authorUrl
			err = q.UpdateAuthor(ctx, author)
			if err != nil {
				log.Logger.Error().Err(err).Str("Author", authorName).Msg("Get Author failed")
				return err
			}
		}
	}
	err = q.BindBookAuthor(ctx, bookId, author.Authorid.Int64)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Book", bookId).Str("Author", author.Name).Int64("AuthorId", author.Authorid.Int64).Msg("Binding Author to Book failed")
		return err
	}
	log.Logger.Debug().Int64("Book", bookId).Str("Author", author.Name).Int64("AuthorId", author.Authorid.Int64).Msg("End Add Author")
	return nil
}

func QueryAuthorsByBookTBLK(q *Queries, ctx context.Context, bookId int64) ([]dbo.Author, error) {
	log.Logger.Debug().Int64("Book", bookId).Msg("Start Query Authors by Book")
	authors, err := q.QueryAuthorsByBook(ctx, bookId)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Book", bookId).Msg("Query Authors by Book failed")
		return nil, err
	}
	log.Logger.Debug().Int64("Book", bookId).Msg("End Query Authors by Book ")
	return authors, nil
}

func UpdateAuthorTBLK(q *Queries, ctx context.Context, author dbo.Author) error {
	log.Logger.Debug().Int64("Author", author.Authorid.Int64).Msg("Start UpdateAuthorTBLK")
	err := q.UpdateAuthor(ctx, author)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Author", author.Authorid.Int64).Msg("UpdateAuthorTBLK failed")
		return err
	}
	log.Logger.Debug().Int64("Author", author.Authorid.Int64).Msg("End UpdateAuthorTBLK ")
	return nil
}

func GetAuthorsById(db *sql.DB, ctx context.Context, authorId int64) (dbo.Author, error) {
	log.Logger.Debug().Int64("Author", authorId).Msg("Start Get Author By Id")
	queries := NewQueries(db)
	author, err := queries.GetAuthorById(ctx, authorId)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Author", authorId).Msg("Get Author By Id failed")
		return dbo.Author{}, err
	}
	log.Logger.Debug().Int64("Author", authorId).Msg("End Get Author By Id")
	return author, nil
}
func DivideBookAllOtherAuthor(q *Queries, ctx context.Context, bookId int64, authorsId []int64) error {
	log.Logger.Debug().Int64("Book", bookId).Ints64("AuthorIds", authorsId).Msg("Start Divide Book All Other Author")
	err := q.DivideAuthorsFromBook(ctx, bookId, authorsId)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Book", bookId).Ints64("AuthorIds", authorsId).Msg("Divide Book All Other Author Failed")
		return err
	}

	log.Logger.Debug().Int64("Book", bookId).Ints64("AuthorIds", authorsId).Msg("End Divide Book All Other Author")
	return nil
}

func QueryAllAuthor(db *sql.DB, ctx context.Context, query string, from int64, qty int64) ([]dbo.ListAuthor, error) {
	log.Logger.Debug().Str("Query", query).Int64("From", from).Int64("Qty", qty).Msg("Start Query All Author")
	queries := NewQueries(db)
	authors, err := queries.QueryAllAuthor(ctx, query, from, qty)
	if err != nil {
		log.Logger.Error().Str("Query", query).Int64("From", from).Int64("Qty", qty).Err(err).Msg("Query All Author Failed")
		return make([]dbo.ListAuthor, 0), err
	}
	log.Logger.Debug().Str("Query", query).Int64("From", from).Int64("Qty", qty).Msg("End Query All Authors")
	return authors, nil
}

func GetAuthorQty(db *sql.DB, ctx context.Context, query string) (int64, error) {
	log.Logger.Debug().Str("Query", query).Msg("Start Get Author Qty")
	queries := NewQueries(db)
	qty, err := queries.GetAuthorQty(ctx, query)
	if err != nil {
		log.Logger.Error().Err(err).Str("Query", query).Msg("Get Author Qty Failed")
		return 0, err
	}
	log.Logger.Debug().Str("Query", query).Msg("End Get Author Qty")
	return qty, nil
}
