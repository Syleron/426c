package main

import (
    "github.com/labstack/gommon/log"
    "github.com/syleron/426c/common/models"
    "github.com/syleron/426c/common/packet"
    "github.com/syleron/426c/common/utils"
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

// Total Chat's per second (TCS) / total number of users
// For example 1 / 10 = 0.10 *


func blockCalcCost() int {
	//t := math.Round(TCSCount / n)
	return 1//TCSCount / n
}

func blockDistribute(s *Server) {
	for _ = range time.Tick(10 * time.Minute) {
		log.Debug("Issuing blocks..")
        // Take a snapshot of clients under read lock
        s.mu.RLock()
        clients := make([]*Client, 0, len(s.clients))
        for _, c := range s.clients { clients = append(clients, c) }
        s.mu.RUnlock()
        for _, c := range clients {
			// Increase user blocks by pre-configured amount
			blocks, err := dbUserBlockCredit(c.Username, 5)
			if err != nil {
				log.Error(err)
			}
			log.Debug("Blocks total ", blocks, " for ", c.Username)
			// Let the user know of their new block balance
			c.Send(packet.SVR_BLOCK, utils.MarshalResponse(&models.BlockResponseModel{
				Blocks: blocks,
				MsgCost: blockCalcCost(),
			}))
		}
	}
}
