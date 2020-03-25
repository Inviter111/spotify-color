package ws

import (
	"log"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client ws
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Fatalln(err)
	}
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Fatalln(err)
		}
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Fatalf("Error: %v", err)
			}
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Fatalln(err)
			}
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Fatalln(err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Fatalln(err)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Fatalln(err)
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs serving websocket
func ServeWs(hub *Hub, ctx *fasthttp.RequestCtx) {
	var client *Client
	upgrader.CheckOrigin = func(ctx *fasthttp.RequestCtx) bool {
		if origin := string(ctx.Request.Header.Peek("Origin")); origin == "http://localhost:3000" {
			return true
		}
		return false
	}
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		client = &Client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte, 256),
		}

		client.hub.register <- client

		go client.writePump()
		client.readPump()
	})
	if err != nil {
		log.Fatal(err)
		return
	}
}
