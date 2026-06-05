# Deployment & Environment Setup

## Prerequisites
- **Go 1.24.3** (Mandatory)
- **SQLite 3.45+** (With vector extension support)
- **Tavily API Key** (For Research module)
- **Anthropic/OpenAI API Keys** (For LLM Waterfall)

## Local Installation
1.  Clone the repository and submodules:
    ```bash
    git clone --recursive https://github.com/robertpelloni/fully_automated_gay_luxury_space_communism
    cd fully_automated_gay_luxury_space_communism
    ```
2.  Initialize the workspace:
    ```bash
    go work init
    go work use ./orchestrator ./hustle/research ./hustle/social ./hustle/curation ./hustle/trading
    ```
3.  Build the system:
    ```bash
    ./build.sh
    ```

## Running the Machine
- **Daemon Mode (Recommended):**
  ```bash
  ./bin/orchestrator -daemon
  ```
- **Mesh Node (with seed):**
  ```bash
  ./bin/orchestrator -daemon -seed http://peer-address:8080 -api 8081
  ```
- **Interactive Dashboard:**
  ```bash
  ./bin/orchestrator -dashboard
  ```

## Production Considerations
- Ensure `token_twitter.json` and `token_linkedin.json` are populated with valid OAuth2 tokens for automated social posting.
- Monitor `ledger.json` for real-time profit/loss tracking.
