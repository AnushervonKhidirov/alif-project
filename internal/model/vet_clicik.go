package model

import "time"

type VetClinic struct {
	Id        int
	Name      string
	About     *string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
