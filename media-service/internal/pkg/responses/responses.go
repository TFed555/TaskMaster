package responses


type GenerateResponse struct {
	Name string `json:"name"`
	Link string`json:"link"`
}

type SaveRequest struct {
	Name string	`json:"name"`
}

type ErrorResponse struct {
	Status  int   `json:"code"`
	Message string `json:"message"`
}