package academiq

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitTables() error {
	SetsRequest := `
	CREATE TABLE IF NOT EXISTS Sets (
		SetID INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
		SetName          TEXT    NOT NULL UNIQUE,
		SetImagePath     TEXT    NOT NULL
	);
	`
	_, err := db.Exec(SetsRequest)
	if err != nil {
		return err
	}

	SetDataRequest := `
	CREATE TABLE IF NOT EXISTS SetData (
		DataID INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
		Text          TEXT    NOT NULL,
		Answer        TEXT    NOT NULL, 
		SetID INT,
		FOREIGN KEY (SetID) REFERENCES Sets(SetID)
	);
	`
	_, err = db.Exec(SetDataRequest)
	if err != nil {
		return err
	}

	return nil
}

func NameAlreadyUsed(dataBaseName string) bool {
	row := db.QueryRow("SELECT SetID FROM Sets WHERE SetName = ?", dataBaseName)

	var SetID string
	err := row.Scan(&SetID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return true
}
