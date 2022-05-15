package main

import (
	"database/sql"
)

// Up is executed when this migration is applied
func Up_20220503024631(txn *sql.Tx) {
	txn.Exec(
		`CREATE TABLE users (
    id           INT          NOT NULL AUTO_INCREMENT,
    name         VARCHAR(255) NOT NULL,
		email        VARCHAR(255) NOT NULL UNIQUE,
		firebase_uid VARCHAR(255) NOT NULL UNIQUE,
		created      DATETIME     NOT NULL,
		updated      DATETIME     NOT NULL,
    PRIMARY KEY(id)
		);`,
	)
}

// Down is executed when this migration is rolled back
func Down_20220503024631(txn *sql.Tx) {
	txn.Exec(`DROP TABLE users;`)
}
