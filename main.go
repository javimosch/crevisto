package main

import (
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		cmdHelp()
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "auth":
		cmdAuth(args)
	case "whoami":
		cmdWhoami(args)
	case "tools":
		cmdTools(args)
	case "generate":
		cmdGenerate(args)
	case "claim":
		cmdClaim(args)
	case "set-key":
		cmdSetKey(args)
	case "gallery":
		cmdGallery(args)
	case "gallery-download":
		cmdGalleryDownload(args)
	case "gallery-toggle":
		cmdGalleryToggle(args)
	case "guide":
		cmdGuide(args)
	case "help-json":
		cmdHelpJSON(args)
	case "--help-json":
		cmdHelpJSON(args)
	case "feedback":
		cmdFeedback(args)
	case "update":
		cmdUpdate(args)
	case "telemetry":
		cmdTelemetry(args)
	case "serve":
		cmdServe(args)
	case "daemon":
		cmdDaemon(args)
	case "--help", "-h", "help":
		cmdHelp()
	case "--version", "version":
		outOK(map[string]string{"version": version})
	default:
		// Suggest a similar command
		suggestion := suggestCommand(cmd)
		msg := "unknown command: " + cmd
		if suggestion != "" {
			fail(ExitInput, "unknown_command", msg, "did you mean: "+suggestion+"?")
		}
		fail(ExitInput, "unknown_command", msg, "crevisto help-json  — list all commands")
	}
}

// suggestCommand finds the closest matching command.
func suggestCommand(input string) string {
	commands := []string{"auth", "whoami", "tools", "generate", "claim", "set-key",
		"gallery", "gallery-download", "gallery-toggle",
		"guide", "help-json", "feedback", "update", "telemetry", "serve", "daemon"}

	for _, c := range commands {
		if strings.HasPrefix(c, input) || strings.HasPrefix(input, c) {
			return c
		}
	}
	return ""
}
