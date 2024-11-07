package wssession

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type WebsocketConn interface {
	WriteJSON(msgType interface{}) error
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type ResponseMsg struct {
	ID      string `json:"id"`
	ConnID  string `json:"conn_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Represents a single session with a client
// A session is from the user's point of view - it should be able to survive a temporary connection dropout
type Session struct {
	ConnID string
	Conn   WebsocketConn // Websocket connection - default is gorilla/websocket.Conn
	Cache  PrunerCache
	mux    sync.Mutex
}

type NoopLogger struct{}

func (n *NoopLogger) Debug(msg string, args ...any) {}
func (n *NoopLogger) Info(msg string, args ...any)  {}
func (n *NoopLogger) Warn(msg string, args ...any)  {}
func (n *NoopLogger) Error(msg string, args ...any) {}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

var globalLogger Logger = &NoopLogger{}

func SetLogger(l Logger) {
	globalLogger = l
}

func logger() Logger {
	return globalLogger
}

func (c *Session) UpdateConnAndReplayCache(wsc WebsocketConn) error {
	logger().Debug(">>>>>>>>>>>>>> Replaying cache", "len", c.Cache.Len())
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Conn = wsc
	for _, msg := range c.Cache.Items() {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			return fmt.Errorf("error writing JSON to websocket: %w", err)
		}
	}
	return nil
}

func (c *Session) writeJSON(msgType string, msg string) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	r := &ResponseMsg{
		ID:      uuid.NewString(),
		ConnID:  c.ID(),
		Type:    msgType,
		Message: msg,
	}
	c.Cache.Add(*r)
	c.Conn.WriteJSON(r)
	return nil // TODO: We aren't returning errors as we're caching the messages in case of a connection timeout - is this the right approach?
}

func (c *Session) ID() string {
	if c.ConnID == "" {
		c.ConnID = uuid.NewString()
	}
	return c.ConnID
}

func (s *Session) Writer() *SessionWriter {
	return &SessionWriter{sess: s}
}

type SessionWriter struct {
	sess *Session
}

func (s *SessionWriter) SendJSON(msgType string, j any) error {
	logger().Debug("Sending JSON:", "msg", j)
	b, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}
	fmt.Println("Sending JSON:", b)
	fmt.Println("Sending JSON:", string(b))
	return s.sess.writeJSON(msgType, string(b))
}

func (s *SessionWriter) SendStr(msgType string, msg string) error {
	return s.sess.writeJSON(msgType, msg)
}

// Mange multiple sessions
type Sessions struct {
	sess map[string]*Session
}

func (s *Sessions) Get(connId string, conn WebsocketConn) (*Session, error) {
	if s.sess == nil {
		s.sess = make(map[string]*Session)
	}
	if connId == "" {
		logger().Info("Creating new connection handler")
		c := &Session{
			Conn:   conn,
			ConnID: uuid.NewString(),
		}
		s.sess[c.ID()] = c
		return c, nil
	}

	logger().Info("Reusing connection handler", "ID", connId)
	if _, exists := s.sess[connId]; !exists {
		return nil, fmt.Errorf("no connection handler found for ID: %s", connId)
	}
	return s.sess[connId], nil
}
