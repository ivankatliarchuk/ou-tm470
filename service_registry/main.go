package main

import (
  "log"
  "github.com/joho/godotenv"
  "web_server"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal(err)
  }

  go func() {
  //  generateGenesisBlock()
  }()
  log.Fatal(web_server.Run())
}
