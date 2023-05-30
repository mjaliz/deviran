package models

type User struct {
	Model
	Email        string `json:"email" validate:"required" gorm:"unique"`
	Password     string `json:"password" validate:"required"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
