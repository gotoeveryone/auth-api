package entity

// RegistrationUser is struct of request data for registration user
type RegistrationUser struct {
	Account     string  `json:"account" binding:"required,min=8,max=20"`
	Name        string  `json:"name" binding:"required,max=50"`
	Gender      string  `json:"gender" binding:"required,oneof=Male Female Unknown"`
	MailAddress string  `json:"mailAddress" binding:"required,email"`
	Birthday    string  `json:"birthday" binding:"required,date"`
	Role        *string `json:"role" binding:"omitempty,oneof=Administrator General"`
}

// Activate is validation struct of using during activate user
type Activate struct {
	Authenticate
	NewPassword string `json:"newPassword" binding:"required,password"`
}

// Authenticate is validation struct of using during authentication
type Authenticate struct {
	Account  string `json:"account" binding:"required,min=8,max=20"`
	Password string `json:"password" binding:"required,password"`
}
