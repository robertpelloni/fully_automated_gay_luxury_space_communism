# Changelog

## [1.0.0-alpha.11] - 2026-05-30
### Added
- Integrated Revenue, Expenses, and Profit into `STATUS.json` reporting.
- Pluggable `LLMProvider` interface with `MockLLM` implementation.
- Orchestrator-wide LLM access for all hustle modules.
- LLM-driven synthesis in Research and Social Media modules.
- Parameterized execution via `-params` flag in Orchestrator CLI.

## [1.0.0-alpha.10] - 2026-05-30
### Added
- Financial Ledger implementation for revenue and expense tracking.
- CLI flags (`-hustle`, `-sync`) for the Hustle Orchestrator.
- Provider interface for Social Media module with Twitter/LinkedIn mocks.
- Automated cost integration for hustle tasks.

## [1.0.0-alpha.9] - 2026-05-30
### Added
- Functional Healer loop with retry logic and memory logging.
- Memory Search with keyword filtering.
- Initial Social Media hustle module scaffolding.

## [1.0.0-alpha.8] - 2026-05-30
### Added
- JSON persistence for tiered memory.
- Standalone entry points for Orchestrator and Research modules.
- Automatic integration between research results and orchestrator memory.

## [1.0.0-alpha.7] - 2026-05-30
### Added
- Utility-based heat scoring for memories with temporal decay.
- Enhanced Research Hustle pipeline with multi-provider support.

## [1.0.0-alpha.6] - 2026-05-30
### Added
- Real Git rollback implementation.
- Failure trap in `sync.sh`.

## [1.0.0-alpha.5] - 2026-05-30
### Added
- GitHub Actions CI pipeline.
- End-to-End (E2E) tests for the research hustle pipeline.

## [1.0.0-alpha.4] - 2026-05-29
### Added
- Enhanced `sync.sh` with branch reconciliation logic.
- Real git status and submodule monitoring.

## [1.0.0-alpha.3] - 2026-05-29
### Added
- Automated rollback handling skeleton.
- Real-time status reporting to `STATUS.json`.

## [1.0.0-alpha.2] - 2026-05-29
### Added
- Repository synchronization orchestration scripts.
- Sync monitoring and health check logic.

## [1.0.0-alpha.1] - 2026-05-29
### Added
- Repository initialization with knowledge base.
- Basic project structure scaffolding.
