# Architecture Overview

This document describes the architecture for the Shell History Sync Client, a cross-platform Go-based tool for synchronizing shell history files using filebrowser as a backend.

---

## High-Level Design

# TODO

- add testing ( unit / in)

### 1. Platform Agnostic Daemon
- The sync client runs as a background process (daemon/service) on Linux, macOS, and Windows.
- Startup integration is provided via systemd (Linux), LaunchAgent/Startup App (macOS), and Task Scheduler/Startup (Windows).

### 2. Configuration
- Uses a YAML configuration file.
- Configurable options include:
  - List of history files to sync (with shell type and path)
  - Backend connection info (filebrowser URL, credentials)
  - Secret filtering rules
  - Sync interval (default: 15 seconds)
  - Logging options

### 3. Main Sync Loop
- On startup:
  - Loads configuration and authenticates with filebrowser.
  - For each history file:
    - Downloads the latest remote version from filebrowser.
    - Merges remote and local histories (preserving order as much as possible, deduplicating lines).
    - Writes the merged result locally, filtering secrets as needed.
    - Uploads the merged file if there are changes.
- The process repeats every configurable interval (default: 15 seconds).

### 4. Merging & Conflict Resolution
- On first run, always fetches remote state and merges with local file.
- Merging is line-based, but attempts to preserve order when possible.
- Deduplicates lines to avoid history bloat.
- Uses file locks and atomic operations to ensure concurrency safety.

### 5. Secret Filtering
- Filters out or redacts sensitive entries before upload.
- Built-in and user-configurable regex patterns.

### 6. Logging
- Logs to stdout and/or file as configured.
- Logs sync events, errors, and filtering actions.

### 7. Extensibility
- Easy to add new shells or history file formats.
- Pluggable backend support (future-proofing for other storage backends).

---

## Backend: filebrowser
- filebrowser is used as the central backend for storing and syncing history files.
- The client interacts with filebrowser via its REST API, authenticating with username/password or JWT.
- Each user/device has a dedicated directory or namespace for their history files.

---

## Platform Startup Integration
- **Linux:** systemd service unit file.
- **macOS:** LaunchAgent or login item.
- **Windows:** Task Scheduler or Startup folder shortcut.
- Documentation/scripts will be provided for each platform.

---

## Security & Best Practices
- Always use HTTPS for backend communication.
- Use strong credentials and/or JWT tokens for authentication.
- Regularly update both the sync client and filebrowser backend.

---

## Future Enhancements
- Support for additional backends (S3, Git, etc.).
- GUI for configuration and monitoring.
- More advanced merge strategies.

---

For more details on setup and usage, see the README.md.
