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

func IsUserCoinExists(db *sql.DB, rollno int) bool {
	row := db.QueryRow("SELECT rollno  from UserData where rollno= ? ", rollno)
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

	//when that user is added in the main table also create
	//its entry in USerData table
	insertcoin_info := `INSERT INTO UserData(rollno, coins) VALUES(?, ?)`
	statement, err := db.Prepare(insertcoin_info)
	if err != nil {
		panic(err)
	}
	statement.Exec(rollno, 0)

	fmt.Println("inserted Student with 0 coins in the table")
}

func UpdateUserCoins(db *sql.DB, rollno int, coins int) {
	fmt.Println("inside InsertCoins Function")
	fmt.Println(coins)
	// updateCoins := `UPDATE USERDATA SET coins = coins + ? WHERE rollno = ?`
	statement, err := db.Exec("UPDATE USERDATA SET coins = coins + ? WHERE rollno = ?", coins, rollno)
	if err != nil {
		panic(err)
	}
	fmt.Println(statement)
}
