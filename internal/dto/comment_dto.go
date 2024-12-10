package dto

type CommentDTO struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	PostID uint   `json:"post_id"`
	Text   string `json:"text"`
}
