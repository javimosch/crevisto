package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Exit code ranges per cli-output-spec
const (
	ExitOK         = 0
	ExitInput      = 80 // 80-89: input/validation
	ExitResource   = 90 // 90-99: precondition/resource
	ExitIntegration = 100 // 100-109: external/integration
	ExitInternal   = 110 // 110-119: internal/bug
)

// outOK prints {"ok":true,...} to stdout and exits 0.
func outOK(data interface{}) {
	j, _ := json.Marshal(map[string]interface{}{"ok": true, "data": data})
	fmt.Println(string(j))
	os.Exit(0)
}

// outData prints raw data to stdout (no ok wrapper) and exits 0.
func outData(data interface{}) {
	j, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(j))
	os.Exit(0)
}

// fail prints a typed error to stderr and exits with the given code.
func fail(code int, typ string, msg string, suggestions ...string) {
	recoverable := code >= 100
	err := map[string]interface{}{
		"ok": false,
		"error": map[string]interface{}{
			"code":        code,
			"type":        typ,
			"message":     msg,
			"recoverable": recoverable,
		},
	}
	if len(suggestions) > 0 {
		err["error"].(map[string]interface{})["suggestions"] = suggestions
	}
	j, _ := json.Marshal(err)
	fmt.Fprintln(os.Stderr, string(j))
	os.Exit(code)
}

// logCtx prints context to stderr (never to stdout).
func logCtx(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
