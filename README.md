# ClusterMind

**ClusterMind** is a terminal-native AI operator for Kubernetes designed specifically for SREs and Platform Engineers. It accelerates cluster troubleshooting by combining deep, real-time Kubernetes context gathering with generative AI to offer instant root-cause analysis and safe, one-click remediation.

## Key Features

- 🔍 **Intent-Based Search:** Use natural language queries (like "pods restarting" or "unhealthy deploy") inside the terminal to instantly filter and surface problematic resources.
- 🤖 **Automated Root Cause Analysis:** Select any failing resource and ClusterMind will build a rich "Evidence Package" (aggregating pod statuses, events, exit codes, and logs) to feed into a local LLM for a deterministic Root Cause Analysis.
- 🛠️ **Remediation Engine:** Instead of just explaining the problem, ClusterMind generates actionable fixes (shell commands or Kubernetes API patches) complete with risk-level assessments.
- 🛡️ **Safety First:** Actions are strictly read-only by default. Remediation proposals require explicit human `Y/N` confirmation through an interactive terminal prompt before any commands are executed against your cluster.
- 💻 **Terminal-Native TUI:** Built entirely in the terminal using `Bubble Tea` and `Lip Gloss`, providing a snappy, low-friction experience directly in your standard CLI workflow.

## Architecture

The project is decoupled into cleanly separated internal modules:

- `internal/ui`: The main Bubble Tea event loop, layout handling, and pane rendering (`ResourcePane`, `ExplainPane`, `QueryBar`).
- `internal/context`: The `EvidenceBuilder` that aggregates real-time metrics, logs, and events from Kubernetes.
- `internal/diagnostics`: Built-in deterministic rules that catch obvious errors (e.g., `CrashLoopBackOff`, `OOMKilled`) prior to AI inference.
- `internal/ai`: The Orchestrator interface mapping complex tasks (RCA, Remediation) to the underlying LLM provider (currently defaulting to local `Ollama`).
- `internal/intent`: Parsers to map natural language intents to strict Kubernetes API field selectors.
- `internal/remediation`: An execution engine capable of safely running shell commands or applying API patches.

*A detailed architectural Mermaid diagram is available in [docs/architecture.md](docs/architecture.md).*

## Prerequisites

- Go 1.23 or higher
- A valid Kubernetes cluster config (`~/.kube/config`)
- [Ollama](https://ollama.com/) running locally with the `llama3.2` model pulled (`ollama run llama3.2`).

*Note: If no Kubernetes cluster is detected, ClusterMind will automatically fall back to a mock, in-memory client to let you preview the UI and AI integration features.*

## Installation & Quick Start

1. **Clone the repository:**
   ```bash
   git clone https://github.com/vignesh245/ClusterMind.git
   cd ClusterMind
   ```

2. **Build the binary:**
   ```bash
   make build
   ```
   This will output the primary binary to `bin/clustermind`.

3. **Run the CLI:**
   ```bash
   ./bin/clustermind
   ```

## Keybindings

- `:` - Focus the Query Bar to enter intent-based searches.
- `Enter` - Execute search intent in the Query Bar.
- `y` / `n` - Approve or reject an AI-proposed remediation action.
- `Esc` - Dismiss active prompts or overlays.
- `q` / `ctrl+c` - Quit the application.

## Development

- **Linting:** Requires `golangci-lint` to be installed locally. Run `make lint`.
- **Pre-commit Hooks:** Run `make install-hooks` to auto-format and lint code on `git commit`.

## License

MIT License