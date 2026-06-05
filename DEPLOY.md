# Deployment & Environment Setup

## Prerequisites
- **Go 1.24.3** (Mandatory toolchain)
- **SQLite 3.45+** (With vector extension or using internal fallback)
- **API Keys**: Tavily (Research), Anthropic/OpenAI (LLM Waterfall)

## Production Installation
1.  **Clone & Update**:
    ```bash
    git clone --recursive https://github.com/robertpelloni/fully_automated_gay_luxury_space_communism
    cd fully_automated_gay_luxury_space_communism
    ```
2.  **Environment Check**:
    ```bash
    go version # Ensure 1.24.3
    ./build.sh
    ```

## Mesh Deployment (Cluster Setup)
To launch a federated luxury mesh:
1.  **Seed Node**:
    ```bash
    ./bin/orchestrator -daemon -api 8080 -real-prices
    ```
2.  **Peer Node**:
    ```bash
    ./bin/orchestrator -daemon -api 8081 -seed http://seed-ip:8080 -real-prices
    ```

## Performance Monitoring
- **Global Ledger**: Monitor `ledger.json` across all nodes.
- **Terminal Dashboard**: Run `./bin/orchestrator -dashboard` to see aggregated mesh profits.
- **Wealth Preservation**: The system automatically audits ROI and terminates underperforming hustles. No manual intervention required.

## Automated Maintenance
The Orchestrator includes a self-healing `sync` protocol. To trigger manually:
```bash
./bin/orchestrator -uri "hustle://sync"
```
