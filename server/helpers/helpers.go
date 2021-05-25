package helper

import (
	"database/sql"

	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// GetUser retrieve a user based on username provided as param
// if user does not exist,
// error "user not found" will be return alongside a nil obj

func VerifyPassword(hashedPassword []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		log.Println(err)
		return false
	} else {
		return true
	}
}

func DeleteUser(db *sql.DB, username string) (bool, error) {
	results, err := db.Exec("DELETE FROM Users WHERE UserName=?;", username)
	if err != nil {
		return false, errors.New("Error deleting record")
	}

	rows, _ := results.RowsAffected()
	fmt.Println(rows)
	if rows < 1 {
		return false, errors.New("record not foudn")
	}

	return true, nil

}
