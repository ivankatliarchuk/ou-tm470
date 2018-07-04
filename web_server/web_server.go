package web_server

import (
	"encoding/json"
	"log"
  "io"
	"os"
  "net/http"
	"time"
	"fmt"
  "github.com/gorilla/mux"
  "github.com/joho/godotenv"
)

// Message structure for inbound PUT and POST
type Message struct {
  Payload string // Q: Is this the inbound format?
}


// Run a web server and listen on localhost
func Run() error {
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
  muxRouter.HandleFunc("/blockchain/", handleGetBlockchain).Methods("GET")
  muxRouter.HandleFunc("/block/new/", handleWriteBlockchain).Methods("POST")
  return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request){
  bytes, err := json.MarshalIndent("{Data : Stub}", "", " ")
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
