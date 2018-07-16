package entity

// Activate is validation struct of using during activate user
type Activate struct {
	Authenticate
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// Authenticate is validation struct of using during authentication
type Authenticate struct {
	Account  string `json:"account" binding:"required,min=6,max=10"`
	Password string `json:"password" binding:"required,min=8"`
}
