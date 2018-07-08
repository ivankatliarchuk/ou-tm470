// Core application for Service Record Blockchain
package main

import (
	"fmt"

	"github.com/mmd93ee/ou-tm470/dataPersist"
	"github.com/mmd93ee/ou-tm470/web"
)

var blockchain []Block         // Core Blockchain
var serviceEntryIdentifier int // Service Entry ID value

// Block struct service chain
type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string

	Record ServiceRecord
}

// ServiceRecord to represent the service record data itself
type ServiceRecord struct {
	Identifier int
	Data       string // PLACEHOLDER FOR DATA MODEL ENTRIES
}

func main() {
	fmt.Println("Starting web server...") // Web Server for blockchain interaction
	web.ServerStart("8000")
}

func loadBlockchain() {

}

func saveBlockchain() {
	dataPersist.Save()
}

func generateGenesisBlock() {
	// STUB
}
