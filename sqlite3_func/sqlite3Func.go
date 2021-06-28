package sqlite3Func

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// type transactionHistory struct {
// 	Sender             int    //roll no of sender
// 	Receiver           int    //roll no of receiver
// 	transaction_amount int    //coins transfered
// 	isAward            int    //if it is a  award or another person transfered to it
// 	dateAndTime        string //timestanp
// 	redeems            int    //redeems
// }

// type row struct {
// 	Count int `json:"count"`
// }

var m sync.Mutex

func IsUserExists(db *sql.DB, rollno int) bool {
	m.Lock()
	defer m.Unlock()

	fmt.Println("inside isUser exists function")
	row := db.QueryRow("SELECT rollno  from User where rollno= ? ", rollno)
	temp := -1
	row.Scan(&temp)
	return temp != -1
}

func IsUserCoinExists(db *sql.DB, rollno int) bool {
	m.Lock()
	defer m.Unlock()
	row := db.QueryRow("SELECT rollno  from UserData where rollno= ? ", rollno)
	temp := -1
	row.Scan(&temp)
	return temp != -1
}

func InsertIntoTable(db *sql.DB, name string, rollno int, batch int, password string) {
	m.Lock()
	defer m.Unlock()
	fmt.Println("inside insetintotable")

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

	fmt.Println("inserted Student with 0 coins in the table")
}

func UpdateUserCoins(db *sql.DB, rollno int, coins int) {
	m.Lock()
	defer m.Unlock()
	fmt.Println("inside InsertCoins Function")
	fmt.Println(coins)
	// updateCoins := `UPDATE USERDATA SET coins = coins + ? WHERE rollno = ?`
	statement, err := db.Exec("UPDATE USERDATA SET coins = coins + ? WHERE rollno = ?", coins, rollno)
	if err != nil {
		panic(err)
	}
	fmt.Println(statement)
}

func UpdateTransactionHistory(db *sql.DB, Sender int, Receiver int, transaction_Amount int, isAward int, redeems int, dateAndTime string) {
	fmt.Println(Sender)
	fmt.Println(Receiver)
	fmt.Println(transaction_Amount)
	fmt.Println(isAward)
	fmt.Println(redeems)
	fmt.Println(dateAndTime)
	fmt.Println("inside UpdateTransactionHistory Function")
	fmt.Println(db)
	_, err := db.Exec(`INSERT INTO Transactionhistory(rollno, isReward, transfered_to, transfered_amount, redeems, date) VALUES(?, ?, ?, ?, ?, ?)`, Sender, isAward, Receiver, transaction_Amount, redeems, dateAndTime)
	// `INSERT INTO UserData(rollno, coins) VALUES(?, ?)`
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// statement.Exec(Sender, isAward, Receiver, transaction_Amount, redeems, dateAndTime)
	fmt.Println("successfully updates transaction history Table")
}

func IsAwardExist(db *sql.DB, rollno int) bool {

	m.Lock()
	defer m.Unlock()

	fmt.Println("inside IsAwardExist Function")

	rows, err := db.Query(`SELECT isReward from Transactionhistory WHERE rollno = ? `, rollno)
	if err != nil {
		panic(err)
	}
	fmt.Println("here")
	var isAward int
	for rows.Next() {
		rows.Scan(&isAward)
		if isAward == 1 {
			return true
		}
	}

	// response, err := db.Query("SELECT count(isReward) as count from Transaction_history WHERE rollno = ? AND isReward = ?", rollno, 1)
	// if err != nil {
	// 	panic(err)
	// }
	// var rows []row
	// err2 := response.Scan(&rows)

	// if err2 != nil {
	// 	panic(err)
	// }

	// count := rows[0].Count

	// fmt.Println("is award returned some value")
	// return count != 0
	fmt.Println("here")
	return false
}

func FindBatch(db *sql.DB, rollno int) int {
	m.Lock()
	defer m.Unlock()

	rows, err := db.Query(`SELECT batch from User WHERE rollno = ?`, rollno)
	if err != nil {
		panic(err)
	}
	var batch int
	rows.Scan(&batch)

	return batch
}

func DisplayTransactionTable(db *sql.DB) {
	m.Lock()
	defer m.Unlock()

	rows, err := db.Query(`SELECT rollno, transfered_to, transfered_amount, isReward, date FROM Transactionhistory`)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var rollno int
	var receiver_rollno int
	var coins_transfered int
	var date string

	for rows.Next() {
		rows.Scan(&rollno, &receiver_rollno, &coins_transfered, &date)
		fmt.Println(rollno, " ", receiver_rollno, " ", coins_transfered, " ", date, " ")
	}
}
