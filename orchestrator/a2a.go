package orchestrator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	Orchestrator  *Orchestrator
	mu            sync.RWMutex
	Subscriptions map[string]chan Message
	Peers         map[string]string // TargetID -> RemoteURL
}

func NewA2ABroker(orch *Orchestrator) *A2ABroker {
	return &A2ABroker{
		Orchestrator:  orch,
		Subscriptions: make(map[string]chan Message),
		Peers:         make(map[string]string),
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

// RegisterPeer mapping a target ID to a remote URL
func (b *A2ABroker) RegisterPeer(targetID, url string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Peers[targetID] = url
	fmt.Printf("[A2A] Registered remote peer %s at %s\n", targetID, url)
}

// Route delivers a message to the appropriate target (local or remote)
func (b *A2ABroker) Route(msg Message) error {
	b.mu.RLock()
	localCh, isLocal := b.Subscriptions[msg.Target]
	remoteURL, isRemote := b.Peers[msg.Target]
	b.mu.RUnlock()

	fmt.Printf("[A2A] Routing %s from %s to %s\n", msg.Type, msg.Source, msg.Target)

	// Persist message to L1 memory for auditing
	b.Orchestrator.L1.Add(MemoryEntry{
		ID:        fmt.Sprintf("msg-%s", msg.ID),
		Content:   fmt.Sprintf("A2A %s: %s -> %s | %s", msg.Type, msg.Source, msg.Target, msg.Payload),
		Timestamp: time.Now(),
		Tags:      []string{"a2a", msg.Source, msg.Target},
	})

	if isLocal {
		select {
		case localCh <- msg:
			return nil
		case <-time.After(1 * time.Second):
			return fmt.Errorf("target %s local subscription buffer full", msg.Target)
		}
	}

	if isRemote {
		return b.forwardToRemote(remoteURL, msg)
	}

	return fmt.Errorf("target %s not found in local mesh or known peers", msg.Target)
}

func (b *A2ABroker) forwardToRemote(url string, msg Message) error {
	fmt.Printf("[A2A] Forwarding message to remote peer: %s\n", url)
	data, _ := json.Marshal(msg)

	client := &http.Client{Timeout: 5 * time.Second}
	// We assume remote has a /message endpoint for A2A
	resp, err := client.Post(url+"/message", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to forward to remote: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote peer returned error status: %d", resp.StatusCode)
	}

	return nil
}

// Broadcast sends a message to all local subscribers and known peers
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
			fmt.Printf("[A2A] Warning: Failed to deliver local broadcast to %s\n", id)
		}
	}

	for id, url := range b.Peers {
		if id == msg.Source {
			continue
		}
		go b.forwardToRemote(url, msg)
	}
}
