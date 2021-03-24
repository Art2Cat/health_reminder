package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gen2brain/beeep"
	"os"
	"os/signal"
	"syscall"

	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var clientId string

func main() {
	client, err := NewWebSocketClient("socket.art2cat.com", "ws/health")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connecting")

	// Close connection correctly on exit
	sigs := make(chan os.Signal, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// The program will wait here until it gets the
	<-sigs

	client.Stop()
	fmt.Println("Goodbye")
}

type Message struct {
	ClientId  string `json:"clientId"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// WebSocketClient return websocket client connection
type WebSocketClient struct {
	configStr string
	sendBuf   chan []byte
	ctx       context.Context
	ctxCancel context.CancelFunc

	mu     sync.RWMutex
	wsconn *websocket.Conn
}

// NewWebSocketClient create new websocket connection
func NewWebSocketClient(host, channel string) (*WebSocketClient, error) {
	conn := WebSocketClient{
		sendBuf: make(chan []byte, 1),
	}
	conn.ctx, conn.ctxCancel = context.WithCancel(context.Background())

	u := url.URL{Scheme: "ws", Host: host, Path: channel}
	conn.configStr = u.String()

	go conn.listen()
	go conn.listenWrite()
	return &conn, nil
}

func (conn *WebSocketClient) Connect() *websocket.Conn {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.wsconn != nil {
		return conn.wsconn
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		select {
		case <-conn.ctx.Done():
			return nil
		default:
			ws, _, err := websocket.DefaultDialer.Dial(conn.configStr, nil)
			if err != nil {
				conn.log("connect", err, fmt.Sprintf("Cannot connect to websocket: %s", conn.configStr))
				continue
			}
			conn.log("connect", nil, fmt.Sprintf("connected to websocket to %s", conn.configStr))
			conn.wsconn = ws
			clientId = uuid.New().String()
			now := time.Now() // current local time
			sec := now.Unix()
			msg := &Message{clientId, "handshake", "", "", sec}
			// msg := map[string]string{"clientId": clientId, "type": "handshake", "timestamp": "12345678"}
			data, _ := json.Marshal(msg)
			ws.WriteMessage(websocket.TextMessage, []byte(data))
			return conn.wsconn
		}
	}
}

func (conn *WebSocketClient) listen() {
	conn.log("listen", nil, fmt.Sprintf("listen for the messages: %s", conn.configStr))
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-conn.ctx.Done():
			return
		case <-ticker.C:
			for {
				ws := conn.Connect()
				if ws == nil {
					return
				}
				_, bytMsg, err := ws.ReadMessage()
				if err != nil {
					conn.log("listen", err, "Cannot read websocket message")
					conn.closeWs()
					break
				}
				conn.log("listen", nil, fmt.Sprintf("websocket msg: %s\n", string(bytMsg)))
				msg := &Message{}
				json.Unmarshal(bytMsg, &msg)
				if msg.Type == "heartbeat" {
					msg.ClientId = clientId
					data, _ := json.Marshal(msg)
					ws.WriteMessage(websocket.TextMessage, []byte(data))
				} else if msg.Type == "showNotice" {
					go showNotify(msg.Title, msg.Message)
				}

			}
		}
	}
}

// Write data to the websocket server
func (conn *WebSocketClient) Write(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	for {
		select {
		case conn.sendBuf <- data:
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context canceled")
		}
	}
}

func (conn *WebSocketClient) listenWrite() {
	for data := range conn.sendBuf {
		ws := conn.Connect()
		if ws == nil {
			err := fmt.Errorf("conn.ws is nil")
			conn.log("listenWrite", err, "No websocket connection")
			continue
		}

		if err := ws.WriteMessage(
			websocket.TextMessage,
			data,
		); err != nil {
			conn.log("listenWrite", nil, "WebSocket Write Error")
		}
		conn.log("listenWrite", nil, fmt.Sprintf("send: %s", data))
	}
}

// Close will send close message and shutdown websocket connection
func (conn *WebSocketClient) Stop() {
	conn.ctxCancel()
	conn.closeWs()
}

// Close will send close message and shutdown websocket connection
func (conn *WebSocketClient) closeWs() {
	conn.mu.Lock()
	if conn.wsconn != nil {
		conn.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.wsconn.Close()
		conn.wsconn = nil
	}
	conn.mu.Unlock()
}

// Log print log statement
// In real word I would recommend to use zerolog or any other solution
func (conn *WebSocketClient) log(f string, err error, msg string) {
	if err != nil {
		fmt.Printf("Error in func: %s, err: %v, msg: %s\n", f, err, msg)
	} else {
		fmt.Printf("Log in func: %s, %s\n", f, msg)
	}
}

func showNotify(title string, message string) {
	err := beeep.Notify(title, message, "assets/information.png")
	if err != nil {
		panic(err)

	}
}
