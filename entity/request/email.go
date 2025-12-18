package entity

type RequestEmailCode struct {
	Email string `json:"email" binding:"required,email"`
}
