package rdbms

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "github.com/duckdb/duckdb-go/v2"
)

type DuckDBClient struct {
	db *sql.DB
}

func NewDuckDBClient() (*DuckDBClient, error) {
	db, err := sql.Open("duckdb", "memory.db")
	if err != nil {
		return nil, err
	}
	err = createSequences(db)
	if err != nil {
		return nil, err
	}
	err = createMemoryTable(db)
	if err != nil {
		return nil, err
	}
	err = createConversationTable(db)
	if err != nil {
		return nil, err
	}
	err = createMemoryMetaTable(db)
	if err != nil {
		return nil, err
	}
	return &DuckDBClient{db: db}, nil

}

func (r *DuckDBClient) GetDB() *sql.DB {
	return r.db
}

func (r *DuckDBClient) Mount(dir string, files map[string]io.Reader) error {
	for fileName, reader := range files {
		if fileName == "memory.parquet" {
			// write the reader to the file memory.parquet inside dir
			bytes, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", fileName, err)
			}
			err = os.WriteFile(filepath.Join(dir, "memory.parquet"), bytes, 0644)
			if err != nil {
				return fmt.Errorf("error writing file %s: %w", fileName, err)
			}
			// now use the sql command of duckdb to copy to the table memories
			_, err = r.db.Exec(fmt.Sprintf("COPY memories FROM '%s' (FORMAT PARQUET)", filepath.Join(dir, "memory.parquet")))
			if err != nil {
				return fmt.Errorf("error copying file %s: %w", fileName, err)
			}
		} else if fileName == "conversations.parquet" {
			// write the reader to the file conversations.parquet inside dir
			bytes, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", fileName, err)
			}
			err = os.WriteFile(filepath.Join(dir, "conversations.parquet"), bytes, 0644)
			if err != nil {
				return fmt.Errorf("error writing file %s: %w", fileName, err)
			}
			// now use the sql command of duckdb to copy to the table conversations
			_, err = r.db.Exec(fmt.Sprintf("COPY conversations FROM '%s' (FORMAT PARQUET)", filepath.Join(dir, "conversations.parquet")))
			if err != nil {
				return fmt.Errorf("error copying file %s: %w", fileName, err)
			}
		} else if fileName == "memories_meta.parquet" {
			// write the reader to the file memories_meta.parquet inside dir
			bytes, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", fileName, err)
			}
			err = os.WriteFile(filepath.Join(dir, "memories_meta.parquet"), bytes, 0644)
			if err != nil {
				return fmt.Errorf("error writing file %s: %w", fileName, err)
			}
			// now use the sql command of duckdb to copy to the table memories_meta
			_, err = r.db.Exec(fmt.Sprintf("COPY memories_meta FROM '%s' (FORMAT PARQUET)", filepath.Join(dir, "memories_meta.parquet")))
			if err != nil {
				return fmt.Errorf("error copying file %s: %w", fileName, err)
			}
		}
	}
	return nil
}

func (r *DuckDBClient) Export(dir string) ([]os.File, error) {
	var files []os.File
	memoryPath := filepath.Join(dir, "memory.parquet")
	conversationsPath := filepath.Join(dir, "conversations.parquet")
	memoriesMetaPath := filepath.Join(dir, "memories_meta.parquet")
	_, err := r.db.Exec(fmt.Sprintf("COPY (SELECT * FROM memories) TO '%s' (FORMAT PARQUET)", memoryPath))
	if err != nil {
		return nil, err
	}
	memoryFile, err := os.Open(memoryPath)
	if err != nil {
		return nil, err
	}
	files = append(files, *memoryFile)
	_, err = r.db.Exec(fmt.Sprintf("COPY (SELECT * FROM conversations) TO '%s' (FORMAT PARQUET)", conversationsPath))
	if err != nil {
		return nil, err
	}
	conversationsFile, err := os.Open(conversationsPath)
	if err != nil {
		return nil, err
	}
	files = append(files, *conversationsFile)
	_, err = r.db.Exec(fmt.Sprintf("COPY (SELECT * FROM memories_meta) TO '%s' (FORMAT PARQUET)", memoriesMetaPath))
	if err != nil {
		return nil, err
	}
	memoriesMetaFile, err := os.Open(memoriesMetaPath)
	if err != nil {
		return nil, err
	}
	files = append(files, *memoriesMetaFile)
	return files, nil
}

func createSequences(db *sql.DB) error {
	queries := []string{
		"CREATE SEQUENCE IF NOT EXISTS memories_id_seq START 1",
		"CREATE SEQUENCE IF NOT EXISTS conversations_id_seq START 1",
	}
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func createMemoryTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS memories (
			id INTEGER PRIMARY KEY DEFAULT nextval('memories_id_seq'),
			uuid UUID NOT NULL,
			conversation_id UUID NOT NULL,
			query TEXT NOT NULL,
			response TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func createMemoryMetaTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS memories_meta (
			id INTEGER PRIMARY KEY DEFAULT nextval('memories_id_seq'),
			uuid UUID NOT NULL,
			conversation_id UUID NOT NULL,
			query TEXT NOT NULL,
			response TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func createConversationTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS conversations (
			id INTEGER PRIMARY KEY DEFAULT nextval('conversations_id_seq'),
			uuid UUID NOT NULL,
			agent TEXT NOT NULL,
			user TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
