package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	helper "github.com/heartziq/course-mgmt/server/helpers"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func genAPIKey(hashedPwd string) string {

	secret := hashedPwd + time.Now().String() // append current time

	hash3 := sha256.New()
	hash3.Write([]byte(secret))

	return hex.EncodeToString(hash3.Sum(nil))
}

func Home(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	fmt.Fprintf(w, "Welcome to the REST API %v", key)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request - Unacceptable Method GET"))
		return
	}

	// Get username and pwd from Body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply course information " +
			"in JSON format"))

		return
	}

	// if empty body
	if len(data) < 1 {
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply body "))

		return
	} else {
		newUser := new(user)
		json.Unmarshal(data, newUser)
		if newUser.Password == "" || newUser.UserName == "" {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply username & password "))

			return
		}

		if err = PopulateNewUser(newUser, newUser.Password); err != nil {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Error hashing Pwd "))

			return
		}
		fmt.Println(newUser)
		_, err := InsertUser(newUser)
		if err != nil {
			panic(err)
		}

		cookie := &http.Cookie{
			Name:     "apiKey",
			Value:    newUser.APIKey,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		w.WriteHeader(http.StatusOK)
		// w.Write([]byte(newUser.APIKey))
		w.Write([]byte("status 200 - Registration successful"))

	}

}

func AllCourses(w http.ResponseWriter, r *http.Request) {
	c := GetRecords()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(c)
}

func Course(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	if r.Method == "GET" {

		// w.Header().Set("Access-Control-Allow-Origin", "*")
		courseId := params["courseid"]
		c, err := GetOneCourse(courseId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
			return
		}

		json.NewEncoder(w).Encode(c)
	}
	// DELETE a course
	if r.Method == "DELETE" {
		courseId := params["courseid"]
		if ok, err := DeleteRecord(courseId); !ok {
			log.Printf("[handlers.go]: Error DeleteRecord() %v", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new course
		if r.Method == "POST" {
			// read the string sent to the service
			var newCourse course
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newCourse)
				if newCourse.Title == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply course " +
							"information " + "in JSON format"))
					return
				}

				courseId := params["courseid"]
				// check if course exists; add only if
				// course does not exist

				if _, err := GetOneCourse(courseId); err != nil {
					InsertRecord(&courseId, newCourse.Title, newCourse.Details, newCourse.Trainer)
					w.WriteHeader(http.StatusCreated)
					msg := fmt.Sprintf("201 - Course added: %s\n",
						courseId)
					w.Write([]byte(msg))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(
						"409 - Duplicate course ID"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information " +
					"in JSON format"))
			}
		}

		//---PUT is for creating or updating
		// existing course---
		if r.Method == "PUT" {
			var newCourse course
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &newCourse)
				if newCourse.Title == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply course " +
							" information " +
							"in JSON format"))
					return
				}
				// check if course exists; add only if
				// course does not exist
				courseId := params["courseid"]
				if _, err := GetOneCourse(courseId); err != nil {

					InsertRecord(&courseId, newCourse.Title, newCourse.Details, newCourse.Trainer)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " +
						courseId))
				} else {
					// update course
					EditRecord(courseId, newCourse.Title, newCourse.Details, newCourse.Trainer)
					w.WriteHeader(http.StatusNoContent)
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " +
					"course information " +
					"in JSON format"))
			}
		}

	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	// Get username and pwd from Body
	data, err := ioutil.ReadAll(r.Body)

	// No user/pwd supplied
	if err != nil {
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply user information " +
			"in JSON format"))

		return
	}

	newUser := new(user)
	json.Unmarshal(data, newUser) // let FrontEnd do the check of empty input

	// retrieve user
	retrievedUser, err := GetOneUser(newUser.UserName)
	// User does not exist in DB
	if err != nil {
		log.Printf("error retrieving user: %v\n", err)
		// Fake user account to simulate user's existence in DB
		retrievedUser = &user{
			Password: "awfep9384utw5vu2m~~",
		}
	}

	if !helper.VerifyPassword([]byte(retrievedUser.Password), newUser.Password) {
		http.Error(w, "Wrong user/password combo", http.StatusUnauthorized) // 401
		return
	}

	// ON LOGIN SUCCESS //

	// Create jwt token
	jwt, err := helper.GenToken(helper.KEY, retrievedUser.Id)
	if err != nil {
		panic(err)
	}

	// Send jwt via cookie
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	// generate new apikey and edit user rrecord
	// if genNewKey := mux.Vars(r)["NewKey"]; genNewKey == "True" {
	// 	newAPIKey := genAPIKey(retrievedUser.Password)

	// 	// set expiry date
	// 	newExpiryDate := time.Now().Add(time.Hour * 24 * 7).Format(FORMAT)
	// 	EditUser(retrievedUser.Password, newAPIKey, newExpiryDate)
	// 	cookie.Value = newAPIKey
	// 	responseMsg = "New api key generated"
	// } else {
	// 	cookie.Value = retrievedUser.APIKey
	// }

	w.WriteHeader(http.StatusAccepted)

	// Send APIKey via json
	json.NewEncoder(w).Encode(map[string]string{"access_token": retrievedUser.APIKey})

}

func TestDraftCookie(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// cookie := &http.Cookie{
	// 	Name:     "token",
	// 	Value:    "afoiwfoninvawv",
	// 	HttpOnly: true,
	// }
	// http.SetCookie(w, cookie)
	// w.Header().Set("Content-Type", "application/json")

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "somejwttoken",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusAccepted)
	// w.Write([]byte("cookie set"))

	data := map[string]string{"token": "afoiwfoninvawv"}

	json.NewEncoder(w).Encode(data)
}

func TestGetToken(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// if c, err := r.Cookie("token"); err == nil {
	// 	log.Printf("c: %v\n", c.Value)
	// 	if c.Value == "afoiwfoninvawv" {
	// 		w.WriteHeader(http.StatusAccepted)
	// 		w.Write([]byte("logged in OK"))

	// 		return
	// 	}
	// }

	if c, err := r.Cookie("jwt"); err == nil {
		log.Printf("jwt: %v\n", c.Value)

		mechanism := strings.Split(r.Header.Get("Authorization"), " ")
		log.Printf("token found %v\n", mechanism)
		if len(mechanism) > 1 && mechanism[0] == "Bearer" {
			if token := mechanism[1]; token == "afoiwfoninvawv" {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("ok"))

				return
			}
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("logged in OK"))

		return

	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Not found 401"))

}
