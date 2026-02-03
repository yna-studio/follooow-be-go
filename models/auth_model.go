package models

type LoginRequest struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type LoginResponse struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Token     string `json:"token,omitempty"`
	Message   string `json:"message"`
}
