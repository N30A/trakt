package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func buildConnString(config DBConfig) string {
	params := config.Params
	if params != "" && !strings.HasPrefix(params, "?") {
		params = "?" + params
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s",
		url.QueryEscape(config.User),
		url.QueryEscape(config.Password),
		config.Host,
		config.Port,
		config.Name,
		params,
	)
}

func connectToDB(ctx context.Context, config DBConfig) (*pgxpool.Pool, error) {
	connString := buildConnString(config)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	log.Printf("connected to database %s on %s:%s", config.Name, config.Host, config.Port)
	return pool, nil
}
