package wssession

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ReceivedMsg struct {
	ConnID  string          `json:"conn_id"`
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
	ReplyTo string          `json:"reply_to"`
}

type SessionGetter interface {
	Get(connID string, conn WebsocketConn) (*Session, error)
}

type OnConnectFn func(m ReceivedMsg) error
type OnDisconnectFn func() error

type Mgr struct {
	Handlers     map[string]MessageHandler
	Sessions     SessionGetter
	onConnFns    []OnConnectFn
	onDisconnFns []OnDisconnectFn
}

// When receiving the initial "connect" message, the OnConnect function is called
// This is a good place to do auth, etc.  Return an error to close the connection.
func (m *Mgr) OnConnect(f OnConnectFn) {
	m.onConnFns = append(m.onConnFns, f)
}

func (m *Mgr) callOnConnectFns(r ReceivedMsg) error {
	for _, f := range m.onConnFns {
		if err := f(r); err != nil {
			return err
		}
	}
	return nil
}

// When the connection is closed, the OnDisconnect function is called
func (m *Mgr) OnDisconnect(f OnDisconnectFn) {
	m.onDisconnFns = append(m.onDisconnFns, f)
}

func (m *Mgr) callOnDisconnectFns() error {
	for _, f := range m.onDisconnFns {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// HandleWebSocket upgrades the HTTP connection to a WebSocket and processes messages
func (m *Mgr) ServeSession(conn WebsocketConn) error {
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

	// defer calling OnDisconnect functions
	defer func() {
		// TODO: shoudl we refactor so it doesn't use a defer? This will let us return errs.
		// This whole func probably needs breaking down.
		if err := m.callOnDisconnectFns(); err != nil {
			logger().Error("Error in OnDisconnectFn", "error", err)
		}
	}()

	// Call OnConnect functions
	if err := m.callOnConnectFns(receivedMsg); err != nil {
		return fmt.Errorf("error in OnConnectFn: %w", err)
	}

	// Get the session
	sess, err := m.Sessions.Get(receivedMsg.ConnID, conn)
	if err != nil {
		return fmt.Errorf("error getting connection handler: %w", err)
	}

	// If the connection ID is not empty then this is a reconnect, update the conn and replay the cache
	if receivedMsg.ConnID != "" {
		err := sess.UpdateConnAndReplayCache(conn)
		if err != nil {
			return fmt.Errorf("error replaying cache: %w", err)
		}
		logger().Debug("Replay complete")
	}

	wg := sync.WaitGroup{}

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

		// Check for any sessions "waiting" for a response
		if ok := sess.CompleteWaiterIfMatch(receivedMsg); ok {
			// If the message was a response to a waiter, don't process it further
			// It's being handled by the sess now
			continue
		}

		// Otherwise get the handler for the message type
		handler, exists := m.Handlers[receivedMsg.Type]
		if exists {
			// Run handler
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := handler.WSHandle(sess.Writer(), receivedMsg.Message)
				if err != nil {
					logger().Error("Non-Fatal Error handling message", "error", err.Error())
				}
			}()
		} else {
			logger().Debug("No handler registered for message type: %s", receivedMsg.Type)
		}
	}
	wg.Wait()
	return nil
}

type Writer interface {
	SendJSON(msgType string, j any) error
	SendJSONAndWait(msgType string, j any, timeout time.Duration) (*json.RawMessage, error)
	SendStr(msgType string, msg string) error
	SendStrAndWait(msgType string, msg string, timeout time.Duration) (*json.RawMessage, error)
}

type MessageHandler interface {
	WSHandle(w Writer, msg json.RawMessage) error
}

func (m *Mgr) RegisterHandler(msgType string, handler MessageHandler) {
	if m.Handlers == nil {
		m.Handlers = make(map[string]MessageHandler)
	}
	m.Handlers[msgType] = handler
}
