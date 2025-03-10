package models

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
    _ "github.com/joho/godotenv/autoload"
)

type database struct {
    *sql.DB
}

var (
    dbname     = os.Getenv("SGUBLOGSITE_DB_DATABASE")
    password   = os.Getenv("SGUBLOGSITE_DB_PASSWORD")
    username   = os.Getenv("SGUBLOGSITE_DB_USERNAME")
    port       = os.Getenv("SGUBLOGSITE_DB_PORT")
    host       = os.Getenv("SGUBLOGSITE_DB_HOST")
    dbInstance *database
)

func New() database {
    // Reuse Connection
    if dbInstance != nil {
        return *dbInstance
    }

    // Opening a driver typically will not attempt to connect to the database.
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
    if err != nil {
        // This will not be a connection error, but a DSN parse error or
        // another initialization error.
        log.Fatal(err)
    }
    db.SetConnMaxLifetime(0)
    db.SetMaxIdleConns(50)
    db.SetMaxOpenConns(50)

    dbInstance = &database{db}
    return *dbInstance
}
