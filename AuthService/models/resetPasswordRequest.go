package models

type ResetPasswordRequest struct {
	Email     string `json:"email"`
	Token  string `json:"token"`
	NewPassword  string `json:"new_password"`
}

  