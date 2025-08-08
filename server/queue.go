package main

import (
    "crypto/md5"
    "encoding/hex"
    "github.com/labstack/gommon/log"
    "github.com/syleron/426c/common/models"
    plib "github.com/syleron/426c/common/packet"
    "github.com/syleron/426c/common/utils"
    "strconv"
    "sync"
    "time"
)

type Queue struct {
    server       *Server
    queues       map[string][]*QueueItem // recipient -> queue items
    senderCounts map[string]int           // sender -> count across all queues
    rrRecipients []string                 // round-robin recipient keys
    rrIndex      int
    ttl          time.Duration
    sync.Mutex
}

type QueueItem struct {
	uid string
	msg *models.MsgToRequestModel
    enq time.Time
}

func newQueue(s *Server) *Queue {
	log.Debug("Message queue initialised")
	// Instantiate our object
    q := &Queue{
        server:       s,
        queues:       make(map[string][]*QueueItem),
        senderCounts: make(map[string]int),
        ttl:          30 * time.Minute,
    }
    // start the queue in a go routine with shutdown support
    go func() {
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                _ = q.process()
            case <-s.shutdownCh:
                return
            }
        }
    }()
    // Token TTL sweeper
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                q.sweepExpired()
            case <-s.shutdownCh:
                return
            }
        }
    }()
	// return our queue
	return q
}

const (
    perRecipientCap = 100
    perSenderCap    = 100
)

// Add enqueues a message and returns whether it was queued (true) or immediately rejected.
// It may also send immediate state updates to the sender.
func (q *Queue) Add(msgObj *models.MsgToRequestModel) bool {
	log.Debug("Queueing new message")
	// Generate uid
	hasher := md5.New()
	hasher.Write([]byte(msgObj.From + ":" + msgObj.To + ":" + strconv.Itoa(msgObj.ID)))
	uid := hex.EncodeToString(hasher.Sum(nil))
    // Check to see if our UID exists (dedupe)
    if mt, err := dbMessageTokenGet(uid); err == nil {
        log.Debug("queue item already exists, sending update chat message state")
        // Token exists; notify sender with current state and avoid duplicate enqueue
        q.server.mu.RLock()
        sender := q.server.clients[msgObj.From]
        q.server.mu.RUnlock()
        if sender != nil {
            sender.Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
            Success: mt.Success,
            MsgID:   msgObj.ID,
            To:      msgObj.To,
            }))
        }
        // Ensure it is in queue if not yet processed
        if !q.queueItemExists(uid) {
            q.Lock()
            q.queues[msgObj.To] = append(q.queues[msgObj.To], &QueueItem{msg: msgObj, uid: uid, enq: time.Now()})
            q.ensureRecipientListed(msgObj.To)
            q.senderCounts[msgObj.From]++
            q.Unlock()
        }
        return true
    }
    // No token exists → create one
    _ = dbMessageTokenAdd(&models.MessageTokenModel{UID: uid, Success: false})
    // Enforce caps
    q.Lock()
    if len(q.queues[msgObj.To]) >= perRecipientCap {
        q.Unlock()
        q.notifyMsgToFailure(msgObj, "recipient_queue_full")
        _ = dbMessageTokenDelete(uid)
        return false
    }
    if q.senderCounts[msgObj.From] >= perSenderCap {
        q.Unlock()
        q.notifyMsgToFailure(msgObj, "sender_queue_full")
        _ = dbMessageTokenDelete(uid)
        return false
    }
    q.queues[msgObj.To] = append(q.queues[msgObj.To], &QueueItem{ msg: msgObj, uid: uid, enq: time.Now() })
    q.ensureRecipientListed(msgObj.To)
    q.senderCounts[msgObj.From]++
    q.Unlock()
    return true
}

func (q *Queue) process() bool {
    q.Lock()
    defer q.Unlock()
    // Round-robin over recipients, delivering at most one message per recipient per tick
    if len(q.rrRecipients) == 0 {
        // rebuild list
        for r := range q.queues {
            if len(q.queues[r]) > 0 {
                q.rrRecipients = append(q.rrRecipients, r)
            }
        }
        q.rrIndex = 0
    }
    if len(q.rrRecipients) == 0 {
        metricQueueLength.Set(0)
        return false
    }
    processed := false
    start := q.rrIndex
    for i := 0; i < len(q.rrRecipients); i++ {
        idx := (start + i) % len(q.rrRecipients)
        recipientKey := q.rrRecipients[idx]
        items := q.queues[recipientKey]
        if len(items) == 0 {
            continue
        }
        m := items[0]
        // Attempt send if recipient online
        q.server.mu.RLock()
        recipient, ok := q.server.clients[recipientKey]
        q.server.mu.RUnlock()
        if !ok {
            continue
        }
        _, err := recipient.Send(plib.SVR_MSG, utils.MarshalResponse(&models.MsgResponseModel{
            Message: m.msg.Message,
        }))
        if err != nil {
            log.Debug("Failed to send message")
            metricMessagesSent.WithLabelValues("fail").Inc()
            continue
        }
        // Success
        q.server.mu.RLock()
        sender, ok := q.server.clients[m.msg.From]
        q.server.mu.RUnlock()
        if ok {
            user, err := dbUserGet(m.msg.From)
            if err == nil {
                if !m.enq.IsZero() {
                    metricMessageDeliverySeconds.Observe(time.Since(m.enq).Seconds())
                }
                sender.Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
                    Success: true,
                    MsgID:   m.msg.ID,
                    To:      m.msg.To,
                    Blocks:  user.Blocks,
                }))
            }
        }
        metricMessagesSent.WithLabelValues("success").Inc()
        _ = dbMessageTokenDelete(m.uid)
        // pop from recipient queue
        q.queues[recipientKey] = q.queues[recipientKey][1:]
        q.senderCounts[m.msg.From]--
        if len(q.queues[recipientKey]) == 0 {
            // remove from rr list
            q.rrRecipients = append(q.rrRecipients[:idx], q.rrRecipients[idx+1:]...)
            // adjust rrIndex to current idx
            if idx < len(q.rrRecipients) {
                q.rrIndex = idx
            } else {
                q.rrIndex = 0
            }
        } else {
            // advance rrIndex
            q.rrIndex = (idx + 1) % len(q.rrRecipients)
        }
        processed = true
        break
    }
    // Update gauge with total queued items
    total := 0
    for _, arr := range q.queues { total += len(arr) }
    metricQueueLength.Set(float64(total))
    return processed
}

func (q *Queue) queueItemExists(uid string) bool {
    for _, items := range q.queues {
        for _, it := range items {
            if it.uid == uid { return true }
        }
    }
    return false
}

func (q *Queue) ensureRecipientListed(recipient string) {
    for _, r := range q.rrRecipients {
        if r == recipient { return }
    }
    q.rrRecipients = append(q.rrRecipients, recipient)
}

// DrainFor attempts to deliver all queued messages for a specific recipient immediately.
func (q *Queue) DrainFor(recipient string) {
    q.Lock()
    defer q.Unlock()
    items := q.queues[recipient]
    if len(items) == 0 { return }
    q.server.mu.RLock()
    rcpt, ok := q.server.clients[recipient]
    q.server.mu.RUnlock()
    if !ok { return }
    // deliver in-order
    delivered := 0
    for len(items) > 0 {
        m := items[0]
        if _, err := rcpt.Send(plib.SVR_MSG, utils.MarshalResponse(&models.MsgResponseModel{ Message: m.msg.Message })); err != nil {
            break
        }
        // notify sender if online
        q.server.mu.RLock()
        sender, ok := q.server.clients[m.msg.From]
        q.server.mu.RUnlock()
        if ok {
            if user, err := dbUserGet(m.msg.From); err == nil {
                sender.Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
                    Success: true,
                    MsgID:   m.msg.ID,
                    To:      m.msg.To,
                    Blocks:  user.Blocks,
                }))
            }
        }
        _ = dbMessageTokenDelete(m.uid)
        q.senderCounts[m.msg.From]--
        items = items[1:]
        delivered++
    }
    q.queues[recipient] = items
    // Update rrRecipients if necessary
    if len(items) == 0 {
        for i, r := range q.rrRecipients {
            if r == recipient {
                q.rrRecipients = append(q.rrRecipients[:i], q.rrRecipients[i+1:]...)
                break
            }
        }
    }
    // update gauge
    total := 0
    for _, arr := range q.queues { total += len(arr) }
    metricQueueLength.Set(float64(total))
    if delivered > 0 {
        metricMessagesSent.WithLabelValues("success").Add(float64(delivered))
    }
}

func (q *Queue) notifyMsgToFailure(msgObj *models.MsgToRequestModel, code string) {
    q.server.mu.RLock()
    sender := q.server.clients[msgObj.From]
    q.server.mu.RUnlock()
    if sender != nil {
        sender.Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
            Success: false,
            MsgID:   msgObj.ID,
            To:      msgObj.To,
            Code:    code,
        }))
    }
}

func (q *Queue) sweepExpired() {
    _ = time.Now().Add(-q.ttl)
    // Remove expired tokens whose messages are no longer in queues (best effort)
    // Since we do not persist payloads, we only delete tokens; queue items are removed by delivery or manual removal
    // No direct index of tokens -> we cannot list; so no-op here without a token index.
}
