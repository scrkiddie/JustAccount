package model

type UserResponse struct {
	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	Username       string `json:"username,omitempty"`
	Email          string `json:"email,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty"`
	Token          string `json:"token,omitempty"`
}

type RegisterUserRequest struct {
	FirstName string `json:"firstName" mod:"normalize_spaces" validate:"required,max=50"`
	LastName  string `json:"lastName,omitempty" mod:"normalize_spaces" validate:"omitempty,max=50"`
	Username  string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Email     string `json:"email" mod:"lcase" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,min=8,max=255,is_password_format,not_contain_space"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password string `json:"password" validate:"required,min=8,max=255,is_password_format,not_contain_space"`
}

type VerifyUserRequest struct {
	Token string `validate:"required"`
}

type GetUserRequest struct {
	ID int `json:"id" validate:"required,numeric"`
}

type UpdateUserRequest struct {
	ID        int    `validate:"required,numeric"`
	FirstName string `json:"firstName" mod:"normalize_spaces" validate:"required,max=50"`
	LastName  string `json:"lastName,omitempty" mod:"normalize_spaces" validate:"omitempty,max=50"`
	Email     string `json:"email" mod:"lcase" validate:"required,email,max=100"`
}

type UpdatePasswordRequest struct {
	ID              int    `validate:"required,numeric"`
	CurrentPassword string `json:"currentPassword" validate:"required,min=8,max=255,is_password_format,not_contain_space"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=255,is_password_format,not_contain_space"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=255,eqfield=NewPassword"`
}
