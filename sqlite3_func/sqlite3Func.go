package sqlite3Func

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func IsUserExists(db *sql.DB, rollno int) bool {
	row := db.QueryRow("SELECT rollno  from User where rollno= ? ", rollno)
	temp := ""
	row.Scan(&temp)
	return temp != ""
}

func InsertIntoTable(db *sql.DB, name string, rollno int, password string) {
	fmt.Println("inside insetintotable")

	insertStudent_info := `INSERT INTO User(rollno, name, password) VALUES(?, ?, ?)`

	insertStudent, err := db.Prepare(insertStudent_info)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	insertStudent.Exec(rollno, name, password)
	fmt.Println("inserted Student in the table")
}
