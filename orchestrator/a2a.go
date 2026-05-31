package orchestrator

import (
	"fmt"
	"sync"
	"time"
)

// MessageType defines the nature of the A2A communication
type MessageType string

const (
	Query    MessageType = "Query"
	Command  MessageType = "Command"
	Response MessageType = "Response"
	Event    MessageType = "Event"
)

// Message is the standard unit of communication in the A2A Mesh
type Message struct {
	ID        string      `json:"id"`
	Source    string      `json:"source"`
	Target    string      `json:"target"`
	Type      MessageType `json:"type"`
	Payload   string      `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

// A2ABroker manages message routing between agents
type A2ABroker struct {
	Orchestrator *Orchestrator
	mu           sync.RWMutex
	Subscriptions map[string]chan Message
}

func NewA2ABroker(orch *Orchestrator) *A2ABroker {
	return &A2ABroker{
		Orchestrator:  orch,
		Subscriptions: make(map[string]chan Message),
	}
}

// Subscribe allows an agent to listen for messages targeting it
func (b *A2ABroker) Subscribe(agentID string) chan Message {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Message, 100)
	b.Subscriptions[agentID] = ch
	fmt.Printf("[A2A] Agent %s subscribed to mesh\n", agentID)
	return ch
}

// Route delivers a message to the appropriate target
func (b *A2ABroker) Route(msg Message) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	fmt.Printf("[A2A] Routing %s from %s to %s\n", msg.Type, msg.Source, msg.Target)

	// Persist message to L1 memory for auditing
	b.Orchestrator.L1.Add(MemoryEntry{
		ID:        fmt.Sprintf("msg-%s", msg.ID),
		Content:   fmt.Sprintf("A2A %s: %s -> %s | %s", msg.Type, msg.Source, msg.Target, msg.Payload),
		Timestamp: time.Now(),
		Tags:      []string{"a2a", msg.Source, msg.Target},
	})

	if ch, ok := b.Subscriptions[msg.Target]; ok {
		select {
		case ch <- msg:
			return nil
		case <-time.After(1 * time.Second):
			return fmt.Errorf("target %s subscription buffer full", msg.Target)
		}
	}

	return fmt.Errorf("target %s not found in mesh", msg.Target)
}

// Broadcast sends a message to all subscribers
func (b *A2ABroker) Broadcast(msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for id, ch := range b.Subscriptions {
		if id == msg.Source {
			continue
		}
		select {
		case ch <- msg:
		default:
			fmt.Printf("[A2A] Warning: Failed to deliver broadcast to %s\n", id)
		}
	}
}
