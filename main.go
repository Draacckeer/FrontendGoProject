package main

import (
	"fmt"
	"net/http"
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
	http.HandleFunc("/create", Create)
	fmt.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)
}

func Start(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hola Dracker")
	templates.ExecuteTemplate(w, "start", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create", nil)
}
