# INSTALL.md

A document on how to install the `zengate` CLI.

---

## Windows

To install on Windows operating systems, you need to remotely execute `./scripts/install.ps1` using the following one-line command:

```sh
powershell -c "irm https://raw.githubusercontent.com/myferr/zengate/refs/heads/main/cli/scripts/install.ps1 | iex"
```

That command runs the script and adds the `zengate` CLI to your PATH.

## Linux / macOS

To install on Unix-based operating systems (e.g. Linux or macOS), you need to remotely execute `./scripts/install.sh` using the following one-line command:

```sh
curl -fsSL https://raw.githubusercontent.com/myferr/zengate/refs/heads/main/cli/scripts/install.sh | bash
```

That command runs the script and adds the `zengate` CLI to your PATH.
