package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	/*response, err := http.Get(`https://gettingbetterapi.azurewebsites.net/api/v1/coaches`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.Body)
	bytes, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		fmt.Println(errRead)
	}
	fmt.Println(string(bytes))*/

	http.HandleFunc("/", Start)
	http.HandleFunc("/create-user", CreateUser)
	http.HandleFunc("/confirm-credentials", ConfirmCredentials)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	fmt.Println("Servidor corriendo...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Money    int    `json:"money"`
}

type Data struct {
	Id           string
	UserResponse []UserResponse
}

func Start(w http.ResponseWriter, r *http.Request) {

	endpoint := "https://go-project-backend.herokuapp.com/api/v1/users"
	resp, err := http.Get(endpoint)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var data Data
	errUnmarshal := json.Unmarshal(body, &data.UserResponse)
	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
	}

	data.Id = r.URL.Query().Get("id")
	if data.Id == "" {
		data.Id = "-1"
	}

	log.Println(r.FormValue("username"))

	templates.ExecuteTemplate(w, "start", data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login", nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register", nil)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	values := map[string]string{
		"username": r.FormValue("username"),
		"password": r.FormValue("password"),
		"money":    "0"}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}
	// Insert The Backend Url Here
	endpoint := "https://go-project-backend.herokuapp.com/api/v1/users"
	resp, err := http.Post(endpoint, "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res["json"])

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func ConfirmCredentials(w http.ResponseWriter, r *http.Request) {
	endpoint := "https://go-project-backend.herokuapp.com/api/v1/users"
	resp, err := http.Get(endpoint)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
	var userResponse []UserResponse
	errUnmarshal := json.Unmarshal(body, &userResponse)
	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
	}
	log.Printf("%v", userResponse)
	id := -1
	for i, s := range userResponse {
		fmt.Println(i, s.Username)
		if s.Username == r.FormValue("username") && s.Password == r.FormValue("password") {
			id = s.Id
		}
	}
	fmt.Println(r.FormValue("username"))

	http.Redirect(w, r, "/?id="+strconv.Itoa(id), http.StatusMovedPermanently)
}
