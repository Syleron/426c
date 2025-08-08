package main

import (
    "fmt"
    "github.com/labstack/gommon/log"
    "github.com/syleron/426c/common/database"
    "github.com/syleron/426c/common/security"
    "os"
    "os/signal"
    "syscall"
)

var (
	port = 9000
	db *database.Database
)

func main() {
	var err error
	fmt.Println(`
  ____ ___  ____    
 / / /|_  |/ __/____
/_  _/ __// _ \/ __/
 /_//____/\___/\__/.net v` + VERSION + "\n")
	// Set our logging level
	log.SetLevel(1) // 1) DEBUG 2) INFO
	// Load our database
	db, err = database.New("426c")
	if err != nil {
		panic(err)
	}
    // Make sure we have our buckets
    _ = db.CreateBucket("users")
    _ = db.CreateBucket("message_tokens")
	// Generate new RSA keys
	if err := security.GenerateKeys("proteus.426c.net"); err != nil {
		panic(err)
	}
	// Create new instance of server
	server := setupServer(fmt.Sprintf(":%v", port))
	defer server.shutdown()
    // Handle OS signals for graceful shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigCh
        log.Info("signal received, shutting down...")
        server.shutdown()
    }()
	// Handle incoming connections
	server.connectionHandler()
}

