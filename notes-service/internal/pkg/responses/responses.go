package responses


type TodoResponse struct {
	Todos *[]OneTodoResponse `json:"todos"`
}

type CreateRequest struct {
	Title       string  `json:"title"`
	Priority    string  `json:"priority"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"createdAt"`
	CompletedAt *string `json:"completedAt,omitempty"`
	UserID      *int    `json:"userID"`
}

type ArchiveRequest struct {
	ID int `json:"id"`
}

type CreateResponse struct {
	ID int `json:"id"`
}

type UpdateRequest struct {
	Title       string  `json:"title,omitempty"`
	Priority    string  `json:"priority,omitempty"`
	Category    string  `json:"category,omitempty"`
	Description *string  `json:"description,omitempty"`
	CreatedAt   string  `json:"createdAt,omitempty"`
	CompletedAt string `json:"completedAt,omitempty"`
	ID          int    `json:"id"`
}

type OneTodoResponse struct {
	Title       string       `json:"title,omitempty"`
	Priority    string       `json:"priority,omitempty"`
	Category    string       `json:"category,omitempty"`
	Description string       `json:"description,omitempty"`
	CreatedAt   string       `json:"createdAt,omitempty"`
	CompletedAt *string      `json:"completedAt,omitempty"`
	ID          int          `json:"id"`
	Tags        []OneTagResponse `json:"tags,omitempty"`
}

type ErrorResponse struct {
	Status  uint   `json:"code"`
	Message string `json:"message"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

type CreateTagRequest struct {
	Name string `json:"name"`
}

type TagResponse struct {
	Tags *[]OneTagResponse
}

type OneTagResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UpdateTagRequest struct {
	Name string `json:"name"`
}

type AddTagRequest struct {
	TagId int `json:"tagId"`
}

type OneAuditTrial struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type AuditTrial struct {
	Changes *[]OneAuditTrial `json:"changes"`
}

type SearchRequest struct {
	SearchString string `json:"search"`
	UserID       uint
}

type PlanRequest struct {
	TodoID int    `json:"todoid"`
	Steps  []string `json:"steps"`
	UserID	int 	`json:"userid"`
}

type CountResponse struct {
	Count int `json:"count"`
}