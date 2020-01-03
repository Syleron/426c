package main

var (
	blocks int = 0
)

// creditBlocks - increases the total blocks we have available
func creditBlocks(b int) error {
	blocks += b
	return nil
}

// debitBlocks - decreases the total blocks we have available
func debitBlocks() error {
	return nil
}

// getBlocks - returns the total number of blocks available
func getBlocks() int {
	return blocks
}
