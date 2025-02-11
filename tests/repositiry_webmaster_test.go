package tests

import (
	"fmt"
	"rest-crud/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestCreateWebMaster_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Преобразуем *sql.DB в *sqlx.DB
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := repository.NewRepository(sqlxDB)

	wm := &repository.Webmaster{
		Name:     "John",
		LastName: "Doe",
		Email:    "john.doe@example.com",
		Status:   "active",
	}

	expectedID := 1

	mock.ExpectQuery(`INSERT INTO webmasters`).
		WithArgs(wm.Name, wm.LastName, wm.Email, wm.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	id, err := repo.CreateWebMaster(wm)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if id != expectedID {
		t.Errorf("expected id %d, got %d", expectedID, id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteWebmaster_ErrorOnZeroID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := repository.NewRepository(sqlxDB)

	zeroID := 0

	mock.ExpectExec(`DELETE FROM webmasters WHERE id = \$1`).
		WithArgs(zeroID).
		WillReturnError(fmt.Errorf("invalid ID"))

	err = repo.DeleteWebmaster(zeroID)

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
