package auth

type LoginRequest struct {
	Username string `json:"username" validate:"required,notblank"`
	Password string `json:"password" validate:"required,notblank,min=6"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,notblank,min=2,max=15"`
	Password string `json:"password" validate:"required,notblank,min=6"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
