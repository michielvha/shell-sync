# Shell History Sync Client

A cross-platform CLI tool to synchronize your shell history files (bash, zsh, PowerShell, etc.) across devices, with secret filtering and flexible configuration.

This project uses [filebrowser/filebrowser](https://github.com/filebrowser/filebrowser) as its backend for file storage and sync. This approach ensures easy setup, minimal maintenance, and leverages a mature, open-source solution for secure file management.

---

## Features

- **Platform-agnostic:** Syncs history files from different shells and OSes.
- **Secret filtering:** Automatically or optionally removes sensitive entries before syncing.
- **Configurable:** Uses a YAML configuration file for easy setup.
- **No server maintenance:** Uses filebrowser as a backend, so you don't need to build or maintain your own server.

---

## How It Works

1. You deploy or use an existing [filebrowser](https://github.com/filebrowser/filebrowser) instance as your central history storage backend.
2. The sync client authenticates to filebrowser (with username/password or JWT token) and uploads/downloads your shell history files.
3. The client can filter secrets and supports custom sync rules per user/config.

---

## Setting Up filebrowser as the Backend

### 1. Deploy filebrowser

You can run filebrowser in several ways:

#### Docker (recommended)

```bash
docker run -d \
  -v /path/to/data:/srv \
  -v /path/to/config:/config \
  -p 8080:80 \
  filebrowser/filebrowser
```
- `/path/to/data` is where user files will be stored.
- `/path/to/config` is where filebrowser's config and database will be stored.
- The web UI and API will be available at `http://localhost:8080` by default.

#### Standalone Binary

Download the latest release from [filebrowser releases](https://github.com/filebrowser/filebrowser/releases) and run:

```bash
./filebrowser -r /path/to/data -p 8080
```

### 2. Initial Setup

- Access the web UI at `http://localhost:8080`.
- Login with the default credentials (`admin`/`admin`), then change the password.
- Create a user for each device or person who will sync history.
- Set user home directories as needed for isolation.

### 3. API Usage

- The sync client will authenticate using username/password or JWT token.
- Files can be uploaded, listed, and downloaded via the REST API.
- Refer to [filebrowser API docs](https://filebrowser.org/api) for details.

---

## Example Workflow

1. Configure your sync client with your filebrowser URL, credentials, and the list of history files to sync.
2. On each sync, the client will upload/download files via the filebrowser API.
3. Secret filtering is applied before upload.

---

## Security Notes

- Use HTTPS for your filebrowser instance in production.
- Use strong, unique passwords and/or JWT tokens.
- Regularly update filebrowser to the latest version.

---

## Resources

- [filebrowser Documentation](https://filebrowser.org/)
- [filebrowser API Reference](https://filebrowser.org/api)
- [filebrowser GitHub](https://github.com/filebrowser/filebrowser)

---

## Next Steps

- Implement the sync client in Go.
- Add support for YAML config and secret filtering.
- Integrate with filebrowser's API for file operations.

For questions or contributions, open an issue or pull request!
