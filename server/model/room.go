package model

import "time"

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
