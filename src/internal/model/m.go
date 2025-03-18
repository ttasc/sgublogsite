package model

import (
    "context"
    "database/sql"
    "sgublogsite/src/internal/model/repos"
    "sgublogsite/src/internal/utils"
)

type Model struct {
    ctx   context.Context
    db    *sql.DB
    query *repos.Queries
}

func New() *Model {
    var (
        ctx     = context.Background()
        db      = utils.NewDB()
        query   = repos.New(db)
    )
    return &Model{
        ctx:   ctx,
        db:    db,
        query: query,
    }
}
