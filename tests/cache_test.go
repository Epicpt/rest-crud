package tests

import (
	"rest-crud/repository"
	services "rest-crud/services/cache"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestLoadCacheFromDB_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewRepository(sqlxDB)

	cache := services.NewCache()

	webmasters := []repository.Webmaster{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	placements := []repository.Placement{
		{ID: 1, UserID: 1, Name: "Placement 1", Description: "Description 1"},
		{ID: 2, UserID: 2, Name: "Placement 2", Description: "Description 2"},
	}

	mock.ExpectQuery(`SELECT id, name, last_name, email, status FROM webmasters`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(webmasters[0].ID, webmasters[0].Name).
			AddRow(webmasters[1].ID, webmasters[1].Name))

	mock.ExpectQuery(`SELECT id, user_id, name, description FROM placements`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description"}).
			AddRow(placements[0].ID, placements[0].UserID, placements[0].Name, placements[0].Description).
			AddRow(placements[1].ID, placements[1].UserID, placements[1].Name, placements[1].Description))

	err = cache.LoadCacheFromDB(repo)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	cacheWebmasters := cache.GetCacheWebmasters()
	cachePlacementsByID := cache.GetCachePlacementsByID()

	if len(cacheWebmasters) != len(webmasters) {
		t.Errorf("expected %d webmasters, got %d", len(webmasters), len(cacheWebmasters))
	}

	if len(cachePlacementsByID) != len(placements) {
		t.Errorf("expected %d placements, got %d", len(placements), len(cachePlacementsByID))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
