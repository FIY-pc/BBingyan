package dto

type LoginDTO struct {
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Nickname string `json:"nickname,omitempty" validate:"omitempty"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type RegisterDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Nickname string `json:"nickname,omitempty" validate:"omitempty"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	Captcha  string `json:"captcha" validate:"required"`
}
