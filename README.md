# vaultpipe

A CLI tool to sync secrets from HashiCorp Vault into local `.env` files safely.

---

## Installation

```bash
go install github.com/yourname/vaultpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/vaultpipe.git
cd vaultpipe
go build -o vaultpipe .
```

---

## Usage

Authenticate with Vault and sync secrets to a `.env` file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-vault-token"

vaultpipe sync --path secret/myapp --output .env
```

This will pull all key-value pairs from the specified Vault path and write them to `.env`:

```
DB_HOST=db.example.com
DB_PASSWORD=supersecret
API_KEY=abc123
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | *(required)* |
| `--output` | Output file path | `.env` |
| `--overwrite` | Overwrite existing file | `false` |
| `--mask` | Mask secret values in logs | `true` |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with KV secrets engine enabled

---

## License

[MIT](LICENSE) © 2024 yourname