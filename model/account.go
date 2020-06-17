package model

import "time"

type Account struct {
	ID         int64     `json:"id,omitempty"`
	Username   string    `json:"username,omitempty"`
	Password   string    `json:"password,omitempty"`
	Email      string    `json:"email,omitempty"`
	Name       string    `json:"name,omitempty"`
	Birthday   time.Time `json:"birthday,omitempty"`
	Gender     int       `json:"gender,omitempty"`
	Address    string    `json:"address,omitempty"`
	Hometown   string    `json:"hometown,omitempty"`
	University string    `json:"university,omitempty"`
	StartDate  time.Time `json:"start_date,omitempty"`
	RefUser    string    `json:"ref_user,omitempty"`
	IsAdmin    string    `json:"is_admin,omitempty"`
	Avatar     string    `json:"avatar,omitempty"`
}

// type Accounts struct {
// 	Accounts []Account `json:"employees"`
// }
