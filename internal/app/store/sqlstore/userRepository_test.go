package sqlstore_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rysya2000/http-rest-api/internal/app/model"
	"github.com/rysya2000/http-rest-api/internal/app/store/sqlstore"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func Test_Create(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	u := &model.User{
		Email:             "user@example.org",
		EncryptedPassword: "encryptedpassword",
	}
	s := sqlstore.New(db)

	rs := mock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO users").WithArgs(u.Email, u.EncryptedPassword).WillReturnRows(rs)

	if err := s.User().Create(u); err != nil {
		t.Errorf("error occured at creating:\n %s", err)
	}
}

func Test_Find(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	s := sqlstore.New(db)

	rs := mock.NewRows([]string{"id", "email", "encryptedpassword"}).
		AddRow(1, "user@example.org", "encryptedpassword")
	mock.ExpectQuery("SELECT (.+) FROM users").WithArgs(1).WillReturnRows(rs)

	if _, err := s.User().Find(1); err != nil {
		t.Errorf("error occured at selecting:\n %s", err)
	}
}

func Test_FindByEmail(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	s := sqlstore.New(db)

	rs := mock.NewRows([]string{"id", "email", "encryptedpassword"}).
		AddRow(1, "user@example.org", "encryptedpassword")

	mock.ExpectQuery("SELECT (.+) FROM users").WithArgs("user@example.org").WillReturnRows(rs)

	if _, err := s.User().FindByEmail("user@example.org"); err != nil {
		t.Errorf("error occured at selecting:\n %s", err)
	}
}
