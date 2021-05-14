package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// const DUMMY_PWD = "password123uSE#4r:0;L0v3~~~"

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

	// body
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
		if newUser.Password == "" {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply password "))

			return
		}

		if err = PopulateNewUser(newUser, newUser.Password); err != nil {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Error hashing Pwd "))

			return
		}

		lastInsertedId, err := InsertUser(newUser)
		if err != nil {
			panic(err)
		}

		fmt.Printf("lastInsertedId: %v\n", lastInsertedId)

	}

}

func AllCourses(w http.ResponseWriter, r *http.Request) {
	c := GetRecords(db)

	json.NewEncoder(w).Encode(c)
}

func Course(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	if r.Method == "GET" {
		courseId := params["courseid"]
		c, err := GetOneCourse(db, courseId)
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
		if ok, err := DeleteRecord(db, courseId); !ok {
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

				if _, err := GetOneCourse(db, courseId); err != nil {
					InsertRecord(db, &courseId, newCourse.Title, newCourse.Details, newCourse.Trainer)
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
				if _, err := GetOneCourse(db, courseId); err != nil {

					InsertRecord(db, &courseId, newCourse.Title, newCourse.Details, newCourse.Trainer)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " +
						courseId))
				} else {
					// update course
					EditRecord(db, courseId, newCourse.Title, newCourse.Details, newCourse.Trainer)
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
