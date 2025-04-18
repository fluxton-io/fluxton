package repositories

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	"os"
	"path/filepath"
	"strings"
)

type DatabaseRepository struct {
	db *sqlx.DB
}

func NewDatabaseRepository(injector *do.Injector) (*DatabaseRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &DatabaseRepository{db: db}, nil
}

func (r *DatabaseRepository) Create(name string, userUUID uuid.NullUUID) error {
	_, err := r.db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, name))
	if err != nil {
		return err
	}

	if userUUID.Valid {
		return r.importSeedFiles(name, userUUID.UUID)
	}

	return nil
}

func (r *DatabaseRepository) DropIfExists(name string) error {
	_, err := r.db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, name))
	if err != nil {
		return err
	}

	return nil
}

func (r *DatabaseRepository) Recreate(name string) error {
	err := r.DropIfExists(name)
	if err != nil {
		return err
	}

	err = r.Create(name, uuid.NullUUID{})
	if err != nil {
		return err
	}

	return nil
}

func (r *DatabaseRepository) List() ([]string, error) {
	var databases []string
	err := r.db.Select(&databases, "SELECT datname FROM pg_database WHERE datistemplate = false")
	if err != nil {
		return []string{}, err
	}

	return databases, nil
}

func (r *DatabaseRepository) Exists(name string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM pg_database WHERE datname = $1", name)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Connect TODO: create actual user for using here
func (r *DatabaseRepository) Connect(name string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=%s port=5432",
		os.Getenv("DATABASE_USER"),
		name,
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_SSL_MODE"),
	)

	return sqlx.Connect("postgres", connStr)
}

func (r *DatabaseRepository) importSeedFiles(databaseName string, userUUID uuid.UUID) error {
	connection, err := r.Connect(databaseName)
	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}
	defer connection.Close()

	seedDir := "database/seeders/client"

	// Read all files in the directory
	files, err := os.ReadDir(seedDir)
	if err != nil {
		return fmt.Errorf("could not read seed directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filePath := filepath.Join(seedDir, file.Name())

		// Load the contents of the SQL file
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("DB: %s => Skipping file %s: could not read file: %v", databaseName, filePath, err)

			continue
		}

		// Why split: If whole file is executed at once, and there is an error in one of the queries,
		// the whole file will be skipped. This way, we can execute the queries one by one.
		sqlCommands := strings.Split(string(sqlContent), ";")
		for _, query := range sqlCommands {
			if strings.Contains(query, "{{USER_ROLE}}") {
				query = strings.ReplaceAll(query, "{{USER_ROLE}}", fmt.Sprintf(`usr_%s`, strings.ReplaceAll(userUUID.String(), "-", "_")))
			}

			if _, err := connection.Exec(query); err != nil {
				log.Printf("DB: %s => Skipping file %s: could not execute SQL: %v", databaseName, filePath, err)

				continue
			}
		}

		log.Printf("DB: %s => Successfully executed seed file %s", databaseName, filePath)
	}

	return nil
}
