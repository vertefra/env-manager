![Static Badge](https://img.shields.io/badge/version-0.1.0-blue?style=flat&label=version&labelColor=darkblue&color=black)

[![test](https://github.com/vertefra/env-manager/actions/workflows/test.yaml/badge.svg)](https://github.com/vertefra/env-manager/actions/workflows/test.yaml)

# Env Manager

Securely store and manage encrypted environment configurations.

## Setup

**Install** (Linux only)
```bash
curl -sSL https://raw.githubusercontent.com/thinktwiceco/env-manager/master/install.sh | bash
```

**Create encryption key**
```bash
openssl rand -hex 16 > .secret
```
> Secret can also be set via `ENV_MANAGER_SECRET` environment variable

## Commands

### `add` - Import file with headers
For files that already have env-manager headers:
```bash
env-manager add -f .env.production
```

**File format:**
```bash
#- identifier: production
#- restore-as: .env
DATABASE_URL=postgres://localhost/db
API_KEY=secret123
```

### `create` - Import plain file
For plain environment files without headers:
```bash
env-manager create -f secrets.txt -i production -r .env
```

### `get` - Restore configuration
```bash
env-manager get -i production
```
Decrypts and restores the configuration file.

### `list` - Show all configurations
```bash
env-manager list
```

### `remove` - Delete configuration
```bash
env-manager remove -i production
```

## How It Works

1. Files are encrypted using AES and stored in `.env-manager/`
2. A `manifest.json` tracks all configurations
3. Identifiers map to encrypted files for easy retrieval
4. On restore, files are decrypted and written with their original name

## Security

- ⚠️ **Keep `.secret` secure** - anyone with this key can decrypt your files
- ✅ Add `.secret` and `.env-manager/` to `.gitignore`
- ✅ Valid keys: 16, 24, or 32 bytes (32, 48, or 64 hex characters)
