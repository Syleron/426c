package main

import (
	"github.com/syleron/426c/common/models"
)

// Setup message queue for failed/pending messages

type MessageQueue struct {
	queue      []*models.Message
	processing bool
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		queue:      make([]*models.Message, 0),
		processing: false,
	}
}

func (mq *MessageQueue) Add(m *models.Message) {
	mq.queue = append(mq.queue, m)
}

//
func (mq *MessageQueue) Process() {
	if !mq.processing {
		mq.processing = true
		for _, m := range mq.queue {
			client.cmdMsgTo(m)
		}
		mq.processing = false
	}
}
