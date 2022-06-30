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
	"strconv"
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

type User struct {
	Id       int
	Username string
	Password string
	Tokens   int
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
	idReceiver := r.URL.Query().Get("id")
	if id == "-1" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	fmt.Println(id + " Sends to " + idReceiver + " " + r.FormValue("tokens") + " tokens")

	// GET BLOCKS
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
	// END GET BLOCKS

	// GET USERS
	endpointUser := "https://go-project-backend.herokuapp.com/api/v1/users"
	resp2, err := http.Get(endpointUser)

	if err != nil {
		log.Fatal(err)
	}

	defer resp2.Body.Close()

	body, err := ioutil.ReadAll(resp2.Body)

	if err != nil {
		log.Fatal(err)
	}

	var userResponse []UserResponse
	errUnmarshal := json.Unmarshal(body, &userResponse)
	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
	}
	log.Printf("%v", userResponse)

	tokensInt, err := strconv.Atoi(r.FormValue("tokens"))
	if err != nil {
		log.Fatal(err)
	}
	var userSender UserResponse
	var userReceiver UserResponse
	canDoTransaction := false

	for i, s := range userResponse {
		fmt.Println(i, s.Id)
		if fmt.Sprint(s.Id) == id && s.Tokens >= tokensInt {
			userSender = s
			canDoTransaction = true
		}
		if fmt.Sprint(s.Id) == idReceiver {
			userReceiver = s
		}
	}

	if !canDoTransaction || userSender.Id == userReceiver.Id {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	userSender.Tokens -= tokensInt
	userReceiver.Tokens += tokensInt
	// END GET USERS
	// CREATE BLOCK
	chain.AddBlock(id + " sends to " + idReceiver + ", " + r.FormValue("tokens") + " tokens")

	prevHash := fmt.Sprintf("%x", chain.blocks[len(chain.blocks)-1].PrevHash)
	data := string(chain.blocks[len(chain.blocks)-1].Data)
	hash := fmt.Sprintf("%x", chain.blocks[len(chain.blocks)-1].Hash)

	// POST BLOCK
	values := map[string]string{
		"hash":     hash,
		"data":     data,
		"prevHash": prevHash}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	endpoint := "https://go-project-backend.herokuapp.com/api/v1/blocks"
	resp, err := http.Post(endpoint, "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	// END POST BLOCK
	// PUT USER SENDER

	// initialize http client
	client := &http.Client{}

	// marshal User to json
	json2, err := json.Marshal(userReceiver)
	if err != nil {
		panic(err)
	}
	json, err := json.Marshal(userSender)
	if err != nil {
		panic(err)
	}

	urlSender := "https://go-project-backend.herokuapp.com/api/v1/users/" + id
	fmt.Println(urlSender)

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, urlSender, bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp3, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp3.StatusCode)

	// END PUT USER SENDER
	// PUT USER RECEIVER

	// marshal User to json

	// set the HTTP method, url, and request body
	urlReceiver := "https://go-project-backend.herokuapp.com/api/v1/users/" + idReceiver
	req2, err := http.NewRequest(http.MethodPut, urlReceiver, bytes.NewBuffer(json2))
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req2.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp4, err := client.Do(req2)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp4.StatusCode)
	// END PUT USER RECEIVER

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Login(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login", nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register", nil)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// POST USERS
	values := map[string]string{
		"username": r.FormValue("username"),
		"password": r.FormValue("password"),
		"tokens":   r.FormValue("tokens")}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	endpoint := "https://go-project-backend.herokuapp.com/api/v1/users"
	resp, err := http.Post(endpoint, "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}
	// END POST USERS

	// Read post response
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

	// GET BLOCKS
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
	// END GET BLOCKS

	// CREATE BLOCK
	chain.AddBlock("0" + " sends to " + id + ", " + r.FormValue("tokens") + " tokens")

	prevHash := fmt.Sprintf("%x", chain.blocks[len(chain.blocks)-1].PrevHash)
	data := string(chain.blocks[len(chain.blocks)-1].Data)
	hash := fmt.Sprintf("%x", chain.blocks[len(chain.blocks)-1].Hash)

	// POST BLOCK
	values2 := map[string]string{
		"hash":     hash,
		"data":     data,
		"prevHash": prevHash}
	json_data2, err := json.Marshal(values2)

	if err != nil {
		log.Fatal(err)
	}

	endpoint2 := "https://go-project-backend.herokuapp.com/api/v1/blocks"
	resp2, err := http.Post(endpoint2, "application/json",
		bytes.NewBuffer(json_data2))

	if err != nil {
		log.Fatal(err)
	}

	defer resp2.Body.Close()
	// END POST BLOCK

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func ConfirmCredentials(w http.ResponseWriter, r *http.Request) {
	// GET USERS
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
	// END GET USERS
	for i, s := range userResponse {
		fmt.Println(i, s.Username)
		if s.Username == r.FormValue("username") && s.Password == r.FormValue("password") {
			id = fmt.Sprint(s.Id)
		}
	}
	fmt.Println(r.FormValue("username"))

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
