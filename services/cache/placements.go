package services

import (
	"rest-crud/repository"
	"sort"
)

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
