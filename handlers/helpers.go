package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

// getPaginationParams – получение параметров пагинации из URL запроса
func getPaginationParams(r *http.Request) (int, int, error) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, 0, fmt.Errorf("invalid page parameter")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return 0, 0, fmt.Errorf("invalid limit parameter")
	}
	return page, limit, nil
}
