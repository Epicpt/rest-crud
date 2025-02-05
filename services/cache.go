package services

import (
	"rest-crud/repository"
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

// GetPlacements - получает список размещений с пагинацией
func (c *Cache) GetPlacements(page, limit int) ([]repository.Placement, int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var allPlacements []repository.Placement
	for _, pls := range c.placements {
		allPlacements = append(allPlacements, pls...) // Собираем все размещения в один список
	}

	total := len(allPlacements)
	start := (page - 1) * limit
	if start > total {
		return []repository.Placement{}, total
	}

	end := start + limit
	if end > total {
		end = total
	}

	return allPlacements[start:end], total
}

// GetWebmasters - получает список вебмастеров с вложенными размещениями и пагинацией
func (c *Cache) GetWebmasters(page, limit int) ([]WebmasterWithPlacements, int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var allWebmasters []WebmasterWithPlacements
	for _, wm := range c.webmasters {
		// Создаем новую структуру для вебмастера с размещениями
		wmWithPlacements := WebmasterWithPlacements{
			Webmaster:  wm,
			Placements: c.placements[wm.ID],
		}
		allWebmasters = append(allWebmasters, wmWithPlacements)
	}

	total := len(allWebmasters)
	start := (page - 1) * limit
	if start > total {
		return []WebmasterWithPlacements{}, total
	}

	end := start + limit
	if end > total {
		end = total
	}

	return allWebmasters[start:end], total
}

func (c *Cache) LoadCacheFromDB(repo *repository.Repository) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	webmasters, err := repo.GetAllWebmasters() // Функция для загрузки всех веб-мастеров из БД
	if err != nil {
		return err
	}
	for _, wm := range webmasters {
		c.webmasters[wm.ID] = wm
	}

	placements, err := repo.GetAllPlacements() // Функция для загрузки всех размещений
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
	c.placements[p.UserID] = append(c.placements[p.UserID], p)
	c.placementsByID[p.ID] = p
}

func (c *Cache) UpdateCacheWhenUpdatePlacement(p repository.Placement) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Обновляем размещение в placements
	for i, placement := range c.placements[c.placementsByID[p.ID].UserID] {
		if placement.ID == p.ID {
			c.placements[c.placementsByID[p.ID].UserID][i] = p
			break
		}
	}
	// Обновляем размещение в placementsByID
	c.placementsByID[p.ID] = p
}
