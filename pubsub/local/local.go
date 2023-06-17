package localps

import (
	"context"
	"log"
	"sync"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/pubsub"
)

// In-memory
// Buffer channel as queue
// Transmission of messages with specific topic to all subscribers within a group
type localPubSub struct {
	name         string
	messageQueue chan *pubsub.Message
	mapTopic     map[pubsub.Topic][]chan *pubsub.Message
	locker       *sync.RWMutex
}

func NewLocalPubSub(name string) *localPubSub {
	return &localPubSub{
		name:         name,
		messageQueue: make(chan *pubsub.Message, 10000),
		mapTopic:     make(map[pubsub.Topic][]chan *pubsub.Message),
		locker:       new(sync.RWMutex),
	}
}

func (ps *localPubSub) Publish(ctx context.Context, topic pubsub.Topic, msg *pubsub.Message) error {
	msg.SetTopic(topic)

	go func() {
		defer common.Recovery()

		ps.messageQueue <- msg
		log.Println("New message published :", msg.String())
	}()

	return nil
}

func (ps *localPubSub) Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, unsubscribe func()) {
	c := make(chan *pubsub.Message)

	ps.locker.Lock()

	val, ok := ps.mapTopic[topic]
	if ok {
		val = append(ps.mapTopic[topic], c)
		ps.mapTopic[topic] = val
	} else {
		ps.mapTopic[topic] = []chan *pubsub.Message{c}
	}

	ps.locker.Unlock()

	return c, func() {
		log.Println("Unsubscribe :", topic)

		if chans, ok := ps.mapTopic[topic]; ok {
			for index := range chans {
				if chans[index] == c {
					chans = append(chans[:index], chans[index+1:]...)

					ps.locker.Lock()
					ps.mapTopic[topic] = chans
					ps.locker.Unlock()
				}
			}
		}
	}
}

// Send message from message queue to subscribed topic channels
func (ps *localPubSub) run() error {
	go func() {
		defer common.Recovery()

		for {
			msg := <-ps.messageQueue
			log.Println("Message dequeue :", msg.String())

			ps.locker.RLock()

			if subs, ok := ps.mapTopic[msg.Topic()]; ok {
				for index := range subs {
					go func(c chan *pubsub.Message) {
						defer common.Recovery()
						c <- msg
					}(subs[index])
				}
			}

			ps.locker.RUnlock()
		}
	}()

	return nil
}

func (ps *localPubSub) GetPrefix() string {
	return ps.name
}

func (ps *localPubSub) Get() interface{} {
	return ps
}

func (ps *localPubSub) Name() string {
	return ps.name
}

func (*localPubSub) InitFlags() {
}

func (*localPubSub) Configure() error {
	return nil
}

func (ps *localPubSub) Run() error {
	return ps.run()
}

func (*localPubSub) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()

	return c
}
