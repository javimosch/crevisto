# Crevisto CLI

AI image generation from your terminal. Agent-first, CLI-first.

## Quick start

```bash
# Build
go build -ldflags "-s -w -X main.version=0.1.0" -o crevisto-cli

# Set your token (invent any 16+ char string — auto-provisions a trial account)
./crevisto-cli auth my-secret-agent-token-12345

# Check your status (5 free images per trial)
./crevisto-cli whoami

# List available image generation tools
./crevisto-cli tools

# Generate an image
./crevisto-cli generate family-hug-generator --input photo=./my-photo.jpg --output result.png

# Attach your own OpenRouter key for unlimited free generations (BYOK)
./crevisto-cli set-key sk-or-v1-xxx
```

## How it works

1. **Invent a bearer token** — any string ≥16 chars. No signup, no email.
2. **Get 5 free images** — the platform auto-provisions a trial account on first `whoami` call.
3. **Bind your email** after 5 images — `crevisto-cli claim --email you@example.com` (verification email sent).
4. **Bring your own key** — attach an OpenRouter API key for unlimited zero-credit generations.

The default image model is **meta/muse-image**.

## CLI specs alignment

This CLI is aligned to all 7 [cli-specs](https://cli-specs.intrane.fr/):

| Spec | Implementation |
|---|---|
| [cli-output-spec](https://cli-specs.intrane.fr/) | JSON stdout, typed errors on stderr, exit codes 80-119, `help-json` |
| [cli-guide-spec](https://cli-specs.intrane.fr/) | `guide` (JSON default, `--human`), embedded in binary |
| [cli-feedback-spec](https://cli-specs.intrane.fr/) | `feedback` command, best-effort, idempotency key |
| [cli-update-spec](https://cli-specs.intrane.fr/) | `update --check`, content-hash, verify-then-swap, `.bak` rollback |
| [cli-telemetry-spec](https://cli-specs.intrane.fr/) | `telemetry` prints payload, `DO_NOT_TRACK=1` respected |
| [cli-daemon-spec](https://cli-specs.intrane.fr/) | `serve` (loopback default), `/_health`, `/_shutdown`, `daemon start\|stop\|status` |
| [cli-trial-spec](https://cli-specs.intrane.fr/) | `auth`, `whoami`, `claim`, `set-key` — bearer-as-tenant auto-provisioning |

## Commands

```
auth <token>            Set your bearer token
whoami                  Check trial status
tools                   List available tools
generate <slug>         Generate an image (--input k=v, --output path)
claim --email <email>   Bind your email
set-key <key>           Attach OpenRouter API key (BYOK, zero credits)
guide [--human]         Print the embedded guide
help-json               Machine-readable command catalog
feedback <message>      Send feedback
update [--check]        Self-update
telemetry               Show telemetry info
serve [--host --port]   Run daemon in foreground
daemon start|stop|status  Manage background daemon
```

## Exit codes

| Range | Meaning |
|---|---|
| 0 | Success |
| 80-89 | Input/validation error |
| 90-99 | Resource/precondition error |
| 100-109 | External/integration error (recoverable) |
| 110-119 | Internal/bug error (recoverable) |

## Build

```bash
go build -ldflags "-s -w -X main.version=0.1.0" -o crevisto-cli
```

Stdlib-only, no external dependencies. Static binary.

## Platform

The Crevisto platform is a separate private project. This CLI talks to it via HTTP.
Default API base: `https://crevisto.intrane.fr` (override with `CREVISTO_API_BASE` env var).

## License

MIT
