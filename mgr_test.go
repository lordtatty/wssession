package wssession_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/lordtatty/wssession"
	mocks "github.com/lordtatty/wssession/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWSMgr_Serve_Success(t *testing.T) {
	// Setup
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)

	sut := &wssession.Mgr{}

	// Prepare the connect message
	receivedMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "connect",
		Message: json.RawMessage(`{}`),
	}
	connectMessage, _ := json.Marshal(receivedMsg)

	// Mock connection behavior
	mockConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	mockConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	// Call the method
	err := sut.Serve(mockConn)

	// Assertions
	assert.NoError(t, err)
}

func TestWSMgr_Serve_InvalidFirstMessageType(t *testing.T) {
	// Setup
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)

	sut := &wssession.Mgr{}

	// Prepare an invalid first message type
	invalidMsg := wssession.ReceivedMsg{
		ConnID:  "",
		Type:    "invalid",
		Message: json.RawMessage(`{}`),
	}
	invalidMessage, _ := json.Marshal(invalidMsg)

	// Mock connection behavior
	mockConn.EXPECT().ReadMessage().Return(websocket.TextMessage, invalidMessage, nil).Once()
	mockConn.EXPECT().WriteMessage(websocket.CloseMessage, mock.Anything).Return(nil).Once()

	// Call the method
	err := sut.Serve(mockConn)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "first message must be of type 'connect'")
}

func TestWSMgr_Serve_ReadMessageError(t *testing.T) {
	// Setup
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)

	sut := &wssession.Mgr{}

	// Mock connection behavior
	mockConn.EXPECT().ReadMessage().Return(0, nil, errors.New("read error")).Once()

	// Call the method
	err := sut.Serve(mockConn)

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
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)

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
	mockConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	mockConn.EXPECT().ReadMessage().Return(websocket.TextMessage, testMessage, nil).Once()
	mockConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	// Mock handler behavior
	mockHandler.EXPECT().WSHandle(mock.AnythingOfType("*wssession.SessionWriter"), testMsg.Message).Return().Once()

	// Call the method
	err := sut.Serve(mockConn)

	// Assertions
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestWSMgr_Serve_NoHandlerRegistered(t *testing.T) {
	// Setup
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)

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
	mockConn.EXPECT().ReadMessage().Return(websocket.TextMessage, connectMessage, nil).Once()
	mockConn.EXPECT().ReadMessage().Return(websocket.TextMessage, unhandledMessage, nil).Once()
	mockConn.EXPECT().ReadMessage().Return(websocket.TextMessage, []byte{}, &websocket.CloseError{Code: websocket.CloseGoingAway}).Once()

	// Call the method
	err := sut.Serve(mockConn)

	// Assertions
	assert.NoError(t, err)
}
