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

var blockchain []Block                         // Core Blockchain
var index int                                  // Block Index
var persistentFilename = "./md5589_blockchain" // What to persist the blockchain to disk as

// Block struct service chain
type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string
	Event     ServiceEvent // The Service Record
}

// ServiceRecord to represent the service record data itself
type ServiceEvent struct {
	Identifier         int
	EventDetails       EventDescription // Describe what has happened
	EventAuthorisor    string           // Who allowed this event to happen.  Store as a public key.
	PerformedOnVehicle Vehicle          // What was the work performed on
	PerformedBy        Garage           // Who did the work
}

// Show the milage and the collection of events associated with it.
type EventDescription struct {
	EventItem     []EventType
	VehicleMilage int
}

// An event identifier and the associated description of the event.  This is a
// predefined set of activities that a garage can select from.
type EventType struct {
	EventId          int
	EventDescription string
}

// The vehicle the activities occur on.
type Vehicle struct {
	V5c                 string   // Reference v5 number, do not record further details
	VehicleMake         string   // Plain text at the moment, should probably lookup
	VehicleModel        string   // Plain text at the moment, should probably lookup
	VehicleColour       []string // Set as a slice with latest colour at last entry on array
	VehicleRegistration []string // Current registration is last entry on array
}

// Garages
type Garage struct {
	GarageId int
	Owner    string
	Name     string
	Location []string
	Type     string
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
	err := dataPersist.Load(persistentFilename, &blockchain)
	if err != nil {
		// Problem with lack of data file, lets create a genesis and return err
		log.Println("INFO: loadBlockchain(): Loading blockchain failed, generating Genesis.  " + err.Error())
		blockchain = append(blockchain, generateGenesisBlock())
		blockChainSize := strconv.Itoa(len(blockchain))
		log.Println("INFO: serviceChain.loadBlockchain(): Created genesis and added to blockchain, now of size " + blockChainSize)
		return nil
	}
	blockChainSize := strconv.Itoa(len(blockchain))
	log.Println("INFO: serviceChain.loadBlockchain(): Loaded blockchain, total records = " + blockChainSize)
	return nil
}

func saveBlockchain() error {
	log.Println("INFO: serviceChain.saveBlockchain(): Persisting blockchain as " + persistentFilename)
	return dataPersist.Save(persistentFilename, blockchain)
}

// generateGenesisBlock will create the first block
func generateGenesisBlock() Block {

	var genesisBlock Block
	var genesisRecord ServiceEvent
	var genesisRecordEventDescription EventDescription
	var genesisRecordEventDescriptionType EventType
	var genesisRecordVehicle Vehicle
	var genesisRecordGarage Garage

	// Seed values for Garage, Vehicle, EventType and EventDescription
	genesisRecordGarage.GarageId = 0
	genesisRecordGarage.Location[0] = "genesis location"
	genesisRecordGarage.Name = "genesis inc."
	genesisRecordGarage.Owner = "genesis and co."
	genesisRecordGarage.Type = "main dealer"

	genesisRecordVehicle.V5c = "63ne515"
	genesisRecordVehicle.VehicleColour[0] = "colour"
	genesisRecordVehicle.VehicleMake = "genesis make"
	genesisRecordVehicle.VehicleModel = "genesis model"
	genesisRecordVehicle.VehicleRegistration[0] = "GEN 351 S"

	genesisRecordEventDescriptionType.EventId = 0
	genesisRecordEventDescriptionType.EventDescription = "genesis event"

	genesisRecordEventDescription.EventItem[0] = genesisRecordEventDescriptionType
	genesisRecordEventDescription.VehicleMilage = 10000000

	// Set the values for the Block
	genesisBlock.Index = 1
	genesisBlock.Hash = "0"
	genesisBlock.PrevHash = "0"
	genesisBlock.Timestamp = time.Now().String()
	genesisBlock.Event = genesisRecord

	blockString, err := json.MarshalIndent(genesisBlock, "", "\t")
	if err != nil {
		log.Println("INFO: serviceChain.createGenesisBlock(): Problem creating the JSON output of the genesis block.  Continuing...")
	}

	log.Println("INFO: serviceChain.generateGenesisBlock(): Created block with contents: " + string(blockString))

	return genesisBlock
}
