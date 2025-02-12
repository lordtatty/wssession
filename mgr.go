package wssession

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ReceivedMsg represents a message received over the WebSocket connection
type ReceivedMsg struct {
	ConnID  string          `json:"conn_id"`
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
	ReplyTo string          `json:"reply_to"`
}

type NoCustomState struct{}

// SessionGetter gets a session, either existing or new
type SessionGetter interface {
	Get(connID string, conn WebsocketConn, cache Cache) (*Session, error)
}

// Writer defines the interface for sending messages over a WebSocket connection
type Writer interface {
	SendJSON(msgType string, j any) error
	SendJSONAndWait(msgType string, j any, timeout time.Duration) (*json.RawMessage, error)
	SendStr(msgType string, msg string) error
	SendStrAndWait(msgType string, msg string, timeout time.Duration) (*json.RawMessage, error)
}

// MessageHandler defines how different message types are processed
type MessageHandler[T any] interface {
	WSHandle(w Writer, msg json.RawMessage, customState *T) error
}

// OnConnectFn is called when a new connection is established
type OnConnectFn[T any] func(m ReceivedMsg, customState T) (T, error)

// OnDisconnectFn is called when a connection is terminated
type OnDisconnectFn func() error

// Mgr is the main WebSocket session manager
type Mgr[T any] struct {
	Handlers     map[string]MessageHandler[T]
	Sessions     SessionGetter
	onConnFns    []OnConnectFn[T]
	onDisconnFns []OnDisconnectFn
}

// RegisterHandler associates a message type with its handler
func (m *Mgr[T]) RegisterHandler(msgType string, handler MessageHandler[T]) {
	if m.Handlers == nil {
		m.Handlers = make(map[string]MessageHandler[T])
	}
	m.Handlers[msgType] = handler
}

func (m *Mgr[T]) OnConnect(fn OnConnectFn[T]) {
	m.onConnFns = append(m.onConnFns, fn)
}

func (m *Mgr[T]) OnDisconnect(fn OnDisconnectFn) {
	m.onDisconnFns = append(m.onDisconnFns, fn)
}

// waitForConnect waits for and validates the initial connect message
func (m *Mgr[T]) waitForConnect(conn WebsocketConn) (*ReceivedMsg, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("error reading connect message: %w", err)
	}

	var receivedMsg ReceivedMsg
	if err := json.Unmarshal(message, &receivedMsg); err != nil {
		return nil, fmt.Errorf("error unmarshalling connect message: %w", err)
	}

	if receivedMsg.Type != "connect" {
		return nil, fmt.Errorf("first message must be of type 'connect'")
	}

	return &receivedMsg, nil
}

// establishSession creates or restores a session for the connection
func (m *Mgr[T]) establishSession(conn WebsocketConn, cache Cache, msg *ReceivedMsg) (*Session, *T, error) {
	// Call OnConnect handler if configured
	var err error
	customState := *new(T)
	for _, fn := range m.onConnFns {
		customState, err = fn(*msg, customState)
		if err != nil {
			return nil, nil, fmt.Errorf("connection handler error: %w", err)
		}
	}

	// Create new session
	sess, err := m.Sessions.Get(msg.ConnID, conn, cache)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating session: %w", err)
	}

	// Handle reconnection if ConnID is present
	if msg.ConnID != "" {
		if err := sess.UpdateConnAndReplayCache(conn); err != nil {
			return nil, nil, fmt.Errorf("error replaying cache: %w", err)
		}
	}

	return sess, &customState, nil
}

// handleMessages processes incoming messages for an established session
func (m *Mgr[T]) handleMessages(ctx context.Context, sess *Session, conn WebsocketConn, customState *T) error {
	wg := sync.WaitGroup{}
	errCh := make(chan error, 2) // Capture errors from handlers

	for {
		select {
		case <-ctx.Done(): // Context canceled
			errCh <- nil
		case err := <-errCh:
			// Stop processing if an error occurred or context is canceled
			wg.Wait() // Ensure all existing goroutines finish
			if err != nil {
				return fmt.Errorf("error handling messages: %w", err)
			}
			return nil
		default:
			// Read message from WebSocket
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					errCh <- nil
					continue
				}
				errCh <- fmt.Errorf("error reading message: %w", err)
				continue
			}

			var receivedMsg ReceivedMsg
			if err := json.Unmarshal(message, &receivedMsg); err != nil {
				logger().Debug("Error unmarshalling message:", err)
				continue
			}

			// Handle waiting responses
			if ok := sess.CompleteWaiterIfMatch(receivedMsg); ok {
				continue
			}

			// Process message with registered handler
			if handler, exists := m.Handlers[receivedMsg.Type]; exists {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := handler.WSHandle(sess.Writer(), receivedMsg.Message, customState)
					if err != nil {
						logger().Error("Non-Fatal Error handling message", "error", err.Error())
					}
				}()
			}
		}
	}
}

// ServeSession handles the main WebSocket session lifecycle
func (m *Mgr[T]) ServeSession(conn WebsocketConn, cache Cache) error {
	// Handle disconnect
	defer func() {
		for _, f := range m.onDisconnFns {
			if err := f(); err != nil {
				logger().Error("Error in disconnect handler", "error", err)
			}
		}
	}()

	// Wait for connect message
	connectMsg, err := m.waitForConnect(conn)
	if err != nil {
		usrMsg := "Error waiting for connect message. Ensure the first message is of type 'connect'"
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, usrMsg))
		return fmt.Errorf("connection setup failed: %w", err)
	}

	// Set up session
	sess, customState, err := m.establishSession(conn, cache, connectMsg)
	if err != nil {
		return fmt.Errorf("session establishment failed: %w", err)
	}

	// Start message handling
	return m.handleMessages(ctx, sess, conn, customState)
}
