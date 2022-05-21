package model

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type Room struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type RoomUser struct {
	ID        int       `db:"id"`
	UserId    int       `db:"user_id"`
	RoomId    int       `db:"room_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type RoomCreate struct {
	Name    string `json:"name" validate:"required,max=20"`
	UserIds []int  `json:"userIds" validate:"required,gt=1"`
}

func (r *RoomCreate) ValidationErrors(err error) []string {
	var errMessages []string

	for _, err := range err.(validator.ValidationErrors) {
		var message string

		switch err.Field() {
		case "Name":
			switch err.Tag() {
			case "required":
				message = "ルーム名は必須です。"
			case "max":
				message = "ルーム名は最大50文字です。"
			}
		case "UserIds":
			switch err.Tag() {
			case "required":
				message = "ユーザーは必須です。"
			case "gt":
				message = "ユーザーは１人以上選択してください。"
			}
		}

		if message != "" {
			errMessages = append(errMessages, message)
		}
	}

	return errMessages
}
