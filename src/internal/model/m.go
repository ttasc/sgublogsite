package model

import (
    "context"
    "database/sql"
    "github.com/ttasc/sgublogsite/src/internal/model/repos"
)

type Model struct {
    ctx   context.Context
    db    *sql.DB
    query *repos.Queries
}

func New(db *sql.DB) *Model {
    var (
        ctx     = context.Background()
        query   = repos.New(db)
    )
    return &Model{
        ctx:   ctx,
        db:    db,
        query: query,
    }
}
