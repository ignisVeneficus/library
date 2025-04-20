package dao

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

//go:embed schema.sql
var createDatabase string

func (q *Queries) CreateDatabase(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, createDatabase)
	return err
}

const getLastId = `select LAST_INSERT_ID()`

func (q *Queries) GetLastId(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLastId)
	var last_insert_id int64
	err := row.Scan(&last_insert_id)
	return last_insert_id, err
}

func GetTx(db *sql.DB, ctx context.Context) (*sql.Tx, error) {
	return db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func NewQueries(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

var ErrDataNotFound = errors.New("data not found")

func GetDataNotFoundError(table string) error {
	return fmt.Errorf("%w, table: %s", ErrDataNotFound, table)
}

func CreateDatabase(db *sql.DB, ctx context.Context) error {
	log.Logger.Debug().Msg("Start Create Database")
	queries := NewQueries(db)
	err := queries.CreateDatabase(ctx)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Create Database Failed")
		return err
	}
	log.Logger.Debug().Msg("End Create Database")
	return nil
}
