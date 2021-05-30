package handlers

import (
	"database/sql"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	cLog "github.com/heartziq/course-mgmt/log"
)

const FORMAT = "2006-01-02 15:04:05"

var (
	db       *sql.DB
	sqlError error
)

func init() {
	db, sqlError = sql.Open("mysql", "user1:password@tcp(127.0.0.1:63579)/my_db")
}

type course struct {
	ID string `json:"id"`

	Title string `json:"Title"`

	Details string `json:"Details"`

	Trainer string `json:"Trainer"`
}

type user struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	APIKey   string `json:"api_key"`
	Count    int    `json:"count"`
	Expiry   string `json:"expiry"`
}

func IsKeyValid(key string) (string, error) {
	if key == "" {
		return "", errors.New("empty APIKey provided")
	}
	query := "SELECT id, expiry from my_db.Users where api_key=?"
	result := db.QueryRow(query, key)

	var u user
	if err := result.Scan(&u.Id, &u.Expiry); err != nil {
		cLog.Info.Printf("error scanning %v", err)
		return "", errors.New("Invalid API_KEY")
	}

	// Check expiry
	timeWhenKeyExpires, _ := time.Parse(FORMAT, u.Expiry)
	// time.Now().Add(time.Hour * 24 * 10) // set this as param to force expire token
	if timeWhenKeyExpires.Before(time.Now()) {
		return "", errors.New("Expired API_KEY")
	}

	return u.Id, nil
}

func PopulateNewUser(u *user, pwd string) error {

	hashedPwd, err := hashPassword(pwd)
	if err != nil {

		return errors.New("Error creating (hashing) password")
	}

	u.Id = uuid.NewString()
	u.Password = hashedPwd
	u.APIKey = genAPIKey(hashedPwd)

	// set expiry date
	u.Expiry = time.Now().Add(time.Hour * 24 * 7).Format(FORMAT)

	return nil
}

func GetOneUser(username string) (*user, error) {
	c := user{}
	row := db.QueryRow("SELECT id, password, api_key From my_db.Users WHERE username=?", username)
	if err := row.Scan(&c.Id, &c.Password, &c.APIKey); err != nil {
		return nil, errors.New("Error scanning")
	}
	return &c, nil
}

func InsertUser(u *user) (int64, error) {

	query := "INSERT INTO my_db.Users VALUES (?,?,?,?,?,?)"
	result, err := db.Exec(query, u.Id, u.UserName, u.Password, u.APIKey, u.Count, u.Expiry)

	if err != nil {
		cLog.Error.Printf("Error: %v\n", err)
		return 0, errors.New("error inserting into my_db.Users")
	}

	rows, _ := result.RowsAffected()
	cLog.Info.Printf("%d row(s) affected\n", rows)

	lastInsertedId, err := result.LastInsertId()

	if err != nil {
		return 0, errors.New("Error getting last inserted ID")
	}

	return lastInsertedId, nil

}

func EditUser(pwd, ApiKey, expiry string) {
	results, err := db.Exec("UPDATE my_db.Users SET api_key=?, expiry=? WHERE password=?;", ApiKey, expiry, pwd)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		cLog.Info.Println(rows)
	}
}

// self-generate course id; called when user input ID that is too long
// Provides another layer of protection (apart from query param specification)
// courseId = 2 Capital letters + 4 random digit number
// e.g. CS1012 or AT3301
func generateCourseId(usrInput *string) (courseId string) {

	rand.Seed(time.Now().Unix())
	fourRandomDigits := strconv.Itoa(1000 + rand.Intn(9999-1000))

	twoCapLetters := (*usrInput)[:1] + (*usrInput)[len((*usrInput))-1:]
	courseId = strings.ToUpper(twoCapLetters) + fourRandomDigits

	return
}

func InsertRecord(ID *string, Title string, Details string, Trainer string) {
	results, err := db.Exec("INSERT INTO my_db.Courses VALUES (?,?,?,?)", ID, Title, Details, Trainer)
	defer func() {
		if r := recover(); r != nil {
			if errorMsg, ok := r.(error); ok {
				if cLog.GetErrorCode(errorMsg.Error()) == cLog.ERR_VARCHAR_TOO_LONG {
					cLog.Error.Printf("Error: %s\n", "Column input too long. DB column VARCHAR(7)")
					// Generate an acceptable column id
					*ID = generateCourseId(&Title)
					InsertRecord(ID, Title, Details, Trainer)
				}
			}

		}
	}()
	if err != nil {
		cLog.Error.Printf("Error: %v\n", err)
		panic(err)

	} else {
		rows, _ := results.RowsAffected()
		cLog.Info.Println(rows)

	}
}

func DeleteRecord(ID string) (bool, error) {
	results, err := db.Exec("DELETE FROM my_db.Courses WHERE ID=?;", ID)
	if err != nil {
		return false, errors.New("Error deleting record")
	}

	rows, _ := results.RowsAffected()
	cLog.Info.Printf("%d row(s) affected\n", rows)
	if rows < 1 {
		return false, errors.New("record not found")
	}

	return true, nil

}

func EditRecord(ID string, Title string, Details string, Trainer string) {
	results, err := db.Exec("UPDATE my_db.Courses SET Title=?, Details=?, Trainer=? WHERE ID=?;", Title, Details, Trainer, ID)
	if err != nil {
		panic(err)
	} else {
		rows, _ := results.RowsAffected()
		cLog.Info.Println(rows)
	}
}

func GetOneCourse(ID string) (*course, error) {
	c := course{}
	row := db.QueryRow("SELECT * From my_db.Courses WHERE ID=?", ID)
	if err := row.Scan(&c.ID, &c.Title, &c.Details, &c.Trainer); err != nil {
		return nil, errors.New("Error scanning")
	}
	return &c, nil
}

func GetRecords() (courses []*course) {

	results, err := db.Query("Select * FROM my_db.Courses")

	if err != nil {

		panic(err.Error())

	}

	for results.Next() {

		// map this type to the record in the table

		var c course

		err = results.Scan(&c.ID, &c.Title,

			&c.Details, &c.Trainer)

		if err != nil {

			panic(err.Error())

		}

		courses = append(courses, &c)

	}
	return
}
