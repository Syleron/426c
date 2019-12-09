package main

import (
	"fmt"
)

var (
	port = 9000
)

func main () {
	fmt.Println("426c Server")
	// Generate new RSA keys
	// Create new instance of server
	server := setupServer(fmt.Sprintf(":%v", port))
	defer server.shutdown()
}

