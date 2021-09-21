package http

type AuthRequest struct {
	Username string `json:"username" validate:"required,gte=4,lt=30"`
	Password string `json:"password" validate:"required,gte=4,lt=40"`
}
