package input

type SignUp struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8"`
	PhoneNumber     string `json:"phone_number"`
	Photo           string `json:"photo,omitempty"`
}
