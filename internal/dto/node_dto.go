package dto

type NodeDTO struct {
	ID     uint   `json:"id" validate:"required" form:"id" query:"id"`
	Name   string `json:"name" validate:"required" form:"name" query:"name"`
	Intro  string `json:"intro" validate:"max=100" form:"intro" query:"intro"`
	Avatar string `json:"avatar" form:"avatar" query:"avatar"`
}
