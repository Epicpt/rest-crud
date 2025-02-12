package repository

type Placement struct {
	ID          int    `db:"id"`
	UserID      int    `db:"user_id" json:"user_id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

func (p Placement) GetID() int {
	return p.ID
}

// Добавление размещения в БД
func (r *Repository) CreatePlacement(p Placement) (int, error) {
	var id int
	query := `INSERT INTO placements (user_id, name, description) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, p.UserID, p.Name, p.Description).Scan(&id)
	return id, err
}

// Обновление размещения в БД
func (r *Repository) UpdatePlacement(p Placement) error {
	query := `UPDATE placements SET name = $1, description = $2 WHERE id = $3`
	_, err := r.db.Exec(query, p.Name, p.Description, p.ID)
	return err
}

// Удаление размещения из БД
func (r *Repository) DeletePlacement(id int) error {
	query := `DELETE FROM placements WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Получение всех размещений из БД
func (r *Repository) GetAllPlacements() ([]Placement, error) {
	query := "SELECT id, user_id, name, description FROM placements"
	var placements []Placement
	err := r.db.Select(&placements, query)
	return placements, err
}
