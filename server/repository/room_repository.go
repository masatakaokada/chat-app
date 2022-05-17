package repository

import (
	"app/model"
	"database/sql"
	"time"
)

func RoomList(currentUserId int) ([]*model.Room, error) {
	query := `SELECT rooms.id, rooms.name FROM rooms INNER JOIN room_users ON rooms.id = room_users.room_id WHERE room_users.user_id = ?;`

	var rooms []*model.Room
	if err := db.Select(&rooms, query, currentUserId); err != nil {
		return nil, err
	}

	return rooms, nil
}

func RoomCreate(name string, current_user *model.User, userIds []int) (sql.Result, error) {
	now := time.Now()
	room := &model.Room{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	query := `INSERT INTO rooms (name, created_at, updated_at)
  VALUES (:name, :created_at, :updated_at);`

	roomUserIds := append(userIds, current_user.ID)

	room_user_query := `INSERT INTO room_users (user_id, room_id, created_at, updated_at)
  VALUES (:user_id, :room_id, :created_at, :updated_at);`

	// トランザクションを開始
	tx := db.MustBegin()

	res, err := tx.NamedExec(query, room)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	roomId, _ := res.LastInsertId()
	castedValue := int(roomId)

	for _, userId := range roomUserIds {
		roomUser := &model.RoomUser{
			UserId:    userId,
			RoomId:    castedValue,
			CreatedAt: now,
			UpdatedAt: now,
		}
		_, err := tx.NamedExec(room_user_query, roomUser)
		if err != nil {
			tx.Rollback()

			return nil, err
		}
	}

	tx.Commit()

	return res, nil
}
