package server

import (
	"crud/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUser insert a user on database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error reading request body!"))
		return
	}

	var user user
	if err = json.Unmarshal(requestBody, &user); err != nil {
		w.Write([]byte("Error converting user to struct!"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Error connecting to database!"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("insert into users (name, email) values ($1, $2) returning id")
	if err != nil {
		w.Write([]byte("Error creating statement!"))
		return
	}
	defer statement.Close()

	var userId int
	err = statement.QueryRow(user.Name, user.Email).Scan(&userId)
	if err != nil {
		w.Write([]byte("Error executing statement!"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User created successfully! Id: %d", userId)))
}

// GetUsers get all users stored in the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Error connecting with the database!"))
		return
	}
	defer db.Close()

	lines, err := db.Query("select * from users order by id")
	if err != nil {
		w.Write([]byte("Error getting users!"))
		return
	}
	defer lines.Close()

	users := []user{}
	for lines.Next() {
		var user user

		if err := lines.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Error getting user!"))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("Error converting users to JSON!"))
		return
	}
}

// GetUser get a specific user stored in the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error converting param to integer!"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Error connecting with database!"))
		return
	}

	line, err := db.Query("select * from users where id = $1", ID)
	if err != nil {
		w.Write([]byte("Error getting user!"))
		return
	}

	var user user
	if line.Next() {
		if err := line.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			if err != nil {
				w.Write([]byte("Error getting user!"))
				return
			}
		}
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write([]byte("Error converting user to JSON!"))
		return
	}
}

// UpdateUser change a user data in the database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error converting param to integer!"))
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error reading request body!"))
		return
	}

	var user user
	if err = json.Unmarshal(requestBody, &user); err != nil {
		w.Write([]byte("Error converting user to struct!"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Error connecting to database!"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update users set name = $1, email = $2 where id = $3")
	if err != nil {
		w.Write([]byte("Error creating statement!"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(user.Name, user.Email, ID); err != nil {
		w.Write([]byte("Error updating user!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
