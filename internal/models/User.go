package models

import (
	"time"
)

type User struct {
	Model
	Name        string  `gorm:"type:varchar(100);not null"`
	Email       string  `gorm:"type:varchar(100);uniqueIndex:idx_email;not null"`
	Password    string  `gorm:"type:varchar(100);not null"`
	Role        *string `gorm:"type:varchar(50);default:'user';not null"`
	Photo       *string `gorm:"not null;default:'default.png'"`
	PhoneNumber string  `gorm:"type:varchar(10)"`
	Verified    *bool   `gorm:"not null;default:false"`
}

type UserResponse struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      *user.Role,
		Photo:     *user.Photo,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}
}
