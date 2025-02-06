package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"rest-crud/repository"

	"github.com/go-chi/chi/v5"
)

// POST /webmasters – создание веб-мастера

func (h *Handler) CreateWebmaster(w http.ResponseWriter, r *http.Request) {
	var wm repository.Webmaster
	if err := json.NewDecoder(r.Body).Decode(&wm); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	log.Printf("Создаём вебмастера: %+v", wm) // удалить

	id, err := h.repo.CreateWebMaster(&wm)
	if err != nil {
		log.Printf("Ошибка создания веб-мастера: %v", err)
		http.Error(w, "Ошибка создания веб-мастера", http.StatusInternalServerError)
		return
	}

	// Обновляем кеш
	h.cache.UpdateCacheWhenCreateWebmaster(wm, id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) UpdateWebmaster(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var wm repository.Webmaster
	if err := json.NewDecoder(r.Body).Decode(&wm); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	wm.ID = id

	if err := h.repo.UpdateWebmaster(wm); err != nil {
		http.Error(w, "Ошибка обновления веб-мастера", http.StatusInternalServerError)
		return
	}

	// Обновляем кеш
	if err := h.cache.UpdateCacheWhenUpdateWebmaster(wm); err != nil {
		log.Printf("Ошибка обновления кеша: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Webmaster обновлён"})
}

func (h *Handler) DeleteWebmaster(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteWebmaster(id); err != nil {
		http.Error(w, "Ошибка удаления веб-мастера", http.StatusInternalServerError)
		return
	}

	// Удаляем из кеша
	if err := h.cache.UpdateCacheWhenDeleteWebmaster(id); err != nil {
		log.Printf("Ошибка удаления из кеша: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Webmaster удалён"})
}

func (h *Handler) GetWebmasters(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationParams(r)

	webmasters := h.cache.GetWebmasters(page, limit)

	response := map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"webmasters": webmasters,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
