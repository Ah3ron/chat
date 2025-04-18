package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Client represents a single chatting user.
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// newHub constructs a new Hub.
func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// run starts the hub's main loop.
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if h.clients[client] {
				close(client.send)
				delete(h.clients, client)
			}

		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					// if client buffer is full, drop it
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// readPump reads messages from the WebSocket and forwards them to the hub.
func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		h.broadcast <- msg
	}
}

// writePump writes messages from the hub to the client WebSocket.
func (c *Client) writePump() {
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				// hub closed client
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}
}

func main() {
	hub := newHub()
	go hub.run()

	app := fiber.New()

	// WebSocket endpoint
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		client := &Client{
			conn: c,
			send: make(chan []byte, 256),
		}
		hub.register <- client

		// launch pumps
		go client.writePump()
		client.readPump(hub)
	}))

	// Simple health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK " + time.Now().Format(time.RFC3339))
	})

	log.Println("Starting chat server on :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
