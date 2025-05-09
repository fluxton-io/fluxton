package repositories

import (
	"database/sql"
	"errors"
	"fluxton/errs"
	"fluxton/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
)

type TableRepository struct {
	connection       *sqlx.DB
	columnRepository *ColumnRepository
}

func NewTableRepository(connection *sqlx.DB) (*TableRepository, error) {
	columnRepository, err := NewColumnRepository(connection)
	if err != nil {
		return nil, err
	}

	return &TableRepository{
		connection:       connection,
		columnRepository: columnRepository,
	}, nil
}

func (r *TableRepository) Exists(name string) (bool, error) {
	var count int
	err := r.connection.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1", name)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *TableRepository) Create(name string, columns []models.Column) error {
	var defs []string
	var foreignConstraints []string

	for _, column := range columns {
		defs = append(defs, r.columnRepository.BuildColumnDefinition(column))

		if fkQuery, ok := r.columnRepository.BuildForeignKeyConstraint(name, column); ok {
			foreignConstraints = append(foreignConstraints, fkQuery)
		}
	}

	createQuery := fmt.Sprintf("CREATE TABLE %s (\n%s\n);", pq.QuoteIdentifier(name), strings.Join(defs, ",\n"))

	if _, err := r.connection.Exec(createQuery); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	for _, fk := range foreignConstraints {
		if _, err := r.connection.Exec(fk); err != nil {
			return fmt.Errorf("failed to add foreign key constraint: %w", err)
		}
	}

	return nil
}

func (r *TableRepository) Duplicate(existingTable string, newTable string) error {
	_, err := r.connection.Exec(fmt.Sprintf("CREATE TABLE %s AS TABLE %s", pq.QuoteIdentifier(newTable), pq.QuoteIdentifier(existingTable)))
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) List() ([]models.Table, error) {
	var tables []models.Table
	query := `
		SELECT
			c.oid AS id,
			c.relname AS name,
			n.nspname AS schema,
			c.reltuples AS estimated_rows,  -- Approximate row count
			pg_size_pretty(pg_total_relation_size(c.oid)) AS total_size -- Table size (including indexes)
		FROM pg_class c
				 JOIN pg_namespace n ON c.relnamespace = n.oid
		WHERE n.nspname = 'public'  -- Only list tables in the "public" schema
		  AND c.relkind = 'r'  -- 'r' means regular table (excludes views, indexes, etc.)
		ORDER BY c.relname;
	`
	err := r.connection.Select(&tables, query)
	if err != nil {
		return []models.Table{}, err
	}

	return tables, nil
}

func (r *TableRepository) GetByNameInSchema(schema, name string) (models.Table, error) {
	var table models.Table
	query := `
		SELECT
			c.oid AS id,
			c.relname AS name,
			n.nspname AS schema,
			c.reltuples AS estimated_rows,  -- Approximate row count
			pg_size_pretty(pg_total_relation_size(c.oid)) AS total_size -- Table size (including indexes)
		FROM pg_class c
		JOIN pg_namespace n ON c.relnamespace = n.oid
		WHERE n.nspname = $1  -- Filter by schema
		  AND c.relname = $2  -- Filter by table name
		LIMIT 1;
	`

	err := r.connection.Get(&table, query, schema, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Table{}, errs.NewNotFoundError("table.error.notFound")
		}

		return models.Table{}, err
	}

	return table, nil
}

func (r *TableRepository) DropIfExists(name string) error {
	_, err := r.connection.Exec("DROP TABLE IF EXISTS " + pq.QuoteIdentifier(name))
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) Rename(oldName string, newName string) error {
	query := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", pq.QuoteIdentifier(oldName), pq.QuoteIdentifier(newName))
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
