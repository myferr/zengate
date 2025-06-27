# Contributing to Zengate

Thank you for taking the time to contribute to Zengate!

This project is made up of three main components:
- **CLI** (`cli/`): Command-line password manager written in Go
- **Backend** (`backend/`): Encrypted password API using Bun + Elysia + Docker
- **Desktop** (`desktop/`): Tauri desktop application built with Bun

---

## Project Setup

### Prerequisites

- [Go](https://go.dev/dl/)
- [Bun](https://bun.sh/docs/installation)
- [Docker](https://www.docker.com/)
- [Rust](https://www.rust-lang.org/tools/install)
- [Node.js](https://nodejs.org) (for Tauri prerequisites)
- [Tauri](https://tauri.app/v1/guides/getting-started/prerequisites)

---

## CLI (`cli/`)

The CLI is written in Go and uses local config/encryption to manage credentials.

### Build

```bash
make cli
```

### Run

```bash
./zengate
```

### Common Dev Tasks

* `make cli`: Build the CLI binary
* `make clean`: Remove old builds
* `go run .` inside `cli/` for fast testing

---

## Backend (`backend/`)

The backend is written in Bun (Elysia) and exposed via Cloudflare tunnels.

### Run Locally

```bash
cd backend
bun install
bun run dev
```

### With Tunnel

```bash
docker compose up --build
```

* Make sure you have a valid `.env` set up.
* API runs at: `http://localhost:3002`
* Cloudflare Tunnel is set up (if you run `docker compose up --build`)

---

## Desktop App (`desktop/`)

The desktop app is powered by Tauri, Bun, Svelte, and TypeScript.

### Setup

```bash
cd desktop
bun install
```

### Build

```bash
bunx tauri build
```

### Development

```bash
bunx tauri dev
```

---

## Testing Your Changes

* Follow the command line prompts and test all encryption flows
* For the backend, use something like `curl` or Postman to test endpoints
* On desktop, make sure AppImage/.dmg/.exe builds correctly

---

## Contributing Guidelines

* Keep PRs focused
* Use consistent formatting and naming
* Add comments for complex code
* Write meaningful commit messages (I suggest using [conventional-commit](https://github.com/myferr/conventional-commit))

---

## Submitting a PR

1. Fork this repository
2. Create a branch: `git checkout -b my-feature`
3. Commit your changes: `git add . && git commit -m "feat: added new feature"`
4. Push to your fork: `git push origin my-feature`
5. Open a pull request!

---

If you have questions, feel free to [open an issue](https://github.com/myferr/zengate/issues/new).
