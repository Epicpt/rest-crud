package main

import (
	"log"
	"log/slog"
	"net/http"

	"rest-crud/config"
	"rest-crud/handlers"
	"rest-crud/repository"
	"rest-crud/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключаемся к БД
	db, err := sqlx.Open("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Проверяем и создаём таблицы
	if err := ensureTables(db); err != nil {
		slog.Error("Ошибка инициализации таблиц", "error", err.Error())
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	cache := services.NewCache()

	// Загружаем кеш из БД
	if err := cache.LoadCacheFromDB(repo); err != nil {
		log.Fatalf("Ошибка загрузки кеша: %v", err)
	}
	log.Printf("Кэш успешно загружен %v", cache)

	handler := handlers.NewHandler(repo, cache)

	r := chi.NewRouter()

	// POST-эндпоинты
	r.Post("/webmasters", handler.CreateWebmaster)
	r.Post("/placements", handler.CreatePlacement)

	// Put-эндпоинты
	r.Put("/webmasters/{id}", handler.UpdateWebmaster)
	r.Put("/placements/{id}", handler.UpdatePlacement)

	// Delete-эндпоинты
	r.Delete("/webmasters/{id}", handler.DeleteWebmaster)
	r.Delete("/placements/{id}", handler.DeletePlacement)

	// GET-эндпоинты
	r.Get("/placements", handler.GetPlacements)
	r.Get("/webmasters", handler.GetWebmasters)

	log.Printf("Сервер запущен на порту %s", cfg.Server.Port)
	http.ListenAndServe(":"+cfg.Server.Port, r)

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
		log.Printf("Таблица успешно проверена или создана: %s", query)
	}

	return nil
}
