// Core application for Service Record Blockchain
package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/mmd93ee/ou-tm470/dataPersist"
	"github.com/mmd93ee/ou-tm470/web"
)

var Blockchain []Block                         // Core Blockchain
var index int                                  // Block Index
var persistentFilename = "./md5589_blockchain" // What to persist the blockchain to disk as

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
	if err := loadBlockchain(); err != nil {
		log.Fatalln(err)
	}
	if err := saveBlockchain(); err != nil {
		log.Fatalln(err)
	}

	log.Println("INFO: serviceChain.main(): Starting web server...") // Web Server for blockchain interaction
	web.ServerStart("8000")
}

func loadBlockchain() error {
	err := dataPersist.Load(persistentFilename, &Blockchain)
	if err != nil {
		// Problem with lack of data file, lets create a genesis and return err
		log.Println("INFO: loadBlockchain(): Loading blockchain failed, generating Genesis.  " + err.Error())
		Blockchain := append(Blockchain, generateGenesisBlock())
		blockChainSize := strconv.Itoa(len(Blockchain))
		log.Println("INFO: serviceChain.loadBlockchain(): Created genesis and added to blockchain, now of size " + blockChainSize)
		return nil
	}
	return nil
}

func saveBlockchain() error {
	log.Println("INFO: serviceChain.saveBlockchain(): Persisting blockchain as " + persistentFilename)
	return dataPersist.Save(persistentFilename, Blockchain)
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

	blockString, err := json.MarshalIndent(genesisBlock, "", "/t")

	log.Println("INFO: serviceChain.generateGenesisBlock(): Created block with contents: " + dataPersist.StructToJSONString(genesisBlock))

	return genesisBlock
}
