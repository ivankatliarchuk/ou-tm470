package main

import (
	"crypto/sha256"
	"encoding/hex"
  "encoding/json"
	"log"
  "io"
	"os"
  "net/http"
	"time"
	"fmt"

  "github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
  "github.com/joho/godotenv"
)

// Data type for the block
type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string
	Payload   string
}

type Message struct {
  Payload string
}

var Blockchain []Block

// Hash function to take key block data and return a SHA256 hash as a string
func calculateHash(block Block) string {
	// Update this to reflect unique block has data
	record := string(string(block.Index) + block.Timestamp + block.Payload + block.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Generate a new block based on the old block and new payload
func generateBlock(oldBlock Block, newPayload string) (Block, error) {

	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Payload = newPayload
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

  return newBlock, nil
}

// Generate the Genesis Block
func generateGenesisBlock(){

  var genesis Block
  t := time.Now()

  genesis.Index = 0
  genesis.Timestamp = t.String()
  genesis.Payload = "payload"
  genesis.PrevHash = ""
  genesis.Hash = ""

  Blockchain = append(Blockchain, genesis)

  spew.Dump(genesis)

}

// Check that a block is valid by comparing to previous blocks attributes and
// validating stored hash is correct.  Return a boolean as to the validation
// being correct (true) or incorrect (false)
func isBlockValid(newBlock, oldBlock Block) bool {

	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// Overwrite the blockchain if a longer chain is presented from other nodes
func replaceChain(newChain []Block) {
	if len(newChain) > len(Blockchain) {
		Blockchain = newChain
	}
}

// Run a web server and listen on localhost
func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on: ", httpAddr)
  s := &http.Server{
    Addr:           ":" + httpAddr,
    Handler:        mux,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  // Start the listener and error out if a problem
  if err := s.ListenAndServe(); err != nil {
    return err
  }
  return nil
}

// Web Handlers for Get and Set Blockchain.  All GET requests
// will return the Blockchain and all POST requests will put a new
// block to the chain
func makeMuxRouter() http.Handler {
  muxRouter := mux.NewRouter()
  muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
  muxRouter.HandleFunc("/", handleWriteBlockchain).Methods("POST")
  return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request){
  bytes, err := json.MarshalIndent(Blockchain, "", " ")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  io.WriteString(w, string(bytes))
}

func handleWriteBlockchain(w http.ResponseWriter, r *http.Request){
  var m Message

  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&m); err != nil {
    respondWithJSON(w, r, http.StatusBadRequest, r.Body)
    return
  }
  defer r.Body.Close()

  newBlock, err := generateBlock(Blockchain[len(Blockchain) - 1], m.Payload)
  if err != nil {
    respondWithJSON(w, r, http.StatusInternalServerError, m)
    return
  }

  if isBlockValid(newBlock, Blockchain[len(Blockchain) - 1]){
    newBlockchain := append(Blockchain, newBlock)
    replaceChain(newBlockchain)
    spew.Dump(Blockchain)
  }

  respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}){
  response, err := json.MarshalIndent(payload, "", "  ")
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("HTTP 500: Internal Server Error"))
    return
  }

  w.WriteHeader(code)
  w.Write(response)
}

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal(err)
  }

  go func() {
    generateGenesisBlock()
  }()
  log.Fatal(run())

}
