package entity

type Task struct {
	Name      string `json:"name"`
	Completed string `json:"completed"`
	ID        int    `json:"id"`

}

func (*Task) TableName () string {
	return "task"
}