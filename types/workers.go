package types

import "time"

// Workers' Status Data

type Leader struct {
	ID                  int       `json:"id"`
	APICallCountSuccess int       `json:"apiCallCountSuccess"`
	APICallCountFail    int       `json:"apiCallCountFail"`
	LastAPICallAt       time.Time `json:"lastAPICallTime"`
	Active              bool      `json:"active"`
	ExitMsg             string    `json:"exitMsg"`
}

type Worker struct {
	ID                  int       `json:"id"`
	APICallCountSuccess int       `json:"apiCallCountSuccess"`
	APICallCountFail    int       `json:"apiCallCountFail"`
	LastAPICallAt       time.Time `json:"lastAPICallTime"`
	Active              bool      `json:"active"`
	ExitMsg             string    `json:"exitMsg"`
}

type TaskSender struct {
	CurrentBatch int    `json:"clock"`
	Active       bool   `json:"active"`
	ExitMsg      string `json:"exitMsg"`
}
