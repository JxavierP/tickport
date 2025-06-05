package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations(db *sql.DB) error  {
	migrationsDir := "internal/database/migrations/"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("Failed to read migrations directory: %w", err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".sql" {
			continue	
		}
		
		path := filepath.Join(migrationsDir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Failed to read migrations: %w", err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("Failed to execute migration %s: %w", file.Name(), err)
		}

		fmt.Printf("✅ Applied migration: %s\n", file.Name())
	}

	return nil
}
