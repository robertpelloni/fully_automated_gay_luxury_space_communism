# Changelog

## [1.0.0-alpha.3] - 2026-05-29
### Added
- Automated rollback handling in `orchestrator/rollback.go`.
- Real-time status reporting to `STATUS.json` via `orchestrator/status.go`.
- Integration of rollback into the `Healer` loop.

## [1.0.0-alpha.2] - 2026-05-29
### Added
- Repository synchronization orchestration scripts (`sync.sh`, `build.sh`).
- Sync monitoring and health check logic in `orchestrator/sync_monitor.go`.
- Unit tests for sync health and edge case handling (mock merge conflicts).

## [1.0.0-alpha.1] - 2026-05-29
### Added
- Repository initialization with knowledge base.
- Basic project structure scaffolding (Orchestrator, Hustle).
- Core documentation files (`VERSION.md`, `MEMORY.md`, `CHANGELOG.md`, `TODO.md`, `HANDOFF.md`).
