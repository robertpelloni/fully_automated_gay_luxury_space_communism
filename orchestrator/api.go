package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type API struct {
	Orchestrator *Orchestrator
	Protocol     *HustleProtocol
}

func NewAPI(orch *Orchestrator, protocol *HustleProtocol) *API {
	return &API{
		Orchestrator: orch,
		Protocol:     protocol,
	}
}

func (a *API) Start(port string) error {
	http.HandleFunc("/dispatch", a.handleDispatch)
	http.HandleFunc("/status", a.handleStatus)

	fmt.Printf("[API] Orchestrator listening on port %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}

func (a *API) handleDispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URI string `json:"uri"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("[API] Received dispatch request for URI: %s\n", req.URI)

	if err := a.Protocol.HandleURI(req.URI); err != nil {
		http.Error(w, fmt.Sprintf("Protocol error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "URI %s dispatched successfully", req.URI)
}

func (a *API) handleStatus(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "Active",
		"version":  "v1.0.0-alpha.27",
		"profit":   a.Orchestrator.Ledger.Profit(),
		"memories": len(a.Orchestrator.L1.Entries),
	})
}
