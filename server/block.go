package main

import (
	"time"
)

var (
	// The cost to send a message through the 426c network
	msgCost int

	// Total chat messages per second
	TCSCount int

	// How often to distribute the block total
	distBlockPeriod time.Duration

	// Total amount to be distributed per time period
	distBlockTotal float32
)

func calcBlockCost(n int) int {
	//t := math.Round(TCSCount / n)
	return 0//TCSCount / n
}
