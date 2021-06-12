package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOk)
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
			panic(err)
		} else {
			fmt.Println("Connected with database")
		}

		if sqlite3Func.IsUserExists(database, rollno) {
			// w.Write([]byte("USer with this rollno created"))
			fmt.Fprintf(w, "User exists with this rollno")
			return
		}

		sqlite3Func.InsertIntoTable(database, name, rollno, hashed_password)
		w.Write([]byte("USer with this rollno created"))
	}
}

// /**********************/

// func isUserExists(db *sql.DB, rollno int) bool {
// 	row := db.QueryRow("SELECT rollno  from User where rollno= ? ", rollno)
// 	temp := ""
// 	row.Scan(&temp)
// 	return temp != ""
// }

// func insertIntoTable(db *sql.DB, name string, rollno int, password string) {
// 	fmt.Println("inside insetintotable")

// 	insertStudent_info := `INSERT INTO User(rollno, name, password) VALUES(?, ?, ?)`

// 	insertStudent, err := db.Prepare(insertStudent_info)

// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}
// 	insertStudent.Exec(rollno, name, password)
// 	fmt.Println("inserted Student in the table")
// }
