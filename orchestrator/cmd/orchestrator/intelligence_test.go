package main

import (
	"github.com/robertpelloni/hustle/hustle/research"
	"github.com/robertpelloni/hustle/hustle/trading"
	"github.com/robertpelloni/hustle/orchestrator"
	"strings"
	"testing"
)

// TestIntelligenceConfluenceE2E verifies the end-to-end intelligence flow:
// Research (Sentiment) -> Memory -> Trading (Confluence Decision)
func TestIntelligenceConfluenceE2E(t *testing.T) {
	orch := orchestrator.NewOrchestrator()
	broker := orchestrator.NewA2ABroker(orch)

	// 1. Setup Research module
	searcher := research.NewResearchSearch(research.Tavily, orch, broker)

	// 2. Setup Trading module
	trader := &trading.TradingModule{
		Orchestrator: orch,
		Broker:       broker,
		Symbol:       "BTC",
		Fetcher:      &trading.MockPriceFetcher{},
		// Give enough history for technical indicators
		History: make([]float64, 15),
		RSIHistory: make([]float64, 15),
	}
	// Initialize history with values that would trigger a signal
	for i := range trader.History {
		trader.History[i] = 40000.0 // Default price
		trader.RSIHistory[i] = 50.0 // Neutral RSI
	}

	// Case A: Technical Oversold (RSI < 30) but Sentiment is BEARISH
	// Expected: Decision should remain HOLD because confluence fails
	// We force the indicators into an oversold state
	for i := 0; i < 15; i++ {
		trader.History[i] = 40000.0
		trader.RSIHistory[i] = 25.0
	}
	trader.Fetcher = &fixedPriceFetcher{price: 35000.0} // Price < SMA

	orch.L1.Add(orchestrator.MemoryEntry{
		ID: "bear-sentiment",
		Content: "Sentiment for BTC: BEARISH",
		Tags: []string{"research", "sentiment", "BTC"},
	})

	err := trader.ExecuteStrategy()
	if err != nil {
		t.Fatalf("failed to execute strategy: %v", err)
	}

	lastTrade := getLastTradeDecision(orch, "BTC")
	if !strings.Contains(lastTrade, "HOLD") {
		t.Errorf("Expected HOLD due to bear sentiment mismatch, got: %s", lastTrade)
	}

	// Case B: Technical Oversold (RSI < 30) and Sentiment is BULLISH
	// Expected: Decision should be BUY
	orch.L1.Entries = []orchestrator.MemoryEntry{} // Clear L1
	orch.L1.Add(orchestrator.MemoryEntry{
		ID: "bull-sentiment",
		Content: "Sentiment for BTC: BULLISH",
		Tags: []string{"research", "sentiment", "BTC"},
	})

	err = trader.ExecuteStrategy()
	if err != nil {
		t.Fatalf("failed to execute strategy: %v", err)
	}

	lastTrade = getLastTradeDecision(orch, "BTC")
	if !strings.Contains(lastTrade, "BUY") {
		t.Errorf("Expected BUY due to bullish confluence, got: %s", lastTrade)
	}

	// Case C: Research query generates sentiment automatically
	results, _ := searcher.Query("What is the current $BTC sentiment?")
	if len(results) == 0 {
		t.Skip("Mock search results empty")
	}

	// Verify sentiment was logged to L1
	foundSentiment := false
	for _, e := range orch.L1.Entries {
		for _, tag := range e.Tags {
			if tag == "sentiment" {
				foundSentiment = true
				break
			}
		}
	}
	if !foundSentiment {
		t.Errorf("Research query failed to generate sentiment in L1 memory")
	}
}

type fixedPriceFetcher struct {
	price float64
}

func (f *fixedPriceFetcher) GetPrice(symbol string) (float64, error) {
	return f.price, nil
}

func getLastTradeDecision(o *orchestrator.Orchestrator, symbol string) string {
	entries := o.L1.Search("Trade Decision for " + symbol)
	if len(entries) == 0 {
		return "NONE"
	}
	return entries[len(entries)-1].Content
}
