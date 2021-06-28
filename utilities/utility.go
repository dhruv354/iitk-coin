package utility

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	Name     string `json:"name"`
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
}

var jwt_key = []byte("dhruv_singhal")

type Claims struct {
	Username string `json:"username"`
	Rollno   int    `json:"rollno"`
	jwt.StandardClaims
}

//function to hash a user enetered password and returning the hashed as a string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func DoesPasswordsMatch(existing_password string, entered_password []byte) bool {
	//first convert the password into byte that is already existing in the database
	existing_password_hashed := []byte(existing_password)
	//compare password
	err := bcrypt.CompareHashAndPassword(existing_password_hashed, entered_password)
	return err == nil
}

func IsLoggedin(w http.ResponseWriter, r *http.Request) (bool, *Claims) {

	fmt.Println("inside IsLogged In function")
	cookie, err := r.Cookie("my_json_token")
	fmt.Println("cookie : ", cookie)
	//if some error occured
	if err != nil {
		if err == http.ErrNoCookie {
			//if cookie is not there
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "you are not authorized to access this page")

		}
		// handling other types of errors
		w.WriteHeader(http.StatusBadRequest)
		return false, nil
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

			return false, nil
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)

		return false, nil
	}

	if !tkn.Valid {
		print("5th")
		w.WriteHeader(http.StatusUnauthorized)

		return false, nil
	}
	fmt.Println("IsLogged In function retured true")

	return true, claims
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
