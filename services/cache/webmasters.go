package services

import (
	"rest-crud/repository"
	"sort"
)

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
