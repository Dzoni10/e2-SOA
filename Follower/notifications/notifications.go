package notifications

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
}

var (
	clients   = make(map[string]*Client) // userID -> client
	clientsMu sync.RWMutex
)

// Add a new client
func Register(userID string, conn *websocket.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients[userID] = &Client{UserID: userID, Conn: conn}
}

// Remove a client
func Unregister(userID string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	delete(clients, userID)
}

// Send notification to a specific user
func Notify(userID string, message interface{}) error {
	clientsMu.RLock()
	defer clientsMu.RUnlock()
	client, ok := clients[userID]
	if !ok {
		return nil // user offline
	}
	return client.Conn.WriteJSON(message)
}
