package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"rest-crud/repository"

	"github.com/go-chi/chi/v5"
)

// POST /placements  – создание размещения
func (h *Handler) CreatePlacement(w http.ResponseWriter, r *http.Request) {
	var p repository.Placement
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	id, err := h.repo.CreatePlacement(p)
	if err != nil {
		log.Printf("Ошибка создания размещения: %v", err)
		http.Error(w, "Ошибка создания размещения", http.StatusInternalServerError)
		return
	}

	// Обновляем кеш
	h.cache.UpdateCacheWhenCreatePlacement(p, id)

	w.WriteHeader(http.StatusCreated)
	if errEncode := json.NewEncoder(w).Encode(map[string]int{"id": id}); errEncode != nil {
		log.Printf("Ошибка json encode: %v", errEncode)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		return
	}

}

// PUT /placements/:id  – изменение размещения
func (h *Handler) UpdatePlacement(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var p repository.Placement
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	p.ID = id

	if err := h.repo.UpdatePlacement(p); err != nil {
		http.Error(w, "Ошибка обновления размещения", http.StatusInternalServerError)
		return
	}

	// Обновляем кеш
	err = h.cache.UpdateCacheWhenUpdatePlacement(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if errEncode := json.NewEncoder(w).Encode(map[string]string{"message": "Placement обновлён"}); errEncode != nil {
		log.Printf("Ошибка json encode: %v", errEncode)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		return
	}
}

// DELETE /placements/:id  – удаление размещения
func (h *Handler) DeletePlacement(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeletePlacement(id); err != nil {
		http.Error(w, "Ошибка удаления размещения", http.StatusInternalServerError)
		return
	}

	// Обновляем кеш
	err = h.cache.UpdateCacheWhenDeletePlacement(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if errEncode := json.NewEncoder(w).Encode(map[string]string{"message": "Placement удалён"}); errEncode != nil {
		log.Printf("Ошибка json encode: %v", errEncode)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		return
	}
}

// GET /placements  – получение списка размещений
func (h *Handler) GetPlacements(w http.ResponseWriter, r *http.Request) {
	page, limit, err := getPaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	placements := h.cache.GetPlacements(page, limit)

	response := map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"placements": placements,
	}

	w.Header().Set("Content-Type", "application/json")
	if errEncode := json.NewEncoder(w).Encode(response); errEncode != nil {
		log.Printf("Ошибка json encode: %v", errEncode)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		return
	}
}
