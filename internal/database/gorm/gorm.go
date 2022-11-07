package database

import (
	"challenge/internal/entity"

	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const(
	DbUrl = "DATABASE"
)
type GormRepo struct {
	db *gorm.DB
}

//estabelece conexao com o gorm
func New() (*GormRepo, error) {
	url := os.Getenv(DbUrl)
	dsn := fmt.Sprintf("root:password@tcp(%s:3306)/task", url)
	DB, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	
	return &GormRepo{
		db: DB,
	}, nil
}

// é realizada uma busca por meio da função do gorm, find, para pegar todas as tarefas e adiliciona-las num slice que serve de retorno para a função 
func (g *GormRepo) ListTask() ([]entity.Task, error) {
	var task []entity.Task
	if err := g.db.Find(&task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// é realizada uma busca por meio da função do gorm, first, para pegar a tarefa, de acordo com o id de parametro e retorna tarefa especifica
func (g *GormRepo) GetTask(id int) (*entity.Task, error) {
	task:=  &entity.Task{}
	if err := g.db.First(&task, id).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// a função do gorm, create, cria uma tarefa e adiciona-a ao banco de dados
// a função retorna um inteiro como confirmação
func (g *GormRepo) NewTask(t entity.Task) (int64, error) {
	if err := g.db.Create(&t).Error; err != nil {
		return 0, err
	}

	return 1, nil

}

//é realizada uma busca por meio da função do gorm, model.where.update, para atualizar a tarefa, de acordo com o id de parametro
// e retorna o erro, caso exista
func (g *GormRepo) UpdateTask(id int) error {
	var task []entity.Task
	if err := g.db.Model(&task).Where("ID = ?", id).Update("completed", "yes").Error; err != nil {
		return err
	}

	return nil
}

// é realizada uma busca por meio da função do gorm, find, para pegar todas as tarefas que sao
//identificadas com a coluna completed = parametro da funçao e adiliciona-las num slice que serve de retorno para a função 
func (g *GormRepo) ListComp(yn string) ([]entity.Task, error) {
	var task []entity.Task
	if err := g.db.Where("completed = ?", yn).Find(&task).Error; err != nil {
		return nil, err
	}

	return task, nil

}

