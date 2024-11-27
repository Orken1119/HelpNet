package pkg

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewConn() (*pgxpool.Pool, error) {

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:Akatsuki2005@localhost:5433/postgres?sslmode=disable"
	}

	cfg, err := pgxpool.ParseConfig(connStr)

	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	cfg.MaxConns = int32(10)
	cfg.MinConns = int32(2)
	cfg.HealthCheckPeriod = 1 * time.Minute
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		Timeout:   cfg.ConnConfig.ConnectTimeout,
	}).DialContext

	pool, err := pgxpool.ConnectConfig(context.Background(), cfg)

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	fmt.Println("success connection")

	return pool, nil
}
