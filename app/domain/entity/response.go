package entity

// Error is struct of error object
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   error  `json:"-"`
}

// State is struct of Application state
type State struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	LogLevel    string `json:"logLevel"`
	TimeZone    string `json:"timeZone"`
}

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
