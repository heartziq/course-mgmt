// curl http://localhost:5000/api/v1/courses
// curl http://localhost:5000/api/v1/courses/AND101?key=2c78afaf-97da-4816-bbee-9ad239abb296

// curl -H "Content-Type: application/json" -X POST http://localhost:5000/api/v1/courses/AND101 -d "{\"title\":\"Android Development\"}"
// curl -H "Content-Type: application/json" -X PUT http://localhost:5000/api/v1/courses/IOT101 -d "{\"title\":\"Internet of Things Development\"}"

// curl -X DELETE http://localhost:5000/api/v1/courses/IOT101?key=2c78afaf-97da-4816-bbee-9ad239abb296

// curl -H "Content-Type: application/json" -X PUT http://localhost:5000/api/v1/courses/AND101 -d "{\"title\":\"How to Cook a Beef\"}"

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseURL = "http://localhost:5000/api/v1/courses"

const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

type course struct {
	ID string `json:"id"`

	Title string `json:"Title"`

	Details string `json:"Details"`

	Trainer string `json:"Trainer"`
}

func getCourse(code string) {

	url := baseURL

	if code != "" {

		url = baseURL + "/" + code + "?key=" + key

	}

	response, err := http.Get(url)

	if err != nil {

		fmt.Printf("The HTTP request failed with error %s\n", err)

	} else {

		data, _ := ioutil.ReadAll(response.Body)

		fmt.Println(response.StatusCode)

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

		fmt.Println(response.StatusCode)

		fmt.Println(string(data))

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

func main() {

	var usrInput int
	for {
		fmt.Print("Enter choice (1-4): ")
		fmt.Scanln(&usrInput)
		var c string
		switch usrInput {
		case 1:

			fmt.Print("Enter course code: ")
			fmt.Scanln(&c)
			getCourse(c)
		case 2:

			// {
			// 	"ID": "78g2ss",
			// 	"Title": "How to be Software Engineer",
			// 	"Details": "a renewed ASPD financial advisor",
			// 	"Trainer": "Dr. Nasrullah Sejeong"
			// }
			fmt.Println("Enter course title: ")
			fmt.Scanln(&c)
			// jsonData := map[string]string{"title": c}

			var code string
			fmt.Println("Enter course code: ")
			fmt.Scanln(&code)

			// Create struct obj
			newCourse := course{
				ID:      code,
				Title:   c,
				Details: "a renewed ASPD financial advisor",
				Trainer: "Dr. Nasrullah Sejeong",
			}

			addCourse(code, newCourse)
			// case 3:
			// 	fmt.Println("Enter (NEW) course title: ")
			// 	fmt.Scanln(&c)
			// 	var code string
			// 	fmt.Println("Enter course code: ")
			// 	fmt.Scanln(&code)
			// 	updateCourse(code, map[string]string{"Title": c})
			// case 4:
			// 	fmt.Println("Enter course code: ")
			// 	fmt.Scanln(&c)
			// 	deleteCourse(c)
		}

		// flush c
		c = ""
	}

}
