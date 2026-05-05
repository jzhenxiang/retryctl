# retryctl

> CLI utility for wrapping flaky commands with configurable retry logic, backoff strategies, and structured logging.

---

## Installation

```bash
go install github.com/yourusername/retryctl@latest
```

Or download a pre-built binary from the [Releases](https://github.com/yourusername/retryctl/releases) page.

---

## Usage

```bash
retryctl [flags] -- <command> [args...]
```

### Examples

```bash
# Retry a flaky curl command up to 5 times with exponential backoff
retryctl --attempts 5 --backoff exponential -- curl https://example.com/api

# Retry with a fixed 2-second delay between attempts
retryctl --attempts 3 --backoff fixed --delay 2s -- ./deploy.sh

# Retry with jitter and structured JSON logging
retryctl --attempts 4 --backoff exponential --jitter --log-format json -- make test
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--attempts` | `3` | Maximum number of retry attempts |
| `--backoff` | `fixed` | Backoff strategy: `fixed`, `exponential`, `linear` |
| `--delay` | `1s` | Initial delay between retries |
| `--jitter` | `false` | Add random jitter to backoff delay |
| `--timeout` | none | Per-attempt timeout |
| `--log-format` | `text` | Log output format: `text`, `json` |

---

## Backoff Strategies

- **fixed** — constant delay between each attempt
- **linear** — delay increases linearly with each attempt
- **exponential** — delay doubles with each attempt (with optional jitter)

---

## License

[MIT](LICENSE)