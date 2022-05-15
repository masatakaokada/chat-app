package repository

import (
	"app/model"
)

func RoomUserGetByUserIdAndRoomId(userId int, roomId int) (*model.RoomUser, error) {
	query := `SELECT *
	FROM room_users
	WHERE user_id = ?
	AND room_id = ?;`

	var roomUser model.RoomUser

	if err := db.Get(&roomUser, query, userId, roomId); err != nil {
		return nil, err
	}

	return &roomUser, nil
}
