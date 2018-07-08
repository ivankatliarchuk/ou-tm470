// Core application for Service Record Blockchain
package main

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
	server := webServer // Web Server for blockchain interaction

}

func generateGenesisBlock() {
	// STUB
}
