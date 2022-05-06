package model

import "time"

type User struct {
	ID          int       `db:"id"`
	Name        string    `db:"name" json:"name"`
	Email       string    `db:"email"`
	FirebaseUid string    `db:"firebase_uid"`
	Created     time.Time `db:"created"`
	Updated     time.Time `db:"updated"`
}
