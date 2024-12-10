package dto

type UserCreateDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Nickname string `json:"nickname" validate:"omitempty"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type UserUpdateDTO struct {
	UID      uint   `json:"uid" validate:"required"`
	Email    string `json:"email" validate:"omitempty,email"`
	Nickname string `json:"nickname" validate:"omitempty"`
	Password string `json:"password" validate:"omitempty"`
}
