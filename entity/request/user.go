package entity

type RegisterUserRequest struct {
	UserName string `json:"userName" binding:"required,max=10"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Avatar   string `json:"avatar"`
	Code     string `json:"code" binding:"required"`
}

type RegisterUserResponse struct {
	ID             uint64 `json:"id"`
	UserName       string `json:"userName"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
	CreatedAt      int64  `json:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt"`
	Token          string `json:"token"`
	TokenExpiresAt int64  `json:"tokenExpiresAt"`
	LastActiveAt   int64  `json:"lastActiveAt"`
}

type LoginUserRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}
