package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	handler "github.com/dhruv354/iitk-coin/Handlers"
	_ "github.com/mattn/go-sqlite3"
)

//function to create table which will store user rollno
//name and password
func createSqliteTable(db *sql.DB) {
	//creating a string with table info
	UserTable_info := `CREATE TABLE IF NOT EXISTS User(
		"rollno" INTEGER UNSIGNED NOT NULL,
		"name" TEXT NOT NULL,
		"password" TEXT NOT NULL,
		"batch" INT UNSIGNED  NOT NULL,
		"isAdmin" INT UNSIGNED NOT NULL,
		"events" INT UNSIGNED NOT NULL
		);`
	//create table with above info
	UserTable, err := db.Prepare(UserTable_info)
	if err != nil {
		fmt.Println(err)
	}

	UserTable.Exec()
	fmt.Println("User table created or not altered if already created")

}

//function to create table for user_coin data
func UserCoinTable(db *sql.DB) {
	UserCoin_info := `CREATE TABLE IF NOT EXISTS UserData(
		"rollno" INTEGER NOT NULL,
		"coins" INTEGER NOT NULL
		);`

	statement, err := db.Prepare(UserCoin_info)
	if err != nil {
		panic(err)
	}
	statement.Exec()
	fmt.Println("user coins table created")
}

func createTransactionTable(db *sql.DB) {
	transaction := `CREATE TABLE IF NOT EXISTS Transactionhistory(
		"rollno" INTEGER UNSIGNED NOT NULL,
		"isReward" INTEGER ,
		"transfered_to" INTEGER UNSIGNED NOT NULL,
		"transfered_amount" INTEGER NOT NULL,
		"redeems" INTEGER,
		"date" TEXT NOT NULL
		);`
	statement, err := db.Prepare(transaction)
	if err != nil {
		panic(err)
	}
	statement.Exec()
	fmt.Println("trasaction  table created")
}

// func createTransactionTable(db *sql.DB) {
// 	transaction := `CREATE TABLE IF NOT EXISTS Transaction(
// 		"sender" INT UNSIGNED NOT NULL,
// 		"isReward" INT ,
// 		"receiver" INT UNSIGNED NOT NULL,
// 		"amount" INT NOT NULL,
// 		"redeem" INT ,
// 		"datetime" TEXT NOT NULL
// 		);`
// 	statement, err := db.Prepare(transaction)
// 	if err != nil {
// 		print("here")
// 		fmt.Println(err)
// 		// panic(err)
// 	}
// 	statement.Exec()
// 	fmt.Println("trasaction  table created")
// }

func main() {
	database, err := sql.Open("sqlite3", "./Student_info.db")
	if err != nil {
		panic(err)
	}
	fmt.Println("created my database")

	createSqliteTable(database)
	UserCoinTable(database)
	createTransactionTable(database)
	// "Signin" and "Welcome" are the handlers that we will implement
	http.HandleFunc("/login", handler.LoginRoute)
	http.HandleFunc("/secretpage", handler.Secretpage)
	http.HandleFunc("/signup", handler.SignupRoute)
	http.HandleFunc("/logout", handler.Logout)
	http.HandleFunc("/getcoins", handler.GetUserCoins)
	http.HandleFunc("/addcoins", handler.AddCoins)
	http.HandleFunc("/transfercoins", handler.TransferCoin)
	// start the server on port 8080

	// sqlite3Func.DisplayTransactionTable(database)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
