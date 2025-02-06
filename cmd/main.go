package main

import (
	"log"
	"net/http"

	"rest-crud/config"
	"rest-crud/handlers"
	"rest-crud/repository"
	"rest-crud/services"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключаемся к БД
	db, err := repository.Init(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

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
