package entity

type Task struct {
	ID       int
	Username string
	Title    string
	Complete bool `db:"is_completed"`
}
