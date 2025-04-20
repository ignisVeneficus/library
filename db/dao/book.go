package dao

import (
	"context"
	"database/sql"
	"ignis/library/server/db/dbo"
	"strings"

	"github.com/rs/zerolog/log"
)

/* duplikacio kereses

SELECT b.title,
   b.file
FROM library.book as b
   INNER JOIN (SELECT title, url
               FROM library.book
               GROUP BY title, url
               HAVING COUNT(bookid) > 1) dup
           ON b.title = dup.title and b.url = dup.url;

*/

const (
	ALL_BOOK_MODE_ALL = iota
	ALL_BOOK_MODE_QUERY
	ALL_BOOK_MODE_NEW
	ALL_BOOK_MODE_FILE
)

const deleteBook = `DELETE FROM book WHERE bookId = ?`
const deleteBookAuthors = `DELETE FROM bookauthors WHERE bookId = ?`
const deleteBookSeries = `DELETE FROM bookseries WHERE bookId = ?`
const deleteBookTags = `DELETE FROM booktags WHERE bookId = ?`

func (q *Queries) DeleteBook(ctx context.Context, bookId int64) error {
	_, err := q.db.ExecContext(ctx, deleteBook, bookId)
	return err
}
func (q *Queries) DeleteBookAuthors(ctx context.Context, bookId int64) error {
	_, err := q.db.ExecContext(ctx, deleteBookAuthors, bookId)
	return err
}
func (q *Queries) DeleteBookSeries(ctx context.Context, bookId int64) error {
	_, err := q.db.ExecContext(ctx, deleteBookSeries, bookId)
	return err
}
func (q *Queries) DeleteBookTags(ctx context.Context, bookId int64) error {
	_, err := q.db.ExecContext(ctx, deleteBookTags, bookId)
	return err
}

const createBook = `INSERT INTO book (title, format, file, hash, updateTs, url, isbn, hasCover, coverColor,coverType, edited, blurb) VALUES
(?,?,?,?,NOW(),?,?,?,?,?,?,?)`

func (q *Queries) CreateBook(ctx context.Context, book dbo.Book) error {
	_, err := q.db.ExecContext(ctx, createBook,
		book.Title,
		book.Format,
		book.File,
		book.Hash,
		book.Url,
		book.Isbn,
		book.Hascover,
		book.CoverColor,
		book.Edited,
		book.CoverType,
		book.Blurb,
	)
	return err
}

const createNewBook = `INSERT INTO book (title, format, file, hash,  updateTs, isbn, hascover, coverColor,coverType, edited,blurb) VALUES
 (?,?,?,?,NOW(),?,?,?,?,0,?)`

func (q *Queries) CreateNewBook(ctx context.Context, book dbo.Book) error {
	_, err := q.db.ExecContext(ctx, createNewBook,
		book.Title,
		book.Format,
		book.File,
		book.Hash,
		book.Isbn,
		book.Hascover,
		book.CoverColor,
		book.CoverType,
		book.Blurb,
	)
	return err
}

const updateNewBook = `UPDATE book SET format=?, file=?, updateTs = NOW(), isbn=?, hascover=?, coverColor=?,coverType=? WHERE hash=?`

func (q *Queries) UpdateNewFile(ctx context.Context, book dbo.Book) error {
	_, err := q.db.ExecContext(ctx, updateNewBook,
		book.Format,
		book.File,
		book.Isbn,
		book.Hascover,
		book.CoverColor,
		book.CoverType,
		book.Hash,
	)
	return err
}

const moveBook = `UPDATE book SET file = ?, updateTs = NOW() WHERE bookId = ?`

func (q *Queries) MoveBook(ctx context.Context, bookId int64, newFileName string) error {
	_, err := q.db.ExecContext(ctx, moveBook, newFileName, bookId)
	return err
}

const updateBookMetadata = `UPDATE book SET title=?, url=?, blurb =?, updateTs = NOW(), edited = 1 WHERE bookId=?`

func (q *Queries) UpdateBookMetadata(ctx context.Context, book dbo.Book) error {
	_, err := q.db.ExecContext(ctx, updateBookMetadata,
		book.Title,
		book.Url,
		book.Blurb,
		book.Bookid,
	)
	return err
}

const selectBook = `b.bookId, b.title, b.format, b.file, b.hash, b.updatets, b.url, b.isbn, b.covercolor, b.edited, b.hascover, b.covertype, b.blurb`

const getBookByHash = `SELECT ` + selectBook + ` FROM book AS b WHERE hash = ?`

func (q *Queries) GetBookByHash(ctx context.Context, hash string) (dbo.Book, error) {
	row := q.db.QueryRowContext(ctx, getBookByHash, hash)
	var i dbo.Book
	err := row.Scan(
		&i.Bookid,
		&i.Title,
		&i.Format,
		&i.File,
		&i.Hash,
		&i.Updatets,
		&i.Url,
		&i.Isbn,
		&i.CoverColor,
		&i.Edited,
		&i.Hascover,
		&i.CoverType,
		&i.Blurb,
	)
	return i, err
}

const getBookById = `SELECT ` + selectBook + ` FROM book AS b WHERE bookId = ?`

func (q *Queries) GetBookById(ctx context.Context, id int64) (dbo.Book, error) {
	row := q.db.QueryRowContext(ctx, getBookById, id)
	var i dbo.Book
	err := row.Scan(
		&i.Bookid,
		&i.Title,
		&i.Format,
		&i.File,
		&i.Hash,
		&i.Updatets,
		&i.Url,
		&i.Isbn,
		&i.CoverColor,
		&i.Edited,
		&i.Hascover,
		&i.CoverType,
		&i.Blurb,
	)
	return i, err
}

const queryBookFilter = `SELECT ` + selectBook + `, COALESCE(GROUP_CONCAT(a.name ORDER BY a.name SEPARATOR ', '), '') AS allauthor 
FROM (
	(SELECT b1.bookId, b1.title, b1.format, b1.file, b1.hash, b1.updatets, b1.url, b1.isbn, b1.covercolor, b1.edited, b1.hascover, b1.covertype, b1.blurb
		FROM author AS a1
		JOIN bookauthors AS ba1 ON a1.authorID = ba1.authorID
		JOIN book AS b1 ON ba1.bookID = b1.bookID
		WHERE a1.name like concat('%',?,'%')
	)
	UNION
	(SELECT b2.bookId, b2.title, b2.format, b2.file, b2.hash, b2.updatets, b2.url, b2.isbn, b2.covercolor, b2.edited, b2.hascover, b2.covertype, b2.blurb
		FROM series AS s2
		JOIN bookseries AS bs2 ON s2.seriesID = bs2.seriesID
		JOIN book AS b2 ON bs2.bookID = b2.bookID
		WHERE s2.title like concat('%',?,'%')
	)
	UNION
	(SELECT b3.bookId, b3.title, b3.format, b3.file, b3.hash, b3.updatets, b3.url, b3.isbn, b3.covercolor, b3.edited, b3.hascover, b3.covertype, b3.blurb
		FROM tag AS t3
		JOIN booktags AS bt3 ON t3.tagID = bt3.tagID
		JOIN book AS b3 ON bt3.bookID = b3.bookID
		WHERE t3.name like concat('%',?,'%')
	)
	UNION
	(SELECT b4.bookId, b4.title, b4.format, b4.file, b4.hash, b4.updatets, b4.url, b4.isbn, b4.covercolor, b4.edited, b4.hascover, b4.covertype, b4.blurb
		FROM book AS b4
		WHERE b4.title like concat('%',?,'%')
	)
) as b
LEFT OUTER JOIN bookauthors AS ba ON b.bookId = ba.bookId
LEFT OUTER JOIN author AS a ON a.authorId = ba.authorId
GROUP BY b.bookId
ORDER BY allauthor, b.title
LIMIT ?,?`

const getBookQtyFilter = `SELECT count(*) FROM (
	(SELECT b1.bookId
		FROM author AS a1
		JOIN bookauthors AS ba1 ON a1.authorID = ba1.authorID
		JOIN book AS b1 ON ba1.bookID = b1.bookID
		WHERE a1.name like concat('%',?,'%')
	)
	UNION
	(SELECT b2.bookId
		FROM series AS s2
		JOIN bookseries AS bs2 ON s2.seriesID = bs2.seriesID
		JOIN book AS b2 ON bs2.bookID = b2.bookID
		WHERE s2.title like concat('%',?,'%')
	)
	UNION
	(SELECT b3.bookId
		FROM tag AS t3
		JOIN booktags AS bt3 ON t3.tagID = bt3.tagID
		JOIN book AS b3 ON bt3.bookID = b3.bookID
		WHERE t3.name like concat('%',?,'%')
	)
	UNION
	(SELECT b4.bookId
		FROM book AS b4
		WHERE b4.title like concat('%',?,'%')
	)
) as cnt`

const queryBook = `SELECT ` + selectBook + `, COALESCE(GROUP_CONCAT(a.name ORDER BY a.name SEPARATOR ', '), '') AS allauthor FROM book AS b
LEFT OUTER JOIN bookauthors AS ba ON b.bookId = ba.bookId
LEFT OUTER JOIN author AS a ON a.authorId = ba.authorId
GROUP BY b.bookId
ORDER BY allauthor, b.title
LIMIT ?,?`

const queryBookNew = `SELECT ` + selectBook + `, COALESCE(GROUP_CONCAT(a.name ORDER BY a.name SEPARATOR ', '), '') AS allauthor FROM book AS b
LEFT OUTER JOIN bookauthors AS ba ON b.bookId = ba.bookId
LEFT OUTER JOIN author AS a ON a.authorId = ba.authorId
WHERE b.edited=0
GROUP BY b.bookId
ORDER BY allauthor, b.title
LIMIT ?,?`

const queryBookFile = `SELECT ` + selectBook + `, COALESCE(GROUP_CONCAT(a.name ORDER BY a.name SEPARATOR ', '), '') AS allauthor FROM book AS b
LEFT OUTER JOIN bookauthors AS ba ON b.bookId = ba.bookId
LEFT OUTER JOIN author AS a ON a.authorId = ba.authorId
WHERE b.file LIKE concat('%',?,'%')
GROUP BY b.bookId
ORDER BY allauthor, b.title
LIMIT ?,?`

func (q *Queries) QueryBook(ctx context.Context, query string, from int64, qty int64) ([]dbo.Book, error) {
	var (
		rows *sql.Rows
		err  error
	)
	switch GetBookQueryType(query) {
	case ALL_BOOK_MODE_ALL:
		rows, err = q.db.QueryContext(ctx, queryBook, from, qty)
	case ALL_BOOK_MODE_QUERY:
		rows, err = q.db.QueryContext(ctx, queryBookFilter, query, query, query, query, from, qty)
	case ALL_BOOK_MODE_NEW:
		rows, err = q.db.QueryContext(ctx, queryBookNew, from, qty)
	case ALL_BOOK_MODE_FILE:
		rows, err = q.db.QueryContext(ctx, queryBookFile, query[2:], from, qty)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dbo.Book
	for rows.Next() {
		var i dbo.Book
		var tmp string
		if err := rows.Scan(
			&i.Bookid,
			&i.Title,
			&i.Format,
			&i.File,
			&i.Hash,
			&i.Updatets,
			&i.Url,
			&i.Isbn,
			&i.CoverColor,
			&i.Edited,
			&i.Hascover,
			&i.CoverType,
			&i.Blurb,
			&tmp,
		); err != nil {
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

const queryAllBooks = `SELECT ` + selectBook + ` FROM book as b`

func (q *Queries) QueryAllBooks(ctx context.Context) ([]dbo.Book, error) {
	rows, err := q.db.QueryContext(ctx, queryAllBooks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dbo.Book
	for rows.Next() {
		var i dbo.Book
		if err := rows.Scan(
			&i.Bookid,
			&i.Title,
			&i.Format,
			&i.File,
			&i.Hash,
			&i.Updatets,
			&i.Url,
			&i.Isbn,
			&i.CoverColor,
			&i.Edited,
			&i.Hascover,
			&i.CoverType,
			&i.Blurb,
		); err != nil {
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

const getBookQty = `SELECT count(*) FROM book`

const getBookQtyNew = `SELECT count(*) FROM book
WHERE edited=0`

const getBookQtyFile = `SELECT count(*) FROM book
WHERE file LIKE concat('%',?,'%')`

func (q *Queries) GetBookQty(ctx context.Context, query string) (int64, error) {
	var row *sql.Row
	switch GetBookQueryType(query) {
	case ALL_BOOK_MODE_ALL:
		row = q.db.QueryRowContext(ctx, getBookQty)
	case ALL_BOOK_MODE_QUERY:
		row = q.db.QueryRowContext(ctx, getBookQtyFilter, query, query, query, query)
	case ALL_BOOK_MODE_NEW:
		row = q.db.QueryRowContext(ctx, getBookQtyNew)
	case ALL_BOOK_MODE_FILE:
		row = q.db.QueryRowContext(ctx, getBookQtyFile, query[2:])
	}
	var count int64
	err := row.Scan(&count)
	return count, err
}

const queryBookByAuthorId = `SELECT ` + selectBook + ` FROM bookauthors ba1
JOIN book as b on b.bookId = ba1.bookId
JOIN bookauthors AS ba ON b.bookId = ba.bookId
JOIN author AS a ON a.authorId = ba.authorId
WHERE ba1.authorId = ?
GROUP BY b.bookId
ORDER BY a.name, b.title
LIMIT ?,?`

func (q *Queries) QueryBookByAuthorId(ctx context.Context, authorId int64, from int64, qty int64) ([]dbo.Book, error) {
	rows, err := q.db.QueryContext(ctx, queryBookByAuthorId, authorId, from, qty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dbo.Book
	for rows.Next() {
		var i dbo.Book
		if err := rows.Scan(
			&i.Bookid,
			&i.Title,
			&i.Format,
			&i.File,
			&i.Hash,
			&i.Updatets,
			&i.Url,
			&i.Isbn,
			&i.CoverColor,
			&i.Edited,
			&i.Hascover,
			&i.CoverType,
			&i.Blurb,
		); err != nil {
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

const getBookByAuthorIdQty = `SELECT count(*) FROM (select 1 from bookauthors ba1
JOIN book as b on b.bookId = ba1.bookId
WHERE ba1.authorId = ?
GROUP BY b.bookId) as cnt;`

func (q *Queries) GetBookByAuthorIdQty(ctx context.Context, authorId int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getBookByAuthorIdQty, authorId)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const queryBookBySeriesId = `SELECT ` + selectBook + ` FROM bookseries bs1
JOIN book as b on b.bookId = bs1.bookId
JOIN bookauthors AS ba ON b.bookId = ba.bookId
JOIN author AS a ON a.authorId = ba.authorId
WHERE bs1.seriesId = ?
GROUP BY b.bookId
ORDER BY bs1.sequence, a.name, b.title
LIMIT ?,?`

func (q *Queries) QueryBookBySeriesId(ctx context.Context, authorId int64, from int64, qty int64) ([]dbo.Book, error) {
	rows, err := q.db.QueryContext(ctx, queryBookBySeriesId, authorId, from, qty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dbo.Book
	for rows.Next() {
		var i dbo.Book
		if err := rows.Scan(
			&i.Bookid,
			&i.Title,
			&i.Format,
			&i.File,
			&i.Hash,
			&i.Updatets,
			&i.Url,
			&i.Isbn,
			&i.CoverColor,
			&i.Edited,
			&i.Hascover,
			&i.CoverType,
			&i.Blurb,
		); err != nil {
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

const getBookBySeriesIdQty = `SELECT count(*) FROM (select 1 from bookseries bs1
JOIN book as b on b.bookId = bs1.bookId
WHERE bs1.seriesId = ?
GROUP BY b.bookId) as cnt;`

func (q *Queries) GetBookBySeriesIdQty(ctx context.Context, seriesId int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getBookBySeriesIdQty, seriesId)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const queryAllBookId = `SELECT bookId FROM book ORDER BY bookId`

func (q *Queries) QueryAllBookId(ctx context.Context) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, queryAllBookId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var i int64
		if err := rows.Scan(
			&i,
		); err != nil {
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

const queryBookByTagId = `SELECT ` + selectBook + `, COALESCE(GROUP_CONCAT(a.name ORDER BY a.name SEPARATOR ', '), '') AS allauthor FROM bookTags bt1
JOIN book as b on b.bookId = bt1.bookId
JOIN bookauthors AS ba ON b.bookId = ba.bookId
JOIN author AS a ON a.authorId = ba.authorId
WHERE bt1.tagId = ?
GROUP BY b.bookId
ORDER BY allauthor, b.title
LIMIT ?,?`

func (q *Queries) QueryBookByTagId(ctx context.Context, tagId int64, from int64, qty int64) ([]dbo.Book, error) {
	rows, err := q.db.QueryContext(ctx, queryBookByTagId, tagId, from, qty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dbo.Book
	for rows.Next() {
		var i dbo.Book
		var tmp string
		if err := rows.Scan(
			&i.Bookid,
			&i.Title,
			&i.Format,
			&i.File,
			&i.Hash,
			&i.Updatets,
			&i.Url,
			&i.Isbn,
			&i.CoverColor,
			&i.Edited,
			&i.Hascover,
			&i.CoverType,
			&i.Blurb,
			&tmp,
		); err != nil {
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

const getBookByTagIdQty = `SELECT count(*) FROM (SELECT 1 FROM bookTags bt1
JOIN book as b on b.bookId = bt1.bookId
JOIN bookauthors AS ba ON b.bookId = ba.bookId
JOIN author AS a ON a.authorId = ba.authorId
WHERE bt1.tagId = ?
GROUP BY b.bookId) as cnt;`

func (q *Queries) GetBookByTagIdQty(ctx context.Context, tagId int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getBookByTagIdQty, tagId)
	var count int64
	err := row.Scan(&count)
	return count, err
}

// <======================[ Public functions ]==============================>

func GetBookQueryType(query string) int {
	if query == "" {
		return ALL_BOOK_MODE_ALL
	}
	if query == ":n" || query == "n:" {
		return ALL_BOOK_MODE_NEW
	}
	if strings.HasPrefix(query, "f:") {
		return ALL_BOOK_MODE_FILE
	}
	return ALL_BOOK_MODE_QUERY
}

func readAllMetadata(q *Queries, ctx context.Context, book dbo.Book) dbo.Book {
	authors, err := QueryAuthorsByBookTBLK(q, ctx, book.Bookid.Int64)
	if err == nil {
		book.Authors = authors
	} else {
		log.Logger.Error().Err(err).Msg("Fill author")
	}
	series, err := QuerySeriesByBookTBLK(q, ctx, book.Bookid.Int64)
	if err == nil {
		book.Series = series
	} else {
		log.Logger.Error().Err(err).Msg("Fill series")
	}
	tags, err := QueryTagsByBookTBLK(q, ctx, book.Bookid.Int64)
	if err == nil {
		book.Tags = tags
	} else {
		log.Logger.Error().Err(err).Msg("Fill tags")
	}
	return book
}

func insertBook(db *sql.DB, ctx context.Context, book dbo.Book, full bool) (int, error) {
	log.Logger.Debug().Str("Book", book.Title.String).Msg("Start Insert Book")
	tx, err := GetTx(db, ctx)
	if err != nil {
		log.Logger.Error().Str("Book", book.Title.String).Err(err).Msg("Insert Book failed")
		return 0, err
	}
	defer tx.Rollback()
	queries := NewQueries(tx)
	if full {
		err = queries.CreateBook(ctx, book)
	} else {
		err = queries.CreateNewBook(ctx, book)
	}
	if err != nil {
		log.Logger.Error().Str("Book", book.Title.String).Err(err).Msg("Insert Book failed")
		return 0, err
	}
	lastRow, err := queries.GetLastId(ctx)
	if err != nil {
		log.Logger.Error().Str("Book", book.Title.String).Err(err).Msg("Insert Book failed")
		return 0, err
	}
	book.Bookid = sql.NullInt64{Int64: lastRow, Valid: true}
	for _, author := range book.Authors {
		err = AddAuthorToBookTBLK(queries, ctx, book.Bookid.Int64, author)
		if err != nil {
			log.Logger.Error().Str("Book", book.Title.String).Str("Author", author.Name).Err(err).Msg("Insert Author failed")
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Logger.Error().Str("Book", book.Title.String).Err(err).Msg("Insert Book failed")
		return 0, err
	}
	log.Logger.Debug().Str("Book", book.Title.String).Int("bookId", int(lastRow)).Msg("Book Inserted")
	return int(lastRow), nil

}

func CreateBook(db *sql.DB, ctx context.Context, book dbo.Book) (int, error) {
	lastId, err := insertBook(db, ctx, book, true)
	return lastId, err
}

func CreateNewBook(db *sql.DB, ctx context.Context, book dbo.Book) (int, error) {
	lastId, err := insertBook(db, ctx, book, false)
	return lastId, err
}
func GetBookByHash(db *sql.DB, ctx context.Context, hash string, needOther bool) (dbo.Book, error) {
	log.Logger.Debug().Str("Book Hash", hash).Msg("Start Get Book by Hash")
	queries := NewQueries(db)
	book, err := queries.GetBookByHash(ctx, hash)
	if err == sql.ErrNoRows {
		return dbo.Book{}, GetDataNotFoundError("Book")
	} else if err != nil {
		log.Logger.Error().Str("Book Hash", hash).Err(err).Msg("Get Book by Hash Failed")
		return dbo.Book{}, err
	}
	if !needOther {
		log.Logger.Debug().Str("Book Hash", hash).Msg("End Get Book by Hash, No other metadata Readed")
		return book, nil
	}
	book = readAllMetadata(queries, ctx, book)
	log.Logger.Debug().Str("Book Hash", hash).Msg("End Get Book by Hash, All metadata Readed")
	return book, nil
}
func MoveBook(db *sql.DB, ctx context.Context, bookId int64, newPath string) error {
	log.Logger.Debug().Int64("Book Id", bookId).Str("New Path", newPath).Msg("Start Move Book")
	queries := NewQueries(db)
	err := queries.MoveBook(ctx, bookId, newPath)
	if err != nil {
		log.Logger.Error().Int64("Book Id", bookId).Str("New Path", newPath).Err(err).Msg("Move Book Failed")
		return err
	}
	log.Logger.Debug().Int64("Book Id", bookId).Str("New Path", newPath).Msg("End Move Book")
	return nil
}

func QueryBook(db *sql.DB, ctx context.Context, query string, from int64, qty int64) ([]dbo.Book, error) {
	log.Logger.Debug().Int64("From", from).Int64("Qty", qty).Msg("Start Query Book")
	queries := NewQueries(db)
	books, err := queries.QueryBook(ctx, query, from, qty)
	if err != nil {
		log.Logger.Error().Int64("From", from).Int64("Qty", qty).Err(err).Msg("Query Book Failed")
		return make([]dbo.Book, 0), err
	}
	log.Logger.Debug().Int64("From", from).Int64("Qty", qty).Int("Got books", len(books)).Msg("Start Query Book")
	for i := range books {
		books[i] = readAllMetadata(queries, ctx, books[i])
	}
	log.Logger.Debug().Int64("From", from).Int64("Qty", qty).Msg("End Query Book")
	return books, nil
}
func GetBookQty(db *sql.DB, ctx context.Context, query string) (int64, error) {
	log.Logger.Debug().Msg("Start Get Book Qty")
	queries := NewQueries(db)
	qty, err := queries.GetBookQty(ctx, query)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Get Book Qty Failed")
		return 0, err
	}
	log.Logger.Debug().Msg("End Get Book Qty")
	return qty, nil
}
func GetBookById(db *sql.DB, ctx context.Context, id int64) (dbo.Book, error) {
	log.Logger.Debug().Int64("Book Id", id).Msg("Start Get Book by Id")
	queries := NewQueries(db)
	book, err := queries.GetBookById(ctx, id)
	if err == sql.ErrNoRows {
		return dbo.Book{}, GetDataNotFoundError("Book")
	} else if err != nil {
		log.Logger.Error().Int64("Book Id", id).Err(err).Msg("Get Book by Id Failed")
		return dbo.Book{}, err
	}
	book = readAllMetadata(queries, ctx, book)
	log.Logger.Debug().Int64("Book Id", id).Msg("End Get Book by Id")
	return book, nil

}
func QueryAllBookByAuthorId(db *sql.DB, ctx context.Context, authorId int64, from int64, qty int64) ([]dbo.Book, error) {
	log.Logger.Debug().Int64("Author", authorId).Int64("From", from).Int64("Qty", qty).Msg("Start Query Book All By AuthorId")
	queries := NewQueries(db)
	books, err := queries.QueryBookByAuthorId(ctx, authorId, from, qty)
	if err != nil {
		log.Logger.Error().Int64("Author", authorId).Int64("From", from).Int64("Qty", qty).Err(err).Msg("Query Book All By AuthorId Failed")
		return make([]dbo.Book, 0), err
	}
	for i := range books {
		books[i] = readAllMetadata(queries, ctx, books[i])
	}
	log.Logger.Debug().Int64("Author", authorId).Int64("From", from).Int64("Qty", qty).Msg("End Query Book All By AuthorId")
	return books, nil
}
func GetBookByAuthorIdQty(db *sql.DB, ctx context.Context, authorId int64) (int64, error) {
	log.Logger.Debug().Msg("Start Get Book BY AuthorId Qty")
	queries := NewQueries(db)
	qty, err := queries.GetBookByAuthorIdQty(ctx, authorId)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Get Book BY AuthorId Qty Failed")
		return 0, err
	}
	log.Logger.Debug().Msg("End Get Book BY AuthorId Qty")
	return qty, nil
}

func QueryAllBookBySeriesId(db *sql.DB, ctx context.Context, seriesId int64, from int64, qty int64) ([]dbo.Book, error) {
	log.Logger.Debug().Int64("Series", seriesId).Int64("From", from).Int64("Qty", qty).Msg("Start Query Book All By SeriesId")
	queries := NewQueries(db)
	books, err := queries.QueryBookBySeriesId(ctx, seriesId, from, qty)
	if err != nil {
		log.Logger.Error().Int64("Series", seriesId).Int64("From", from).Int64("Qty", qty).Err(err).Msg("Query Book All By SeriesId Failed")
		return make([]dbo.Book, 0), err
	}
	for i := range books {
		books[i] = readAllMetadata(queries, ctx, books[i])
	}
	log.Logger.Debug().Int64("Series", seriesId).Int64("From", from).Int64("Qty", qty).Msg("End Query Book All By SeriesId")
	return books, nil
}
func GetBookBySeriesIdQty(db *sql.DB, ctx context.Context, seriesId int64) (int64, error) {
	log.Logger.Debug().Msg("Start Get Book BY SeriesId Qty")
	queries := NewQueries(db)
	qty, err := queries.GetBookBySeriesIdQty(ctx, seriesId)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Get Book BY SeriesId Qty Failed")
		return 0, err
	}
	log.Logger.Debug().Msg("End Get Book BY SeriesId Qty")
	return qty, nil
}

func QueryAllBookByTagId(db *sql.DB, ctx context.Context, tagId int64, from int64, qty int64) ([]dbo.Book, error) {
	log.Logger.Debug().Int64("Tag", tagId).Int64("From", from).Int64("Qty", qty).Msg("Start Query Book All By TagId")
	queries := NewQueries(db)
	books, err := queries.QueryBookByTagId(ctx, tagId, from, qty)
	if err != nil {
		log.Logger.Error().Int64("Tag", tagId).Int64("From", from).Int64("Qty", qty).Err(err).Msg("Query Book All By TagId Failed")
		return make([]dbo.Book, 0), err
	}
	for i := range books {
		books[i] = readAllMetadata(queries, ctx, books[i])
	}
	log.Logger.Debug().Int64("Tag", tagId).Int64("From", from).Int64("Qty", qty).Msg("End Query Book All By TagId")
	return books, nil
}
func GetBookByTagIdQty(db *sql.DB, ctx context.Context, tagId int64) (int64, error) {
	log.Logger.Debug().Msg("Start Get Book By TagId Qty")
	queries := NewQueries(db)
	qty, err := queries.GetBookByTagIdQty(ctx, tagId)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Get Book By TagId Qty Failed")
		return 0, err
	}
	log.Logger.Debug().Msg("End Get Book By TagId Qty")
	return qty, nil
}

func UpdateBook(db *sql.DB, ctx context.Context, book dbo.Book) error {
	log.Logger.Debug().Int64("Book", book.Bookid.Int64).Msg("Start Update Book")
	tx, err := GetTx(db, ctx)
	if err != nil {
		log.Logger.Error().Int64("Book", book.Bookid.Int64).Err(err).Msg("Update Book failed")
		return err
	}
	defer tx.Rollback()

	queries := NewQueries(tx)
	err = queries.UpdateBookMetadata(ctx, book)
	if err != nil {
		log.Logger.Error().Int64("Book", book.Bookid.Int64).Err(err).Msg("Update Book failed")
		return err
	}
	// authors
	authorsId := make([]int64, 0)
	for _, author := range book.Authors {
		if author.Authorid.Valid {
			authorsId = append(authorsId, author.Authorid.Int64)
			err = UpdateAuthorTBLK(queries, ctx, author)
			if err != nil {
				log.Logger.Error().Int64("Book", book.Bookid.Int64).Str("Author", author.Name).Int64("AuthorId", author.Authorid.Int64).Err(err).Msg("Update Author failed")
				return err
			}
		}
	}
	err = DivideBookAllOtherAuthor(queries, ctx, book.Bookid.Int64, authorsId)
	if err != nil {
		log.Logger.Error().Int64("Book", book.Bookid.Int64).Err(err).Msg("Delete Author failed")
		return err
	}
	for _, author := range book.Authors {
		if !author.Authorid.Valid {
			err = AddAuthorToBookTBLK(queries, ctx, book.Bookid.Int64, author)
			if err != nil {
				log.Logger.Error().Int64("Book", book.Bookid.Int64).Err(err).Str("Author", author.Name).Msg("Inser Author failed")
				return err
			}
		}
	}
	// series
	seriesId := make([]int64, 0)
	for _, series := range book.Series {
		if series.SeriesId.Valid {
			seriesId = append(seriesId, series.SeriesId.Int64)
			err = UpdateBookSeriesTBLK(queries, ctx, book.Bookid.Int64, series)
			if err != nil {
				log.Logger.Error().Int64("Book", book.Bookid.Int64).Str("Series", series.Title).Int64("SeriesId", series.SeriesId.Int64).Err(err).Msg("Update Series failed")
				return err
			}
		}
	}

	err = DivideBookAllOtherSeries(queries, ctx, book.Bookid.Int64, seriesId)
	if err != nil {
		log.Logger.Error().Int64("Book", book.Bookid.Int64).Err(err).Msg("Delete Series failed")
		return err
	}
	for _, series := range book.Series {
		if !series.SeriesId.Valid {
			err = AddSeriesToBookTBLK(queries, ctx, book.Bookid.Int64, series.Series, series.Seqno)
			if err != nil {
				log.Logger.Error().Int64("Book", book.Bookid.Int64).Err(err).Str("Series", series.Title).Msg("Inser Series failed")
				return err
			}
		}
	}

	// tags
	tagsId := make([]int64, 0)
	for _, tag := range book.Tags {
		if tag.TagId.Valid {
			tagsId = append(tagsId, tag.TagId.Int64)
			err = UpdateTagTBLK(queries, ctx, book.Bookid.Int64, tag)
			if err != nil {
				log.Logger.Error().Int64("Book", book.Bookid.Int64).Str("Tag", tag.Name).Int64("TagId", tag.TagId.Int64).Err(err).Msg("Update Tag failed")
				return err
			}
		}
	}

	err = DivideBookAllOtherTags(queries, ctx, book.Bookid.Int64, tagsId)
	if err != nil {
		log.Logger.Error().Int64("Book", book.Bookid.Int64).Err(err).Msg("Delete Tags failed")
		return err
	}
	for _, tag := range book.Tags {
		if !tag.TagId.Valid {
			err = AddTagToBookTBLK(queries, ctx, book.Bookid.Int64, tag)
			if err != nil {
				log.Logger.Error().Int64("Book", book.Bookid.Int64).Err(err).Str("Tag", tag.Name).Msg("Inser Tag failed")
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Logger.Error().Int64("Book", book.Bookid.Int64).Err(err).Msg("Update Book failed")
		return err
	}
	log.Logger.Debug().Int64("Book", book.Bookid.Int64).Msg("Book updated")
	return nil

}
func DeleteBook(db *sql.DB, ctx context.Context, bookId int64) error {
	log.Logger.Debug().Int64("Book", bookId).Msg("Start Delete Book")

	tx, err := GetTx(db, ctx)
	if err != nil {
		log.Logger.Error().Int64("Book", bookId).Err(err).Msg("Delete Book failed")
		return err
	}
	defer tx.Rollback()

	queries := NewQueries(tx)
	err = queries.DeleteBook(ctx, bookId)
	if err != nil {
		log.Logger.Error().Int64("Book", bookId).Err(err).Msg("Delete Book failed")
		return err
	}
	err = queries.DeleteBookAuthors(ctx, bookId)
	if err != nil {
		log.Logger.Error().Int64("Book", bookId).Err(err).Msg("Delete Book failed")
		return err
	}
	err = queries.DeleteBookSeries(ctx, bookId)
	if err != nil {
		log.Logger.Error().Int64("Book", bookId).Err(err).Msg("Delete Book failed")
		return err
	}
	err = queries.DeleteBookTags(ctx, bookId)
	if err != nil {
		log.Logger.Error().Int64("Book", bookId).Err(err).Msg("Delete Book failed")
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Logger.Error().Int64("Book", bookId).Err(err).Msg("Delete Book failed")
		return err
	}
	log.Logger.Debug().Int64("Book", bookId).Msg("End Delete Book")
	return nil
}
func QueryAllBookId(db *sql.DB, ctx context.Context) ([]int64, error) {
	log.Logger.Debug().Msg("Start Query All Book Id")
	queries := NewQueries(db)
	ret, err := queries.QueryAllBookId(ctx)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Query All Book Id Failed")
		return nil, err
	}
	log.Logger.Debug().Msg("End Query All Book Id")
	return ret, nil
}

func UpdateBookFile(db *sql.DB, ctx context.Context, book dbo.Book) error {
	log.Logger.Debug().Str("Hash", book.Hash).Msg("Start UpdateBook File")
	queries := NewQueries(db)
	err := queries.UpdateNewFile(ctx, book)
	if err != nil {
		log.Logger.Error().Err(err).Str("Hash", book.Hash).Msg("UpdateBook File Failed")
		return err
	}

	log.Logger.Debug().Str("Hash", book.Hash).Msg("End UpdateBook File")
	return nil
}
func QueryAllBooks(db *sql.DB, ctx context.Context) ([]dbo.Book, error) {
	log.Logger.Debug().Msg("Start Query All Books")
	queries := NewQueries(db)
	books, err := queries.QueryAllBooks(ctx)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Query All Books")
		return make([]dbo.Book, 0), err
	}
	for i := range books {
		books[i] = readAllMetadata(queries, ctx, books[i])
	}
	log.Logger.Debug().Msg("End Query All Books")
	return books, nil
}
