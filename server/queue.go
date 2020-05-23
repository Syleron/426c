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
	// start the queue in a go routine
	go scheduler(q.process, 1000)
	// return our queue
	return q
}

func (q *Queue) Add(msgObj *models.MsgToRequestModel) {
	log.Debug("Queueing new message")
	// Generate uid
	hasher := md5.New()
	hasher.Write([]byte(msgObj.From + ":" + msgObj.To + ":" + strconv.Itoa(msgObj.ID)))
	uid := hex.EncodeToString(hasher.Sum(nil))
	// Check to see if our UID exists
	msgToken, err := dbMessageTokenGet(uid)
	if err != nil {
		log.Debug("queue item already exists, sending update chat message state")
		// Our token already exists so let the sender know!
		q.server.clients[msgObj.From].Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
			Success: msgToken.Success,
			MsgID:   msgObj.ID,
			To:      msgObj.To, // Needed?
		}))
		// Double check to see if our message still exists in queue
		if !q.queueItemExists(uid) {
			// Add the message into the queue
			q.queue  = append(q.queue, &QueueItem{
				msg: msgObj,
				uid: uid,
			})
		}
		return
	}
	// Add a message token to the database
	dbMessageTokenAdd(&models.MessageTokenModel{
		UID: uid,
		Success: false,
	})
	// Add the message into the queue
	q.queue  = append(q.queue, &QueueItem{
		msg: msgObj,
		uid: uid,
	})
}

func (q *Queue) process() bool {
	q.Lock()
	defer q.Unlock()
	for i, m := range q.queue {
		if _, ok := q.server.clients[m.msg.To]; ok {
			// Send our message to our recipient
			_, err := q.server.clients[m.msg.To].Send(plib.SVR_MSG, utils.MarshalResponse(&models.MsgResponseModel{
				Message: m.msg.Message,
			}))
			if err != nil {
				log.Debug("Failed to send message")
				// Failed to send message, we need to try again in the queue
				return false
			}
			// Message successfully sent
			if _, ok := q.server.clients[m.msg.From]; ok {
				user, err := dbUserGet(m.msg.From)
				if err != nil {
					log.Debug("unable to process message user does not exist")
					return false
				}
				// Let the sender know that the message successfully sent.
				q.server.clients[m.msg.From].Send(plib.SVR_MSGTO, utils.MarshalResponse(&models.MsgToResponseModel{
					Success: true,
					MsgID:   m.msg.ID,
					To:      m.msg.To, // Needed?
					Blocks:  user.Blocks,
				}))
				// Remove the message token from our server DB
				if err := dbMessageTokenDelete(m.uid); err != nil {
					log.Error("failed to delete message token")
				}
				return false
			}
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
