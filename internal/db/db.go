package db

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

type Document struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

func NewDB(ctx context.Context, connString string) (*DB, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func formatVector(embedding []float32) string {
	strs := make([]string, len(embedding))
	for i, v := range embedding {
		strs[i] = strconv.FormatFloat(float64(v), 'f', -1, 32)
	}
	return "[" + strings.Join(strs, ",") + "]"
}

func (db *DB) SaveDocument(ctx context.Context, content string, embedding []float32) (int64, error) {
	vecStr := formatVector(embedding)
	var id int64
	err := db.pool.QueryRow(ctx, "INSERT INTO documents (content, embedding) VALUES ($1, $2::vector) RETURNING id", content, vecStr).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to save document: %w", err)
	}
	return id, nil
}

func (db *DB) SearchNearest(ctx context.Context, embedding []float32, limit int) ([]Document, error) {
	vecStr := formatVector(embedding)
	rows, err := db.pool.Query(ctx, "SELECT id, content FROM documents ORDER BY embedding <-> $1::vector LIMIT $2", vecStr, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}
	defer rows.Close()

	var docs []Document
	for rows.Next() {
		var doc Document
		if err := rows.Scan(&doc.ID, &doc.Content); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}
