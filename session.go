package wssession

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type WebsocketConn interface {
	WriteJSON(msgType interface{}) error
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type ResponseMsg struct {
	ID      string `json:"id"`
	ConnID  string `json:"conn_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
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

var DefaultLogger = &NoopLogger{}

var globalLogger Logger = DefaultLogger

func SetLogger(l Logger) {
	globalLogger = l
}

func logger() Logger {
	return globalLogger
}

type WaiterResp struct {
	Msg json.RawMessage
	Err error
}

type Cache interface {
	Add(connID string, r ResponseMsg) error
	Items(connID string) ([]*ResponseMsg, error)
}

// Represents a single session with a client
// A session is from the user's point of view - it should be able to survive a temporary connection dropout
type Session struct {
	ConnID   string
	Conn     WebsocketConn // Websocket connection - default is gorilla/websocket.Conn
	cache    Cache
	mux      sync.Mutex
	waitsMux sync.Mutex
	waits    map[string]chan WaiterResp
}

func (c *Session) Cache() Cache {
	return c.cache
}

func (c *Session) addToCache(r ResponseMsg) {
	if c.cache == nil {
		return
	}
	c.cache.Add(c.ConnID, r)
}

func (c *Session) cachedItems() ([]*ResponseMsg, error) {
	if c.cache == nil {
		return nil, nil
	}
	return c.cache.Items(c.ConnID)
}

func (c *Session) SetCache(cache Cache) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cache = cache
}

func (c *Session) UpdateConnAndReplayCache(wsc WebsocketConn) error {
	itms, err := c.cachedItems()
	if err != nil {
		return fmt.Errorf("error getting cache length: %w", err)
	}
	logger().Debug("Replaying cache", "len", len(itms))
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Conn = wsc
	for _, msg := range itms {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			return fmt.Errorf("error writing JSON to websocket: %w", err)
		}
	}
	return nil
}

func (c *Session) writeJSONRaw(msgType string, msg string, waiterTimeout time.Duration) (chan WaiterResp, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	r := &ResponseMsg{
		ID:      uuid.NewString(),
		ConnID:  c.ID(),
		Type:    msgType,
		Message: msg,
	}
	c.addToCache(*r)
	var wCh chan WaiterResp
	if waiterTimeout != 0 {
		var err error
		wCh, err = c.startWaiter(r.ID, waiterTimeout)
		if err != nil {
			return nil, fmt.Errorf("error starting waiter: %w", err)
		}
	}
	c.Conn.WriteJSON(r)
	return wCh, nil // TODO: We aren't returning errors as we're caching the messages in case of a connection timeout - is this the right approach?
}

func (c *Session) writeJSON(msgType string, msg string) error {
	_, err := c.writeJSONRaw(msgType, msg, 0)
	return err
}

func (c *Session) writeJSONWithWaiter(msgType string, msg string, timeout time.Duration) (chan WaiterResp, error) {
	return c.writeJSONRaw(msgType, msg, timeout)
}

func (c *Session) writeJSONAndWait(msgType string, msg string, timeout time.Duration) (*json.RawMessage, error) {
	w, err := c.writeJSONWithWaiter(msgType, msg, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to write and start waiter: %w", err)
	}
	r := <-w
	if r.Err != nil {
		return nil, fmt.Errorf("error in waiter response: %w", r.Err)
	}
	return &r.Msg, nil
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

func (s *Session) CompleteWaiterIfMatch(msg ReceivedMsg) bool {
	s.waitsMux.Lock()
	defer s.waitsMux.Unlock()
	id := msg.ReplyTo
	if _, exists := s.waits[id]; exists {
		s.waits[id] <- WaiterResp{
			Msg: msg.Message,
		}
		close(s.waits[id])
		delete(s.waits, id)
		return true
	}
	return false
}

func (s *Session) timeoutWaiter(id string) {
	s.waitsMux.Lock()
	defer s.waitsMux.Unlock()
	if _, exists := s.waits[id]; exists {
		s.waits[id] <- WaiterResp{
			Err: fmt.Errorf("timeout waiting for reply to message: %s", id),
		}
		close(s.waits[id])
		delete(s.waits, id)
	}
}

func (s *Session) startWaiter(id string, timeout time.Duration) (chan WaiterResp, error) {
	s.waitsMux.Lock()
	defer s.waitsMux.Unlock()
	if s.waits == nil {
		s.waits = make(map[string]chan WaiterResp)
	}
	if _, exists := s.waits[id]; exists {
		return nil, fmt.Errorf("already waiting for reply to message: %s", id)
	}
	s.waits[id] = make(chan WaiterResp, 1)
	go func() {
		<-time.After(timeout)
		s.timeoutWaiter(id)
	}()
	return s.waits[id], nil
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
	return s.sess.writeJSON(msgType, string(b))
}

func (s *SessionWriter) SendJSONAndWait(msgType string, j any, timeout time.Duration) (*json.RawMessage, error) {
	b, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("error marshalling message: %w", err)
	}
	msg, err := s.sess.writeJSONAndWait(msgType, string(b), timeout)
	if err != nil {
		return nil, fmt.Errorf("error in waiter response: %w", err)
	}
	return msg, nil
}

func (s *SessionWriter) SendStr(msgType string, msg string) error {
	return s.sess.writeJSON(msgType, msg)
}

func (s *SessionWriter) SendStrAndWait(msgType string, msg string, timeout time.Duration) (*json.RawMessage, error) {
	resp, err := s.sess.writeJSONAndWait(msgType, msg, timeout)
	if err != nil {
		return nil, fmt.Errorf("error in waiter response: %w", err)
	}
	return resp, nil
}

// Mange multiple sessions
type Sessions struct {
	sess map[string]*Session
}

func (s *Sessions) Get(connId string, conn WebsocketConn, cache Cache) (*Session, error) {
	if s.sess == nil {
		s.sess = make(map[string]*Session)
	}
	if connId == "" {
		logger().Info("Creating new connection handler")
		c := &Session{
			Conn:   conn,
			ConnID: uuid.NewString(),
			cache:  cache,
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
