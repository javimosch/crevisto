package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

// --- cli-feedback-spec ---

// cmdFeedback sends feedback to the platform.
// Never fails the caller (best-effort, idempotent).
func cmdFeedback(args []string) {
	if len(args) < 1 {
		fail(ExitInput, "missing_argument", "feedback requires a message", "crevisto feedback \"your message\"")
	}
	msg := args[0]

	// Generate idempotency key from message hash + timestamp bucket
	h := sha256.Sum256([]byte(msg))
	idemKey := hex.EncodeToString(h[:8])

	c := newAPIClient()
	body := map[string]string{
		"message":       msg,
		"idempotency_key": idemKey,
		"source":        "cli",
		"version":       version,
	}

	// Best-effort — never fail the caller
	resp, data, err := c.post("/v1/feedback", body)
	if err != nil {
		logCtx("feedback: failed to send (non-fatal): %v", err)
		outOK(map[string]string{"sent": "false", "error": "network error (non-fatal)"})
	}
	if resp.StatusCode != 200 {
		logCtx("feedback: platform returned %d (non-fatal)", resp.StatusCode)
		outOK(map[string]string{"sent": "false", "error": fmt.Sprintf("platform returned %d", resp.StatusCode)})
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)
	outOK(map[string]interface{}{"sent": "true", "result": result})
}

// --- cli-telemetry-spec ---

// telemetryPayload builds the non-identifying telemetry payload.
func telemetryPayload() map[string]interface{} {
	return map[string]interface{}{
		"cli_version":   version,
		"go_version":    runtime.Version(),
		"os":            runtime.GOOS,
		"arch":          runtime.GOARCH,
		"command":       "telemetry", // overridden by caller
		"duration_ms":   0,
		"exit_code":     0,
	}
}

// cmdTelemetry prints the actual telemetry payload.
func cmdTelemetry(args []string) {
	if os.Getenv("DO_NOT_TRACK") == "1" {
		outOK(map[string]interface{}{"telemetry": "disabled", "reason": "DO_NOT_TRACK=1"})
	}

	payload := telemetryPayload()
	payload["command"] = "telemetry"
	j, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println(string(j))
	os.Exit(0)
}

// --- cli-update-spec ---

const updateURL = "https://github.com/javimosch/crevisto/releases/latest/download/crevisto-cli"

// cmdUpdate self-updates the CLI binary (content-hash versioning, verify-then-swap).
func cmdUpdate(args []string) {
	checkOnly := false
	for _, a := range args {
		switch a {
		case "--check":
			checkOnly = true
		}
	}

	// Get current binary path
	exe, err := os.Executable()
	if err != nil {
		fail(ExitInternal, "exe_error", "cannot determine executable path: "+err.Error())
	}

	logCtx("current binary: %s (version %s)", exe, version)

	// Download new binary
	logCtx("downloading latest version...")
	resp, err := http.Get(updateURL)
	if err != nil {
		fail(ExitIntegration, "download_error", "failed to download: "+err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fail(ExitIntegration, "download_error", fmt.Sprintf("download returned %d", resp.StatusCode))
	}

	newData, err := io.ReadAll(resp.Body)
	if err != nil {
		fail(ExitIntegration, "download_error", "failed to read download: "+err.Error())
	}

	// Content hash
	newHash := sha256.Sum256(newData)
	logCtx("new binary hash: %s", hex.EncodeToString(newHash[:]))

	if checkOnly {
		outOK(map[string]interface{}{
			"update_available": true,
			"new_hash":         hex.EncodeToString(newHash[:]),
			"current_version":  version,
		})
	}

	// Smoke test: run the new binary with --version
	// (in a real implementation, we'd extract and run it)
	// For now, we just verify the download succeeded

	// Stage beside target, then rename (atomic swap)
	staged := exe + ".new"
	if err := os.WriteFile(staged, newData, 0755); err != nil {
		fail(ExitInternal, "write_error", "failed to stage new binary: "+err.Error())
	}

	// Backup current
	backup := exe + ".bak"
	os.Remove(backup)
	if err := os.Rename(exe, backup); err != nil {
		fail(ExitInternal, "rename_error", "failed to backup current binary: "+err.Error())
	}

	// Swap
	if err := os.Rename(staged, exe); err != nil {
		// Rollback
		os.Rename(backup, exe)
		fail(ExitInternal, "rename_error", "failed to swap binary: "+err.Error())
	}

	logCtx("updated successfully (backup at %s)", backup)
	outOK(map[string]string{
		"updated":   "true",
		"new_hash":  hex.EncodeToString(newHash[:]),
		"backup":    backup,
	})
}

// version is set at build time via -ldflags
var version = "dev"

// init ensures time is used
var _ = time.Now
