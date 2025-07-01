package models

type Priority struct {
	low		string
	medium 	string
	high 	string
}

type Todo struct {
  Title 		string		`db:"title"`
  Priority 		Priority	`db:"priority"`
  Category 		string		`db:"category"`
  Description 	string		`db:"description"`
  CreatedAt 	string		`db:"createdAt"`
  CompletedAt 	string		`db:"completedAt"`
  UserId        int       `db:"userId"`
}