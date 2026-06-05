# Release Notes - v1.0.0-alpha.41

## Core Components (bin/)
- `orchestrator`: Central management and protocol routing hub.
- `research`: LLM-driven market sentiment and ticker extraction.
- `trading`: Real-time execution with CoinGecko pricing and technical indicators.
- `social`: Automated content distribution (Twitter/LinkedIn).
- `curator`: RSS/Atom feed ingestion and synthesis.

## Key Features
- **Live Market Data**: Real-world USD pricing via CoinGecko API.
- **Native Vector Search**: High-performance persistence via `sqlite-vec` with Go-level fallback.
- **Executive Protocol**: Hardened repository sync engine for monorepo health.
- **Federated Mesh**: Stable A2ABroker and Memory Swarm for multi-node collaboration.

## External Dependencies
- Go 1.24.3
- Git 2.x
- SQLite 3 with extension loading support (optional, for native vectors).

## API Credentials (ENV)
- `ANTHROPIC_API_KEY`: Required for core intelligence.
- `TAVILY_API_KEY`: Required for web research.
- `COINGECKO_API_KEY`: Optional for Pro price feeds.
