package services

import (
    "context"
    "database/sql"
    "sgublogsite/src/internal/model/repos"
    "sgublogsite/src/internal/utils"
)

type services struct {
    ctx   context.Context
    db    *sql.DB
    query *repos.Queries
}

func new() *services {
    var (
        ctx     = context.Background()
        db      = utils.NewDB()
        query   = repos.New(db)
    )
    return &services{
        ctx:   ctx,
        db:    db,
        query: query,
    }
}

var s = new()
