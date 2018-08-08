// Package web used to create a Web Server to place Handler logic
package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Blockchain struct to provide access into data for httpHandlers
type blockchainReference struct {
	bc *servicechain.BlockchainReference
}

var defaultListenerPort = "8000" // Default listener port

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

	// DEBUG
	fmt.Fprintf(w, "Request made for item ID %v", requestItem)

	if requestItem == 0 { //Request item is invalid so display that blockID only
		fmt.Printf("\nStub behaviour - no block ID hence print entre chain")
	} else {
		fmt.Fprintf(w, "\nStub behaviour - print block number %s.", requestAction)
	}
}
