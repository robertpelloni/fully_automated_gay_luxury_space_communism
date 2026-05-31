package trading

import (
	"fmt"
	"github.com/robertpelloni/hustle/orchestrator"
	"math/rand"
	"time"
)

type PriceFetcher interface {
	GetPrice(symbol string) (float64, error)
}

type MockPriceFetcher struct{}

func (m *MockPriceFetcher) GetPrice(symbol string) (float64, error) {
	// Simple mock price generation
	return 50000.0 + rand.Float64()*1000.0, nil
}

type TradingModule struct {
	Orchestrator *orchestrator.Orchestrator
	Symbol       string
	Fetcher      PriceFetcher
	History      []float64
}

func (t *TradingModule) ExecuteStrategy() error {
	fmt.Printf("[Trading] Executing strategy for Symbol: %s\n", t.Symbol)

	price, err := t.Fetcher.GetPrice(t.Symbol)
	if err != nil {
		return fmt.Errorf("failed to fetch price: %v", err)
	}
	t.History = append(t.History, price)

	fmt.Printf("[Trading] Current Price for %s: $%.2f\n", t.Symbol, price)

	// Technical Indicator: Simple Moving Average (SMA)
	sma := t.calculateSMA(5)
	fmt.Printf("[Trading] Calculated SMA(5): $%.2f\n", sma)

	decision := "HOLD"
	if len(t.History) >= 5 {
		if price > sma*1.01 {
			decision = "SELL"
		} else if price < sma*0.99 {
			decision = "BUY"
		}
	} else {
		fmt.Println("[Trading] Insufficient history for indicators, defaulting to HOLD.")
	}

	fmt.Printf("[Trading] Strategy Decision: %s\n", decision)

	// Persist to memory
	t.Orchestrator.L1.Add(orchestrator.MemoryEntry{
		ID:        fmt.Sprintf("trade-%s-%d", t.Symbol, time.Now().Unix()),
		Content:   fmt.Sprintf("Trade Decision for %s: %s at $%.2f (SMA: $%.2f)", t.Symbol, decision, price, sma),
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

func (t *TradingModule) calculateSMA(period int) float64 {
	if len(t.History) == 0 {
		return 0
	}

	count := period
	if len(t.History) < period {
		count = len(t.History)
	}

	sum := 0.0
	for i := len(t.History) - count; i < len(t.History); i++ {
		sum += t.History[i]
	}
	return sum / float64(count)
}
