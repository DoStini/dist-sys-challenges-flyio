## Fly.io Distributed Systems Challenges

This repository contains my implementations for the Fly.io/Maelstrom distributed systems challenges (echo, broadcast, unique-ids, etc.) and simple scripts to run them against Maelstrom test workloads.

### Prerequisites
- Go 1.25 installed and on your PATH
- Java 11+ (for running Maelstrom)
- Bash (scripts assume macOS/Linux)

### Install Maelstrom (required)
Maelstrom must be installed inside the `maelstrom` folder at the repository root.

Steps:
1. Download the Maelstrom release zip from the official repo.
2. Extract it so that the Maelstrom launcher and `maelstrom.jar` live under the `maelstrom` directory here.

You should end up with something like:
- `./maelstrom/maelstrom` (executable launcher)
- `./maelstrom/lib/maelstrom.jar` (or `./maelstrom/maelstrom.jar` depending on release layout)

If you already see `./maelstrom/maelstrom` in this repo, you’re good.

### Project layout (relevant bits)
- `cmd/echo` – Go implementation for echo
- `cmd/broadcast` – Go implementation for broadcast
- `cmd/unique-id` – Go implementation for unique-ids
- `scenarios/*` – Helper scripts to build and run each workload against Maelstrom
- `maelstrom/` – Maelstrom binary and jar live here (required)

### How to run
All scenarios follow the same pattern: the script will build the scenario binary and then invoke Maelstrom with the correct workload and flags.

From the repo root:

#### Echo
```bash
./scenarios/echo/run.sh
```

#### Broadcast
```bash
./scenarios/broadcast/run-single.sh
```

#### Unique IDs
```bash
./scenarios/unique-id/run.sh
```

These scripts assume the Maelstrom launcher is available at `./maelstrom/maelstrom`.