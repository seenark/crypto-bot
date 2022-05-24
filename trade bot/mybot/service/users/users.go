package usersService

import (
	"mybot/line"
	"mybot/repository/users"
)

type UserResponse struct {
	Id              string          `json:"id" bson:"_id"`
	Name            string          `json:"name" bson:"name,omitempty"` // omitempty คือถ้าไม่มีค่าก็ไม่ต้องใส่ key เข้าไปใน db
	WatchPairs      []string        `json:"watchPair" bson:"watchPair"`
	Equity          float64         `json:"equity" bson:"equity,omitempty"`
	Position        *users.Position `json:"position" bson:"position"`
	Profit          float64         `json:"profit" bson:"profit"`
	Leverage        float64         `json:"Leverage" bson:"Leverage"`
	UseAssetPercent float64         `json:"useAssetPercent" bson:"useAssetPercent"`
	Line            *line.Line      `json:"line" bson:"line"`
}

type UserService interface {
	Create(users.User) (string, error)
	GetAll() ([]UserResponse, error)
	GetUserById(string) (*UserResponse, error)
	UpdateUserById(string, users.User) (*UserResponse, error)
	DeleteUserById(string) (int, error)
}
