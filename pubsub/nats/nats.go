package natsps

import (
	"context"
	"encoding/json"
	"flag"
	"time"

	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/hoangtk0100/social-todo-list/pubsub"
	"github.com/nats-io/nats.go"
)

type natsPubSub struct {
	name       string
	url        string
	connection *nats.Conn
	logger     logger.Logger
}

func NewNatsPubSub(name string) *natsPubSub {
	return &natsPubSub{
		name: name,
	}
}

func (ps *natsPubSub) Publish(ctx context.Context, topic pubsub.Topic, msg *pubsub.Message) error {
	data, err := json.Marshal(msg.Data())
	if err != nil {
		ps.logger.Errorln(err)
		return err
	}

	if err := ps.connection.Publish(string(topic), data); err != nil {
		ps.logger.Errorln(err)
		return err
	}

	return nil
}

func (ps *natsPubSub) Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, unsubscribe func()) {
	msgChan := make(chan *pubsub.Message)

	sub, err := ps.connection.Subscribe(string(topic), func(msg *nats.Msg) {
		data := make(map[string]interface{})

		_ = json.Unmarshal(msg.Data, &data)

		newMsg := pubsub.NewMessage(data)
		newMsg.SetTopic(topic)

		newMsg.SetAckFunc(func() error {
			return msg.Ack()
		})

		msgChan <- newMsg
	})

	if err != nil {
		ps.logger.Errorln(err)
	}

	return msgChan, func() {
		_ = sub.Unsubscribe()
	}
}

func (ps *natsPubSub) GetPrefix() string {
	return ps.name
}

func (ps *natsPubSub) Name() string {
	return ps.name
}

func (ps *natsPubSub) Get() interface{} {
	return ps
}

func (ps *natsPubSub) InitFlags() {
	flag.StringVar(&ps.url, ps.name+"-url", nats.DefaultURL, "NATS URL - Ex: nats://127.0.0.1:4222")
}

func (ps *natsPubSub) setupOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectWait := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectWait))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectWait)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		ps.logger.Infof("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))

	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		ps.logger.Infof("Reconnected [%s]", nc.ConnectedUrl())
	}))

	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		ps.logger.Infof("Exiting: %v", nc.LastError())
	}))

	return opts
}

func (ps *natsPubSub) Configure() error {
	ps.logger = logger.GetCurrent().GetLogger(ps.name)

	conn, err := nats.Connect(ps.url, ps.setupOptions([]nats.Option{})...)
	if err != nil {
		ps.logger.Fatalln(err)
	}

	ps.logger.Infoln("Connected to NATS service.")

	ps.connection = conn

	return nil
}

func (ps *natsPubSub) Run() error {
	return ps.Configure()
}

func (ps *natsPubSub) Stop() <-chan bool {
	ch := make(chan bool)

	go func() {
		ch <- true
	}()

	return ch
}
