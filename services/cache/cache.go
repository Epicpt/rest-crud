package services

import (
	"fmt"
	"rest-crud/repository"
	"sync"
)

// Cache - структура кеша
type Cache struct {
	mu             sync.RWMutex
	webmasters     map[int]repository.Webmaster
	placements     map[int][]repository.Placement
	placementsByID map[int]repository.Placement
}

// NewCache - создаёт новый кеш
func NewCache() *Cache {
	return &Cache{
		webmasters:     make(map[int]repository.Webmaster),
		placements:     make(map[int][]repository.Placement),
		placementsByID: make(map[int]repository.Placement),
	}
}

type Cacheable interface {
	GetID() int
}

func (c *Cache) UpdateCache(action string, entity Cacheable) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch e := entity.(type) {
	case repository.Webmaster:
		return c.updateWebmasterCache(action, e)
	case repository.Placement:
		return c.updatePlacementCache(action, e)
	default:
		return fmt.Errorf("неподдерживаемый тип")
	}
}

func (c *Cache) updateWebmasterCache(action string, wm repository.Webmaster) error {
	switch action {
	case "create":
		// Добавляем в кеш
		c.webmasters[wm.ID] = wm
		// Инициализируем пустой список размещений для нового веб-мастера
		c.placements[wm.ID] = []repository.Placement{}
	case "update":
		// Проверяем, существует ли веб-мастер в кеше
		if _, exists := c.webmasters[wm.ID]; !exists {
			return fmt.Errorf("веб-мастер с ID %d не найден в кеше", wm.ID)
		}
		// Обновляем веб-мастера в кеше
		c.webmasters[wm.ID] = wm
	case "delete":
		// Проверяем, существует ли веб-мастер в кеше
		if _, exists := c.webmasters[wm.ID]; !exists {
			return fmt.Errorf("веб-мастер с ID %d не найден в кеше", wm.ID)
		}

		// Удаляем веб-мастера и его размещения из кеша
		delete(c.webmasters, wm.ID)
		delete(c.placements, wm.ID)
	default:
		return fmt.Errorf("неизвестное действие: %s", action)
	}
	return nil
}

func (c *Cache) updatePlacementCache(action string, p repository.Placement) error {
	switch action {
	case "create":
		c.placements[p.UserID] = append(c.placements[p.UserID], p)
		c.placementsByID[p.ID] = p
	case "update":
		// Проверяем, существует ли размещение в кеше
		oldPlacement, exists := c.placementsByID[p.ID]
		if !exists {
			return fmt.Errorf("размещение с ID %d не найдено в кеше", p.ID)
		}
		for i, placement := range c.placements[oldPlacement.UserID] {
			if placement.ID == p.ID {
				c.placements[oldPlacement.UserID][i] = p
				break
			}
		}
		c.placementsByID[p.ID] = p
	case "delete":
		// Проверяем, существует ли размещение
		placement, exists := c.placementsByID[p.ID]
		if !exists {
			return fmt.Errorf("размещение с ID %d не найдено в кеше", p.ID)
		}
		// Удаляем из кеша по ID
		delete(c.placementsByID, p.ID)
		var updated []repository.Placement
		for _, pl := range c.placements[placement.UserID] {
			if pl.ID != p.ID {
				updated = append(updated, pl)
			}
		}
		c.placements[placement.UserID] = updated
	default:
		return fmt.Errorf("неизвестное действие: %s", action)
	}
	return nil
}

// Метод для получения webmasters
func (c *Cache) GetCacheWebmasters() map[int]repository.Webmaster {
	return c.webmasters
}

// Метод для получения placementsByID
func (c *Cache) GetCachePlacementsByID() map[int]repository.Placement {
	return c.placementsByID
}

// LoadCacheFromDB - загружает данные из БД в кеш
func (c *Cache) LoadCacheFromDB(repo *repository.Repository) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	webmasters, err := repo.GetAllWebmasters()
	if err != nil {
		return err
	}
	for _, wm := range webmasters {
		c.webmasters[wm.ID] = wm
	}

	placements, err := repo.GetAllPlacements()
	if err != nil {
		return err
	}
	for _, p := range placements {
		c.placements[p.UserID] = append(c.placements[p.UserID], p)
		c.placementsByID[p.ID] = p
	}

	return nil
}
