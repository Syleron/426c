package main

import "github.com/syleron/426c/common/security"

var (
	ui *UserInterface
)

func main() {
	// Generate our keys
	if err := security.GenerateKeys("127.0.0.1"); err != nil {
		panic(err)
	}
	client := setupClient()
	defer client.Close()
	// Setup our go routine for connections handlers
	client.connectionHandler()
	// Setup our UI
	//setupUI()
}