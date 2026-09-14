# Crevisto CLI

AI image generation from your terminal. 100 curated tools, agent-first, CLI-first.

## Install

```sh
curl -fsSL https://crevisto.com/install.sh | sh
```

Or download from [releases](https://github.com/javimosch/crevisto/releases):

| Platform | File |
|---|---|
| Linux amd64 | `crevisto-linux-amd64` |
| macOS Apple Silicon | `crevisto-darwin-arm64` |
| macOS Intel | `crevisto-darwin-amd64` |

## Quick start

```sh
# Get a free trial token (5 credits, no email needed)
crevisto trial

# Check your status
crevisto whoami

# List all 100 AI tools
crevisto tools

# Generate an image
crevisto generate linkedin-avatar --input photo=./me.jpg --output ./result.webp

# Browse the public gallery
crevisto gallery

# Download a gallery image
crevisto gallery-download <id>
```

## How it works

1. **Get a trial token** — `crevisto trial` gives you a bearer token with 5 free credits. No signup, no email.
2. **Generate images** — pick a tool, provide inputs, get a WebP/PNG back.
3. **Run out of credits?** Either:
   - `crevisto claim --email you@example.com` — bind your email for more credits
   - `crevisto set-key sk-or-v1-xxx` — BYOK (Bring Your Own Key), unlimited, zero credits
4. **Go public** — free trial generations are auto-public in the gallery. BYOK users can opt-in with `--public`.

## Commands

| Command | Purpose |
|---|---|
| `crevisto trial` | Get a free trial token (5 credits) |
| `crevisto tools` | List all 100 AI tools (JSON) |
| `crevisto generate <slug>` | Generate an image |
| `crevisto whoami` | Check credit balance and trial status |
| `crevisto claim --email <email>` | Bind email to trial account |
| `crevisto set-key <key>` | Attach OpenRouter BYOK key (unlimited, zero credits) |
| `crevisto gallery` | Browse public gallery (JSON) |
| `crevisto gallery-download <id>` | Download a gallery image |
| `crevisto gallery-toggle <id>` | Toggle your generation public/private |

## Tool examples

| Tool | Credits | Inputs | Type |
|---|---|---|---|
| `linkedin-avatar` | 2 | photo + style | Image-to-image |
| `logo-maker` | 2 | description + style | Text-to-image |
| `pet-portrait` | 2 | photo + style | Image-to-image |
| `photo-restoration` | 3 | photo | Image-to-image |
| `book-cover` | 2 | title + description + genre | Text-to-image |
| `character-art` | 2 | description + style | Text-to-image |

Run `crevisto tools` for the full list of 100 tools.

## BYOK (Bring Your Own Key)

Have an OpenRouter API key? Use it for unlimited generations at zero credit cost:

```sh
crevisto set-key sk-or-v1-xxxxx
crevisto generate logo-maker --input description="A tech startup logo" --output ./logo.webp
```

Your key is encrypted with AES-256-GCM and never exposed.

## Agent integration

The CLI is designed for AI agents. All output is JSON-parseable. No interactive prompts. Deterministic exit codes.

```sh
# Agent workflow
TOKEN=$(crevisto trial --json | jq -r .token)
crevisto auth $TOKEN
crevisto tools --json | jq '.tools[].slug'
crevisto generate linkedin-avatar --input photo=./user.jpg --output ./avatar.webp --json
```

## Build from source

```sh
go build -ldflags "-s -w" -o crevisto .
```

## Popular tools

| Tool | What it does | Link |
|---|---|---|
| `linkedin-avatar` | Professional headshots from selfies | [Try it](https://crevisto.com/tools/linkedin-avatar) |
| `ghibli-style` | Studio Ghibli-style art | [Try it](https://crevisto.com/tools/ghibli-style) |
| `logo-maker` | Custom logos in seconds | [Try it](https://crevisto.com/tools/logo-maker) |
| `photo-restoration` | Fix damaged old photos | [Try it](https://crevisto.com/tools/photo-restoration) |
| `anime-avatar` | Anime-style portraits | [Try it](https://crevisto.com/tools/anime-avatar) |
| `pixel-art` | Retro game sprites | [Try it](https://crevisto.com/tools/pixel-art) |
| `book-cover` | Book cover design | [Try it](https://crevisto.com/tools/book-cover) |
| `tattoo-design` | Custom tattoo designs | [Try it](https://crevisto.com/tools/tattoo-design) |
| `meme-generator` | Create custom memes | [Try it](https://crevisto.com/tools/meme-generator) |
| `fantasy-map` | Fantasy world maps | [Try it](https://crevisto.com/tools/fantasy-map) |

Browse all 100+ tools: https://crevisto.com/tools

## Links

- [Crevisto.com](https://crevisto.com) — full web app with 100 AI image tools
- [Gallery](https://crevisto.com/gallery)
- [Tools](https://crevisto.com/tools)
- [Pricing](https://crevisto.com/pricing)
- [Blog](https://crevisto.com/blog)
- [GitHub Releases](https://github.com/javimosch/crevisto/releases)
