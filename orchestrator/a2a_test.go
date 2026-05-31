package orchestrator

import (
	"testing"
	"time"
)

func TestA2ABroker(t *testing.T) {
	orch := NewOrchestrator()
	broker := NewA2ABroker(orch)

	agent1 := "researcher-1"
	agent2 := "curator-1"

	broker.Subscribe(agent1)
	ch2 := broker.Subscribe(agent2)

	msg := Message{
		ID:      "test-1",
		Source:  agent1,
		Target:  agent2,
		Type:    Query,
		Payload: "hustle://research?query=trends",
		Timestamp: time.Now(),
	}

	// Test Routing
	err := broker.Route(msg)
	if err != nil {
		t.Fatalf("Failed to route message: %v", err)
	}

	select {
	case received := <-ch2:
		if received.ID != msg.ID {
			t.Errorf("Received wrong message ID: got %s, want %s", received.ID, msg.ID)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for message on ch2")
	}

	// Test Topic Pub/Sub
	topic := "alerts"
	chTopic := broker.SubscribeTopic(topic)

	pubMsg := Message{
		ID:      "pub-1",
		Source:  agent1,
		Type:    Event,
		Topic:   topic,
		Payload: "New Alpha Found",
	}
	broker.Publish(pubMsg)

	select {
	case received := <-chTopic:
		if received.Payload != pubMsg.Payload {
			t.Errorf("Received wrong payload: got %s, want %s", received.Payload, pubMsg.Payload)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for topic message")
	}

	// Test Broadcast
	broadcastMsg := Message{
		ID:     "b-1",
		Source: agent1,
		Type:   Event,
		Payload: "System heartbeating",
	}
	broker.Broadcast(broadcastMsg)

	select {
	case <-ch2:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for broadcast on ch2")
	}

	// Verify L1 Memory persistence
	memories := orch.L1.Search("a2a")
	if len(memories) == 0 {
		t.Error("Expected A2A message to be saved to L1 memory")
	}
}

func TestA2ARoutingFailure(t *testing.T) {
	orch := NewOrchestrator()
	broker := NewA2ABroker(orch)

	msg := Message{
		ID:      "fail-1",
		Source:  "a",
		Target:  "non-existent",
		Type:    Command,
		Payload: "do something",
	}

	err := broker.Route(msg)
	if err == nil {
		t.Error("Expected error when routing to non-existent target")
	}
}
