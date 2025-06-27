# Zengate

[![releases ci/cd](https://github.com/myferr/zengate/actions/workflows/release.yml/badge.svg)](https://github.com/myferr/zengate/actions/workflows/release.yml)
![](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![](https://img.shields.io/badge/Uses-Docker-blue?style=flat-square&logo=docker)
![](https://img.shields.io/badge/Tauri-v2-8A2BE2?style=flat-square&logo=tauri)
[![](https://img.shields.io/badge/Contributing-CONTRIBUTING.md-blue?style=flat-square&logo=github)](CONTRIBUTING.md)

Self-host your own end-to-end encrypted password manager.

---

Zengate is an open-source, self-hostable, end-to-end encrypted password manager that uses Docker.

# Requirements.
Before installing the desktop app, CLI, or self-hosting the backend, make sure you have the following technologies installed on your machine:
- [git](https://git-scm.com/)
- [Docker](https://docker.com/)

# Installation
To install Zengate, you can either get the CLI and follow [the CLI installation guide](https://github.com/myferr/zengate/blob/main/cli/INSTALL.md) or get the desktop app.

## Installing the CLI
Go check out [the CLI installation guide](https://github.com/myferr/zengate/blob/main/cli/INSTALL.md)

## Installing the desktop app
The desktop app is built with Tauri, Svelte, TypeScript, and Docker. There are two ways of getting the desktop app on your machine

### 1. GitHub Releases
Go to [the releases page](https://github.com/myferr/zengate/releases/latest) and get the latest binary for your machine. Install that and pin it to your dock or task-bar.

### 2. Dockerfile
Inside the [`desktop/`](https://github.com/myferr/zengate/tree/main/desktop) directory there is a [Dockerfile](https://github.com/myferr/zengate/blob/main/desktop/Dockerfile) that builds the app inside of a Docker container. There is a thorough guide on using Docker to install the desktop app, [see here.](https://github.com/myferr/zengate/blob/main/desktop/DOCKER.md)
