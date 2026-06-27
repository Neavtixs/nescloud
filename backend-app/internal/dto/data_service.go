package dto

import "context"

type InputAuthRegister struct {
	Ctx      context.Context
	Name     string
	Email    string
	Password string
}

type InputAuthLogin struct {
	Ctx      context.Context
	Email    string
	Password string
}

type InputAuthRefresh struct {
	Ctx          context.Context
	RefreshToken string
}

type ResultAuthLogin struct {
	AccessToken  string
	RefreshToken string
}

type ResultAuthMe struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type ResultAuthRegister struct {
	ID           string
	AccessToken  string
	RefreshToken string
}

type InputAuthLogout struct {
	Ctx          context.Context
	RefreshToken string
}

type InputAuthMe struct {
	Ctx    context.Context
	UserID string
}

type ResultAuthRefresh struct {
	AccessToken  string
	RefreshToken string
}
