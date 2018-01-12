package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"apla_test_work/model"

	"github.com/gorilla/mux"
)

var (
	// Auth - yet logined users
	Auth map[string]string
	// Work - work
	Work map[string]int32
)

func main() {
	Auth = make(map[string]string)
	Work = make(map[string]int32)

	err := model.GormInit()
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	u := model.User{
		Login:      "eugene",
		Pass:       "123",
		WorkNumber: 123,
	}
	u.CreateWhenNotExits()

	defer func() {
		model.GormClose()
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", RootPage)
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/login/pass", ChangePass).Methods("POST")
	r.HandleFunc("/do_work", DoWork).Methods("POST")

	// http.Handle("/", r)
	log.Println("Start server on http://localhost:5000")
	if err := http.ListenAndServe("localhost:5000", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// RootPage - main page
func RootPage(w http.ResponseWriter, r *http.Request) {
	body := `<!DOCTYPE html>
		<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<meta name="theme-color" content="#375EAB">

			<title>main page</title>
		</head>
		<body>
			<h3>Page body and some more content</h3>
			<p>Hi there</p>
		</body>
		</html>`
	writeResponse(w, body, http.StatusOK)
}

// Login - end point to user login
func Login(w http.ResponseWriter, r *http.Request) {
	login := strings.Trim(r.FormValue("login"), " ")
	pass := strings.Trim(r.FormValue("pass"), " ")
	log.Println("login: ", login, " pass: ", pass)

	// Check params
	if login == "" || pass == "" {
		writeResponse(w, "Login and password required\n", http.StatusBadRequest)
		return
	}

	// Already authorized
	if savedPass, OK := Auth[login]; OK && savedPass == pass {
		writeResponse(w, "You are already authorized\n", http.StatusOK)
		return
	} else if OK && savedPass != pass {
		// it is not neccessary
		writeResponse(w, "Wrong pass\n", http.StatusBadRequest)
		return
	}

	user := model.User{}
	err := user.Get(login, pass)
	if err == nil {
		Auth[login], Work[login] = pass, user.WorkNumber
		writeResponse(w, "Succesfull authorization\n", http.StatusOK)
		return
	}

	writeResponse(w, "User with same login not found\n", http.StatusNotFound)
}

// ChangePass - end point to change password
func ChangePass(w http.ResponseWriter, r *http.Request) {
	login := strings.Trim(r.FormValue("login"), " ")
	pass := strings.Trim(r.FormValue("pass"), " ")

	// Check params
	if login == "" || pass == "" {
		writeResponse(w, "Login and password required\n", http.StatusBadRequest)
		return
	}

	newPass := strings.Trim(r.FormValue("newPass"), " ")
	if newPass == "" {
		writeResponse(w, "newPass required\n", http.StatusBadRequest)
		return
	}

	if Auth[login] != pass {
		writeResponse(w, "Wrong pass\n", http.StatusBadRequest)
		return
	}

	user := model.User{}
	err := user.Get(login, pass) // в Auth можно сохранять id, чтобы не делать этот запрос
	log.Println("user", user)
	if err != nil {
		writeResponse(w, "Something wrong\n", http.StatusInternalServerError)
		return
	}

	user.Pass = newPass
	err = user.Save()
	if err != nil {
		writeResponse(w, "Error update password in database\n", http.StatusInternalServerError)
		return
	}

	Auth[user.Login] = user.Pass

	writeResponse(w, "Password changed\n", http.StatusOK)
}

// DoWork - do user work
func DoWork(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		BigNumber int64 `json:"number"`
		// SmallNumber int32  `json:"smallNumber"`
		Text string `json:"text"`
	}

	var value resp
	login := r.FormValue("login")
	if Work[login] <= 0 { // idk what is it
		writeResponse(w, "Work not found\n", http.StatusBadRequest)
		return
	}

	err := json.Unmarshal([]byte(r.FormValue("value")), &value)
	if err != nil {
		log.Println("Error parse JSON ", err)
		writeResponse(w, "Error parse json\n", http.StatusBadRequest)
		return
	}

	v := reflect.ValueOf(value)
	tv := reflect.TypeOf(value)
	var res string

	for i := 0; i < v.NumField(); i++ {
		res += tv.Field(i).Name + "  " + reverse(v.Field(i)) + "\n"
	}
	writeResponse(w, res, http.StatusOK)
}

func reverse(val reflect.Value) string {
	switch val.Type().String() {
	case "int64":
		v := val.Interface().(int64)
		return strconv.FormatInt(9223372036854775807-v, 10)
	case "int32":
		v := val.Interface().(int32)
		return strconv.FormatInt(int64(2147483647-v), 10)
	case "string":
		v := val.Interface().(string)

		runes := []rune(v)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}
	return ""
}

func writeResponse(w http.ResponseWriter, data string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(data))
}
