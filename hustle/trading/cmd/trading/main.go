package main

import (
	"fmt"
	"github.com/robertpelloni/hustle/hustle/trading"
	"github.com/robertpelloni/hustle/orchestrator"
)

func main() {
	fmt.Println("=== AI Trading Hustle Module ===")
	orch := orchestrator.NewOrchestrator()

	trader := &trading.TradingModule{
		Orchestrator: orch,
		Symbol:       "BTC",
	}

	if err := trader.ExecuteStrategy(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Module execution complete. Profit: $%.2f\n", orch.Ledger.Profit())
}
