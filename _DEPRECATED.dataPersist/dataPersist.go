// Package dataPersist is used to write and load the blockchain to disk to support persistence
package dataPersist

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

var lock sync.Mutex // Used to synchronise WRITE activities

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
func Save(path string, structIn interface{}) error {

	// Create a lock and then defer the unlock until function exit
	lock.Lock()
	defer lock.Unlock()

	//Create os.File and defer close
	file, err := os.Create(path)
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
func Load(path string, structOut interface{}) error {

	// Lock and defer the unlock until function exit
	lock.Lock()
	defer lock.Unlock()

	log.Println("INFO: Loading " + path)
	fileOut, err := os.Open(path)

	// DEPRECATED - PERSISTENT FILE WILL ALWAYS exist
	/*	if !os.IsNotExist(err) { // Check if it does not exist or if it is a real error
			log.Println("INFO: dataPersist.Load(): File does not exist.")
			return nil
		}
	*/

	if err != nil {
		log.Println("ERROR: dataPersist.Load() " + err.Error() + " while openning file " + path)
		return err
	}

	defer fileOut.Close()
	Unmarshal(fileOut, structOut)

	return nil
}
