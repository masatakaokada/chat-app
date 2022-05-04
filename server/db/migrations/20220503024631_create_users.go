package main

import (
	"database/sql"
)

// Up is executed when this migration is applied
func Up_20220503024631(txn *sql.Tx) {
	txn.Exec(
		`CREATE TABLE users (
    id int AUTO_INCREMENT,
    name varchar(100),
		email varchar(100),
		firebase_uid varchar(100),
		created datetime,
		updated datetime,
    PRIMARY KEY(id)
		);`,
	)
}

// Down is executed when this migration is rolled back
func Down_20220503024631(txn *sql.Tx) {
	txn.Exec(`DROP TABLE users;`)
}
