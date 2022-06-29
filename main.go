package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*"))
var id = 0

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
	http.HandleFunc("/create", Create)
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

func Start(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "start", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login", nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	values := map[string]string{
		"firstName":    "golang1",
		"lastName":     "go2testx",
		"selectedGame": "go3",
		"nickName":     "go4",
		"email":        "go5",
		"password":     "go6",
		"userImage":    "gox",
		"bibliography": "goy"}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}
	//endpoint := "https://gettingbetterapi.azurewebsites.net/api/v1/coaches"
	endpoint := "Dont Click"
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
	endpoint := "https://ksero.herokuapp.com/api/v1/users/auth/get-all"
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
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
