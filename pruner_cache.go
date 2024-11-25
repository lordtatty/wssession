package wssession

import (
	"sync"
	"time"
)

type PrunerCache struct {
	items             []*ResponseMsg
	mu                sync.Mutex
	pruning           bool
	AutoPruneDuration time.Duration // If not set, defaults to 1 minute
}

type PruneResult struct {
	Duration       time.Duration
	ItemsRemoved   int
	ItemsRemaining int
}

// Prune every d duration
// Removes all items from the cache that are older than d (the prior tick)
// Cheaper than a TTL cache
func (c *PrunerCache) startPruning(d time.Duration, doneCh chan bool) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	c.mu.Lock()
	if c.pruning {
		c.mu.Unlock()
		return
	}
	c.pruning = true
	defer func() { c.pruning = false }()
	l := len(c.items)
	c.mu.Unlock()
	for {
		select {
		case <-ticker.C:
			logger().Debug("Pruning cache after", "duration", d)
			c.mu.Lock()
			c.items = c.items[l:]
			l = len(c.items)
			if l == 0 {
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
		case <-doneCh:
			logger().Debug("Received done signal, stopping pruning")
			return
		}
	}
}

func (c *PrunerCache) Add(r ResponseMsg) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.pruning {
		if c.AutoPruneDuration == 0 {
			c.AutoPruneDuration = time.Minute
		}
		go c.startPruning(c.AutoPruneDuration, nil)
	}
	c.items = append(c.items, &r)
}

func (c *PrunerCache) Items() []*ResponseMsg {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.items
}

func (c *PrunerCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

func (c *PrunerCache) PrunerIsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pruning
}
