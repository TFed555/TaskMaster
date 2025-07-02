package models

type Priority struct {
	Low		string    `db:"low"`
	Medium 	string  `db:"medium"`
	High 	string    `db:"high"`
}

type Todo struct {
  ID        int     `db:"id"`
  Title 		string		`db:"title"`
  Priority 		string	`db:"priority"`
  Category 		string		`db:"category"`
  Description 	string		`db:"description"`
  CreatedAt 	string		`db:"createdat"`
  CompletedAt 	string		`db:"completedat"`
  UserId        int       `db:"userid"`
}