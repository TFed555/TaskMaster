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
  UserID uint
  Name string
}