package repository

import (
	"fmt"
	"log"
	"rest-crud/config"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

// NewRepository создает новый экземпляр репозитория
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Init инициализирует репозиторий
func Init(cfg *config.Config) (*sqlx.DB, error) {

	// Подключаемся к БД
	db, err := sqlx.Open("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %v", err)
	}

	// Проверяем и создаём таблицы
	if err := ensureTables(db); err != nil {
		return nil, fmt.Errorf("Ошибка инициализации таблиц: %v", err)
	}

	return db, nil
}

// ensureTables проверяет и создаёт таблицы
func ensureTables(db *sqlx.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS webmasters (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL,
            status TEXT NOT NULL CHECK (status IN ('active', 'banned'))
        );`,
		`CREATE TABLE IF NOT EXISTS placements (
            id SERIAL PRIMARY KEY,
            user_id INT REFERENCES webmasters(id) ON DELETE CASCADE,
            name TEXT NOT NULL,
            description TEXT
        );
		CREATE INDEX IF NOT EXISTS idx_placements_user_id ON placements(user_id);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Ошибка выполнения запроса: %s, ошибка: %v", query, err)
			return err
		}

	}

	return nil
}
