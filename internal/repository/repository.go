package repository

import "challenge/internal/entity"

type Repository interface {
	ListTask() ([]entity.Task, error)
	GetTask(id int) (*entity.Task, error)
	NewTask(entity.Task) (int64, error)
	UpdateTask(n int) error
	ListComp(yn string) ([]entity.Task, error)
}
