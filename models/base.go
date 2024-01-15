package models

import "time"

type Base struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateBaseObject() *Base {
	currentTime := time.Now().UTC()
	return &Base{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}
func UpdatedBaseObject() *Base {
	currentTime := time.Now().UTC()
	return &Base{
		UpdatedAt: currentTime,
	}
}
