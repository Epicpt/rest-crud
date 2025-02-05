package handlers

import (
	"rest-crud/repository"
	"rest-crud/services"
)

type Handler struct {
	repo  *repository.Repository
	cache *services.Cache
}

func NewHandler(repo *repository.Repository, cache *services.Cache) *Handler {
	return &Handler{repo: repo, cache: cache}
}
