package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//function to create user table in a database
func createSqliteTable(db *sql.DB) {
	//creating a string with table info
	UserTable_info := `CREATE TABLE IF NOT EXISTS User(
		"rollno" INTEGER NOT NULL,
		"name" TEXT NOT NULL
		);`
	//create table with above info
	UserTable, err := db.Prepare(UserTable_info)
	if err != nil {
		panic(err)
	}
	// fmt.Println(reflect.TypeOf(UserTable))
	UserTable.Exec()
	fmt.Println("User table created ")
}

func insertIntoTable(db *sql.DB, student *student_details) {

	insertStudent_info := `INSERT INTO User(rollno, name) VALUES(?, ?)`

	insertStudent, err := db.Prepare(insertStudent_info)

	if err != nil {
		panic(err)
	}
	insertStudent.Exec(student.rollno, student.name)
	fmt.Println("inserted Student in the table")
}

func main() {

	my_database, err := sql.Open("sqlite3", "./Student_info.db")

	//if some error in creating the database
	if err != nil {
		panic(err)
	}

	type student_details struct {
		name   string
		rollno int
	}
	student2 := student_details{name: "shyam", rollno: 1111}

	createSqliteTable(my_database)
	insertIntoTable(my_database, student2)
}
