package handlers

import (
	"rest-crud/repository"
	services "rest-crud/services/cache"
)

type Handler struct {
	repo  *repository.Repository
	cache *services.Cache
}

func NewHandler(repo *repository.Repository, cache *services.Cache) *Handler {
	return &Handler{repo: repo, cache: cache}
}
