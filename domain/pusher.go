package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"terminal/logger"

	"github.com/gorilla/websocket"
	"github.com/pusher/pusher-http-go/v5"
)

type PusherEvent string

const (
	PUSHER_REPORT PusherEvent = "Report"
)

func (p PusherEvent) IsValid() bool {
	switch p {
	case PUSHER_REPORT:
		return true
	}

	return false
}

type Pusher struct {
	client    *pusher.Client
	wsConn    *websocket.Conn
	listenCtx context.Context
	cancel    context.CancelFunc
}

func NewPusher() *Pusher {
	pusherClient := pusher.Client{
		AppID:   os.Getenv("PUSHER_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: os.Getenv("PUSHER_CLUSTER"),
	}

	client := &Pusher{
		client: &pusherClient,
	}

	return client
}

func (p *Pusher) chanelName(terminal *Terminal) string {
	return fmt.Sprintf("trm_%s", terminal.UserName)
}

func (p *Pusher) Send(terminal *Terminal, event PusherEvent, message string) error {
	if !event.IsValid() {
		return fmt.Errorf("o evento %s não existe", string(event))
	}
	data := map[string]string{"message": message}

	chanell := p.chanelName(terminal)

	err := p.client.Trigger(chanell, string(event), data)
	if err != nil {
		return fmt.Errorf("erro para disparar o evento: %s", err)
	}

	return nil
}

func (p *Pusher) Read(terminal *Terminal) {
	channel, chanelErr := p.client.Channel(terminal.Slug(), pusher.ChannelParams{})
	if chanelErr != nil {
		logger.Println("Falha ao ler socket")
	}
	logger.Println("channel", channel)
}

func (p *Pusher) Listen(terminal *Terminal, eventName PusherEvent, handler func(data map[string]interface{})) error {
	// Build Pusher websocket URL (protocol 7)
	urlStr := url.URL{
		Scheme:   "wss",
		Host:     fmt.Sprintf("ws-%s.pusher.com", os.Getenv("PUSHER_CLUSTER")),
		Path:     fmt.Sprintf("/app/%s", os.Getenv("PUSHER_KEY")),
		RawQuery: "protocol=7&client=go&version=0.1&flash=false",
	}

	conn, _, err := websocket.DefaultDialer.Dial(urlStr.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect websocket: %w", err)
	}

	p.wsConn = conn

	ctx, cancel := context.WithCancel(context.Background())
	p.listenCtx = ctx
	p.cancel = cancel

	// Subscribe to channel
	subscribeMessage := map[string]interface{}{
		"event": "pusher:subscribe",
		"data": map[string]interface{}{
			"channel": p.chanelName(terminal),
		},
	}

	if err := conn.WriteJSON(subscribeMessage); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	// Listen loop in goroutine
	go func() {
		defer conn.Close()

		for {
			select {
			case <-ctx.Done():
				log.Println("Pusher listen canceled")
				return
			default:
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("Error reading websocket message: %v", err)
					return
				}

				var msg map[string]interface{}
				if err := json.Unmarshal(message, &msg); err != nil {
					log.Printf("Error parsing websocket message: %v", err)
					continue
				}

				// Check event name and call handler
				if ev, ok := msg["event"].(string); ok && ev == string(eventName) {
					if dataRaw, ok := msg["data"].(string); ok {
						var data map[string]interface{}
						if err := json.Unmarshal([]byte(dataRaw), &data); err != nil {
							log.Printf("Error parsing event data: %v", err)
							continue
						}
						handler(data)
					}
				}
			}
		}
	}()

	return nil
}
