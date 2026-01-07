package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Http를 WebSocket으로 업그레이드
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var hub *Hub

// initWebSocket WebSocket Hub 초기화
func initWebSocket() {
	hub = &Hub{
		clients:  make(map[*Client]bool),
		bChannal: make(chan []byte),
		rChannel: make(chan *Client),
		uChannel: make(chan *Client),
	}
	go hub.run()
}

// Hub WebSocket 연결을 관리하는 Hub
/*
	새 클라이언트 → register → clients에 추가
	메시지 수신 → broadcast → 모든 clients에 전송
	연결 종료 → unregister → clients에서 제거
*/
type Hub struct {
	clients  map[*Client]bool // 연결된 클라이언트 목록
	bChannal chan []byte      // 메시지 브로드캐스트 채널
	rChannel chan *Client     // 클라이언트 등록 채널
	uChannel chan *Client     // 클라이언트 해제 채널
}

// Client WebSocket 클라이언트
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// run Hub 실행 (클라이언트 등록/해제 및 메시지 브로드캐스트 처리)
func (h *Hub) run() {
	for {
		select {
		case client := <-h.rChannel:
			h.clients[client] = true

		case client := <-h.uChannel:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.bChannal:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// WebSocket 핸들러
/*
	클라이언트 연결을 받아 WebSocket으로 전환하고,
	Hub에 등록한 뒤 메시지 송수신을 시작
*/
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.rChannel <- client

	go client.writeMessage()
	go client.readMessage()
}

// 클라이언트로부터 메시지 수신
func (c *Client) readMessage() {
	defer func() {
		c.hub.uChannel <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.hub.bChannal <- message
	}
}

// writePump 클라이언트로 메시지 전송
func (c *Client) writeMessage() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
