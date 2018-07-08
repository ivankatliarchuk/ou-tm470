// Core application for Service Record Blockchain
package main

var Blockchain []Block         // Core Blockchain
var ServiceEntryIdentifier int // Service Entry ID value
var WebServer web              // Web server for application entry point

// Core block for service chain
type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string

	Record ServiceRecord
}

// A struct to represent the service record data itself
type ServiceRecord struct {
	Identifier int
	Data       string // PLACEHOLDER FOR DATA MODEL ENTRIES
}

func main() {

}

func generateGenesisBlock() {
	// STUB
}
