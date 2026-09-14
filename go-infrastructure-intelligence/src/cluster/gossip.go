package cluster

import (
	"fmt"
	"sync"
	"time"
)

type NodeState string

const (
	StateAlive   NodeState = "ALIVE"
	StateSuspect NodeState = "SUSPECT"
	StateDead    NodeState = "DEAD"
)

type Node struct {
	ID        string
	Address   string
	State     NodeState
	Heartbeat uint64
	LastSeen  time.Time
}

type GossipCluster struct {
	mu       sync.RWMutex
	SelfID   string
	Nodes    map[string]*Node
	seqCount uint64
}

func NewGossipCluster(selfID string, address string) *GossipCluster {
	c := &GossipCluster{
		SelfID: selfID,
		Nodes:  make(map[string]*Node),
	}
	c.Nodes[selfID] = &Node{
		ID:        selfID,
		Address:   address,
		State:     StateAlive,
		Heartbeat: 1,
		LastSeen:  time.Now(),
	}
	return c
}

func (c *GossipCluster) AddPeer(id string, address string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Nodes[id] = &Node{
		ID:        id,
		Address:   address,
		State:     StateAlive,
		Heartbeat: 1,
		LastSeen:  time.Now(),
	}
}

func (c *GossipCluster) RecordHeartbeat(id string, seq uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if n, ok := c.Nodes[id]; ok {
		if seq > n.Heartbeat {
			n.Heartbeat = seq
			n.LastSeen = time.Now()
			n.State = StateAlive
		}
	}
}

func (c *GossipCluster) DetectFailures(timeout time.Duration) []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var failed []string
	now := time.Now()

	for id, n := range c.Nodes {
		if id == c.SelfID {
			continue
		}
		if now.Sub(n.LastSeen) > timeout && n.State == StateAlive {
			n.State = StateSuspect
			failed = append(failed, fmt.Sprintf("Node %s marked SUSPECT (No heartbeat in %v)", id, timeout))
		}
	}
	return failed
}
