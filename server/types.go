package srv

import "go-gin-boilerplate/db"

type Result struct {
	Code    int         `json:"code"`
	Results interface{} `json:"results,omitempty"`
	Err     error       `json:"message,omitempty"`
}

type Job func(*db.Core) *Result
