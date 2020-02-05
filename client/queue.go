package main

import (
	"github.com/syleron/426c/common/models"
	"sync"
)

type MessageQueue struct {
	queue      []*models.Message
	processing bool
	sync.Mutex
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		queue:      make([]*models.Message, 0),
		processing: false,
	}
}

func (mq *MessageQueue) Add(m *models.Message) {
	mq.Lock()
	defer mq.Unlock()
	mq.queue = append(mq.queue, m)
}

func (mq *MessageQueue) Remove(m *models.Message) {
	mq.Lock()
	defer mq.Unlock()
	if len(mq.queue) == 1 {
		mq.queue = make([]*models.Message, 0)
	} else {
		for i, message := range mq.queue {
			if message == m {
				mq.queue = append(mq.queue[:i], mq.queue[i+1:]...)
			}
		}
	}
}

func (mq *MessageQueue) Process() {
	if !mq.processing {
		mq.processing = true
		for _, m := range mq.queue {
			client.cmdMsgTo(m)
			mq.Remove(m)
		}
		mq.processing = false
	}
}
