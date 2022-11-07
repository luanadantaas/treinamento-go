package database

import (
	"challenge/internal/entity"
	"challenge/internal/logger"
	"database/sql"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var u = &entity.Task{
	Name: "Luana Dantas",
	Completed: "no",
	ID:    2,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Log().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

//test list task
func TestListTask(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer func() {
		repo.Close()
	}()

	query := "SELECT ID, name, completed FROM task"

	rows := sqlmock.NewRows([]string{"ID", "name", "completed"}).
	AddRow(u.ID, u.Name, u.Completed)

	mock.ExpectQuery(query).WillReturnRows(rows)

	_, err := repo.ListTask()
	if err != nil {
		t.Errorf("not expected error here: %v", err)
	}

	if err:= mock.ExpectationsWereMet();err != nil{
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

//teste get task
func TestGetTask(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer func() {
		repo.Close()
	}()

	query := "SELECT ID, name, completed FROM task WHERE ID = ?"

	rows := sqlmock.NewRows([]string{"ID", "name", "completed"}).
		AddRow(u.ID, u.Name, u.Completed)

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	user, err := repo.GetTask(u.ID)
	if err != nil {
		t.Errorf("not expected error here: %v", err)
	}
	if !reflect.DeepEqual(user, u) {
		t.Errorf("not expected entity to be nil")
	}
	if err:= mock.ExpectationsWereMet();err != nil{
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

//teste get task dando erro

func TestNewTask(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer func() {
		repo.Close()
	}()

	query := "INSERT INTO task (name, completed) VALUES ( ?, 'no')"

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(u.Name).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err := repo.NewTask(*u)
	if err != nil {
		t.Errorf("not expected error here: %v", err)
	}
	if err:= mock.ExpectationsWereMet();err != nil{
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestUpdateTask(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer func() {
		repo.Close()
	}()

	query := "UPDATE task SET completed = 'yes' WHERE ID = ?"

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateTask(u.ID)
	if err != nil {
		t.Errorf("not expected error here: %v", err)
	}

	if err:= mock.ExpectationsWereMet();err != nil{
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestListComp(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}
	defer func() {
		repo.Close()
	}()

	query := "SELECT ID, name, completed FROM task WHERE completed = ?"

	rows := sqlmock.NewRows([]string{"ID", "name", "completed"}).
		AddRow(u.ID, u.Name, u.Completed)

	mock.ExpectQuery(query).WithArgs(u.Completed).WillReturnRows(rows)

	_, err := repo.ListComp(u.Completed)

	if err != nil {
		t.Errorf("not expected error here: %v", err)
	}
	if err:= mock.ExpectationsWereMet();err != nil{
		t.Errorf("unfulfilled expectations: %v", err)
	}
	
}
