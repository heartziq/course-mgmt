package log

import (
	"io"
	"log"
	"os"
	"regexp"
)

var (
	Info, Error          *log.Logger
	ERR_VARCHAR_TOO_LONG errCode = 1406
)

type errCode uint16

func init() {
	file, err := os.Create("errors.log")
	if err != nil {
		log.Fatal(err)
	}

	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stdout, file), "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Gets an error code from a string - usually an error message
func GetErrorCode(errorMsg string) errCode {
	pattern := "^Error\\s{1}(\\d+):"
	r, err := regexp.Compile(pattern)

	if err != nil {
		log.Printf("Error compiling regexp: %v", err)
		return 0
	}

	result := r.FindStringSubmatch(errorMsg)
	if len(result) > 1 && result[len(result)-1] == "1406" {
		return ERR_VARCHAR_TOO_LONG
	}

	return 0

}
