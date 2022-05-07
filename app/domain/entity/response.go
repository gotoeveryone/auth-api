package entity

// Error is struct of error object
type Error struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Error   error       `json:"-"`
}

// State is struct of Application state
type State struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	LogLevel    string `json:"logLevel"`
	TimeZone    string `json:"timezone"`
}

// GeneratedPassword is struct of generated password
type GeneratedPassword struct {
	Password string `json:"password"`
}

// Claim is struct of logged in user claim data
type Claim struct {
	Expire string `json:"expire"`
	Token  string `json:"token"`
}
