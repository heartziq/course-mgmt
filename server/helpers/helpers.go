package helper

import (
	"database/sql"
	"io"
	"os"
	"strings"
	"time"

	"errors"
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	SECRET_FILE = "secret.pem"
	KEY         string
)

func init() {
	KEY = initSecretKey()

}

func initSecretKey() string {

	//
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	// log.Fatal(path) // for example /home/user
	//
	file, err := os.Open(path + "\\helpers\\" + SECRET_FILE)
	if err != nil {
		panic(err)
	}

	content, _ := io.ReadAll(file)

	secret := strings.TrimSpace(string(content))
	secret = strings.TrimPrefix(string(secret), "-----BEGIN RSA PRIVATE KEY-----")
	secret = strings.TrimSuffix(string(secret), "-----END RSA PRIVATE KEY-----")
	return secret
}

// GetUser retrieve a user based on username provided as param
// if user does not exist,
// error "user not found" will be return alongside a nil obj

func GenToken(secret, userid string) (string, error) {

	mySigningKey := []byte(secret)
	expiryDate := time.Now().Add(time.Hour * 24 * 7).Unix()

	// get userid

	// Create the Claims
	claims := &jwt.StandardClaims{
		Audience:  userid,
		ExpiresAt: expiryDate,
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		panic(err)
	}

	return ss, nil
}

func VerifyToken(tokenString string) (bool, error) {
	// Verify
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(KEY), nil
	})

	if err != nil {

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token has expired
				log.Println("token expired")
				return false, errors.New("token expired")
			}
		}

		log.Println("invalid token")
		return false, errors.New("invalid token")

	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Welcome, ", claims["aud"])
		log.Println("Token is valid.")
		return true, nil
	}

	return false, errors.New("Unknown error")

}

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
