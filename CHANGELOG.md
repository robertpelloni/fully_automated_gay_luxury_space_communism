# Changelog

## [1.0.0-alpha.8] - 2026-05-30
### Added
- JSON persistence for tiered memory (`Save`/`Load` methods).
- Standalone entry points for Orchestrator and Research modules.
- Automatic integration between research results and orchestrator memory.
- Root-level `bin/` for compiled executables.

## [1.0.0-alpha.7] - 2026-05-30
### Added
- Utility-based heat scoring for memories with exponential temporal decay.
- Enhanced Research Hustle pipeline with multi-provider support (Tavily, Brave, Google).
- Improved report synthesis logic.

## [1.0.0-alpha.6] - 2026-05-30
### Added
- Real Git rollback implementation (merge abortion, hard reset, clean).
- Failure trap in `sync.sh` for automated error handling during synchronization.

## [1.0.0-alpha.5] - 2026-05-30
### Added
- GitHub Actions CI pipeline in `.github/workflows/ci.yml`.
- End-to-End (E2E) tests for the research hustle pipeline.
- Go workspace (`go.work`) configuration.

## [1.0.0-alpha.4] - 2026-05-29
### Added
- Enhanced `sync.sh` with "Executive Protocol" reconciliation logic.
- Real git status and submodule monitoring in `orchestrator/sync_monitor.go`.
- Comprehensive documentation updates for `VISION.md`, `DEPLOY.md`, and `IDEAS.md`.

## [1.0.0-alpha.3] - 2026-05-29
### Added
- Automated rollback handling skeleton.
- Real-time status reporting to `STATUS.json`.

## [1.0.0-alpha.2] - 2026-05-29
### Added
- Repository synchronization orchestration scripts (`sync.sh`, `build.sh`).
- Sync monitoring and health check logic.

## [1.0.0-alpha.1] - 2026-05-29
### Added
- Repository initialization with knowledge base.
- Basic project structure scaffolding (Orchestrator, Hustle).
- Core documentation files.
