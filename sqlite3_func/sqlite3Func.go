package sqlite3Func

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

func IsUserExists(db *sql.DB, rollno int) bool {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	row := db.QueryRow("SELECT rollno  from User where rollno= ? ", rollno)
	temp := -1
	row.Scan(&temp)
	mutex.Unlock()
	return temp != -1
}

func IsUserCoinExists(db *sql.DB, rollno int) bool {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	row := db.QueryRow("SELECT rollno  from UserData where rollno= ? ", rollno)
	temp := -1
	row.Scan(&temp)
	mutex.Unlock()
	return temp != -1
}

func InsertIntoTable(db *sql.DB, name string, rollno int, batch int, password string) {
	fmt.Println("inside insetintotable")

	var mutex = &sync.Mutex{}
	mutex.Lock()
	insertStudent_info := `INSERT INTO User(rollno, name, password, batch, isAdmin, events) VALUES(?, ?, ?, ?, ?, ?)`

	insertStudent, err := db.Prepare(insertStudent_info)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if rollno == 190294 {
		insertStudent.Exec(rollno, name, password, batch, 1, 0)
	} else {
		insertStudent.Exec(rollno, name, password, batch, 0, 0)
	}

	//when that user is added in the main table also create
	//its entry in USerData table
	insertcoin_info := `INSERT INTO UserData(rollno, coins) VALUES(?, ?)`
	statement, err := db.Prepare(insertcoin_info)
	if err != nil {
		panic(err)
	}
	statement.Exec(rollno, 0)
	mutex.Unlock()

	fmt.Println("inserted Student with 0 coins in the table")
}

func UpdateUserCoins(db *sql.DB, rollno int, coins int) {
	fmt.Println("inside InsertCoins Function")
	fmt.Println(coins)
	// updateCoins := `UPDATE USERDATA SET coins = coins + ? WHERE rollno = ?`
	var mutex = &sync.Mutex{}
	mutex.Lock()
	statement, err := db.Exec("UPDATE USERDATA SET coins = coins + ? WHERE rollno = ?", coins, rollno)
	mutex.Unlock()
	if err != nil {
		panic(err)
	}
	fmt.Println(statement)

}
