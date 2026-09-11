package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const defaultDaemonPort = 9192
const daemonLogDir = ".crevisto"

// cmdServe runs the daemon in the foreground (cli-daemon-spec).
func cmdServe(args []string) {
	host := "127.0.0.1"
	port := defaultDaemonPort
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--host":
			if i+1 < len(args) {
				host = args[i+1]
				i++
			}
		case "--port":
			if i+1 < len(args) {
				port, _ = strconv.Atoi(args[i+1])
				i++
			}
		}
	}

	mux := http.NewServeMux()

	// GET /_health — open, fast, no auth
	mux.HandleFunc("/_health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":      true,
			"service": appName,
			"pid":     os.Getpid(),
		})
	})

	// POST /_shutdown — Bearer-gated when off-loopback
	mux.HandleFunc("/_shutdown", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", 405)
			return
		}
		// If bound off-loopback, require token
		if host != "127.0.0.1" && host != "localhost" {
			token := readShutdownToken()
			auth := r.Header.Get("Authorization")
			if auth != "Bearer "+token {
				http.Error(w, "unauthorized", 401)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true,"shutdown":true}`))
		go func() {
			time.Sleep(100 * time.Millisecond)
			os.Exit(0)
		}()
	})

	// GET /guide — same JSON as `crevisto guide`
	mux.HandleFunc("/guide", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, guideJSON)
	})

	// GET /llms.txt — breadcrumb
	mux.HandleFunc("/llms.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "# Crevisto CLI\n\nGuide: GET /guide\nCLI command: crevisto guide\n")
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	logCtx("serve: listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fail(ExitInternal, "serve_error", err.Error())
	}
}

// cmdDaemon manages the background daemon lifecycle.
func cmdDaemon(args []string) {
	if len(args) < 1 {
		fail(ExitInput, "missing_argument", "daemon requires start|stop|status", "crevisto daemon start")
	}

	sub := args[0]
	port := defaultDaemonPort
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--port":
			if i+1 < len(args) {
				port, _ = strconv.Atoi(args[i+1])
				i++
			}
		}
	}

	switch sub {
	case "start":
		daemonStart(port)
	case "stop":
		daemonStop(port)
	case "status":
		daemonStatus(port)
	default:
		fail(ExitInput, "invalid_argument", "daemon subcommand must be start|stop|status")
	}
}

func daemonLogPath(port int) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, daemonLogDir, fmt.Sprintf("serve-%d.log", port))
}

func daemonStart(port int) {
	// Check if already running
	if daemonIsUp(port) {
		outOK(map[string]interface{}{"daemon": "already_running", "port": port})
	}

	// Spawn ourselves detached
	exe, _ := os.Executable()
	logPath := daemonLogPath(port)
	cmd := exec.Command(exe, "serve", "--port", strconv.Itoa(port))
	cmd.Stdout = openLog(logPath)
	cmd.Stderr = cmd.Stdout
	cmd.SysProcAttr = setSysProcAttr()

	if err := cmd.Start(); err != nil {
		fail(ExitInternal, "spawn_error", "failed to start daemon: "+err.Error())
	}

	// Poll /_health until up (never fixed sleep)
	for i := 0; i < 50; i++ {
		if daemonIsUp(port) {
			outOK(map[string]interface{}{
				"daemon": "started",
				"host":   "127.0.0.1",
				"port":   port,
				"log":    logPath,
			})
		}
		time.Sleep(100 * time.Millisecond)
	}

	fail(ExitIntegration, "daemon_timeout", "daemon did not become healthy within 5s")
}

func daemonStop(port int) {
	if !daemonIsUp(port) {
		outOK(map[string]interface{}{"daemon": "not_running"})
	}

	// Send /_shutdown
	resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d/_shutdown", port), "application/json", nil)
	if err != nil {
		fail(ExitIntegration, "shutdown_error", "failed to send shutdown: "+err.Error())
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	// Poll until down
	for i := 0; i < 50; i++ {
		if !daemonIsUp(port) {
			outOK(map[string]interface{}{"daemon": "stopped"})
		}
		time.Sleep(100 * time.Millisecond)
	}

	fail(ExitIntegration, "shutdown_timeout", "daemon did not stop within 5s")
}

func daemonStatus(port int) {
	if daemonIsUp(port) {
		outOK(map[string]interface{}{"daemon": "running", "host": "127.0.0.1", "port": port})
	}
	// exit 90 = not_running (resource error per spec)
	fail(ExitResource, "not_running", "daemon is not running")
}

func daemonIsUp(port int) bool {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/_health", port))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func openLog(path string) *os.File {
	os.MkdirAll(filepath.Dir(path), 0700)
	f, _ := os.Create(path)
	return f
}

func readShutdownToken() string {
	home, _ := os.UserHomeDir()
	data, _ := os.ReadFile(filepath.Join(home, daemonLogDir, "shutdown-token"))
	return strings.TrimSpace(string(data))
}

// init ensures net is used
var _ = net.Listen
