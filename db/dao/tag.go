package dao

import (
	"context"
	"strings"

	"database/sql"
	"ignis/library/server/db/dbo"

	"github.com/rs/zerolog/log"
)

const queryTagsByBook = `SELECT t.tagId, t.name, t.color FROM tag AS t
JOIN booktags AS bt ON t.tagId = bt.tagId
WHERE bt.bookId = ?
ORDER BY t.name`

func (q *Queries) QueryTagsByBook(ctx context.Context, bookid int64) ([]dbo.Tag, error) {
	rows, err := q.db.QueryContext(ctx, queryTagsByBook, bookid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dbo.Tag
	for rows.Next() {
		var i dbo.Tag
		if err := rows.Scan(&i.TagId, &i.Name, &i.Color); err != nil {
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

const updateTag = `UPDATE tag SET name=?, color=? WHERE tagId =?`

func (q *Queries) UpdateTag(ctx context.Context, tag dbo.Tag) error {
	_, err := q.db.ExecContext(ctx, updateTag, tag.Name, tag.Color, tag.TagId)
	return err
}

const getTagById = `SELECT tagId, name, color FROM tag WHERE tagId = ?`

func (q *Queries) GetTagById(ctx context.Context, tagId int64) (dbo.Tag, error) {
	row := q.db.QueryRowContext(ctx, getTagById, tagId)
	var i dbo.Tag
	err := row.Scan(&i.TagId, &i.Name, &i.Color)
	return i, err
}

const getTag = `SELECT tagId, name, color FROM tag WHERE name = ?`

func (q *Queries) GetTag(ctx context.Context, name string) (dbo.Tag, error) {
	row := q.db.QueryRowContext(ctx, getTag, name)
	var i dbo.Tag
	err := row.Scan(&i.TagId, &i.Name, &i.Color)
	return i, err
}

const createTag = `INSERT INTO tag (name, color) VALUES (?, ?)`

func (q *Queries) CreateTag(ctx context.Context, arg dbo.Tag) error {
	_, err := q.db.ExecContext(ctx, createTag, arg.Name, arg.Color)
	return err
}

// bookSeries query func
const bindBookTag = `INSERT INTO booktags (bookId,tagId) VALUES (?,?)`

func (q *Queries) BindBookTag(ctx context.Context, bookId int64, tagId int64) error {
	_, err := q.db.ExecContext(ctx, bindBookTag, bookId, tagId)
	return err
}

const divideAllTagFromBook = "DELETE FROM booktags WHERE bookid=?"
const divideTagsFromBookStart = "DELETE FROM booktags WHERE bookid=? and tagId not in (?"
const divideTagsFromBookEnd = ")"

func (q *Queries) DivideTagsFromBook(ctx context.Context, bookId int64, tagsIds []int64) error {
	if len(tagsIds) == 0 {
		_, err := q.db.ExecContext(ctx, divideAllTagFromBook, bookId)
		return err
	}
	args := make([]interface{}, len(tagsIds)+1)
	args[0] = bookId
	for i, id := range tagsIds {
		args[i+1] = id
	}
	query := divideTagsFromBookStart + strings.Repeat(",?", len(tagsIds)-1) + divideTagsFromBookEnd
	_, err := q.db.ExecContext(ctx, query, args...)
	return err
}

// ///////////////////////////////
func UpdateTagTBLK(q *Queries, ctx context.Context, bookId int64, tag dbo.Tag) error {
	log.Logger.Debug().Int64("Tag", tag.TagId.Int64).Msg("Start UpdateTagTBLK")
	err := q.UpdateTag(ctx, tag)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Tag", tag.TagId.Int64).Msg("UpdateTagTBLK failed")
		return err
	}
	log.Logger.Debug().Int64("Tag", tag.TagId.Int64).Msg("End UpdateTagTBLK")
	return nil
}

func QueryTagsByBookTBLK(q *Queries, ctx context.Context, bookId int64) ([]dbo.Tag, error) {
	log.Logger.Debug().Int64("Book", bookId).Msg("Start Query Tags by Book")
	tags, err := q.QueryTagsByBook(ctx, bookId)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Book", bookId).Msg("Query Tags by Book failed")
		return nil, err
	}
	log.Logger.Debug().Int64("Book", bookId).Msg("End Query Tags by Book ")
	return tags, nil
}
func GetTagById(db *sql.DB, ctx context.Context, tagId int64) (dbo.Tag, error) {
	log.Logger.Debug().Int64("Tag", tagId).Msg("Start Get Tag By Id")
	queries := NewQueries(db)
	series, err := queries.GetTagById(ctx, tagId)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Tag", tagId).Msg("Get Tag By Id failed")
		return dbo.Tag{}, err
	}
	log.Logger.Debug().Int64("Tag", tagId).Msg("End Get Tag By Id")
	return series, nil
}
func AddTagToBookTBLK(q *Queries, ctx context.Context, bookId int64, tag dbo.Tag) error {
	log.Logger.Debug().Int64("Book", bookId).Str("Tag", tag.Name).Int64("TagId", tag.TagId.Int64).Bool("TagId is Valid", tag.TagId.Valid).Msg("Start Add Tag")
	var err error
	if !tag.TagId.Valid {
		tagName := tag.Name
		tagColor := tag.Color
		tag, err = q.GetTag(ctx, tagName)
		if err == sql.ErrNoRows {
			log.Logger.Trace().Str("Tag", tag.Name).Msg("Tag not found, create")
			tag = dbo.Tag{Name: tagName, Color: tagColor}
			err = q.CreateTag(ctx, tag)
			if err != nil {
				log.Logger.Error().Err(err).Str("Tag", tag.Name).Msg("Create Tag failed")
				return err
			}
			tagId, err := q.GetLastId(ctx)
			if err != nil {
				log.Logger.Error().Err(err).Str("Tag", tag.Name).Msg("Create Tag failed")
				return err
			}
			tag.TagId = sql.NullInt64{Int64: tagId, Valid: true}
		} else if err != nil {
			log.Logger.Error().Err(err).Str("Tag", tag.Name).Msg("Get Tag failed")
			return err
		}
	}
	err = q.BindBookTag(ctx, bookId, tag.TagId.Int64)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Book", bookId).Str("Tag", tag.Name).Int64("TagId", tag.TagId.Int64).Msg("Binding Tag to Book failed")
		return err
	}
	log.Logger.Debug().Int64("Book", bookId).Str("Tag", tag.Name).Int64("TagId", tag.TagId.Int64).Msg("End Add Tag")
	return nil
}
func DivideBookAllOtherTags(q *Queries, ctx context.Context, bookId int64, tagIds []int64) error {
	log.Logger.Debug().Int64("Book", bookId).Ints64("TagIds", tagIds).Msg("Start Divide Book All Other Tags")
	err := q.DivideTagsFromBook(ctx, bookId, tagIds)
	if err != nil {
		log.Logger.Error().Err(err).Int64("Book", bookId).Ints64("TagIds", tagIds).Msg("Divide Book All Other Tags Failed")
		return err
	}

	log.Logger.Debug().Int64("Book", bookId).Ints64("TagIds", tagIds).Msg("End Divide Book All Other Tags")
	return nil
}
