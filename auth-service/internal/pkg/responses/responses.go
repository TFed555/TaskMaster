package responses

type RegRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateRequest struct {
	Email	string	`json:"email,omitempty"`
	Name	string 	`json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	ImgPath	 string `json:"img_path,omitempty"`	
}

type RegResponse struct {
	UserID int `json:"id"`
}

type AuthResponse struct {
	UserID uint   `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	ImgURL	string	`json:"img"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Status  int   `json:"code"`
	Message string `json:"message"`
}