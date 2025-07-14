package entity

type Task struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Text     string `json:"text"`
	Complete bool   `json:"complete"`
}
