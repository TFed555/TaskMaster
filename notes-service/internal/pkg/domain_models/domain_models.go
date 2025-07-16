package domain_models

type Todo struct {
  ID			*int
  Title 		string
  Priority 		string
  Category 		string
  Description 	string
  CreatedAt 	string
  CompletedAt *string
  UserId        *int
}

type Tag struct {
  ID     *int
  UserID uint
  Name string
}

type TagTodo struct {
  TodoID int
  TagID int
}