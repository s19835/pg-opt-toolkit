package connector

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/s19835/pg-opt-toolkit/pkg/models"
)

type PGConnector struct {
	db *sql.DB
}

func NewPGConnector(cfg models.PGConfig) (*PGConnector, error) {
	// create a database connection
	db, err := sql.Open("pgx", cfg.URL)

	if err != nil {
		return nil, fmt.Errorf("failed to connect database %w", err)
	}

	// Limits active connections to prevent overloading DB
	db.SetMaxOpenConns(5)

	// Maintains 2 idle connections for quick reuse
	db.SetMaxIdleConns(2)

	// Recycles connections after 30 minutes to prevent stale connections
	db.SetConnMaxLifetime(30 * time.Minute)

	// validate connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error establishing connections %w", err)
	}

	fmt.Println("Successfully Connected to PostgreSQL")
	return &PGConnector{db: db}, nil
}

// Analyze and out put result as json
func (c *PGConnector) ExplainAnalyze(query string) (string, error) {
	explainQuery := "EXPLAIN (ANALYZE, COSTS, VERBOSE, BUFFERS, FORMAT JSON)" + query

	var result string
	err := c.db.QueryRow(explainQuery).Scan(&result)
	if err != nil {
		return "", fmt.Errorf("fail to execute EXPLAIN ANALYZE: %w", err)
	}

	return result, nil
}

func (c *PGConnector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
