package main

var (
	ui *UserInterface
)

func main() {
	// Generate our keys
	//if err := security.GenerateKeys("127.0.0.1"); err != nil {
	//	panic(err)
	//}
	generateKeys()
	//client := setupClient()
	//defer client.Close()
	// Setup our go routine for connections handlers
	//client.connectionHandler()
	// Setup our UI
	//setupUI()
}