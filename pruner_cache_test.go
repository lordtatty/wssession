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

	sut := &wssession.PrunerCache{
		// Do not set auto prune time, so we can test it defaults later to 1 minute
	}
	assert.Equal(0, sut.Len())

	// Add two items
	sut.Add(*want[0])
	sut.Add(*want[1])

	assert.Equal(2, sut.Len())
	assert.Equal(want[0], sut.Items()[0])
	assert.Equal(want[1], sut.Items()[1])

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

	sut := &wssession.PrunerCache{
		AutoPruneDuration: time.Second,
	}
	assert.Equal(0, sut.Len())

	// Add two items
	sut.Add(*want[0])
	sut.Add(*want[1])

	assert.Equal(2, sut.Len())
	assert.Equal(want[0], sut.Items()[0])
	assert.Equal(want[1], sut.Items()[1])

	// wait 200 ms and add another item
	time.Sleep(200 * time.Millisecond)
	sut.Add(*want[2])
	assert.True(sut.PrunerIsRunning())

	// Wait for the cache to be pruned
	// The cache should have only the last item
	time.Sleep(1 * time.Second)
	assert.Equal(1, sut.Len())
	assert.True(sut.PrunerIsRunning())

	// wait one more second and the cache should be empty
	time.Sleep(1 * time.Second)
	assert.Equal(0, sut.Len())
	assert.False(sut.PrunerIsRunning())
}
