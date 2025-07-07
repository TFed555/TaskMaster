package domain_models

type User struct {
    ID       uint
    Email    string
    Name     string
    Password string
}

type RegisterParams struct {
	Name	string
	Email	string
	Password string
}

type AuthorizeParams struct {
	Email	string
	Password string
}

type UpdateUserParams struct {
	ID        uint
	Name      *string
	Password  *string
	Email     *string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
