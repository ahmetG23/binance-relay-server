package main

import (
	// "client/original"
	"client/server"
)

// Before executing the main function, start the server and wait 
// until see the message "Server started at localhost:8080"

func main() {
	pairs := []string{"BTCUSDT", "ETHUSDT", "BNBUSDT", "XRPUSDT"}
	timeoutSec := 30

	// uncomment this line and the import statement to run the original client
	
	// go original.Init(pairs, true, timeoutSec)

	server.Init(pairs, true, timeoutSec)
}
