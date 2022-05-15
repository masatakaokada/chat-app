package repository

import (
	"app/model"
)

func RoomUserGetByRoomId(firebase_uid string) (*model.User, error) {
	query := `SELECT *
	FROM room_users
	WHERE firebase_uid = ?;`

	var user model.User

	if err := db.Get(&user, query, firebase_uid); err != nil {
		return nil, err
	}

	return &user, nil
}
