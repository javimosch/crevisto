package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// cmdHelpJSON implements `crevisto help-json` (cli-output-spec).
func cmdHelpJSON(args []string) {
	catalog := map[string]interface{}{
		"version": "1.0",
		"name":    appName,
		"commands": []map[string]string{
			{"cmd": "auth", "args": "<token>", "desc": "set your bearer token (invent any 16+ char string)"},
			{"cmd": "whoami", "args": "", "desc": "check trial status (auto-provisions on first call)"},
			{"cmd": "tools", "args": "", "desc": "list available image generation tools"},
			{"cmd": "generate", "args": "<slug> [--input k=v]... [--output path]", "desc": "generate an image"},
			{"cmd": "claim", "args": "--email <email> [--orkey <openrouter-key>]", "desc": "bind email + optional BYOK key"},
			{"cmd": "set-key", "args": "<openrouter-key>", "desc": "attach your OpenRouter API key (BYOK)"},
			{"cmd": "guide", "args": "[--human]", "desc": "print the embedded guide (JSON default, --human for text)"},
			{"cmd": "help-json", "args": "", "desc": "machine-readable command catalog"},
			{"cmd": "feedback", "args": "<message>", "desc": "send feedback to the platform"},
			{"cmd": "update", "args": "[--check] [--force]", "desc": "self-update the CLI binary"},
			{"cmd": "telemetry", "args": "", "desc": "show telemetry info and next payload"},
			{"cmd": "serve", "args": "[--host H] [--port P]", "desc": "run the daemon in foreground (loopback by default)"},
			{"cmd": "daemon", "args": "start|stop|status [--port P]", "desc": "manage the background daemon"},
		},
		"exit_codes": map[string]string{
			"0":          "ok",
			"80-89":      "input",
			"90-99":      "resource",
			"100-109":    "integration",
			"110-119":    "internal",
		},
		"env":      []string{"CREVISTO_API_BASE", "CREVISTO_TOKEN", "DO_NOT_TRACK"},
		"see_also": []string{"crevisto guide", "crevisto --help"},
	}

	j, _ := json.MarshalIndent(catalog, "", "  ")
	fmt.Println(string(j))
	os.Exit(0)
}

// cmdHelp prints human-readable help.
func cmdHelp() {
	fmt.Fprintf(os.Stderr, `crevisto — AI image generation CLI (agent-first, cli-specs aligned)

Usage:
  crevisto <command> [options]

Commands:
  auth <token>            Set your bearer token (invent any 16+ char string)
  whoami                  Check trial status (auto-provisions on first call)
  tools                   List available image generation tools
  generate <slug>         Generate an image (--input k=v, --output path)
  claim --email <email>   Bind your email to unlock more images
  set-key <key>           Attach your OpenRouter API key (BYOK, zero credits)
  guide [--human]         Print the embedded guide
  help-json               Machine-readable command catalog
  feedback <message>      Send feedback
  update [--check]        Self-update the CLI binary
  telemetry               Show telemetry info
  serve [--host --port]   Run daemon in foreground
  daemon start|stop|status  Manage background daemon

Exit codes: 0=ok, 80-89=input, 90-99=resource, 100-109=integration, 110-119=internal

See: crevisto guide --human  for the full operator manual
     crevisto help-json      for the machine-readable catalog
`)
	os.Exit(0)
}
