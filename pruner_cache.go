package wssession

import (
	"sync"
	"time"
)

type PrunerCache struct {
	nextTickIndexes   map[string]int
	items             map[string][]*ResponseMsg
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
	c.nextTickIndexes = make(map[string]int, len(c.items))
	for k := range c.items {
		c.nextTickIndexes[k] = len(c.items[k])
	}
	c.mu.Unlock()
	for {
		select {
		case <-ticker.C:
			logger().Debug("Pruning cache after", "duration", d)
			c.mu.Lock()
			for k, v := range c.items {
				l := c.nextTickIndexes[k]
				c.items[k] = v[l:]
				c.nextTickIndexes[k] = len(c.items[k])
				if c.nextTickIndexes[k] == 0 {
					delete(c.items, k)
					delete(c.nextTickIndexes, k)
				}
			}
			if len(c.items) == 0 {
				logger().Debug("Cache is empty, stopping pruning")
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

func (c *PrunerCache) ensureItems() {
	if c.items == nil {
		c.items = make(map[string][]*ResponseMsg)
	}
}

func (c *PrunerCache) Add(connID string, r ResponseMsg) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ensureItems()
	if !c.pruning {
		if c.AutoPruneDuration == 0 {
			c.AutoPruneDuration = time.Minute
		}
		go c.startPruning(c.AutoPruneDuration, nil)
	}
	if _, ok := c.items[connID]; !ok {
		c.items[connID] = []*ResponseMsg{}
	}
	c.items[connID] = append(c.items[connID], &r)
}

func (c *PrunerCache) Items(connID string) []*ResponseMsg {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ensureItems()
	if _, ok := c.items[connID]; !ok {
		return nil
	}
	return c.items[connID]
}

func (c *PrunerCache) Len(connID string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ensureItems()
	if _, ok := c.items[connID]; !ok {
		return 0
	}
	return len(c.items[connID])
}

func (c *PrunerCache) PrunerIsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pruning
}
