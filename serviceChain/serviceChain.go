// Core application for Service Record Blockchain
package main

import (
	"fmt"
	"time"

	"github.com/mmd93ee/ou-tm470/dataPersist"
	"github.com/mmd93ee/ou-tm470/web"
)

var blockchain []Block                       // Core Blockchain
var Index int                                // Block Index
var persistentFilename = "md5589_blockchain" // What to persist the blockchain to disk as

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
	loadBlockchain()
	saveBlockchain()

	fmt.Println("Starting web server...") // Web Server for blockchain interaction
	web.ServerStart("8000")
}

func loadBlockchain() error {
	if err := dataPersist.Load(persistentFilename, blockchain); err != nil {
		// Problem with lack of data file, lets create a genesis and return err
		fmt.Println("Error: " + err.Error() + ".  Loading blockchain failed, generating Genesis.")
		blockchain = append(blockchain, generateGenesisBlock())
		return err
	}
	return nil
}

func saveBlockchain() {
	fmt.Println("Persisting blockchain as " + persistentFilename)
	dataPersist.Save(persistentFilename, blockchain)
}

// generateGenesisBlock will create the first block
func generateGenesisBlock() Block {

	var genesisBlock Block
	var genesisRecord ServiceRecord

	// Set the values
	genesisBlock.Index = 1
	genesisBlock.Hash = "0"
	genesisBlock.PrevHash = "0"
	genesisBlock.Timestamp = time.Now().String()
	genesisRecord.Identifier = 0
	genesisRecord.Data = "seeded data"
	genesisBlock.Record = genesisRecord

	return genesisBlock
}
