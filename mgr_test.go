package wssession_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lordtatty/wssession"
	mocks "github.com/lordtatty/wssession/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWSMgr_ServeSession_Success(t *testing.T) {
	t.Parallel()
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)
	mCache := &mocks.MockCache{}
	defer mCache.AssertExpectations(t)

	// Prepare the connect message
	receivedMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "connect",
		Message: json.RawMessage(`{}`),
	}
	connectMessage, _ := json.Marshal(receivedMsg)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.EXPECT().Get("", mConn, mCache).Return(&wssession.Session{}, nil).Once()

	// Run test
	sut := &wssession.Mgr[wssession.NoCustomState]{
		Sessions: mSessions,
	}
	err := sut.ServeSession(mConn, mCache)

	// Assertions
	assert.NoError(t, err)
}

func TestWSMgr_ServeSession_InvalidFirstMessageType(t *testing.T) {
	t.Parallel()
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	// Prepare an invalid first message type
	invalidMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "invalid",
		Message: json.RawMessage(`{}`),
	}
	invalidMessage, _ := json.Marshal(invalidMsg)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, invalidMessage, nil).Once()
	mConn.EXPECT().WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Error waiting for connect message. Ensure the first message is of type 'connect'")).Return(nil).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.EXPECT().Get("", mConn, nil).Return(&wssession.Session{}, nil).Once()

	mCache := &mocks.MockCache{}
	defer mCache.AssertExpectations(t)

	// Run test
	sut := &wssession.Mgr[wssession.NoCustomState]{
		Sessions: mSessions,
	}
	err := sut.ServeSession(mConn, mCache)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "first message must be of type 'connect'")
}

func TestWSMgr_ServeSession_ReadMessageError(t *testing.T) {
	t.Parallel()
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(0, nil, errors.New("read error")).Once()
	mConn.EXPECT().WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Error waiting for connect message. Ensure the first message is of type 'connect'")).Return(nil).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.EXPECT().Get("", mConn, nil).Return(&wssession.Session{}, nil).Once()

	mCache := &mocks.MockCache{}
	defer mCache.AssertExpectations(t)

	// Run test
	sut := &wssession.Mgr[wssession.NoCustomState]{
		Sessions: mSessions,
	}
	err := sut.ServeSession(mConn, mCache)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading connect message: read error")
}

func TestWSMgr_RegisterHandler(t *testing.T) {
	t.Parallel()
	sut := &wssession.Mgr[wssession.NoCustomState]{}

	// Create a mock handler
	mockHandler := new(mocks.MockMessageHandler[wssession.NoCustomState])

	// Register the handler
	sut.RegisterHandler("testType", mockHandler)

	// Assertions
	assert.NotNil(t, sut.Handlers)
	assert.Equal(t, mockHandler, sut.Handlers["testType"])
}

func TestWSMgr_ServeSession_HandlerInvocation(t *testing.T) {
	t.Parallel()
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	// Prepare a valid connect message
	connectMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "connect",
		Message: json.RawMessage(`{}`),
	}
	connectMessage, _ := json.Marshal(connectMsg)

	// Prepare a message that should invoke the handler
	testMsg := wssession.ReceivedMsg{
		ConnID:  "some-id",
		Type:    "testType",
		Message: json.RawMessage(`{"key":"value"}`),
	}
	testMessage, _ := json.Marshal(testMsg)

	// Mock cache
	mCache := &mocks.MockCache{}
	defer mCache.AssertExpectations(t)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, testMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.EXPECT().Get("", mConn, mCache).Return(&wssession.Session{}, nil).Once()

	mockHandler := new(mocks.MockMessageHandler[wssession.NoCustomState])
	sut := &wssession.Mgr[wssession.NoCustomState]{
		Handlers: map[string]wssession.MessageHandler[wssession.NoCustomState]{
			"testType": mockHandler,
		},
		Sessions: mSessions,
	}

	// Mock handler behavior
	mockHandler.EXPECT().WSHandle(mock.AnythingOfType("*wssession.SessionWriter"), testMsg.Message, &wssession.NoCustomState{}).Return(nil).Once()

	// Call the method
	err := sut.ServeSession(mConn, mCache)

	// Assertions
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestWSMgr_ServeSession_HandlerInvocationReturnsErrAndLogs(t *testing.T) {
	// DO not parallelise due to changing the logger
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	// Mock cache
	mCache := &mocks.MockCache{}
	defer mCache.AssertExpectations(t)

	// Prepare a valid connect message
	connectMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "connect",
		Message: json.RawMessage(`{}`),
	}
	connectMessage, _ := json.Marshal(connectMsg)

	// Prepare a message that should invoke the handler
	testMsg := wssession.ReceivedMsg{
		ConnID:  "some-id",
		Type:    "testType",
		Message: json.RawMessage(`{"key":"value"}`),
	}
	testMessage, _ := json.Marshal(testMsg)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, testMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.EXPECT().Get("", mConn, mCache).Return(&wssession.Session{}, nil).Once()

	mockHandler := new(mocks.MockMessageHandler[wssession.NoCustomState])
	sut := &wssession.Mgr[wssession.NoCustomState]{
		Handlers: map[string]wssession.MessageHandler[wssession.NoCustomState]{
			"testType": mockHandler,
		},
		Sessions: mSessions,
	}

	// Mock handler behavior
	wantErr := fmt.Errorf("test handler error")
	mockHandler.EXPECT().WSHandle(mock.AnythingOfType("*wssession.SessionWriter"), testMsg.Message, &wssession.NoCustomState{}).Return(wantErr).Once()

	mLogger := &mocks.MockLogger{}
	defer mLogger.AssertExpectations(t)
	mLogger.EXPECT().Error("Non-Fatal Error handling message", "error", "test handler error").Once()

	wssession.SetLogger(mLogger)
	defer func() {
		wssession.SetLogger(wssession.DefaultLogger)
	}()
	// Call the method
	err := sut.ServeSession(mConn, mCache)

	// Assertions
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestWSMgr_ServeSession_NoHandlerRegistered(t *testing.T) {
	t.Parallel()
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	// Mock cache
	mCache := &mocks.MockCache{}
	defer mCache.AssertExpectations(t)

	// Prepare a valid connect message
	connectMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "connect",
		Message: json.RawMessage(`{}`),
	}
	connectMessage, _ := json.Marshal(connectMsg)

	// Prepare a message with an unregistered type
	unhandledMsg := wssession.ReceivedMsg{
		ConnID:  "some-id",
		Type:    "unknownType",
		Message: json.RawMessage(`{}`),
	}
	unhandledMessage, _ := json.Marshal(unhandledMsg)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, unhandledMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.EXPECT().Get("", mConn, mCache).Return(&wssession.Session{}, nil).Once()

	// Run test
	sut := &wssession.Mgr[wssession.NoCustomState]{
		Handlers: map[string]wssession.MessageHandler[wssession.NoCustomState]{},
		Sessions: mSessions,
	}
	err := sut.ServeSession(mConn, mCache)

	// Assertions
	assert.NoError(t, err)
}

func TestWSMgr_ServeSession_PassMsgToWaiter(t *testing.T) {
	t.Parallel()
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	// Mock cache
	mCache := &mocks.MockCache{}
	defer mCache.AssertExpectations(t)

	replyTo := ""
	mConn.EXPECT().WriteJSON(mock.MatchedBy(func(v interface{}) bool {
		// Type assertion and JSON validation
		r, ok := v.(*wssession.ResponseMsg)
		if !ok {
			return false
		}
		replyTo = r.ID
		return r.Type == "testType" && r.Message == "test message"
	})).Return(nil).Once()

	sess := &wssession.Session{
		ConnID: "some-id",
		Conn:   mConn,
	}
	finishedWaitingCh := make(chan bool)
	go func() {
		sess.Writer().SendStrAndWait("testType", "test message", 1*time.Minute)
		close(finishedWaitingCh)
	}()
	<-time.After(100 * time.Millisecond)

	// Manually clear the cache so it doesn't replay
	sess.SetCache(&wssession.PrunerCache{})

	mSessions := &mocks.MockSessionGetter{}
	mSessions.EXPECT().Get("some-id", mConn, mCache).Return(sess, nil).Once()

	// Create handler and SUT
	mockHandler := new(mocks.MockMessageHandler[wssession.NoCustomState])
	defer mockHandler.AssertExpectations(t)
	sut := &wssession.Mgr[wssession.NoCustomState]{
		Handlers: map[string]wssession.MessageHandler[wssession.NoCustomState]{
			"testType": mockHandler,
		},
		Sessions: mSessions,
	}
	assert.NotNil(t, sut)

	// Prepare a valid connect message
	connectMsg := wssession.ReceivedMsg{
		ConnID:  sess.ID(),
		Type:    "connect",
		Message: json.RawMessage(`{}`),
	}
	connectMessage, _ := json.Marshal(connectMsg)

	// Prepare a message that responds to the wait message
	testMsg := wssession.ReceivedMsg{
		ConnID:  sess.ID(),
		Type:    "testType",
		Message: json.RawMessage(`{"key":"value"}`),
		ReplyTo: replyTo,
	}
	testMessage, _ := json.Marshal(testMsg)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, testMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	// IMPORTANT - WSHandle should not be called because the response should go to the waiter

	// Run test
	err := sut.ServeSession(mConn, mCache)
	assert.NoError(t, err)

	// Finally we need to wait for the wait call to complete to know this works
	<-finishedWaitingCh
}

func TestWSMgr_ServeSessionOnConnectDisconnect(t *testing.T) {
	t.Parallel()
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	// Mock cache
	mCache := &mocks.MockCache{}
	defer mCache.AssertExpectations(t)

	// Prepare the connect message
	receivedMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "connect",
		Message: json.RawMessage(`{"auth_token":"123"}`),
	}
	connectMessage, _ := json.Marshal(receivedMsg)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.EXPECT().Get("", mConn, mCache).Return(&wssession.Session{}, nil).Once()

	// Run test
	sut := &wssession.Mgr[wssession.NoCustomState]{
		Sessions: mSessions,
	}
	var hasRun1, hasRun2 bool
	sut.OnConnect(func(r wssession.ReceivedMsg, customState wssession.NoCustomState) (wssession.NoCustomState, error) {
		assert.Equal(t, "", r.ConnID)
		assert.Equal(t, "connect", r.Type)
		assert.Equal(t, json.RawMessage(`{"auth_token":"123"}`), r.Message)
		hasRun1 = true
		return customState, nil
	})
	sut.OnConnect(func(r wssession.ReceivedMsg, customState wssession.NoCustomState) (wssession.NoCustomState, error) {
		assert.Equal(t, "", r.ConnID)
		assert.Equal(t, "connect", r.Type)
		assert.Equal(t, json.RawMessage(`{"auth_token":"123"}`), r.Message)
		hasRun2 = true
		return customState, nil
	})
	var hasRunDiscon1, hasRunDiscon2 bool
	sut.OnDisconnect(func() error {
		hasRunDiscon1 = true
		return nil
	})
	sut.OnDisconnect(func() error {
		hasRunDiscon2 = true
		return nil
	})
	err := sut.ServeSession(mConn, mCache)

	// Assertions
	assert.NoError(t, err)
	assert.True(t, hasRun1)
	assert.True(t, hasRun2)
	assert.True(t, hasRunDiscon1)
	assert.True(t, hasRunDiscon2)
}

func TestWSMgr_ServeSessionOnConnectErr(t *testing.T) {
	t.Parallel()
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	// Prepare the connect message
	receivedMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "connect",
		Message: json.RawMessage(`{"auth_token":"123"}`),
	}
	connectMessage, _ := json.Marshal(receivedMsg)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	// mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.On("Get", "", mConn).Return(&wssession.Session{}, nil).Once()

	// Run test
	sut := &wssession.Mgr[wssession.NoCustomState]{
		Sessions: mSessions,
	}
	var hasRun1, hasRun2 bool
	sut.OnConnect(func(r wssession.ReceivedMsg, customState wssession.NoCustomState) (wssession.NoCustomState, error) {
		assert.Equal(t, "", r.ConnID)
		assert.Equal(t, "connect", r.Type)
		assert.Equal(t, json.RawMessage(`{"auth_token":"123"}`), r.Message)
		hasRun1 = true
		return customState, fmt.Errorf("expected conn rejection error")
	})
	sut.OnConnect(func(r wssession.ReceivedMsg, customState wssession.NoCustomState) (wssession.NoCustomState, error) {
		hasRun2 = true
		return customState, nil
	})
	var hasRunDiscon1, hasRunDiscon2 bool
	sut.OnDisconnect(func() error {
		hasRunDiscon1 = true
		return nil
	})
	sut.OnDisconnect(func() error {
		hasRunDiscon2 = true
		return nil
	})
	cache := &mocks.MockCache{}
	err := sut.ServeSession(mConn, cache)

	// Assertions
	assert.EqualError(t, err, "session establishment failed: connection handler error: expected conn rejection error")
	assert.True(t, hasRun1)
	assert.False(t, hasRun2) /// This should be False - we should bail after the first connection error and not run any more onConnect functions
	assert.True(t, hasRunDiscon1)
	assert.True(t, hasRunDiscon2)
}

func TestWSMgr_CustomStateInOnConectPassedToHandler(t *testing.T) {
	t.Parallel()
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	// Mock cache
	mCache := &mocks.MockCache{}
	defer mCache.AssertExpectations(t)

	// Prepare a valid connect message
	connectMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "connect",
		Message: json.RawMessage(`{}`),
	}
	connectMessage, _ := json.Marshal(connectMsg)

	// Prepare a message that should invoke the handler
	testMsg := wssession.ReceivedMsg{
		ConnID:  "some-id",
		Type:    "testType",
		Message: json.RawMessage(`{"key":"value"}`),
	}
	testMessage, _ := json.Marshal(testMsg)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, testMessage, nil).Once()
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.EXPECT().Get("", mConn, mCache).Return(&wssession.Session{}, nil).Once()

	// Define custom struct
	type CustomState struct {
		Value string
	}

	mHandler := new(mocks.MockMessageHandler[CustomState])
	mHandler.EXPECT().WSHandle(mock.AnythingOfType("*wssession.SessionWriter"), mock.Anything, mock.MatchedBy(func(v interface{}) bool {
		// Type assertion and JSON validation
		r, ok := v.(*CustomState)
		if !ok {
			return false
		}
		return r.Value == "test"
	})).Return(nil).Once()

	// Run test
	sut := &wssession.Mgr[CustomState]{
		Sessions: mSessions,
		Handlers: map[string]wssession.MessageHandler[CustomState]{
			"testType": mHandler,
		},
	}
	sut.OnConnect(func(r wssession.ReceivedMsg, customState CustomState) (CustomState, error) {
		customState.Value = "test"
		return customState, nil
	})
	err := sut.ServeSession(mConn, mCache)

	// Assertions
	assert.NoError(t, err)
}
