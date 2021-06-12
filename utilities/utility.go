package utility

import (
	"golang.org/x/crypto/bcrypt"
)

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
