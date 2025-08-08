package main

import (
    "math"
    "sync/atomic"
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

// Adaptive cost calculator state (sliding window of delivered messages)
const costWindowSeconds = 60
var deliveredBuckets [costWindowSeconds]int64
var deliveredIndex int64
var currentCost int64 = 1
var scaleFactor float64 = 1.0 // tuneable multiplier
var maxCost int64 = 10        // clamp the cost to avoid extremes

// recordMessageDelivered should be called whenever a message is successfully delivered
func recordMessageDelivered() {
    idx := atomic.LoadInt64(&deliveredIndex)
    atomic.AddInt64(&deliveredBuckets[idx], 1)
}

// startCostCalculator runs a ticker to recompute cost every second based on recent throughput and online users
func (s *Server) startCostCalculator() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    var lastTick time.Time
    for {
        select {
        case <-ticker.C:
            // Advance ring buffer
            now := time.Now()
            if lastTick.IsZero() || now.Sub(lastTick) >= time.Second {
                // move index
                next := (atomic.LoadInt64(&deliveredIndex) + 1) % costWindowSeconds
                atomic.StoreInt64(&deliveredIndex, next)
                atomic.StoreInt64(&deliveredBuckets[next], 0)
                lastTick = now
            }
            // Sum window
            var sum int64
            for i := 0; i < costWindowSeconds; i++ {
                sum += atomic.LoadInt64(&deliveredBuckets[i])
            }
            mps := float64(sum) / float64(costWindowSeconds)
            // Get connected users
            s.mu.RLock()
            users := len(s.clients)
            s.mu.RUnlock()
            if users <= 0 {
                atomic.StoreInt64(&currentCost, 1)
                continue
            }
            perUserRate := mps / float64(users)
            adaptive := int64(math.Ceil(perUserRate * scaleFactor))
            if adaptive < 1 { adaptive = 1 }
            if adaptive > maxCost { adaptive = maxCost }
            atomic.StoreInt64(&currentCost, adaptive)
        case <-s.shutdownCh:
            return
        }
    }
}

func (s *Server) blockCalcCost() int {
    cost := atomic.LoadInt64(&currentCost)
    if cost < 1 { return 1 }
    return int(cost)
}

func blockDistribute(s *Server) {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
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
                    MsgCost: s.blockCalcCost(),
                }))
                metricBlocksIssued.Add(float64(5))
            }
        case <-s.shutdownCh:
            return
        }
    }
}
