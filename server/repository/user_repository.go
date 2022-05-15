package repository

import (
	"app/model"
	"database/sql"
	"time"
)

func RoomCreationUsers(user_id int) ([]*model.User, error) {
	query := `SELECT * FROM users WHERE NOT id = ?;`

	var users []*model.User
	if err := db.Select(&users, query, user_id); err != nil {
		return nil, err
	}

	return users, nil
}

func UserCreate(user *model.User) (sql.Result, error) {
	now := time.Now()
	user.Created = now
	user.Updated = now

	query := `INSERT INTO users (name, email, firebase_uid, created, updated)
  VALUES (:name, :email, :firebase_uid, :created, :updated);`

	// トランザクションを開始
	tx := db.MustBegin()

	res, err := tx.NamedExec(query, user)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	tx.Commit()

	return res, nil
}

func UserGetByFirebaseUid(firebase_uid string) (*model.User, error) {
	query := `SELECT *
	FROM users
	WHERE firebase_uid = ?;`

	var user model.User

	if err := db.Get(&user, query, firebase_uid); err != nil {
		return nil, err
	}

	return &user, nil
}
