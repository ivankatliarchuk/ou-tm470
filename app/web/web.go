package web

import (
  "crypto/sha256"
	"encoding/hex"
  "encoding/json"
	"log"
  "io"
	"os"
  "net/http"
	"time"

  "github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
  "github.com/joho/godotenv"


)

// Run the webserver listening on localhost
func Run() error {

  mux := makeMuxRouter()
  httpAddress := os.Getenv("ADDR")



}
