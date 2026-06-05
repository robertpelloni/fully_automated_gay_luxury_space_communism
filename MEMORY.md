# Architectural Observations & Preferences

## Core Traits
- **Protocol-Driven:** The system heavily relies on `hustle://` URIs for decoupling. This allows for easy mesh routing and manual testing.
- **Tiered Memory:** L1 (Volatile/Events), L2 (Successes/Discoveries), L3 (System/Long-term). This separation prevents "context rot" in long-running sessions.
- **Go 1.24.3:** Strict adherence to the latest Go toolchain for performance and security.
- **A2A Mesh:** Distributed by design. No single point of failure; memory and status are replicated across peers.

## Design Preferences
- **Minimal Dependencies:** Prefer standard library or high-signal open source tools (like `sqlite-vec`).
- **Fail-Fast:** Rollback logic is first-class. If a sync or build fails, the system reverts immediately.
- **Autonomous Evolution:** The system should generate its own tasks and code improvements (ChainDiscoverer).

## Discovered Optimizations
- **LLM Sentiment Confluence:** Technical signals are noisy; filtering them through LLM-extracted market sentiment significantly reduces false positives in Trading.
- **Mesh Aggregation:** Centralizing status in L1 memory allows the Dashboard to remain stateless and high-performance.
