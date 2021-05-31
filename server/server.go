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

func addAuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("jwt")
		if err != nil {
			log.Printf("Error: Cookie not found %v", err)
			log.Println("Auth will fail")

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized access. Pls login"))

			return
		}
		bearerToken := fmt.Sprintf("Bearer %s", token.Value)
		r.Header.Set("Authorization", bearerToken)

		next.ServeHTTP(w, r)
	})
}

func verifyAPIKey(next http.Handler) http.Handler {
	newHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("This prints...")
		if cookie, err := r.Cookie("jwt"); err == nil {
			// add to BEARER token http AUTHORIZATION
			log.Printf("[Cookie] name: %s, value: %s\n", cookie.Name, cookie.Value)
		} else {
			log.Println("no jwt cookie found")
		}
		vars := mux.Vars(r)

		if key, exist := vars["key"]; exist {

			userId, err := handlers.IsKeyValid(key)
			if err == nil {
				log.Printf("[server] Key is valid, WElcome %s\n", userId)
				// update session with userId retrieved
				next.ServeHTTP(w, r)
				return
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				if err.Error() == "Invalid API_KEY" {
					http.Error(w, fmt.Sprintf("error - %v", err), http.StatusBadRequest) //400
					// can later change to redirect
					return
				}
				http.Error(w, fmt.Sprintf("error - %v", err), http.StatusForbidden) // 403
				// can later change to redirect
				return

			}

		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.Error(w, "Forbidden Access - No API_KEY provided", http.StatusUnauthorized) // 401

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

	// draft handler - for learning/prac purposes
	handler.HandleFunc("/draft", handlers.TestDraftCookie) // get cookie

	// retrieve token
	handler.Handle("/dashboard/{id}", addAuthHeader(http.HandlerFunc(handlers.TestGetToken)))
	// handler.HandleFunc("/dashboard/{id}", handlers.TestGetToken)
	return handler
}

// NOTE: _test.go will not run main(), hence
// any router config i.e. routers, subrouters, pathprefix etc
// won't be added
func main() {

	c := make(chan os.Signal)
	router := createServer()
	go func() {
		http.ListenAndServeTLS(":5000", "cert/cert.pem", "cert/key.pem", router)
		// http.ListenAndServe(":5000", router)
	}()
	signal.Notify(c, os.Interrupt) // User abruptly quit - Ctrl-C
	<-c

	// Do some cleaning ups before shutdown
	// close connection
	log.Println("INterrupt.. closing connection...")
	log.Println("Doing cleanup...")
	log.Println("done cleaning up")

}
