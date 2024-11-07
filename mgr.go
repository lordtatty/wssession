package wssession

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type ReceivedMsg struct {
	ConnID  string          `json:"conn_id"`
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
}

type Mgr struct {
	Handlers map[string]MessageHandler
	Sessions Sessions
}

// HandleWebSocket upgrades the HTTP connection to a WebSocket and processes messages
func (m *Mgr) Serve(conn WebsocketConn) error {
	// Read first messages - it must be of type "connect"
	_, message, err := conn.ReadMessage()
	if err != nil {
		logger().Debug("Error reading message:", err)
		return fmt.Errorf("error reading message: %w", err)
	}
	var receivedMsg ReceivedMsg
	if err := json.Unmarshal(message, &receivedMsg); err != nil {
		logger().Debug("Error unmarshalling incoming WS message:", err)
		return err
	}
	if receivedMsg.Type != "connect" {
		logger().Debug("First message must be of type 'connect'")
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "First message must be of type 'connect'"))
		return fmt.Errorf("first message must be of type 'connect'")
	}
	logger().Debug("Received connect message", "ConnID", receivedMsg.ConnID)

	// Get the session
	sess, err := m.Sessions.Get(receivedMsg.ConnID, conn)
	if err != nil {
		return fmt.Errorf("error getting connection handler: %w", err)
	}

	sess.Cache.AutoPruneDuration = time.Minute

	// If the connection ID is not empty then this is a reconnect, update the conn and replay the cache
	if receivedMsg.ConnID != "" {
		err := sess.UpdateConnAndReplayCache(conn)
		if err != nil {
			return fmt.Errorf("error replaying cache: %w", err)
		}
		logger().Debug("Replay complete")
	}

	for {
		// Read message from WebSocket
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logger().Debug("Connection closed by client")
				break
			}
			return fmt.Errorf("error reading message: %w", err)
		}
		logger().Debug("Received: %s", message)

		// Unmarshal the received message into ReceivedMsg
		var receivedMsg ReceivedMsg
		if err := json.Unmarshal(message, &receivedMsg); err != nil {
			logger().Debug("Error unmarshalling incoming WS message:", err)
			continue
		}

		// Get the handler for the message type
		handler, exists := m.Handlers[receivedMsg.Type]
		if exists {
			handler.WSHandle(sess.Writer(), receivedMsg.Message)
		} else {
			logger().Debug("No handler registered for message type: %s", receivedMsg.Type)
		}
	}
	return nil
}

type Writer interface {
	SendJSON(msgType string, j any) error
	SendStr(msgType string, msg string) error
}

type MessageHandler interface {
	WSHandle(w Writer, msg json.RawMessage)
}

func (m *Mgr) RegisterHandler(msgType string, handler MessageHandler) {
	if m.Handlers == nil {
		m.Handlers = make(map[string]MessageHandler)
	}
	m.Handlers[msgType] = handler
}