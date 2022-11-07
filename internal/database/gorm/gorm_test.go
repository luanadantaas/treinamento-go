package database

import (
	"challenge/internal/entity"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var u = &entity.Task{
	Name:      "Luana Dantas",
	Completed: "no",
	ID:        2,
}

func NewMock(t *testing.T) (sqlmock.Sqlmock, *GormRepo) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repo := &GormRepo{db: gdb}

	return mock, repo
}

func TestListTaskGorm(t *testing.T) {
	mock, repo := NewMock(t)

	query := "SELECT * FROM `task`"

	rows := sqlmock.NewRows([]string{"ID", "name", "completed"}).
		AddRow(u.ID, u.Name, u.Completed)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	_, err := repo.ListTask()
	if err != nil {
		t.Errorf("not expected error here: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetTaskGorm(t *testing.T) {
	mock, repo := NewMock(t)

	query := "SELECT * FROM `task` WHERE `task`.`id` = ? ORDER BY `task`.`id` LIMIT 1"

	rows := sqlmock.NewRows([]string{"ID", "name", "completed"}).
		AddRow(u.ID, u.Name, u.Completed)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.ID).WillReturnRows(rows)

	user, err := repo.GetTask(u.ID)
	if err != nil {
		t.Errorf("not expected error here: %v", err)
	}
	if !reflect.DeepEqual(user, u) {
		t.Errorf("not expected entity to be nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestNewTaskGorm(t *testing.T) {
	mock, repo := NewMock(t)

	query := "INSERT INTO `task` (`name`,`completed`,`id`) VALUES (?,?,?)"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(u.Name, u.Completed, u.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	_, err := repo.NewTask(*u)
	if err != nil {
		t.Errorf("not expected error here: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestUpdateTaskGorm(t *testing.T) {
	mock, repo := NewMock(t)

	query := "UPDATE `task` SET `completed`=? WHERE ID = ?"

	yes := "yes"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(yes, u.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.UpdateTask(u.ID)
	if err != nil {
		t.Errorf("not expected error here: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestListCompGorm(t *testing.T) {
	mock, repo := NewMock(t)
	query := "SELECT * FROM `task` WHERE completed = ?"

	rows := sqlmock.NewRows([]string{"ID", "name", "completed"}).
		AddRow(u.ID, u.Name, u.Completed)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Completed).WillReturnRows(rows)

	_, err := repo.ListComp(u.Completed)

	if err != nil {
		t.Errorf("not expected error here: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
