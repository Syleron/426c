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
	server *Server
	queue []*QueueItem // Array of Packets (byte arrays)
	sync.Mutex
}

type QueueItem struct {
	uid string
	msg *models.MsgToRequestModel
}

func newQueue(s *Server) *Queue {
	log.Debug("Message queue initialised")
	// Instantiate our object
	q := &Queue{
		server: s,
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
	// return our queue
	return q
}

func (q *Queue) Add(msgObj *models.MsgToRequestModel) {
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
            q.queue = append(q.queue, &QueueItem{msg: msgObj, uid: uid})
            q.Unlock()
        }
        return
    }
    // No token exists → create one
    _ = dbMessageTokenAdd(&models.MessageTokenModel{UID: uid, Success: false})
	// Add the message into the queue
    q.Lock()
    q.queue  = append(q.queue, &QueueItem{
		msg: msgObj,
		uid: uid,
	})
    q.Unlock()
}

func (q *Queue) process() bool {
	q.Lock()
	defer q.Unlock()
	for i, m := range q.queue {
        q.server.mu.RLock()
        recipient, ok := q.server.clients[m.msg.To]
        q.server.mu.RUnlock()
        if ok {
			// Send our message to our recipient
            _, err := recipient.Send(plib.SVR_MSG, utils.MarshalResponse(&models.MsgResponseModel{
				Message: m.msg.Message,
			}))
			if err != nil {
				log.Debug("Failed to send message")
				// Failed to send message, we need to try again in the queue
				return false
			}
            // Message successfully sent
            q.server.mu.RLock()
            sender, ok := q.server.clients[m.msg.From]
            q.server.mu.RUnlock()
            if ok {
				user, err := dbUserGet(m.msg.From)
				if err != nil {
					log.Debug("unable to process message user does not exist")
					return false
				}
				// Let the sender know that the message successfully sent.
                sender.Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
					Success: true,
					MsgID:   m.msg.ID,
					To:      m.msg.To, // Needed?
					Blocks:  user.Blocks,
				}))
            }
            // Remove the message token from our server DB
            if err := dbMessageTokenDelete(m.uid); err != nil {
                log.Error("failed to delete message token")
			}
            // Remove item from queue now that it is processed
            q.queue = append(q.queue[:i], q.queue[i+1:]...)
            // Continue processing next items
            return false
			// Remove item from queue
            q.queue = append(q.queue[:i], q.queue[i+1:]...)
		}
	}
	return false
}

func (q *Queue) queueItemExists(uid string) bool {
	for _, msgI := range q.queue {
		if msgI.uid == uid {
			return true
		}
	}
	return false
}
