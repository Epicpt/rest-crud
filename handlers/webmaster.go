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

	id, err := h.repo.CreateWebMaster(&wm)
	if err != nil {
		log.Printf("Ошибка создания веб-мастера: %v", err)
		http.Error(w, "Ошибка создания веб-мастера", http.StatusInternalServerError)
		return
	}
	// Назначаем ID новому веб-мастеру
	wm.ID = id

	// Обновляем кеш
	if err = h.cache.UpdateCache("create", wm); err != nil {
		log.Printf("Ошибка обновления кеша: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	if errEncode := json.NewEncoder(w).Encode(map[string]int{"id": id}); errEncode != nil {
		log.Printf("Ошибка json encode: %v", errEncode)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		return
	}
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
	if err := h.cache.UpdateCache("update", wm); err != nil {
		log.Printf("Ошибка обновления кеша: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	if errEncode := json.NewEncoder(w).Encode(map[string]string{"message": "Webmaster обновлён"}); errEncode != nil {
		log.Printf("Ошибка json encode: %v", errEncode)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		return
	}
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

	var wm repository.Webmaster
	wm.ID = id

	// Удаляем из кеша
	if err := h.cache.UpdateCache("delete", wm); err != nil {
		log.Printf("Ошибка удаления из кеша: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	if errEncode := json.NewEncoder(w).Encode(map[string]string{"message": "Webmaster удалён"}); errEncode != nil {
		log.Printf("Ошибка json encode: %v", errEncode)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetWebmasters(w http.ResponseWriter, r *http.Request) {
	page, limit, err := getPaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	webmasters := h.cache.GetWebmasters(page, limit)

	response := map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"webmasters": webmasters,
	}

	w.Header().Set("Content-Type", "application/json")
	if errEncode := json.NewEncoder(w).Encode(response); errEncode != nil {
		log.Printf("Ошибка json encode: %v", errEncode)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		return
	}
}
