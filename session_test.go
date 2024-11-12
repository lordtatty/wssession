package wssession_test

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lordtatty/wssession"
	mocks "github.com/lordtatty/wssession/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSession_SendJSONAndWait(t *testing.T) {
	t.Parallel()
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)
	sut := &wssession.Session{
		Conn: mockConn,
	}

	// Test data
	testMsg := map[string]string{"message": "hello"}
	msgType := "testType"
	expectedJSON, err := json.Marshal(testMsg)
	require.NoError(t, err)

	replyTo := ""
	mockConn.EXPECT().WriteJSON(mock.MatchedBy(func(v interface{}) bool {
		// Type assertion and JSON validation
		r, ok := v.(*wssession.ResponseMsg)
		if !ok {
			return false
		}
		replyTo = r.ID
		return r.Type == msgType && r.Message == string(expectedJSON)
	})).Return(nil).Once()

	// Call the method
	go func() {
		<-time.After(2 * time.Second)
		sut.CompleteWaiterIfMatch(wssession.ReceivedMsg{
			ReplyTo: replyTo,
			Message: json.RawMessage(`{"msg":"expected-reply"}`),
		})
	}()
	start := time.Now()
	resp, err := sut.Writer().SendJSONAndWait(msgType, testMsg, time.Minute)
	assert.True(t, time.Since(start) > 1*time.Second)
	assert.True(t, time.Since(start) < 3*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, `{"msg":"expected-reply"}`, string(*resp))
}

func TestSession_SendJSONAndWaitTimesOut(t *testing.T) {
	t.Parallel()
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)
	sut := &wssession.Session{
		Conn: mockConn,
	}

	// Test data
	testMsg := map[string]string{"message": "hello"}
	msgType := "testType"
	expectedJSON, err := json.Marshal(testMsg)
	require.NoError(t, err)

	mockConn.EXPECT().WriteJSON(mock.MatchedBy(func(v interface{}) bool {
		// Type assertion and JSON validation
		r, ok := v.(*wssession.ResponseMsg)
		if !ok {
			return false
		}
		return r.Type == msgType && r.Message == string(expectedJSON)
	})).Return(nil).Once()

	// Call the method
	start := time.Now()
	resp, err := sut.Writer().SendJSONAndWait(msgType, testMsg, 2*time.Second)
	assert.True(t, time.Since(start) > 1*time.Second)
	assert.True(t, time.Since(start) < 3*time.Second)
	assert.ErrorContains(t, err, "error in waiter response: error in waiter response: timeout waiting for reply to message")
	assert.Nil(t, resp)
}

func TestSession_SendStrAndWait(t *testing.T) {
	t.Parallel()
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)
	sut := &wssession.Session{
		Conn: mockConn,
	}

	// Test data
	msgType := "testType"
	testMsg := "hello"

	replyTo := ""
	mockConn.EXPECT().WriteJSON(mock.MatchedBy(func(v interface{}) bool {
		r, ok := v.(*wssession.ResponseMsg)
		if !ok {
			return false
		}
		replyTo = r.ID
		return r.Type == msgType && r.Message == testMsg
	})).Return(nil).Once()

	// Call the method
	go func() {
		<-time.After(2 * time.Second)
		sut.CompleteWaiterIfMatch(wssession.ReceivedMsg{
			ReplyTo: replyTo,
			Message: json.RawMessage(`{"msg":"expected-reply"}`),
		})
	}()
	start := time.Now()
	resp, err := sut.Writer().SendStrAndWait(msgType, testMsg, time.Minute)
	assert.True(t, time.Since(start) > 1*time.Second)
	assert.True(t, time.Since(start) < 3*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, `{"msg":"expected-reply"}`, string(*resp))
}

func TestSession_SendStrAndWaitTimesOut(t *testing.T) {
	t.Parallel()
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)
	sut := &wssession.Session{
		Conn: mockConn,
	}

	// Test data
	msgType := "testType"
	testMsg := "hello"

	mockConn.EXPECT().WriteJSON(mock.MatchedBy(func(v interface{}) bool {
		r, ok := v.(*wssession.ResponseMsg)
		if !ok {
			return false
		}
		return r.Type == msgType && r.Message == testMsg
	})).Return(nil).Once()

	// Call the method
	start := time.Now()
	resp, err := sut.Writer().SendStrAndWait(msgType, testMsg, 2*time.Second)
	assert.True(t, time.Since(start) > 1*time.Second)
	assert.True(t, time.Since(start) < 3*time.Second)
	assert.ErrorContains(t, err, "error in waiter response: error in waiter response: timeout waiting for reply to message")
	assert.Nil(t, resp)
}
func TestSession_SendJSON(t *testing.T) {
	t.Parallel()
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)
	sut := &wssession.Session{
		Conn: mockConn,
	}

	// Test data
	testMsg := map[string]string{"message": "hello"}
	msgType := "testType"
	expectedJSON, err := json.Marshal(testMsg)
	require.NoError(t, err)

	mockConn.EXPECT().WriteJSON(mock.MatchedBy(func(v interface{}) bool {
		// Type assertion and JSON validation
		r, ok := v.(*wssession.ResponseMsg)
		if !ok {
			return false
		}
		return r.Type == msgType && r.Message == string(expectedJSON)
	})).Return(nil).Once()

	// Call the method
	err = sut.Writer().SendJSON(msgType, testMsg)

	// Assertions
	assert.NoError(t, err)
	mockConn.AssertExpectations(t)
}

func TestSession_SendStr(t *testing.T) {
	t.Parallel()
	mockConn := new(mocks.MockWebsocketConn)
	defer mockConn.AssertExpectations(t)
	sut := &wssession.Session{
		Conn: mockConn,
	}

	msgType := "testType"
	testMsg := "hello"

	mockConn.EXPECT().WriteJSON(mock.MatchedBy(func(v interface{}) bool {
		r, ok := v.(*wssession.ResponseMsg)
		if !ok {
			return false
		}
		return r.Type == msgType && r.Message == testMsg
	})).Return(nil).Once()

	// Call the method
	err := sut.Writer().SendStr(msgType, testMsg)

	// Assertions
	assert.NoError(t, err)
}

func TestSession_UpdateConnAndReplayCache(t *testing.T) {
	t.Parallel()
	connID := uuid.NewString()
	wantMsgs := []*wssession.ResponseMsg{
		{
			ID:      uuid.NewString(),
			ConnID:  connID,
			Type:    "chat-response",
			Message: "message 1",
		},
		{
			ID:      uuid.NewString(),
			ConnID:  connID,
			Type:    "chat-response",
			Message: "message 2",
		},
		{
			ID:      uuid.NewString(),
			ConnID:  connID,
			Type:    "chat-response",
			Message: "message 3",
		},
		{
			ID:      uuid.NewString(),
			ConnID:  connID,
			Type:    "chat-response",
			Message: "message 4",
		},
	}

	// Mock conn1 - this is initial and old conn, and should not be sent any messages
	mockConn1 := new(mocks.MockWebsocketConn)
	defer mockConn1.AssertExpectations(t)

	// Mock conn2 - the new conn and should be sent all replayed messages
	mockConn2 := new(mocks.MockWebsocketConn)
	defer mockConn2.AssertExpectations(t)

	mockConn2.EXPECT().WriteJSON(wantMsgs[0]).Return(nil).Once()
	mockConn2.EXPECT().WriteJSON(wantMsgs[1]).Return(nil).Once()
	mockConn2.EXPECT().WriteJSON(wantMsgs[2]).Return(nil).Once()
	mockConn2.EXPECT().WriteJSON(wantMsgs[3]).Return(nil).Once()

	// Initialize the session and send messages

	sut := &wssession.Session{
		ConnID: connID,
		Conn:   mockConn1,
		Cache:  wssession.PrunerCache{},
	}
	for _, msg := range wantMsgs {
		sut.Cache.Add(*msg)
	}

	err := sut.UpdateConnAndReplayCache(mockConn2)

	// Assertions
	assert.NoError(t, err)
}

// Regular expression to match a UUID (version 4).
var uuidRegex = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)

func respMsgMatcher(connID, msgType, msg string) func(resp *wssession.ResponseMsg) bool {
	return func(resp *wssession.ResponseMsg) bool {
		return resp.ConnID == connID &&
			resp.Type == msgType &&
			resp.Message == msg &&
			uuidRegex.MatchString(resp.ID)
	}
}

func TestSession_SessionWriterRespectsUpdatedConn(t *testing.T) {
	t.Parallel()
	connID := uuid.NewString()

	// Mock conn1 - this is the initial conn
	mConn1 := new(mocks.MockWebsocketConn)
	defer mConn1.AssertExpectations(t)

	// Mock conn2 - the new conn
	mConn2 := new(mocks.MockWebsocketConn)
	defer mConn2.AssertExpectations(t)

	mConn2.EXPECT().WriteJSON(mock.MatchedBy(respMsgMatcher(connID, "chat-response", "message 1"))).Return(nil).Times(1)
	mConn2.EXPECT().WriteJSON(mock.MatchedBy(respMsgMatcher(connID, "chat-response", `{"msg":"message 2"}`))).Return(nil).Times(1)

	// Initialize the session and send messages
	sess := &wssession.Session{
		ConnID: connID,
		Conn:   mConn1,
	}

	sut := sess.Writer()
	sess.UpdateConnAndReplayCache(mConn2)
	sut.SendStr("chat-response", "message 1")
	sut.SendJSON("chat-response", struct {
		Msg string `json:"msg"`
	}{Msg: "message 2"})
}

func TestSessionWriter_SendStr(t *testing.T) {
	t.Parallel()
	connID := uuid.NewString()

	// Conn - we expect to see the message written to the conn
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)
	mConn.EXPECT().WriteJSON(mock.MatchedBy(respMsgMatcher(connID, "chat-response", "message 1"))).Return(nil).Times(1)

	// Initialize the session and send messages
	sess := &wssession.Session{
		ConnID: connID,
		Conn:   mConn,
	}

	// Get the SUT
	sut := sess.Writer()
	sut.SendStr("chat-response", "message 1")
}

func TestSessionWriter_SendJSON(t *testing.T) {
	t.Parallel()
	connID := uuid.NewString()

	// Conn - we expect to see the message written to the conn
	mConn := new(mocks.MockWebsocketConn)
	defer mConn.AssertExpectations(t)
	mConn.EXPECT().WriteJSON(mock.MatchedBy(respMsgMatcher(connID, "chat-response", `{"msg":"message 2"}`))).Return(nil).Times(1)

	// Initialize the session and send messages
	sess := &wssession.Session{
		ConnID: connID,
		Conn:   mConn,
	}

	// Get the SUT
	sut := sess.Writer()
	sut.SendJSON("chat-response", struct {
		Msg string `json:"msg"`
	}{Msg: "message 2"})
}

func TestSessions_Get(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Mock conn
	mConn1 := mocks.MockWebsocketConn{}
	mConn2 := mocks.MockWebsocketConn{}

	var sut wssession.Sessions

	sess1, err := sut.Get("", &mConn1)
	assert.NoError(err)
	assert.Equal(sess1.Conn, &mConn1)
	assert.Equal(sess1.Cache.Len(), 0)

	sess2, err := sut.Get("", &mConn2)
	assert.NoError(err)
	assert.Equal(sess2.Conn, &mConn2)
	assert.Equal(sess2.Cache.Len(), 0)

	// Get the sessions again
	rtvSess1, err := sut.Get(sess1.ID(), &mConn1)
	assert.NoError(err)
	assert.Equal(rtvSess1.Conn, &mConn1)
	assert.Equal(rtvSess1.Cache.Len(), 0)
	assert.Equal(rtvSess1.ConnID, sess1.ID())

	rtvSess2, err := sut.Get(sess2.ID(), &mConn2)
	assert.NoError(err)
	assert.Equal(rtvSess2.Conn, &mConn2)
	assert.Equal(rtvSess2.Cache.Len(), 0)
	assert.Equal(rtvSess2.ConnID, sess2.ID())
}
