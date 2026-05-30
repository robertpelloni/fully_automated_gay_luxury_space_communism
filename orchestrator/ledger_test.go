package orchestrator

import (
	"testing"
)

func TestLedger(t *testing.T) {
	ledger := Ledger{}

	ledger.Add(Transaction{Amount: 100.0, Type: Revenue, Hustle: "TestHustle"})
	ledger.Add(Transaction{Amount: 30.0, Type: Expense, Hustle: "TestHustle"})

	if ledger.TotalRevenue() != 100.0 {
		t.Errorf("Expected 100.0 revenue, got %v", ledger.TotalRevenue())
	}
	if ledger.TotalExpenses() != 30.0 {
		t.Errorf("Expected 30.0 expenses, got %v", ledger.TotalExpenses())
	}
	if ledger.Profit() != 70.0 {
		t.Errorf("Expected 70.0 profit, got %v", ledger.Profit())
	}
}
