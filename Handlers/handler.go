package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	sqlite3Func "github.com/dhruv354/iitk-coin/sqlite3_func"
	utility "github.com/dhruv354/iitk-coin/utilities"
	_ "github.com/mattn/go-sqlite3"
)

var jwt_key = []byte("dhruv_singhal")

type UserData struct {
	Name     string `json:"name"`
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type UserCoins struct {
	Rollno int `json:"rollno"`
	Coins  int `json:"coins"`
}

type transferBWUsers struct {
	Rollno1 int `json:"rollno1"` //sender rollno
	Rollno2 int `json:"rollno2"` //receiver rollno
	Coins   int `json:"coins"`   //Coins to be transfered
}

/*****************login route handler ******************************/
func LoginRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, "currently not dealing in get request, please make a post request only")
	} else {
		database, err := sql.Open("sqlite3", "../iitk-coin/Student_info.db")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		var user_data UserData

		json.NewDecoder(r.Body).Decode(&user_data)

		rows, err := database.Query("SELECT password,rollno FROM User")
		if err != nil {
			panic(err)
		}

		//iterate through the database
		var database_password string
		var database_rollno int
		var userFound bool = false

		for rows.Next() {
			rows.Scan(&database_password, &database_rollno)
			isRollnoSame := database_rollno == user_data.Rollno
			isPasswordSame := utility.DoesPasswordsMatch(database_password, []byte(user_data.Password))

			//if roll no and password matches then
			//user is authenticated and send a json token
			//that will expire after some time
			if isRollnoSame && isPasswordSame {
				userFound = true
				//time after which token gets expired
				expirationTime := time.Now().Add(12 * time.Minute)
				//data to be stored in the cookie
				claims := &Claims{
					Username: user_data.Name,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: expirationTime.Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenString, err := token.SignedString(jwt_key)
				if err != nil {
					fmt.Println("error: not able to generate a token")
					//code for the response
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				cookie := http.Cookie{
					Name:    "my_json_token",
					Value:   tokenString,
					Expires: expirationTime,
				}
				http.SetCookie(w, &cookie)
				fmt.Fprintf(w, "you are signed in for 12 minutes after that you have to login again")
				return
			}
		}
		time.Sleep(1 * time.Second)
		if !userFound {
			fmt.Fprintf(w, "Oh No! invalid username or password ")
		}
	}
}

/****************************************secretpage handler *****************************************/

func Secretpage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		fmt.Println("first let me verify you")
	}
	cookie, err := r.Cookie("my_json_token")
	//if some error occured
	if err != nil {
		if err == http.ErrNoCookie {
			//if cookie is not there
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "you are not authorized to access this page")
		}
		// handling other types of errors
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//access cookie value present
	tknStr := cookie.Value

	//address of an empty Claim Struct
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwt_key, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "you are verified %s", claims.Username)
}

/****************************************signup page handler**************************************/

func SignupRoute(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside signup route")
	// fmt.Println(r)
	if r.Method == "GET" {
		fmt.Println("currently not habdling get requests so please make a post request")
	} else {
		var user_data UserData
		// user_data, _ = json.Marshal(user_data)
		json.NewDecoder(r.Body).Decode(&user_data)
		fmt.Println(user_data)
		name := user_data.Name
		password := user_data.Password
		rollno := user_data.Rollno
		hashed_password, _ := utility.HashPassword(password)

		// database_path, _ := filepath.Abs("..")
		// fmt.Println(database_path)

		database, err := sql.Open("sqlite3", "../iitk-coin/Student_info.db")
		if err != nil {
			fmt.Println("error", err)
			panic(err)
		} else {
			fmt.Println("Connected with database")
		}

		if sqlite3Func.IsUserExists(database, rollno) {
			// w.Write([]byte("USer with this rollno created"))
			fmt.Fprintf(w, "User exists with this rollno")
			return
		}

		// rows, err := database.Query("SELECT name, rollno from User")
		// if err != nil {
		// 	fmt.Println("error")
		// 	panic(err)
		// }

		// var database_name string
		// var database_rollno int

		// var flag bool = true
		// var temp1 bool = false
		// var temp2 bool = false

		// for rows.Next() {
		// 	rows.Scan(&database_name, &database_rollno)
		// 	if database_name == user_data.Name {
		// 		temp1 = true
		// 	}
		// 	if database_rollno == user_data.Rollno {
		// 		temp2 = true
		// 	}
		// 	if temp1 || temp2 {
		// 		flag = false
		// 		fmt.Fprintf(w, "user already exists")
		// 		return
		// 	}
		// 	time.Sleep(200 * time.Millisecond)
		// }

		// time.Sleep(3 * time.Second)
		// if flag {
		// 	sqlite3_func.InsertIntoTable(database, name, rollno, hashed_password)
		// 	fmt.Fprintf(w, "signed in")
		// }

		sqlite3Func.InsertIntoTable(database, name, rollno, hashed_password)
		w.Write([]byte("USer with this rollno created"))
	}
}

/**********************route to get user coins ******************************/
func GetUserCoins(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside getUserCoins Route")

	if r.Method == "POST" {
		fmt.Println("this route can only handle GET request")
		return
	}

	var user_coins UserCoins
	json.NewDecoder(r.Body).Decode(&user_coins)

	fmt.Println(user_coins)
	//open the database
	database, err := sql.Open("sqlite3", "../iitk-coin/Student_info.db")
	if err != nil {
		log.Fatal(err.Error())
		return
	} else {
		fmt.Println("connected with Database")
	}

	rows, err := database.Query("SELECT rollno ,coins FROM UserData")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	print("going to scan rows")

	var rollno int
	var coins int

	for rows.Next() {
		rows.Scan(&rollno, &coins)
		fmt.Println("scanning")
		if rollno == user_coins.Rollno {
			fmt.Fprintf(w, "you have coins: %d", coins)
			return
		}
	}

	fmt.Fprintf(w, "this rollno does not exist")
}

/****************************************handler to award coins to the user****************************/

func AddCoins(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside AddCoins function")

	//handling GET requests
	if r.Method == "GET" {
		fmt.Fprintf(w, "this route is only for post request so please make a post request")
		return
	}

	var user_coins UserCoins

	json.NewDecoder(r.Body).Decode(&user_coins)

	//open the database of Student_info
	database, err := sql.Open("sqlite3", "../iitk-coin/Student_info.db")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("successfully connected with database")

	//check if user with that rollno exists
	if !sqlite3Func.IsUserExists(database, user_coins.Rollno) {
		fmt.Fprintf(w, "user with this rollno does not exists")
		return
	}

	sqlite3Func.UpdateUserCoins(database, user_coins.Rollno, user_coins.Coins)
	fmt.Fprintf(w, "added coins in your wallet")

	//display userdata
	rows, err := database.Query(`SELECT rollno, coins from USERDATA`)
	if err != nil {
		panic(err)
	}
	var rollno int
	var coins int

	for rows.Next() {
		rows.Scan(&rollno, &coins)
		log.Println("rollno ", rollno, " ", "coins: ", coins)
	}
}

/************************************Handler to transfer coins between two users********************/

func TransferCoin(w http.ResponseWriter, r *http.Request) {

	//if get request return
	if r.Method == "GET" {
		fmt.Fprintf(w, "only post request is possible at this route")
		return
	}

	var transfer_data transferBWUsers
	json.NewDecoder(r.Body).Decode(&transfer_data)

	database, err := sql.Open("sqlite3", "../iitk-coin/Student_info.db")
	if err != nil {
		panic(err)
	}

	//checking if both users have an account or not
	isUser1Exists := sqlite3Func.IsUserExists(database, transfer_data.Rollno1)

	isUser2Exists := sqlite3Func.IsUserExists(database, transfer_data.Rollno2)

	if !isUser1Exists || !isUser2Exists {
		fmt.Fprintf(w, "Either user1 or user2 does not exists")
		return
	}

	//if both user exists
	//get current balance holded by both the users
	var user1_balance int
	var user2_balance int
	row1 := database.QueryRow(`SELECT coins FROM USERDATA WHERE rollno = ?`, transfer_data.Rollno1)
	row2 := database.QueryRow(`SELECT coins FROM USERDATA WHERE rollno = ?`, transfer_data.Rollno2)

	err1 := row1.Scan(&user1_balance)
	err2 := row2.Scan(&user2_balance)

	if err1 != nil && err1 != sql.ErrNoRows {
		// log the error
		fmt.Println(err1)
		panic(err1)
	}

	if err2 != nil && err2 != sql.ErrNoRows {
		// log the error
		fmt.Println(err2)
		panic(err2)
	}

	// check if sender has sufficient coins in his account
	if user2_balance < transfer_data.Coins {
		fmt.Fprintf(w, "user2 does not have sufficient money")
		return
	}

	//if he have enough coins to be sent
	statement, err := database.Prepare(`UPDATE  UserData SET coins = ? WHERE rollno = ?`)

	if err != nil {
		panic(err)
	}

	statement.Exec(user1_balance+transfer_data.Coins, transfer_data.Rollno1)
	statement.Exec(user2_balance-transfer_data.Coins, transfer_data.Rollno2)
	fmt.Fprintf(w, "transaction is successful")
}
