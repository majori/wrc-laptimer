package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"fmt"

	"github.com/marcboeker/go-duckdb/v2"
)

type Database struct {
	ctx      context.Context
	conn     driver.Conn
	db       *sql.DB
	appender *duckdb.Appender
}

//go:embed schema.sql
var dbSchema string

func NewDatabase(ctx context.Context, path string) (*Database, error) {
	connector, err := duckdb.NewConnector(path, nil)
	if err != nil {
		return nil, err
	}

	dbConnection, err := connector.Connect(ctx)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)

	_, err = db.ExecContext(ctx, dbSchema)
	if err != nil {
		return nil, err
	}

	appender, err := duckdb.NewAppenderFromConn(dbConnection, "", "telemetry")
	if err != nil {
		return nil, fmt.Errorf("could not create new appender for telemetry: %w", err)
	}

	return &Database{
		ctx:      ctx,
		conn:     dbConnection,
		db:       db,
		appender: appender,
	}, nil
}

func (d *Database) Close() {
	if d.appender != nil {
		//nolint:errcheck
		d.appender.Close()
	}
	if d.db != nil {
		//nolint:errcheck
		d.db.Close()
	}
	if d.conn != nil {
		//nolint:errcheck
		d.conn.Close()
	}
}

func (d *Database) exec(query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(d.ctx, query, args...)
}

func (d *Database) query(query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(d.ctx, query, args...)
}

func (d *Database) queryRow(query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(d.ctx, query, args...)
}

// Supports only SELECT queries. Safe(?) to execute from frontend
func (d *Database) ExecuteSelectQuery(query string) (string, error) {
	var result string
	err := d.queryRow(`
		SELECT COALESCE(CAST(to_json(list(t)) AS VARCHAR), '[]') 
		FROM (SELECT * FROM json_execute_serialized_sql(json_serialize_sql(?::STRING))) t
	`, query).Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}
