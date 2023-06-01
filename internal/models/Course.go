package models

type Course struct {
	Model
	UserId      int    `json:"user_id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}
