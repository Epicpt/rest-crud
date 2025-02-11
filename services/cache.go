package services

import (
	"fmt"
	"rest-crud/repository"
	"sort"
	"sync"
)

// WebmasterWithPlacements - структура, содержащая вебмастера и его размещения
type WebmasterWithPlacements struct {
	Webmaster  repository.Webmaster
	Placements []repository.Placement
}

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

// Метод для получения webmasters
func (c *Cache) GetCacheWebmasters() map[int]repository.Webmaster {
	return c.webmasters
}

// Метод для получения placementsByID
func (c *Cache) GetCachePlacementsByID() map[int]repository.Placement {
	return c.placementsByID
}

// GetPlacements - получает список размещений с пагинацией
func (c *Cache) GetPlacements(page, limit int) []repository.Placement {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Извлекаем и сортируем ключи
	keys := make([]int, 0, len(c.placementsByID))
	for k := range c.placementsByID {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	start := (page - 1) * limit
	count := 0
	var paginatedPlacements []repository.Placement

	for _, k := range keys {
		if count >= start && count < start+limit {
			paginatedPlacements = append(paginatedPlacements, c.placementsByID[k])
		}
		count++
		if count >= start+limit {
			break
		}
	}

	return paginatedPlacements
}

// GetWebmasters - получает список вебмастеров с вложенными размещениями и пагинацией
func (c *Cache) GetWebmasters(page, limit int) []repository.Webmaster {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Извлекаем и сортируем ключи
	keys := make([]int, 0, len(c.webmasters))
	for k := range c.webmasters {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	start := (page - 1) * limit
	count := 0
	var paginatedWebmasters []repository.Webmaster

	for _, k := range keys {
		if count >= start && count < start+limit {
			wm := c.webmasters[k]
			// Добавляем вложенные плейсменты
			wm.Placements = c.placements[wm.ID]
			paginatedWebmasters = append(paginatedWebmasters, wm)
		}
		count++
		if count >= start+limit {
			break
		}
	}

	return paginatedWebmasters
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

func (c *Cache) UpdateCacheWhenCreatePlacement(p repository.Placement, id int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	p.ID = id

	// Проверяем, есть ли уже размещения для этого UserID
	if _, exists := c.placements[p.UserID]; !exists {
		c.placements[p.UserID] = []repository.Placement{}
	}

	c.placements[p.UserID] = append(c.placements[p.UserID], p)
	c.placementsByID[p.ID] = p
}

func (c *Cache) UpdateCacheWhenUpdatePlacement(p repository.Placement) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Проверяем, существует ли размещение в кеше
	oldPlacement, exists := c.placementsByID[p.ID]
	if !exists {
		return fmt.Errorf("размещение с ID %d не найдено в кеше", p.ID)
	}

	// Обновляем размещение в placements
	userID := oldPlacement.UserID
	for i, placement := range c.placements[userID] {
		if placement.ID == p.ID {
			c.placements[userID][i] = p
			c.placementsByID[p.ID] = p
			return nil
		}
	}

	return fmt.Errorf("размещение с ID %d не найдено в списке вебмастеров %d", p.ID, userID)
}

func (c *Cache) UpdateCacheWhenDeletePlacement(id int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Проверяем, существует ли размещение
	placement, exists := c.placementsByID[id]
	if !exists {
		return fmt.Errorf("размещение с ID %d не найдено в кеше", id)
	}

	// Удаляем из кеша по ID
	delete(c.placementsByID, id)

	// Удаляем из кеша по UserID
	userID := placement.UserID
	placements, exists := c.placements[userID]
	if !exists {
		return fmt.Errorf("не найдено размещений для пользователя %d", userID)
	}

	var updated []repository.Placement
	found := false
	for _, p := range placements {
		if p.ID != id {
			updated = append(updated, p)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("размещение с ID %d не найдено в списке пользователя %d", id, userID)
	}

	c.placements[userID] = updated
	return nil
}

func (c *Cache) UpdateCacheWhenCreateWebmaster(wm repository.Webmaster, id int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Назначаем ID новому веб-мастеру
	wm.ID = id

	// Добавляем в кеш
	c.webmasters[wm.ID] = wm

	// Инициализируем пустой список размещений для нового веб-мастера
	c.placements[wm.ID] = []repository.Placement{}
}

func (c *Cache) UpdateCacheWhenUpdateWebmaster(wm repository.Webmaster) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Проверяем, существует ли веб-мастер в кеше
	if _, exists := c.webmasters[wm.ID]; !exists {
		return fmt.Errorf("веб-мастер с ID %d не найден в кеше", wm.ID)
	}

	// Обновляем веб-мастера в кеше
	c.webmasters[wm.ID] = wm
	return nil
}

func (c *Cache) UpdateCacheWhenDeleteWebmaster(id int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Проверяем, существует ли веб-мастер в кеше
	if _, exists := c.webmasters[id]; !exists {
		return fmt.Errorf("веб-мастер с ID %d не найден в кеше", id)
	}

	// Удаляем веб-мастера и его размещения из кеша
	delete(c.webmasters, id)
	delete(c.placements, id)

	return nil
}
