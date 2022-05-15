package main

import (
	"database/sql"
)

// Up is executed when this migration is applied
func Up_20220506223437(txn *sql.Tx) {
	txn.Exec(
		`CREATE TABLE rooms (
    id         INT          NOT NULL AUTO_INCREMENT,
    name       VARCHAR(255) NOT NULL,
		created_at DATETIME     NOT NULL,
		updated_at DATETIME     NOT NULL,
    PRIMARY KEY(id)
		);`,
	)
	txn.Exec(
		`CREATE TABLE room_users (
    id         INT      NOT NULL AUTO_INCREMENT,
		user_id    INT      NOT NULL,
		room_id    INT      NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    PRIMARY KEY(id)
		);`,
	)
}

// Down is executed when this migration is rolled back
func Down_20220506223437(txn *sql.Tx) {
	txn.Exec(`DROP TABLE room_users;`)
	txn.Exec(`DROP TABLE rooms;`)
}
