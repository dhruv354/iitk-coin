package sqlite3Func

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

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
	row := db.QueryRow("SELECT rollno  from User where rollno= ?", rollno)
	temp := -1
	row.Scan(&temp)
	return temp != -1
}

func IsUserCoinExists(db *sql.DB, rollno int) bool {

	row := db.QueryRow("SELECT rollno  from UserData where rollno= ? ", rollno)
	temp := -1
	row.Scan(&temp)
	return temp != -1
}

func InsertIntoTable(db *sql.DB, name string, rollno int, batch int, password string) {

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

	context := context.Background()
	tx, err := db.BeginTx(context, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("inside InsertCoins Function")
	fmt.Println(coins)
	// updateCoins := `UPDATE USERDATA SET coins = coins + ? WHERE rollno = ?`
	statement, err := db.ExecContext(context, "UPDATE USERDATA SET coins = coins + ? WHERE rollno = ?", coins, rollno)
	if err != nil {
		tx.Rollback()
	}

	rows_affected, err := statement.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rows_affected != 1 {

		tx.Rollback()

	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	fmt.Println(statement)
}

func UpdateTransactionHistory(db *sql.DB, Sender int, Receiver int, transaction_Amount int, isAward int, redeems int, dateAndTime string) {

	// fmt.Println(Sender)
	// fmt.Println(Receiver)
	// fmt.Println(transaction_Amount)
	// fmt.Println(isAward)
	// fmt.Println(redeems)
	// fmt.Println(dateAndTime)
	// fmt.Println("inside UpdateTransactionHistory Function")
	// fmt.Println(db)

	addtrans, err := db.Prepare(`INSERT INTO EVENTS(sender,receiver,amount,isreward,date, redeem) VALUES(?,?,?,?,?,?)`)

	if err != nil {
		panic(err)
	}
	addtrans.Exec(Sender, Receiver, transaction_Amount, isAward, dateAndTime, redeems)
	fmt.Println("successfully updates transaction history Table")
}

func IsAwardExist(db *sql.DB, rollno int) bool {

	fmt.Println("inside IsAwardExist Function")

	rows, err := db.Query("SELECT isreward from EVENTS WHERE receiver = ? ", rollno)
	if err != nil {
		panic(err)
	}
	fmt.Println("here")
	var isAward int
	defer rows.Close()
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

	rows, err := db.Query("SELECT batch from User WHERE rollno = ?", rollno)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var batch int
	rows.Scan(&batch)

	return batch
}

func DisplayTransactionTable(db *sql.DB) {

	rows, err := db.Query("SELECT sender, receiver, isreward FROM  EVENTS")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var rollno int
	var receiver_rollno int
	var isaward int

	for rows.Next() {
		rows.Scan(&rollno, &receiver_rollno, &isaward)
		fmt.Println(rollno, " ", receiver_rollno, " ", isaward)
	}
}

func GetUserCoins(db *sql.DB, rollno int) int {

	fmt.Println("inside GetUSerCoins function")
	row := db.QueryRow("SELECT coins  from UserData where rollno= ?", rollno)
	var temp int
	row.Scan(&temp)
	return temp
}

func RedeemCoins(db *sql.DB, rollno int, coins int) {

	context := context.Background()
	tx, err := db.BeginTx(context, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("inside RedeemCoins Function")
	fmt.Println(coins)
	// updateCoins := `UPDATE USERDATA SET coins = coins + ? WHERE rollno = ?`
	statement, err := db.ExecContext(context, "UPDATE USERDATA SET coins = coins - ? WHERE rollno = ?", coins, rollno)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	rows_affected, err := statement.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rows_affected != 1 {
		tx.Rollback()
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	fmt.Println(statement)

}

func InsertIntoRedeemRequest(db *sql.DB, rollno int, coins int, item string, dateAndTime string) {

	fmt.Println("inside insertintoredeem request")

	insertRedeemRequest := `INSERT INTO REDEEMREQUESTS(rollno, coins, item, status, date) VALUES(?, ?, ?, ?, ?)`

	statement, err := db.Prepare(insertRedeemRequest)

	if err != nil {
		panic(err)
	}
	print(rollno, coins, item, 0, dateAndTime)
	statement.Exec(rollno, coins, item, 0, dateAndTime)
	fmt.Println("Made  a redeem request of item", item)
}

func DisplayRedeemTable(db *sql.DB) {
	rows, err := db.Query("SELECT rollno, coins, item FROM  REDEEMREQUESTS")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var rollno int
	var coins int
	var item string

	for rows.Next() {
		rows.Scan(&rollno, &coins, &item)
		fmt.Println(rollno, " ", coins, " ", item)
	}
}

// function to check a particular id in redeem request table and if it exists then check pending or unpending request

func ApproveStatus(db *sql.DB, id int) {

	fmt.Println("inside approve status for admin")

	row := db.QueryRow("SELECT rollno, coins, status from REDEEMREQUESTS where id= ?", id)
	coins := -1
	rollno := -1
	status := -1
	row.Scan(&rollno, &coins, &status)
	if rollno == -1 {
		fmt.Println("this id do not exist")
		return
	}

	//now update status of the pending request
	_, err := db.Exec(`UPDATE REDEEMREQUESTS SET status = ? WHERE id = ?`, 1, id)

	if err != nil {
		panic(err)
	}

	//now update the coins of the user
	_, err = db.Exec(`UPDATE UserData SET coins = coins - ? WHERE rollno = ?`, coins, rollno)

	if err != nil {
		panic(err)
	}

	//update transaction history table
	UpdateTransactionHistory(db, rollno, rollno, coins, 0, 1, time.Now().String())
	// DisplayTransactionTable(db)

}
