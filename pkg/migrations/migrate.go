package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// Migrate выполняет все миграции из указанной директории
func Migrate(db *sql.DB, migrationsDir string) error {
	// Получаем список файлов миграций
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	// Сортируем файлы по имени
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Создаем таблицу для отслеживания выполненных миграций
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %v", err)
	}

	// Выполняем каждую миграцию
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Проверяем, была ли миграция уже выполнена
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", file.Name()).Scan(&count)
		if err != nil {
			return fmt.Errorf("error checking migration status: %v", err)
		}

		if count > 0 {
			continue
		}

		// Читаем SQL-файл
		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", file.Name(), err)
		}

		// Выполняем миграцию в транзакции
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction: %v", err)
		}

		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing migration %s: %v", file.Name(), err)
		}

		// Записываем информацию о выполненной миграции
		_, err = tx.Exec("INSERT INTO migrations (name) VALUES ($1)", file.Name())
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error recording migration %s: %v", file.Name(), err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("error committing migration %s: %v", file.Name(), err)
		}

		fmt.Printf("Applied migration: %s\n", file.Name())
	}

	return nil
}
