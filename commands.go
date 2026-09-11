package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// mimeByExt returns the MIME type for common image file extensions.
func mimeByExt(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".webp":
		return "image/webp"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	default:
		return "image/jpeg" // safe fallback
	}
}

// cmdAuth sets the bearer token.
func cmdAuth(args []string) {
	if len(args) < 1 {
		fail(ExitInput, "missing_argument", "auth requires a token", "crevisto auth <token>")
	}
	token := args[0]
	if len(token) < 16 {
		fail(ExitInput, "invalid_argument", "token must be at least 16 characters")
	}
	if err := setToken(token); err != nil {
		fail(ExitInternal, "config_error", err.Error())
	}
	outOK(map[string]string{"token_set": "true", "hint": "run 'crevisto whoami' to provision your trial account"})
}

// cmdWhoami checks trial status (auto-provisions on first call).
func cmdWhoami(args []string) {
	token := requireToken()
	c := newAPIClient()
	c.token = token

	resp, data, err := c.get("/v1/whoami")
	if err != nil {
		fail(ExitIntegration, "api_error", "failed to reach platform: "+err.Error(), "check CREVISTO_API_BASE or network")
	}
	if resp.StatusCode == 400 {
		var e map[string]interface{}
		json.Unmarshal(data, &e)
		fail(ExitInput, "invalid_token", fmt.Sprintf("%v", e["error"]))
	}
	if resp.StatusCode != 200 {
		fail(ExitIntegration, "api_error", fmt.Sprintf("whoami returned %d", resp.StatusCode))
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)
	outData(result)
}

// cmdTools lists available tools.
func cmdTools(args []string) {
	c := newAPIClient()
	resp, data, err := c.get("/v1/tools")
	if err != nil {
		fail(ExitIntegration, "api_error", "failed to reach platform: "+err.Error())
	}
	if resp.StatusCode != 200 {
		fail(ExitIntegration, "api_error", fmt.Sprintf("tools returned %d", resp.StatusCode))
	}
	outData(jsonRaw(data))
}

// cmdGenerate generates an image.
func cmdGenerate(args []string) {
	if len(args) < 1 {
		fail(ExitInput, "missing_argument", "generate requires a tool slug", "crevisto tools  — list available slugs")
	}
	slug := args[0]
	token := requireToken()

	// Parse --input k=v, --output path, --public
	inputs := map[string]string{}
	outputPath := ""
	isPublic := false
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--input", "-i":
			if i+1 < len(args) {
				kv := args[i+1]
				parts := strings.SplitN(kv, "=", 2)
				if len(parts) == 2 {
					// If the value is a file path, read and base64-encode it
					val := parts[1]
					if _, err := os.Stat(val); err == nil {
						data, err := os.ReadFile(val)
						if err != nil {
							fail(ExitInput, "file_error", "cannot read file: "+val)
						}
						mime := mimeByExt(val)
						val = "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data)
					}
					inputs[parts[0]] = val
				}
				i++
			}
		case "--output", "-o":
			if i+1 < len(args) {
				outputPath = args[i+1]
				i++
			}
		case "--public":
			isPublic = true
		}
	}

	c := newAPIClient()
	c.token = token

	body := map[string]interface{}{"inputs": inputs, "isPublic": isPublic}
	resp, data, err := c.post("/v1/tools/"+slug+"/execute", body)
	if err != nil {
		fail(ExitIntegration, "api_error", "failed to reach platform: "+err.Error())
	}
	if resp.StatusCode == 401 {
		fail(ExitResource, "not_authenticated", "invalid or missing token", "crevisto auth <token>")
	}
	if resp.StatusCode == 402 {
		var e map[string]interface{}
		json.Unmarshal(data, &e)
		fail(ExitResource, "trial_limit", fmt.Sprintf("%v", e["error"]), "crevisto claim --email <email>")
	}
	if resp.StatusCode == 404 {
		fail(ExitResource, "not_found", "tool not found: "+slug, "crevisto tools  — list available slugs")
	}
	if resp.StatusCode != 200 {
		fail(ExitIntegration, "api_error", fmt.Sprintf("generate returned %d: %s", resp.StatusCode, string(data)))
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	// If --output is set and there's an image, save it
	if outputPath != "" {
		if r, ok := result["result"].(map[string]interface{}); ok {
			if img, ok := r["image"].(string); ok && img != "" {
				// Strip data URI prefix
				b64 := img
				if idx := strings.Index(img, "base64,"); idx >= 0 {
					b64 = img[idx+7:]
				}
				data, err := base64.StdEncoding.DecodeString(b64)
				if err != nil {
					fail(ExitInternal, "decode_error", "failed to decode image: "+err.Error())
				}
				if err := os.WriteFile(outputPath, data, 0644); err != nil {
					fail(ExitInternal, "write_error", "failed to write file: "+err.Error())
				}
				logCtx("image saved to %s (%d bytes)", outputPath, len(data))
			}
		}
	}

	outData(result)
}

// cmdClaim binds an email to the trial account.
func cmdClaim(args []string) {
	email := ""
	orkey := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--email", "-e":
			if i+1 < len(args) {
				email = args[i+1]
				i++
			}
		case "--orkey", "--key":
			if i+1 < len(args) {
				orkey = args[i+1]
				i++
			}
		}
	}
	if email == "" {
		fail(ExitInput, "missing_argument", "claim requires --email", "crevisto claim --email me@example.com")
	}

	token := requireToken()
	c := newAPIClient()
	c.token = token

	body := map[string]string{"token": token, "email": email}
	if orkey != "" {
		body["orkey"] = orkey
	}

	resp, data, err := c.post("/app/claim", body)
	if err != nil {
		fail(ExitIntegration, "api_error", "failed to reach platform: "+err.Error())
	}
	if resp.StatusCode != 200 {
		fail(ExitIntegration, "api_error", fmt.Sprintf("claim returned %d: %s", resp.StatusCode, string(data)))
	}
	outData(jsonRaw(data))
}

// cmdSetKey attaches an OpenRouter API key (BYOK).
func cmdSetKey(args []string) {
	if len(args) < 1 {
		fail(ExitInput, "missing_argument", "set-key requires an OpenRouter API key", "crevisto set-key sk-or-v1-xxx")
	}
	orkey := args[0]
	token := requireToken()

	c := newAPIClient()
	c.token = token

	body := map[string]string{"token": token, "email": "", "orkey": orkey}
	resp, data, err := c.post("/app/claim", body)
	if err != nil {
		fail(ExitIntegration, "api_error", "failed to reach platform: "+err.Error())
	}
	if resp.StatusCode != 200 {
		fail(ExitIntegration, "api_error", fmt.Sprintf("set-key returned %d: %s", resp.StatusCode, string(data)))
	}
	outData(jsonRaw(data))
}

// jsonRaw wraps raw JSON bytes in a json.RawMessage for output.
func jsonRaw(data []byte) json.RawMessage {
	return json.RawMessage(data)
}

// cmdGallery lists public gallery images.
func cmdGallery(args []string) {
	toolSlug := ""
	limit := "24"
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--tool", "-t":
			if i+1 < len(args) {
				toolSlug = args[i+1]
				i++
			}
		case "--limit", "-l":
			if i+1 < len(args) {
				limit = args[i+1]
				i++
			}
		}
	}

	c := newAPIClient()
	q := "?limit=" + limit
	if toolSlug != "" {
		q += "&toolSlug=" + toolSlug
	}
	resp, data, err := c.get("/v1/gallery" + q)
	if err != nil {
		fail(ExitIntegration, "api_error", "failed to reach platform: "+err.Error())
	}
	if resp.StatusCode != 200 {
		fail(ExitIntegration, "api_error", fmt.Sprintf("gallery returned %d: %s", resp.StatusCode, string(data)))
	}
	outData(jsonRaw(data))
}

// cmdGalleryDownload downloads a gallery image by id.
func cmdGalleryDownload(args []string) {
	if len(args) < 1 {
		fail(ExitInput, "missing_argument", "gallery-download requires an image id", "crevisto gallery  — list ids")
	}
	id := args[0]
	outputPath := ""
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--output", "-o":
			if i+1 < len(args) {
				outputPath = args[i+1]
				i++
			}
		}
	}
	if outputPath == "" {
		outputPath = "gallery-" + id + ".png"
	}

	c := newAPIClient()
	resp, data, err := c.get("/v1/gallery/" + id + "/image")
	if err != nil {
		fail(ExitIntegration, "api_error", "failed to reach platform: "+err.Error())
	}
	if resp.StatusCode == 404 {
		fail(ExitResource, "not_found", "gallery image not found: "+id)
	}
	if resp.StatusCode != 200 {
		fail(ExitIntegration, "api_error", fmt.Sprintf("download returned %d: %s", resp.StatusCode, string(data)))
	}
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		fail(ExitInternal, "write_error", "failed to write file: "+err.Error())
	}
	logCtx("image saved to %s (%d bytes)", outputPath, len(data))
	outData(map[string]string{"saved": outputPath, "bytes": fmt.Sprintf("%d", len(data))})
}

// cmdGalleryToggle toggles public/private on a generation.
func cmdGalleryToggle(args []string) {
	if len(args) < 1 {
		fail(ExitInput, "missing_argument", "gallery-toggle requires a generation id")
	}
	id := args[0]
	token := requireToken()

	c := newAPIClient()
	c.token = token
	resp, data, err := c.post("/v1/gallery/"+id+"/toggle", map[string]interface{}{})
	if err != nil {
		fail(ExitIntegration, "api_error", "failed to reach platform: "+err.Error())
	}
	if resp.StatusCode == 401 {
		fail(ExitResource, "not_authenticated", "invalid or missing token", "crevisto auth <token>")
	}
	if resp.StatusCode == 403 {
		fail(ExitResource, "forbidden", "not your generation")
	}
	if resp.StatusCode == 404 {
		fail(ExitResource, "not_found", "generation not found: "+id)
	}
	if resp.StatusCode != 200 {
		fail(ExitIntegration, "api_error", fmt.Sprintf("toggle returned %d: %s", resp.StatusCode, string(data)))
	}
	outData(jsonRaw(data))
}

// init ensures http is imported (used by api.go)
var _ = http.StatusOK
var _ = io.EOF
