package academiq

import (
	"database/sql"
	"fmt"
)

func Launch() {
	fmt.Println("##### Launching Academiq #####")

	fmt.Println("\nLoading DataBase")
	var err error
	db, err = sql.Open("sqlite3", "./academiq.db")
	if err != nil {
		fmt.Println("Error while opening database ", err)
	}
	defer db.Close()

	err = InitTables()
	if err != nil {
		fmt.Println("Error while initiating tables ", err)
	}

	fmt.Println("\nLaunching Server")
	LaunchServer()
}
