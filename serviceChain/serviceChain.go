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
var BlockchainLength int
var ValidGarages []Garage
var ValidVehicles []Vehicle
var ValidEvents []EventType
var vehicleMap map[string][]int
var handlerStrings []string
var filewritelock sync.Mutex
var blockchainwritelock sync.Mutex

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

	if err := loadBaseData(); err != nil {
		log.Fatalln(err)
	}

	log.Println("INFO: serviceChain.main(): Starting web server (port 8000)...") // Web Server for blockchain interaction
	ServerStart("8000")
}

func loadBlockchain() error {
	err := fileToInterface(persistentFilename, &Blockchain)
	if err != nil {
		// Problem with lack of data file, lets create a genesis and return err
		log.Println("INFO: loadBlockchain(): Loading Blockchain failed, generating Genesis.  " + err.Error())
		Blockchain = append(Blockchain, generateGenesisBlock())
		log.Printf("INFO: serviceChain.loadBlockchain(): Created genesis and added to Blockchain, now of size %s", string(len(Blockchain)))
		return nil
	}
	log.Printf("INFO: serviceChain.loadBlockchain(): Loaded Blockchain, total records = %s", strconv.Itoa(len(Blockchain)))

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

	err := fileToInterface(validGarageDataFile, &ValidGarages)

	if err != nil {
		log.Println("ERROR: Unable to load valid garage data")
		return err
	}

	err = fileToInterface(validVehicleDataFile, &ValidVehicles)
	if err != nil {
		log.Println("ERROR: Unable to load valid vehicles data")
		return err
	}

	err = fileToInterface(validEventDataFile, &ValidEvents)
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
	genesisRecord.Identifier = 1
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

func generateBlock(priorBlock Block, event ServiceEvent) (Block, error) {

	// New block and values
	var b Block
	log.Printf("INFO: Creating new block with previous Block ID set to: %s", strconv.Itoa(priorBlock.Index))
	b.Index = priorBlock.Index + 1
	b.PrevHash = priorBlock.Hash
	b.Event = event
	b.Timestamp = time.Now().String()
	b.Hash = calculateHash(b)

	log.Printf("INFO: Created new block with Index %s, hash %s at %s from eventID %s", strconv.Itoa(b.Index), string(b.Hash), string(b.Timestamp), strconv.Itoa(b.Event.Identifier))
	return b, nil
}

// Perform tests to validate if a specific block is valid by comparing it to its own values
func isBlockValid(prior Block, current Block) bool {
	return true
}

// Replace the current chain by adding a new block
func replaceChain(newBlock Block) bool {

	// make this thread safe
	blockchainwritelock.Lock()
	defer blockchainwritelock.Unlock()

	// Is the block valid and if so then append it to the chain
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)
		BlockchainLength = len(Blockchain)
		var registration string

		// Update vehicle lookups
		log.Printf("INFO: replaceChain(): Adding vehicle to vehicle and blocklookup map.")
		lastregindex := len(newBlock.Event.PerformedOnVehicle.VehicleRegistration)
		registration = newBlock.Event.PerformedOnVehicle.VehicleRegistration[lastregindex-1]
		registration = strings.ToUpper(registration)
		registration = strings.Replace(registration, " ", "", -1)

		blocklist := vehicleMap[registration]

		log.Printf("INFO: replaceChain(): REG %s, BLOCKLIST SIZE %s", registration, strconv.Itoa(len(blocklist)))

		log.Printf("INFO: replaceChain(): Captured registration: %s with %s previous entries", registration, strconv.Itoa(len(blocklist)))

		if (len(blocklist)) > 0 {
			log.Printf("INFO: replaceChain(): Vehicle been added before.  Appending new block id with value %s", strconv.Itoa(newBlock.Index))
			vehicleMap[registration] = append(blocklist, newBlock.Index)
		} else {
			newBlockSlice := []int{newBlock.Index}
			log.Printf("INFO: replaceChain(): created new list of blocks for registration %s, size is %s", registration, strconv.Itoa(len(newBlockSlice)))

			// vehicleMap not initialised so set it up
			if len(vehicleMap) < 1 {
				log.Printf("INFO: replaceChain(): vehicleMap is not initialised, size is %s", strconv.Itoa(len(vehicleMap)))
				vehicleMap = make(map[string][]int)
			}
			// Add the new vehicle to the map
			log.Printf("INFO: replaceChain(): Adding vehicle %s to new vehicleMap", registration)
			vehicleMap[registration] = newBlockSlice
		}

		log.Printf("INFO: replaceChain(): Added vehicle reg %s to block lookup table, blockid %s", registration, strconv.Itoa(newBlock.Index))

		log.Printf("INFO: replaceChain(): Appended new block, writing to disk with ID %s", strconv.Itoa(BlockchainLength))
		err := interfaceToFile("./saved_chains/md5589_blockchain_"+strconv.Itoa(BlockchainLength), Blockchain)
		if err != nil {
			log.Printf("ERROR: Unable to write blockchain to disk: %s", strconv.Itoa(BlockchainLength))
		}
		return true
	}
	return false
}

// ServerStart starts the web server on the specified TCP port.  Blank will default to 8000.
func ServerStart(port string) (string, error) {

	// List of view handlers
	handlerStrings = append(handlerStrings, "/", "/blockchain/view/<ID>", "/garage/view/<ID>", "serviceevent/add/", "/vehicle/view/<ID>")

	http.HandleFunc("/", defaultHandler) // Each call to "/" will invoke defaultHandler
	http.HandleFunc("/blockchain/view/", blockchainViewHandler)
	http.HandleFunc("/garage/view/", garageViewHandler)
	http.HandleFunc("/serviceevent/add/", writeServiceEventHandler)
	http.HandleFunc("/vehicle/view/", vehicleViewHandler)

	//log.Fatal(http.ListenAndServe("localhost:"+port, nil))
	return "Started on: " + port, http.ListenAndServe("localhost:"+port, nil)

}

// Default handler to catch-all
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "INFO: Default Handler called from %s.  Please try alternative methods such as \n %s", r.RemoteAddr, handlerStrings)
}

func writeServiceEventHandler(w http.ResponseWriter, r *http.Request) {

	var newServiceEvent ServiceEvent

	if err := json.NewDecoder(r.Body).Decode(&newServiceEvent); err != nil {
		log.Printf("ERROR: writeServiceEventHandler(): Unable to decode data payload: %s", err.Error())
		http.Error(w, "ERROR: writeServiceEventHandler(): Unable to decode data payload: "+err.Error(), 400)
		r.Body.Close()

		return
	}

	dataPayload, _ := json.MarshalIndent(&newServiceEvent, " ", "  ")
	log.Printf("INFO: Data Payload being written: %s", dataPayload)
	r.Body.Close()

	// Generate block
	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], newServiceEvent)
	if err != nil {
		http.Error(w, "ERROR: writeServiceEventHandler(): Unable to generate new block: "+err.Error(), 400)
		return
	}

	result := replaceChain(newBlock)
	if result {
		http.StatusText(200)
	}
}

// Handler to manage requests to /blockchain/ subchain
func blockchainViewHandler(w http.ResponseWriter, r *http.Request) {

	// Take the URL beyond /blockchain/ and split into request and value strings
	requestAction := strings.Split(r.URL.String(), "/")
	requestItem, err := strconv.Atoi(requestAction[3])
	requestItem = requestItem - 1

	if err != nil {
		log.Println("INFO: Unable to convert argument to integer, assume this is a request for entire chain")
	}

	if requestItem == -1 { //Request item is invalid so display that blockID only
		blockString, err := json.MarshalIndent(Blockchain, "", "\t")
		if err != nil {
			log.Println("ERROR: blockchainViewHandler(): Cannot print Blockchain")
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
		log.Println("INFO: garageViewHandler(): No garage id requested, setting to -1")
		requestItem = -1
	}

	if requestItem < 0 { //no value so display them all
		garageString, err := json.MarshalIndent(ValidGarages, "", "\t")

		if err != nil {
			log.Println("ERROR: garageViewHandler(): Cannot print Garages JSON data")
		}
		fmt.Fprintf(w, "\n %s", garageString)
	} else {
		garageString, _ := json.MarshalIndent(ValidGarages[requestItem], "", "\t") // Do nothing if index too high
		fmt.Fprintf(w, "\n %s.", garageString)
	}
}

func vehicleViewHandler(w http.ResponseWriter, r *http.Request) {

	// Take the URL beyond /vehicle/ and split into request and value strings
	requestAction := strings.Split(r.URL.String(), "/")
	requestReg := requestAction[3]
	requestReg = strings.ToUpper(requestReg)

	blockItemStringJSON, _ := json.MarshalIndent(vehicleMap[requestReg], "", "\t") // Do nothing if index too high
	vehicleBlockString := fmt.Sprintf("Block locations for vehicle %s: %s", requestReg, blockItemStringJSON)
	fmt.Fprintf(w, "\n %s.", vehicleBlockString)
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
	filewritelock.Lock()
	defer filewritelock.Unlock()

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
	filewritelock.Lock()
	defer filewritelock.Unlock()

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
