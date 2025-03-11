package utils

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
    _ "github.com/joho/godotenv/autoload"
)

var (
    dbname     = os.Getenv("SGUBLOGSITE_DB_DATABASE")
    password   = os.Getenv("SGUBLOGSITE_DB_PASSWORD")
    username   = os.Getenv("SGUBLOGSITE_DB_USERNAME")
    port       = os.Getenv("SGUBLOGSITE_DB_PORT")
    host       = os.Getenv("SGUBLOGSITE_DB_HOST")
    params     = os.Getenv("SGUBLOGSITE_DB_PARAMS")
    dbInstance *sql.DB
)

func NewDB() *sql.DB {
    // Reuse Connection
    if dbInstance != nil {
        return dbInstance
    }

    // Opening a driver typically will not attempt to connect to the database.
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", username, password, host, port, dbname, params))
    if err != nil {
        // This will not be a connection error, but a DSN parse error or
        // another initialization error.
        log.Fatal(err)
    }
    db.SetConnMaxLifetime(0)
    db.SetMaxIdleConns(50)
    db.SetMaxOpenConns(50)

    dbInstance = db
    return dbInstance
}
