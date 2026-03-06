# silo

Backup tool with chunking, deduplication, compression, and encryption.

## Installation

Install the latest release (Linux, macOS, Windows).

### Linux / macOS

```bash
curl -sSL https://raw.githubusercontent.com/zhravan/silo/main/scripts/install.sh | sh
```

The script installs the binary to `/usr/local/bin/silo` (uses `sudo` if needed). To install to a different prefix:

```bash
SILO_PREFIX=$HOME/.local curl -sSL https://raw.githubusercontent.com/zhravan/silo/main/scripts/install.sh | sh
```

To install a specific version:

```bash
SILO_VERSION=0.0.1 curl -sSL https://raw.githubusercontent.com/zhravan/silo/main/scripts/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/zhravan/silo/main/scripts/install.ps1 | iex
```

This installs to `%LOCALAPPDATA%\Programs\silo` and adds it to your user PATH. For a specific version:

```powershell
$env:SILO_VERSION="0.0.1"; irm https://raw.githubusercontent.com/zhravan/silo/main/scripts/install.ps1 | iex
```

## Usage

```bash
silo run      # run backup
silo restore  # restore from backup
silo status   # show status
```

Use a config file (default: `backup.yaml`) and optionally pass `-config` to point to another path.
