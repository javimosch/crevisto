package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// guideJSON is the embedded mental model (cli-guide-spec).
// It MUST NOT fetch a URL at runtime — it travels with every binary.
const guideJSON = `{
  "one_liner": "AI image generation platform with a CLI-first, agent-first interface.",
  "model": "Crevisto turns text prompts into images via OpenRouter models. The CLI is the primary interface: agents invent a bearer token, get 5 free images, and can attach their own OpenRouter key (BYOK) for unlimited zero-credit generations. The platform manages tools (image generation pipelines), credits, and billing — but the CLI bypasses all of that for agent use.",
  "loop": "1. crevisto auth <token>    — set your bearer token (invent any 16+ char string)\n2. crevisto whoami           — auto-provisions a trial account, shows status\n3. crevisto tools            — list available image generation tools\n4. crevisto generate <slug>  — generate an image (uses meta/muse-image by default)\n5. crevisto claim --email X  — bind your email after 5 free images\n6. crevisto set-key <or-key> — attach your OpenRouter key for unlimited BYOK generations",
  "concepts": {
    "bearer_token": "Any string ≥16 chars. The token IS your credential — no signup, no email. Auto-provisions a trial account on first use.",
    "trial_limit": "5 free images per trial account. After that, bind an email (with verification) to unlock more.",
    "byok": "Bring Your Own Key — attach an OpenRouter API key. BYOK generations are free (zero credits) because you pay OpenRouter directly.",
    "tools": "Image generation pipelines. Each tool has a slug, inputs, and a model. The default model is meta/muse-image.",
    "credits": "Platform credits are consumed when using the operator's OpenRouter key. BYOK users consume zero credits."
  },
  "commands": {
    "auth <token>": "Set your bearer token (invent any 16+ char string)",
    "whoami": "Check trial status (auto-provisions on first call)",
    "tools": "List available image generation tools",
    "generate <slug>": "Generate an image (--input key=value, --output path)",
    "claim --email X": "Bind your email to unlock more images",
    "set-key <key>": "Attach your OpenRouter API key (BYOK)",
    "guide": "Print this guide (JSON by default, --human for readable)",
    "help-json": "Machine-readable command catalog",
    "feedback <msg>": "Send feedback to the platform",
    "update": "Self-update the CLI binary",
    "telemetry": "Show telemetry info",
    "serve": "Run the local daemon (foreground)",
    "daemon start|stop|status": "Manage the background daemon"
  },
  "examples": [
    "crevisto auth my-secret-token-12345",
    "crevisto whoami",
    "crevisto tools",
    "crevisto generate family-hug-generator --input photo=base64data --input style=warm",
    "crevisto claim --email me@example.com",
    "crevisto set-key sk-or-v1-xxx"
  ],
  "gotchas": [
    "The bearer token is stored in plaintext at ~/.crevisto/config.json — protect it.",
    "Trial accounts get 5 free images. After that, email binding is required.",
    "BYOK keys are validated live against OpenRouter before being accepted.",
    "The default model is meta/muse-image. Other models (nano_banana, nano_banana_pro) are available per-tool.",
    "Image inputs must be base64-encoded data URIs (data:image/jpeg;base64,...)."
  ],
  "version": "1.0",
  "see_also": ["crevisto help-json", "crevisto --help"]
}`

// cmdGuide implements `crevisto guide` (cli-guide-spec).
func cmdGuide(args []string) {
	human := false
	for _, a := range args {
		if a == "--human" || a == "--format" && len(args) > 1 {
			human = true
		}
	}

	if !human {
		// JSON by default (agent-first)
		fmt.Println(guideJSON)
		os.Exit(0)
	}

	// --human: render as readable text
	var g map[string]interface{}
	json.Unmarshal([]byte(guideJSON), &g)

	fmt.Println("# Crevisto CLI Guide")
	fmt.Println()
	fmt.Println(g["one_liner"])
	fmt.Println()
	fmt.Println("## Model")
	fmt.Println(g["model"])
	fmt.Println()
	fmt.Println("## Loop")
	fmt.Println(g["loop"])
	fmt.Println()
	fmt.Println("## Commands")
	if cmds, ok := g["commands"].(map[string]interface{}); ok {
		for cmd, desc := range cmds {
			fmt.Printf("  %-30s %s\n", cmd, desc)
		}
	}
	fmt.Println()
	fmt.Println("## Examples")
	if examples, ok := g["examples"].([]interface{}); ok {
		for _, ex := range examples {
			fmt.Printf("  $ %s\n", ex)
		}
	}
	fmt.Println()
	fmt.Println("## Gotchas")
	if gotchas, ok := g["gotchas"].([]interface{}); ok {
		for _, gt := range gotchas {
			fmt.Printf("  ⚠ %s\n", gt)
		}
	}
	os.Exit(0)
}
