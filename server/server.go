package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/heartziq/course-mgmt/server/handlers"
)

// var router = mux.NewRouter()

func verifyAPIKey(next http.Handler) http.Handler {
	newHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		if key, exist := vars["key"]; exist {

			userId, err := handlers.IsKeyValid(key)
			if err == nil {
				log.Printf("[server] Key is valid, WElcome %s\n", userId)
				// update session with userId retrieved
				next.ServeHTTP(w, r)
				return
			} else {

				http.Error(w, fmt.Sprintf("Forbidden Access - %v", err), http.StatusForbidden)
				// can later change to redirect
				return

			}

		}

		http.Error(w, "Forbidden Access - No API_KEY provided", http.StatusForbidden)

	})

	return newHandlerFunc
}

func createServer() http.Handler {
	handler := mux.NewRouter()

	// router := mux.NewRouter()

	// Protected route - need to supply API_KEY
	subR := handler.NewRoute().Subrouter()

	// Course id must be in [A-Z]{2}\\d{4} format
	// i.e. 2 Capital letters + 4 Random Digits
	// e.g. FB4513 or XZ1142
	subR.
		Methods("GET", "PUT", "POST", "DELETE").
		Path("/api/v1/courses/{courseid:[A-Z]{2}\\d{4}}").
		Queries("key", "{key}").
		HandlerFunc(handlers.Course)

	subR.Use(verifyAPIKey)

	// Public API - No API key is necessary for this
	handler.HandleFunc("/api/v1/courses", handlers.AllCourses)

	// Only allow method POST - else, return Error 404 - Not Found
	handler.HandleFunc("/register", handlers.Register).Methods("POST")
	handler.HandleFunc("/login", handlers.Login).Methods("POST").Queries("NewKey", "{NewKey:True|False}")

	return handler
}

// NOTE: _test.go will not run main(), hence
// any router config i.e. routers, subrouters, pathprefix etc
// wont be added
func main() {

	// // router := mux.NewRouter()

	// // Protected route - need to supply API_KEY
	// subR := router.NewRoute().Subrouter()

	// // Course id must be in [A-Z]{2}\\d{4} format
	// // i.e. 2 Capital letters + 4 Random Digits
	// // e.g. FB4513 or XZ1142
	// subR.
	// 	Methods("GET", "PUT", "POST", "DELETE").
	// 	Path("/api/v1/courses/{courseid:[A-Z]{2}\\d{4}}").
	// 	Queries("key", "{key}").
	// 	HandlerFunc(handlers.Course)

	// subR.Use(verifyAPIKey)

	// // Public API - No API key is necessary for this
	// router.HandleFunc("/api/v1/courses", handlers.AllCourses)

	// // Only allow method POST - else, return Error 404 - Not Found
	// router.HandleFunc("/register", handlers.Register).Methods("POST")
	// router.HandleFunc("/login", handlers.Login).Methods("POST").Queries("NewKey", "{NewKey:True|False}")

	c := make(chan os.Signal)
	router := createServer()
	go func() {
		http.ListenAndServe(":5000", router)
	}()
	signal.Notify(c, os.Interrupt) // User abruptly quit - Ctrl-C
	<-c

	// Do some cleaning ups before shutdown
	// close connection
	log.Println("INterrupt.. closing connection...")
	log.Println("Doing cleanup...")
	log.Println("done cleaning up")

}
