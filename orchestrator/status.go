package orchestrator

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type StatusReport struct {
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
}

func WriteStatusReport(version, status, message string) error {
	report := StatusReport{
		Version:   version,
		Timestamp: time.Now(),
		Status:    status,
		Message:   message,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("STATUS.json", data, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Status report written to STATUS.json: %s\n", status)
	return nil
}
