package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"
)

func ConnectWithRetry(dsn string, maxRetries int) (*sql.DB, error) {
    var db *sql.DB
    var err error

    for attempt := 1; attempt <= maxRetries; attempt++ {

        log.Printf("🔄 Database connection attempt %d/%d...", attempt, maxRetries)

        db, err = sql.Open("postgres", dsn)
        if err == nil {

            pingErr := db.Ping()
            if pingErr == nil {

                log.Println("✅ Database connected successfully!")
                return db, nil
            }

            err = pingErr
            db.Close()
        }

        if attempt < maxRetries {

            backoffDuration := time.Duration(1<<uint(attempt-1)) * time.Second


            log.Printf("⚠️  Connection failed: %v. Retrying in %v...", err, backoffDuration)
            time.Sleep(backoffDuration)
        }
    }

    return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

