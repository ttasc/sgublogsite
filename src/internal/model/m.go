package model

import (
    "context"
    "database/sql"
    "github.com/ttasc/sgublogsite/src/internal/model/repos"
    "github.com/ttasc/sgublogsite/src/internal/utils"
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
