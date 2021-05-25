package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	cLog "github.com/heartziq/course-mgmt/log"
)

var (
	key     string
	scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	msg, _                 = os.ReadFile("main-menu.txt")
)

const ERR_NO_API_KEY = "Forbidden Access - empty string provided"

const baseURL = "http://localhost:5000/api/v1/courses"

type course struct {
	ID string `json:"id"`

	Title string `json:"Title"`

	Details string `json:"Details"`

	Trainer string `json:"Trainer"`
}

func getCourse(code string, key string) {

	url := baseURL

	if code != "" {

		url = baseURL + "/" + code + "?key=" + key
	}

	response, err := http.Get(url)

	if err != nil {

		fmt.Printf("The HTTP request failed with error %s\n", err)

	} else {

		data, _ := ioutil.ReadAll(response.Body)

		fmt.Printf("STATUS %d\n", response.StatusCode)

		fmt.Println(string(data))

		response.Body.Close()

	}

}

func addCourse(code string, c course) {

	jsonValue, _ := json.Marshal(c)

	response, err := http.Post(baseURL+"/"+code+"?key="+key,

		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {

		fmt.Printf("The HTTP request failed with error %s\n", err)

	} else {

		data, _ := ioutil.ReadAll(response.Body)

		fmt.Printf("STATUS %d\n", response.StatusCode)

		if string(data) == ERR_NO_API_KEY {
			fmt.Println("Please provide API_KEY. Login/Register to attain them")
		} else {

			fmt.Println(string(data))
		}

		response.Body.Close()

	}

}

func updateCourse(code string, jsonData map[string]string) {

	jsonValue, _ := json.Marshal(jsonData)

	request, err := http.NewRequest(http.MethodPut,

		baseURL+"/"+code+"?key="+key,

		bytes.NewBuffer(jsonValue))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {

		fmt.Printf("The HTTP request failed with error %s\n", err)

	} else {

		data, _ := ioutil.ReadAll(response.Body)

		fmt.Println(response.StatusCode)

		fmt.Println(string(data))

		response.Body.Close()

	}

}

func deleteCourse(code string) {

	request, err := http.NewRequest(http.MethodDelete,

		baseURL+"/"+code+"?key="+key, nil)

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {

		fmt.Printf("The HTTP request failed with error %s\n", err)

	} else {

		data, _ := ioutil.ReadAll(response.Body)

		fmt.Println(response.StatusCode)

		fmt.Println(string(data))

		response.Body.Close()

	}

}

func getUserInput(sc *bufio.Scanner, promptMsg string) string {
	fmt.Print(promptMsg)

	sc.Scan()
	line := sc.Text()

	return strings.TrimSpace(line)
}

func printBanner(bannerMsg string) {

	fmt.Println(strings.Repeat("-", len(bannerMsg)*2))
	fmt.Printf("%s%s%s\n", strings.Repeat(" ", len(bannerMsg)/2), bannerMsg, strings.Repeat(" ", len(bannerMsg)/2))
	fmt.Println(strings.Repeat("-", len(bannerMsg)*2))
}

func registerUser() {
	printBanner("REGISTER NEW USER")
	username := getUserInput(scanner, "Enter Username: ")
	password := getUserInput(scanner, "Enter Password: ")
	data, err := json.Marshal(map[string]string{"username": username, "password": password})

	if err != nil {
		cLog.Error.Printf("error parsing json req body: %v", err)
	} else {
		res, err := http.Post("http://localhost:5000/register", "application/json", bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}
		if res.StatusCode == http.StatusOK {
			fmt.Println("Successfully registered. Here is your API_KEY:")

			key = res.Cookies()[0].Value
			fmt.Printf("apiKey: %s\n", key)

		}
	}
}

func login(setNewKey bool) {
	url := "http://localhost:5000/login"
	printBanner("Login")
	username := getUserInput(scanner, "Username: ")
	pwd := getUserInput(scanner, "Password: ")
	data, err := json.Marshal(map[string]string{"username": username, "password": pwd})

	if err != nil {
		cLog.Error.Printf("error parsing json req body: %v", err)
		return
	}

	if setNewKey {
		url += "?NewKey=True"
	} else {
		url += "?NewKey=False"
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	if res.StatusCode == http.StatusAccepted {
		key = res.Cookies()[0].Value
		fmt.Printf("Your API_KEY is auto-included in cookies, but for your reference, here is your API_KEY: %s\n", key)
	} else {
		fmt.Println(res.StatusCode)
		fmt.Println("Incorrect User/Password Combo - Login Failed.")
	}
}

func main() {

	var usrInput int

	for {
		fmt.Print(string(msg))
		fmt.Scanln(&usrInput)

		switch usrInput {
		case 1:
			registerUser()
		case 2:
			login(false)
		case 3:
			// Print course detail //
			printBanner("Display Course Details")
			courseId := getUserInput(scanner, "Input course id [Enter to view all]: ")
			getCourse(courseId, key)
		case 4:
			// Add new course //
			printBanner("ADD NEW COURSE")
			code := getUserInput(scanner, "Course Code (2 Uppercase Character, followed by 4 digits e.g. JR3412): ")
			title := getUserInput(scanner, "Course Title: ")
			details := getUserInput(scanner, "Course Details: ")
			trainer := getUserInput(scanner, "Trainer: ")

			addCourse(code, course{"", title, details, trainer})
		case 5:
			// Edit Course //
			printBanner("EDIT COURSE DETAILS")

			code := getUserInput(scanner, "Select a course (Enter Course Code): ")
			title := getUserInput(scanner, "Input (NEW) Course Title - [ENTER] to skip: ")
			details := getUserInput(scanner, "Input (NEW) Course Details - [ENTER] to skip: ")
			trainer := getUserInput(scanner, "Input (NEW) Course Trainer - [ENTER] to skip: ")

			if (title + details + trainer) == "" {
				fmt.Printf("No new inputs detected - Course %s remains unchanged\n", code)
			} else {
				jsonData := map[string]string{"title": title, "details": details, "trainer": trainer}
				updateCourse(code, jsonData)
			}

		case 6:
			// Delete Course //
			printBanner("DELETE A COURSE")
			code := getUserInput(scanner, "Select a course (Enter Course Code): ")

			deleteCourse(code)
		case 7:
			fmt.Println("Login to renew your APIKey")
			login(true)
		default:
			os.Exit(0)
		}

		// flush c
		usrInput = 0
	}

}
