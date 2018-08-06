package entity

import "github.com/gotoeveryone/golib/config"

// Error is struct of error object
type Error struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Error   error       `json:"-"`
}

// State is struct of Application state
type State struct {
	Status      string          `json:"status"`
	Environment string          `json:"environment"`
	LogLevel    config.LogLevel `json:"logLevel"`
	TimeZone    string          `json:"timezone"`
}
