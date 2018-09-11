// Core application for Service Record Blockchain
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	//"github.com/mmd93ee/ou-tm470/dataPersist"
	//	"github.com/mmd93ee/ou-tm470/web"

	"fmt"
	"net/http"
	"strings"
)

// Core Blockchain and persistent data sets
var Blockchain []Block
var ValidGarages []Garage
var ValidVehicles []Vehicle
var ValidEvents []EventType
var lock sync.Mutex

// Data files and persistence variables
var persistentFilename = "./md5589_blockchain"      // What to persist the blockchain to disk as
var validGarageDataFile = "./md5589_garages.json"   // Where is the list of Garages
var validVehicleDataFile = "./md5589_vehicles.json" // Where is the list of Garages
var validEventDataFile = "./md5589_events.json"     // Where is the list of Garages

var index int                    // Block Index
var defaultListenerPort = "8000" // Default listener port

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
	Location string
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
	ServerStart("8000")
}

func loadBlockchain() error {
	err := fileToInterface(persistentFilename, &Blockchain)
	if err != nil {
		// Problem with lack of data file, lets create a genesis and return err
		log.Println("INFO: loadBlockchain(): Loading Blockchain failed, generating Genesis.  " + err.Error())
		Blockchain = append(Blockchain, generateGenesisBlock())
		blockchainSize := strconv.Itoa(len(Blockchain))
		log.Println("INFO: serviceChain.loadBlockchain(): Created genesis and added to Blockchain, now of size " + blockchainSize)
		return nil
	}
	blockchainSize := strconv.Itoa(len(Blockchain))
	log.Println("INFO: serviceChain.loadBlockchain(): Loaded Blockchain, total records = " + blockchainSize)

	if len(Blockchain) < 1 { // Blockchain is too small so is missing genesis data
		log.Println("INFO: Block is too small to hold data, seeding genesis block")
		Blockchain = append(Blockchain, generateGenesisBlock())
		return nil
	}
	return nil
}

func saveBlockchain() error {
	log.Println("INFO: serviceChain.saveBlockchain(): Persisting Blockchain as " + persistentFilename)

	return interfaceToFile(persistentFilename, Blockchain)
}

// Generate a new block based on the old block and new payload
func generateNewBlock(oldBlock Block, dataPayload string) (Block, error) {

	var newBlock Block
	timeNow := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = timeNow.String()

	newEvent, err := dataPayloadtoServiceEvent(dataPayload)

	if err != nil {
		log.Println("ERROR: Unable to convert data payload into ServiceEvent for new block generation.")
	}

	newBlock.Event = newEvent
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

// Load Garages, Vehicles and Events in order to provide baseline data sets
func loadBaseData() error {

	err := fileToInterface(validGarageDataFile, ValidGarages)

	if err != nil {
		log.Println("ERROR: Unable to load valid garage data")
		return err
	}

	err = fileToInterface(validVehicleDataFile, ValidVehicles)
	if err != nil {
		log.Println("ERROR: Unable to load valid vehicles data")
		return err
	}

	err = fileToInterface(validEventDataFile, ValidEvents)
	if err != nil {
		log.Println("ERROR: Unable to load valid events data")
		return err
	}
	return nil
}

// Convert a data payload from the web handler into a ServiceEvent
func dataPayloadtoServiceEvent(dataPayload string) (ServiceEvent, error) {

	var newServiceEvent ServiceEvent
	var newEventDescription EventDescription
	var newEventDescriptionType EventType
	var vehicle Vehicle
	var garage Garage

	log.Print(string(newServiceEvent.Identifier) + string(newEventDescription.VehicleMilage) + string(newEventDescriptionType.EventId) + string(vehicle.V5c) + string(garage.GarageId))

	return newServiceEvent, nil
}

// Hash function to take key block data and return a SHA256 hash as a string
func calculateHash(block Block) string {

	// Time and vehicle identifier (v5c) are the key block items to generate the hash
	record := string(string(block.Index) + block.Timestamp + block.Event.PerformedOnVehicle.V5c + block.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
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
	genesisRecordGarage.Location = "genesis location"
	genesisRecordGarage.Name = "genesis inc."
	genesisRecordGarage.Owner = "genesis and co."
	genesisRecordGarage.Type = "main dealer"

	genesisRecordVehicle.V5c = "63ne515"
	genesisRecordVehicle.VehicleColour = append(genesisRecordVehicle.VehicleColour, "starting colour")
	genesisRecordVehicle.VehicleMake = "genesis make"
	genesisRecordVehicle.VehicleModel = "genesis model"
	genesisRecordVehicle.VehicleRegistration = append(genesisRecordVehicle.VehicleRegistration, "GEN 351 S")

	genesisRecordEventDescriptionType.EventId = 0
	genesisRecordEventDescriptionType.EventDescription = "genesis event"

	genesisRecordEventDescription.EventItem = append(genesisRecordEventDescription.EventItem, genesisRecordEventDescriptionType)
	genesisRecordEventDescription.VehicleMilage = 10000000

	// Pull all the objects into ServiceEvent
	genesisRecord.EventAuthorisor = "Created by serviceChain as the Genesis Block"
	genesisRecord.EventDetails = genesisRecordEventDescription
	genesisRecord.Identifier = 0
	genesisRecord.PerformedBy = genesisRecordGarage
	genesisRecord.PerformedOnVehicle = genesisRecordVehicle

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

// ServerStart starts the web server on the specified TCP port.  Blank will default to 8000.
func ServerStart(port string) (string, error) {

	http.HandleFunc("/", defaultHandler) // Each call to "/" will invoke defaultHandler
	http.HandleFunc("/blockchain/view/", blockchainViewHandler)

	//log.Fatal(http.ListenAndServe("localhost:"+port, nil))
	return "Started on: " + port, http.ListenAndServe("localhost:"+port, nil)

}

// Default handler to catch-all
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Default Handler called from %s.  Please try alternative methods such as /blockchain/view/<id>", r.RemoteAddr)
}

// Handler to manage requests to /blockchain/ subchain
func blockchainViewHandler(w http.ResponseWriter, r *http.Request) {

	// Take the URL beyond /blockchain/ and split into request and value strings
	requestAction := strings.Split(r.URL.String(), "/")
	requestItem, err := strconv.Atoi(requestAction[3])
	if err != nil {
		log.Println("ERROR: Unable to convert argument to integer" + err.Error())
	}

	if requestItem == 0 { //Request item is invalid so display that blockID only
		blockString, err := json.MarshalIndent(Blockchain, "", "\t")
		if err != nil {
			log.Println("ERROR: Cannot print Blockchain")
		}
		fmt.Fprintf(w, "\n %s", blockString)
	} else {
		blockItemString, _ := json.MarshalIndent(Blockchain[requestItem], "", "\t") // Do nothing if index too high
		fmt.Fprintf(w, "\n %s.", blockItemString)
	}
}

// Handler to manage requests to /garage/ subchain
func garageViewHandler(w http.ResponseWriter, r *http.Request) {

	// Take the URL beyond /garage/ and split into request and value strings
	requestAction := strings.Split(r.URL.String(), "/")
	requestItem, err := strconv.Atoi(requestAction[3])
	if err != nil {
		log.Println("ERROR: Unable to convert argument to integer" + err.Error())
	}

	garageString, _ := json.MarshalIndent(ValidGarages[requestItem], "", "\t") // Do nothing if index too high
	fmt.Fprintf(w, "\n %s.", garageString)
}

/*
All data persistence functions and ways of writing and saving to the file
system to be located below here.  FUture state to move this to a different
package for portability.
*/

// Marshal interface to JSON for writing out to file
var Marshal = func(structIn interface{}) (io.Reader, error) {

	bytesIn, err := json.MarshalIndent(structIn, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bytesIn), nil
}

// Unmarshal interface for converting file to struct
var Unmarshal = func(reader io.Reader, structIn interface{}) error {
	return json.NewDecoder(reader).Decode(structIn)
}

// Save will convert the input interface (v) into a JSON formatted object on disk
func interfaceToFile(path string, structIn interface{}) error {

	// Create a lock and then defer the unlock until function exit
	lock.Lock()
	defer lock.Unlock()

	//Create os.File and defer close
	file, err := os.Create(string(path))
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a reader (via marshal) that we can use to copy into the writer
	reader, err := Marshal(structIn)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, reader)
	return err

}

// Load is used to convert a JSON (marshall output formatted) file to a struct (interface)
func fileToInterface(path string, structOut interface{}) error {

	// Lock and defer the unlock until function exit
	lock.Lock()
	defer lock.Unlock()

	log.Println("INFO: Loading " + path)
	fileOut, err := os.Open(path)

	if err != nil {
		log.Println("ERROR: fileToInterface() " + err.Error() + " while openning file " + path)
		return err
	}

	defer fileOut.Close()
	Unmarshal(fileOut, structOut)

	return nil
}
