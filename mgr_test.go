package wssession_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lordtatty/wssession"
	mocks "github.com/lordtatty/wssession/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWSMgr_Serve_Success(t *testing.T) {
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	sut := &wssession.Mgr{}

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
	mSessions.On("Get", "", mConn).Return(&wssession.Session{}, nil).Once()

	// Call the method
	err := sut.Serve(mConn, mSessions)

	// Assertions
	assert.NoError(t, err)
}

func TestWSMgr_Serve_InvalidFirstMessageType(t *testing.T) {
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	sut := &wssession.Mgr{}

	// Prepare an invalid first message type
	invalidMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "invalid",
		Message: json.RawMessage(`{}`),
	}
	invalidMessage, _ := json.Marshal(invalidMsg)

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(websocket.TextMessage, invalidMessage, nil).Once()
	mConn.EXPECT().WriteMessage(websocket.CloseMessage, mock.Anything).Return(nil).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.On("Get", "", mConn).Return(&wssession.Session{}, nil).Once()

	// Call the method
	err := sut.Serve(mConn, mSessions)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "first message must be of type 'connect'")
}

func TestWSMgr_Serve_ReadMessageError(t *testing.T) {
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	sut := &wssession.Mgr{}

	// Mock connection behavior
	mConn.EXPECT().ReadMessage().Return(0, nil, errors.New("read error")).Once()

	mSessions := &mocks.MockSessionGetter{}
	mSessions.On("Get", "", mConn).Return(&wssession.Session{}, nil).Once()

	// Call the method
	err := sut.Serve(mConn, mSessions)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading message")
}

func TestWSMgr_RegisterHandler(t *testing.T) {
	sut := &wssession.Mgr{}

	// Create a mock handler
	mockHandler := new(mocks.MockMessageHandler)

	// Register the handler
	sut.RegisterHandler("testType", mockHandler)

	// Assertions
	assert.NotNil(t, sut.Handlers)
	assert.Equal(t, mockHandler, sut.Handlers["testType"])
}

func TestWSMgr_Serve_HandlerInvocation(t *testing.T) {
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	mockHandler := new(mocks.MockMessageHandler)
	sut := &wssession.Mgr{
		Handlers: map[string]wssession.MessageHandler{
			"testType": mockHandler,
		},
	}

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
	mSessions.On("Get", "", mConn).Return(&wssession.Session{}, nil).Once()

	// Mock handler behavior
	mockHandler.EXPECT().WSHandle(mock.AnythingOfType("*wssession.SessionWriter"), testMsg.Message).Return().Once()

	// Call the method
	err := sut.Serve(mConn, mSessions)

	// Assertions
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestWSMgr_Serve_NoHandlerRegistered(t *testing.T) {
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	sut := &wssession.Mgr{
		Handlers: map[string]wssession.MessageHandler{},
	}

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
	mSessions.On("Get", "", mConn).Return(&wssession.Session{}, nil).Once()

	// Call the method
	err := sut.Serve(mConn, mSessions)

	// Assertions
	assert.NoError(t, err)
}

func TestWSMgr_Serve_PassMsgToWaiter(t *testing.T) {
	// Setup
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)

	mockHandler := new(mocks.MockMessageHandler)
	defer mockHandler.AssertExpectations(t)
	sut := &wssession.Mgr{
		Handlers: map[string]wssession.MessageHandler{
			"testType": mockHandler,
		},
	}

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
	sess.Cache = wssession.PrunerCache{}

	assert.NotNil(t, sut)

	mSessions := &mocks.MockSessionGetter{}
	mSessions.On("Get", "some-id", mConn).Return(sess, nil).Once()

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

	// Call the method
	err := sut.Serve(mConn, mSessions)
	assert.NoError(t, err)

	// Finally we need to wait for the wait call to complete to know this works
	<-finishedWaitingCh
}
