package responses

type SaveRequest struct {
	File string `json:"file"`
}

type SaveResponse struct {
	Link string`json:"link"`
}