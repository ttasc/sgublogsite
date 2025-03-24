package model

import (
    "context"
    "database/sql"
    "github.com/ttasc/sgublogsite/src/internal/model/repos"
)

type Model struct {
    DB    *sql.DB
    ctx   context.Context
    query *repos.Queries
}

func New(db *sql.DB) *Model {
    var (
        ctx     = context.Background()
        query   = repos.New(db)
    )
    return &Model{
        DB:    db,
        ctx:   ctx,
        query: query,
    }
}
