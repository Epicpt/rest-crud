package repository

type Webmaster struct {
	ID       int    `db:"id"`
	Name     string `db:"name" json:"name"`
	LastName string `db:"last_name" json:"last_name"`
	Email    string `db:"email" json:"email"`
	Status   string `db:"status" json:"status"`
}

// Добавление веб-мастера в БД
func (r *Repository) CreateWebMaster(wm *Webmaster) (int, error) {
	var id int
	query := `INSERT INTO webmasters (name, last_name, email, status) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, wm.Name, wm.LastName, wm.Email, wm.Status).Scan(&id)
	return id, err
}

func (r *Repository) UpdateWebmaster(wm Webmaster) error {
	query := `UPDATE webmasters SET name = $1, last_name = $2, email = $3, status = $4 WHERE id = $5`
	_, err := r.db.Exec(query, wm.Name, wm.LastName, wm.Email, wm.Status, wm.ID)
	return err
}

func (r *Repository) DeleteWebmaster(id int) error {
	query := `DELETE FROM webmasters WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository) GetAllWebmasters() ([]Webmaster, error) {
	query := "SELECT id, name, last_name, email, status FROM webmasters"
	var webmasters []Webmaster
	err := r.db.Select(&webmasters, query)
	return webmasters, err
}
