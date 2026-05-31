package trading

import (
	"fmt"
	"github.com/robertpelloni/hustle/orchestrator"
	"math/rand"
	"time"
)

type TradingModule struct {
	Orchestrator *orchestrator.Orchestrator
	Symbol       string
}

func (t *TradingModule) ExecuteStrategy() error {
	fmt.Printf("[Trading] Executing strategy for Symbol: %s\n", t.Symbol)

	// Mock price fetch
	price := 50000.0 + rand.Float64()*1000.0
	fmt.Printf("[Trading] Current Price for %s: $%.2f\n", t.Symbol, price)

	decision := "HOLD"
	if price > 50500 {
		decision = "SELL"
	} else if price < 50200 {
		decision = "BUY"
	}

	fmt.Printf("[Trading] Strategy Decision: %s\n", decision)

	// Persist to memory
	t.Orchestrator.L1.Add(orchestrator.MemoryEntry{
		ID:        fmt.Sprintf("trade-%s-%d", t.Symbol, time.Now().Unix()),
		Content:   fmt.Sprintf("Trade Decision for %s: %s at $%.2f", t.Symbol, decision, price),
		Timestamp: time.Now(),
		Tags:      []string{"trading", t.Symbol, decision},
	})

	if decision != "HOLD" {
		t.Orchestrator.Ledger.Add(orchestrator.Transaction{
			Amount: 0.10, // Simulating execution fee
			Type:   orchestrator.Expense,
			Hustle: "Trading",
			Note:   fmt.Sprintf("%s %s", decision, t.Symbol),
		})
	}

	return nil
}
