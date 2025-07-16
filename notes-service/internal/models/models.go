package models

type Priority struct {
	Low    string `db:"low"`
	Medium string `db:"medium"`
	High   string `db:"high"`
}

type Todo struct {
	ID          *int       `db:"id"`
	Title       string     `db:"title"`
	Priority    string     `db:"priority"`
	Category    string     `db:"category"`
	Description *string     `db:"description"`
	CreatedAt   string     `db:"createdat"`
	CompletedAt *string    `db:"completedat"` //т.к. может быть null
	UserId      *int       `db:"userid"`
	Tags        []Tag
}

type Tag struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	UserID *int    `db:"userid"`
}

type HistoryTodo struct {
  ID	int		`db:"id"`
  Title string `db:"title"`
  Status string `db:"action"`
}