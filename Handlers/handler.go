package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
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
	Batch    int    `json:"batch"`
}

type Claims struct {
	Username string `json:"username"`
	Rollno   int    `json:"rollno"`
	jwt.StandardClaims
}

type UserCoins struct {
	Rollno int `json:"rollno"`
	Coins  int `json:"coins"`
}

type transferBWUsers struct {
	ReceiverRollno int `json:"rollno"` //receiver rollno
	Coins          int `json:"coins"`  //Coins to be transfered
}

// type transactionHistory struct {
// 	Sender             int    //roll no of sender
// 	Receiver           int    //roll no of receiver
// 	transaction_amount int    //coins transfered
// 	isAward            int    //if it is a  award or another person transfered to it
// 	dateAndTime        string //timestanp
// 	redeems            int    //redeems
// }

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
					Rollno:   user_data.Rollno,
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
		batch := user_data.Batch
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

		sqlite3Func.InsertIntoTable(database, name, rollno, batch, hashed_password)
		w.Write([]byte("USer with this rollno created"))
	}
}

/************************logout endpoint**********************************/

func Logout(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "my_json_token",
		MaxAge: -1}
	http.SetCookie(w, &c)

	w.Write([]byte("old cookie deleted!\n"))
}

/**********************route to get user coins ******************************/
func GetUserCoins(w http.ResponseWriter, r *http.Request) {

	// if !utility.IsLoggedin(w, r) {
	// 	fmt.Fprintf(w, "first login in to access this endpoint")
	// 	return
	// }
	isLogged, claims := utility.IsLoggedin(w, r)
	if !isLogged {
		fmt.Fprintf(w, "first login in to access this endpoint")
		return
	}

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

	var rollno int
	var coins int
	print("claims.Rollno: ", claims.Rollno)
	for rows.Next() {
		rows.Scan(&rollno, &coins)
		if rollno == claims.Rollno {
			fmt.Fprintf(w, "you have coins: %d", coins)
			return
		}
	}

	fmt.Fprintf(w, "this rollno does not exist")
}

//****************************************handler to award coins to the user****************************/

func AddCoins(w http.ResponseWriter, r *http.Request) {

	isLogged, claims := utility.IsLoggedin(w, r)
	if !isLogged {
		fmt.Fprintf(w, "first login in to access this endpoint")
		return
	}

	fmt.Println("inside AddCoins function")

	//handling GET requests
	if r.Method == "GET" {
		fmt.Fprintf(w, "this route is only for post request so please make a post request")
		return
	}
	//if current user is not admin
	if claims.Rollno != 190294 {
		fmt.Fprintf(w, "only admin is allowed to add coins into it")
		return
	}

	var user_coins UserCoins

	json.NewDecoder(r.Body).Decode(&user_coins)

	//admin cannot
	if user_coins.Rollno == 190294 {
		fmt.Fprintf(w, "you are a admin so you cannot add coins to your wallet")
		return
	}

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

	// no redeem
	dateAndTime := time.Now().String() //current time

	sqlite3Func.UpdateTransactionHistory(database, 190294, user_coins.Rollno, user_coins.Coins, 1, 0, dateAndTime)

	// sqlite3Func.DisplayTransactionTable(database)

	fmt.Fprintf(w, "added coins in your wallet")

}

/************************************Handler to transfer coins between two users********************/

func TransferCoin(w http.ResponseWriter, r *http.Request) {

	//if get request return
	if r.Method == "GET" {
		fmt.Fprintf(w, "only post request is possible at this route")
		return
	}

	isLogged, claims := utility.IsLoggedin(w, r)
	if !isLogged {
		fmt.Fprintf(w, "first login in to access this endpoint")
		return
	}

	database, err := sql.Open("sqlite3", "../iitk-coin/Student_info.db")
	if err != nil {
		panic(err)
	}

	//check if the user has participated in any of the events

	if !sqlite3Func.IsAwardExist(database, claims.Rollno) {
		fmt.Fprintf(w, "you haven't participated in any of the events so you can't send money to anyone")
		return
	}

	var transfer_data transferBWUsers
	json.NewDecoder(r.Body).Decode(&transfer_data)

	//checking if both users have an account or not
	isUser1Exists := sqlite3Func.IsUserExists(database, transfer_data.ReceiverRollno)

	isUser2Exists := sqlite3Func.IsUserExists(database, claims.Rollno)

	if !isUser1Exists || !isUser2Exists {
		fmt.Fprintf(w, "Either user1 or user2 does not exists")
		return
	}

	tx, err := database.Begin()
	if err != nil {
		fmt.Println("error lies in database.begin()")
		return
	}

	//statement and updates in the same statement to solve problems during concurrency
	sender_batch := sqlite3Func.FindBatch(database, claims.Rollno)
	receiver_batch := sqlite3Func.FindBatch(database, transfer_data.ReceiverRollno)

	var tax float32

	if sender_batch == receiver_batch {
		tax = 0.2
	} else {
		tax = 0.33
	}

	var mutex = &sync.Mutex{}
	mutex.Lock()

	_, err1 := database.Exec(`UPDATE  UserData SET coins = coins + ? WHERE rollno = ?`, int((1-tax)*float32(transfer_data.Coins)), transfer_data.ReceiverRollno)

	mutex.Unlock()

	if err1 != nil {
		//if some error rollback databse to initial condition and print the error

		fmt.Println("error lies in database.Exec() err1")
		print(err1)
		tx.Rollback()
		return
		// panic(err)
	}

	mutex.Lock()

	_, err2 := database.Exec(`UPDATE UserData SET coins = coins - ?  WHERE rollno = ? AND coins - ? >= 0`, transfer_data.Coins, claims.Rollno, transfer_data.Coins)

	mutex.Unlock()

	if err2 != nil {
		//if some error rollback databse to initial condition and print the error
		fmt.Println(err2)
		fmt.Println("error lies in database.Exec() err2")
		tx.Rollback()
		return
	}
	//we are here so this means transaction is successful so Commit this change to the database
	tx.Commit()

	//save this transaction history
	sqlite3Func.UpdateTransactionHistory(database, claims.Rollno, transfer_data.ReceiverRollno, transfer_data.Coins,
		0, 0, time.Now().String())
	fmt.Fprintf(w, "transaction is successful")
}

// c, err := r.Cookie("Tok")
// if err != nil {
// 	if err == http.ErrNoCookie {

// 		w.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprintf(w, "Please login to acess the page")

// 	}

// 	w.WriteHeader(http.StatusBadRequest)
// 	return
// }

// tknStr := c.Value

// claims := &Claims{}

// tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
// 	return jwtKey, nil
// })

// if err != nil {
// 	if err == jwt.ErrSignatureInvalid {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}
// 	w.WriteHeader(http.StatusBadRequest)
// 	return
// }

// if !tkn.Valid {
// 	w.WriteHeader(http.StatusUnauthorized)
// 	return
// }
