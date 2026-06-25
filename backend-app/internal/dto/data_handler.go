package dto

type ResponseWeb[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type AuthRegisterReq struct {
	Name     string `json:"name"     validate:"required,min=1,max=100"`
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthLoginReq struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthRegisterRes struct {
	AccessToken string `json:"access_token"`
}

type AuthLoginRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type AuthMeRes struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type ErrorWeb struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}
