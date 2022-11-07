package database

import (
	"challenge/internal/entity"
	"challenge/internal/logger"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	sqlDB *sql.DB
}

const(
	DbUrl = "DATABASE"	
)

// estabelece conexao com o banco de dados pelo mysql
func New() (*Repository, error) {

	url := os.Getenv(DbUrl)
	path := fmt.Sprintf("root:password@tcp("+url+":3306)/task")
	db, err := sql.Open("mysql", path)

	if err != nil {
		logger.Log().Warn("couldn't open database: %v", err)
		return nil, err
	}

	return &Repository{sqlDB: db}, nil
}

//fecha conexao do mysql
func (d *Repository) Close() {
	d.sqlDB.Close()
}

//seleciona todas as tarefas do banco de dados e adiciona em um slice como retorno da função
func (d *Repository) ListTask() ([]entity.Task, error) {
	rows, err := d.sqlDB.Query("SELECT ID, name, completed FROM task")
	if err != nil {
		logger.Log().Errorf("couldn't list data: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task

	//goes through all the rows and adds the data in the slice
	for rows.Next() {
		var t entity.Task
		if err := rows.Scan(&t.ID, &t.Name,
			&t.Completed); err != nil {
			return tasks, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil

}

//faz uma busca pelo banco de dados de acordo com o id usado como parametro
//retorna a tarefa do determinado id
func (d *Repository) GetTask(id int) (*entity.Task, error) {
	var t *entity.Task
	g := "SELECT ID, name, completed FROM task WHERE ID = ?"
	err := d.sqlDB.QueryRow(g, id).Scan(&t.ID, &t.Name,
		&t.Completed)

	if err != nil {
		logger.Log().Errorf("couldn't get data from task %d: %v", id, err)
		return nil, err
	}

	return t, nil
}

//cria uma nova tarefa e insere valores no banco
//retorna o id da tarefa criada
func (d *Repository) NewTask(t entity.Task) (int64, error) {
	insert := "INSERT INTO task (name, completed) VALUES ( ?, 'no')"
	id, err := d.sqlDB.Exec(insert, t.Name) 
	if err != nil {
		return 0, err
	}

	return id.LastInsertId()
}

//recebe um id de parametro e atualiza a linha na tabela
//de acordo com o id recebido
func (d *Repository) UpdateTask(id int) error {
	update := "UPDATE task SET completed = 'yes' WHERE ID = ?"
	_, err := d.sqlDB.Exec(update, id)
	if err != nil { 
		logger.Log().Errorf("coundn't update task %d: %v", id, err)			
		return err
	}

	return err
}

//pegar todas as tarefas que sao
//identificadas com a coluna completed = parametro da funçao e adiliciona-las num slice que serve de retorno para a função 
func (d *Repository) ListComp(yn string) ([]entity.Task, error) {
	selectTask := "SELECT ID, name, completed FROM task WHERE completed = ?"
	row, err := d.sqlDB.Query(selectTask, yn)
	if err != nil {
		logger.Log().Errorf("couldn't list completed data: %v", err)
		return nil, err
	}
	defer row.Close()

	var tasks []entity.Task

	for row.Next() {
		var t entity.Task
		if err := row.Scan(&t.ID, &t.Name,
			&t.Completed); err != nil {
			return tasks, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
