package wssession_test

import (
	"testing"
	"time"

	"github.com/lordtatty/wssession"
	"github.com/stretchr/testify/assert"
)

func TestCache_Add_Items_Len(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	want := []*wssession.ResponseMsg{
		{
			ID:      "1",
			ConnID:  "conn-1",
			Type:    "test",
			Message: "hello",
		},
		{
			ID:      "2",
			ConnID:  "conn-2",
			Type:    "test",
			Message: "hello",
		},
	}

	connID := "0001"

	sut := &wssession.PrunerCache{
		// Do not set auto prune time, so we can test it defaults later to 1 minute
	}
	itms, err := sut.Items(connID)
	assert.Nil(err)
	assert.Len(itms, 0)

	// Add two items
	sut.Add(connID, *want[0])
	sut.Add(connID, *want[1])

	itms, err = sut.Items(connID)
	assert.Nil(err)
	assert.Len(itms, 2)
	assert.Equal(want[0], itms[0])
	assert.Equal(want[1], itms[1])

	// check that we've defaulted to a minute autoprune
	assert.Equal(time.Minute, sut.AutoPruneDuration)
}

func TestCache_Pruning(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	want := []*wssession.ResponseMsg{
		{
			ID:      "1",
			ConnID:  "conn-1",
			Type:    "test",
			Message: "hello",
		},
		{
			ID:      "2",
			ConnID:  "conn-2",
			Type:    "test",
			Message: "hello",
		},
		{
			ID:      "3",
			ConnID:  "conn-3",
			Type:    "test",
			Message: "hello",
		},
	}

	connID := "0001"

	sut := &wssession.PrunerCache{
		AutoPruneDuration: time.Second,
	}
	itms, err := sut.Items(connID)
	assert.Nil(err)
	assert.Len(itms, 0)

	// Add two items
	sut.Add(connID, *want[0])
	sut.Add(connID, *want[1])

	itms, err = sut.Items(connID)
	assert.Nil(err)
	assert.Len(itms, 2)

	assert.Equal(want[0], itms[0])
	assert.Equal(want[1], itms[1])

	// wait 200 ms and add another item
	time.Sleep(200 * time.Millisecond)
	sut.Add(connID, *want[2])
	assert.True(sut.PrunerIsRunning())

	// Wait for the cache to be pruned
	// The cache should have only the last item
	time.Sleep(1 * time.Second)
	itms, err = sut.Items(connID)
	assert.Nil(err)
	assert.Len(itms, 1)
	assert.True(sut.PrunerIsRunning())

	// wait one more second and the cache should be empty
	time.Sleep(1 * time.Second)
	itms, err = sut.Items(connID)
	assert.Nil(err)
	assert.Len(itms, 0)
	assert.False(sut.PrunerIsRunning())
}
