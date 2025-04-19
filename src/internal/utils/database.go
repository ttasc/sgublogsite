package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_520_ci&parseTime=True&loc=Local", username, password, host, port, dbname))
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

func DBhealth(db *sql.DB) map[string]string {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    stats := make(map[string]string)

    // Ping the database
    err := db.PingContext(ctx)
    if err != nil {
        stats["status"] = "down"
        stats["error"] = fmt.Sprintf("db down: %v", err)
        log.Fatalf("db down: %v", err) // Log the error and terminate the program
        return stats
    }

    // Database is up, add more statistics
    stats["status"] = "up"
    stats["message"] = "It's healthy"

    // Get database stats (like open connections, in use, idle, etc.)
    dbStats := db.Stats()
    stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
    stats["in_use"] = strconv.Itoa(dbStats.InUse)
    stats["idle"] = strconv.Itoa(dbStats.Idle)
    stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
    stats["wait_duration"] = dbStats.WaitDuration.String()
    stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
    stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

    // Evaluate stats to provide a health message
    if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
        stats["message"] = "The database is experiencing heavy load."
    }
    if dbStats.WaitCount > 1000 {
        stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
    }

    if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
        stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
    }

    if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
        stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
    }

    return stats
}
