package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*"))
var id string = "-1"

func main() {

	http.HandleFunc("/", Start)
	http.HandleFunc("/create-user", CreateUser)
	http.HandleFunc("/confirm-credentials", ConfirmCredentials)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/blockchain", Blockchain)
	http.HandleFunc("/send-tokens", SendTokens)
	fmt.Println("Servidor corriendo...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}

type BlockChain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func RequestBlock(hash []byte, data string, prevhash []byte) *Block {
	block := &Block{hash, []byte(data), prevhash}
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

func Genesis(genesis string) *Block {
	return CreateBlock(genesis, []byte{})
}

func InitBlockChain(genesis string) *BlockChain {
	return &BlockChain{[]*Block{Genesis(genesis)}}
}

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Tokens   int    `json:"tokens"`
}

type BlockResponse struct {
	Id       int    `json:"id"`
	Hash     string `json:"hash"`
	Data     string `json:"data"`
	PrevHash string `json:"prevHash"`
}

type Data struct {
	Id           string
	UserResponse []UserResponse
}

func Start(w http.ResponseWriter, r *http.Request) {
	urlBlocks := "https://go-project-backend.herokuapp.com/api/v1/blocks"
	respBlocks, err := http.Get(urlBlocks)

	if err != nil {
		log.Fatal(err)
	}

	defer respBlocks.Body.Close()

	bodyBlocks, err := ioutil.ReadAll(respBlocks.Body)

	if err != nil {
		log.Fatal(err)
	}

	var blockResponse []BlockResponse
	errUnmarshalBlocks := json.Unmarshal(bodyBlocks, &blockResponse)
	if errUnmarshalBlocks != nil {
		log.Fatal(errUnmarshalBlocks)
	}

	var chain BlockChain

	for i, s := range blockResponse {
		//fmt.Println(i, s.Data)
		if i == 0 {
			chain = *InitBlockChain(s.Data)
		} else {
			chain.AddBlock(s.Data)
		}
	}
	/*for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}*/

	/*chain := InitBlockChain("Genesis")
	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block after Genesis")
	chain.AddBlock("Third Block after Genesis")

	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}*/

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

	data.Id = id
	if data.Id == "" {
		data.Id = "-1"
	}

	log.Println(r.FormValue("username"))

	templates.ExecuteTemplate(w, "start", data)
}

func Blockchain(w http.ResponseWriter, r *http.Request) {
	urlBlocks := "https://go-project-backend.herokuapp.com/api/v1/blocks"
	respBlocks, err := http.Get(urlBlocks)

	if err != nil {
		log.Fatal(err)
	}

	defer respBlocks.Body.Close()

	bodyBlocks, err := ioutil.ReadAll(respBlocks.Body)

	if err != nil {
		log.Fatal(err)
	}

	var blockResponse []BlockResponse
	errUnmarshalBlocks := json.Unmarshal(bodyBlocks, &blockResponse)
	if errUnmarshalBlocks != nil {
		log.Fatal(errUnmarshalBlocks)
	}

	var chain BlockChain

	for i, s := range blockResponse {
		if i == 0 {
			chain = *InitBlockChain(s.Data)
		} else {
			chain.AddBlock(s.Data)
		}
	}
	/*for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}*/
	templates.ExecuteTemplate(w, "blockchain", blockResponse)
}

func SendTokens(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("tokens"))
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
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
		"tokens":   r.FormValue("tokens")}
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

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var userResponse UserResponse
	errUnmarshal := json.Unmarshal(body, &userResponse)
	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
	}

	id = fmt.Sprint(userResponse.Id)

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

	var userResponse []UserResponse
	errUnmarshal := json.Unmarshal(body, &userResponse)
	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
	}
	log.Printf("%v", userResponse)
	for i, s := range userResponse {
		fmt.Println(i, s.Username)
		if s.Username == r.FormValue("username") && s.Password == r.FormValue("password") {
			id = fmt.Sprint(s.Id)
		}
	}
	fmt.Println(r.FormValue("username"))

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
