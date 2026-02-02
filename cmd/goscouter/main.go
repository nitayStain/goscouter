package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"goscouter/internal/server"
)

var debugMode bool
var daemonMode bool
var isDaemonChild bool

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Parse flags
	args := os.Args[1:]
	command := args[0]
	flags := args[1:]

	for _, flag := range flags {
		if flag == "--debug" || flag == "--verbose" || flag == "-v" {
			debugMode = true
			server.DebugMode = true
		}
		if flag == "-d" {
			daemonMode = true
		}
		if flag == "--daemon" {
			isDaemonChild = true
		}
	}

	switch command {
	case "run":
		if isDaemonChild {
			// We are the daemon child process, just run normally
			daemonMode = true // Set this so we don't print startup messages
			runCommand()
		} else if daemonMode {
			// User wants to daemonize, so fork
			if err := startDaemon(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Normal foreground run
			runCommand()
		}
	case "stop":
		if err := stopDaemon(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "status":
		if err := statusDaemon(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "build":
		buildCommand()
	case "version", "-v", "--version":
		printVersion()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func runCommand() {
	// Skip startup messages if running as daemon
	if !daemonMode {
		fmt.Println("🚀 Starting GoScouter...")
		if debugMode {
			fmt.Println("🔍 Debug mode enabled")
		}

		// Check for updates in background (not for daemon)
		checkForUpdates()
	}

	if !daemonMode {
		fmt.Println("📦 Checking frontend build...")
	}

	// Find frontend directory (local or installed)
	frontendPath := getFrontendPath()
	if frontendPath == "" {
		log.Fatalf("Frontend not found. Please run 'goscouter build' first or install using 'make install'")
	}

	if !daemonMode {
		fmt.Printf("✅ Frontend found at: %s\n", frontendPath)
		fmt.Println("🌐 Starting server on http://localhost:8080")
		fmt.Println("Press Ctrl+C to stop")
		fmt.Println()
	}

	// Set environment variable for server to find frontend
	os.Setenv("GOSCOUTER_FRONTEND_PATH", frontendPath)

	srv := server.New()
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getFrontendPath() string {
	// Check local development directory first
	localPath := filepath.Join("frontend", "out")
	if _, err := os.Stat(filepath.Join(localPath, "index.html")); err == nil {
		return localPath
	}

	// Check installed directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		installedPath := filepath.Join(homeDir, ".goscouter", "out")
		if _, err := os.Stat(filepath.Join(installedPath, "index.html")); err == nil {
			return installedPath
		}
	}

	return ""
}

func buildCommand() {
	fmt.Println("📦 Building GoScouter...")

	if err := buildFrontend(); err != nil {
		log.Fatalf("Failed to build frontend: %v", err)
	}

	fmt.Println("✅ Build complete!")
	fmt.Println("Run 'goscouter run' to start the application")
}

func buildFrontend() error {
	fmt.Println("📦 Building frontend...")

	frontendPath := "frontend"

	// Check if node_modules exists
	nodeModulesPath := filepath.Join(frontendPath, "node_modules")
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		fmt.Println("📥 Installing frontend dependencies...")
		cmd := exec.Command("npm", "install")
		cmd.Dir = frontendPath

		// Only show npm output if debug mode
		if debugMode {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else {
			// Suppress all output except errors
			cmd.Stdout = io.Discard
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("npm install failed: %w", err)
		}
	}

	// Build frontend
	fmt.Println("⚙️  Compiling...")
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = frontendPath

	// Suppress stdout, capture stderr for errors
	var stderr bytes.Buffer
	if debugMode {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = &stderr
	}

	if err := cmd.Run(); err != nil {
		if !debugMode && stderr.Len() > 0 {
			fmt.Fprintln(os.Stderr, stderr.String())
		}
		return fmt.Errorf("npm build failed: %w", err)
	}

	fmt.Println("✅ Frontend built successfully")
	return nil
}

func printVersion() {
	fmt.Printf("GoScouter v%s\n", Version)
	fmt.Println()

	// Check for updates
	fmt.Println("Checking for updates...")
	latestVersion, err := fetchLatestVersion()
	if err != nil {
		fmt.Printf("Unable to check for updates: %v\n", err)
		return
	}

	if latestVersion == "" {
		fmt.Println("No stable release found")
		return
	}

	if latestVersion == Version {
		fmt.Println("✅ You are running the latest version")
	} else if compareVersions(Version, latestVersion) < 0 {
		fmt.Printf("⚠️  A newer version is available: v%s\n", latestVersion)
		fmt.Println()
		fmt.Println("Update with:")
		fmt.Println("  cd /path/to/goscouter")
		fmt.Println("  git pull")
		fmt.Println("  make install")
	} else {
		fmt.Printf("You are running a development version (v%s > v%s)\n", Version, latestVersion)
	}
}

func printUsage() {
	fmt.Println(`GoScouter - Subdomain Discovery Tool

Usage:
  goscouter <command> [flags]

Commands:
  run       Start the GoScouter web service
  stop      Stop the running daemon
  status    Check if daemon is running
  build     Build the frontend and prepare for production
  version   Show version information and check for updates
  help      Show this help message

Flags:
  -d, --daemon              Run as background daemon
  --debug, --verbose        Enable debug mode (shows DNS lookup logs)

Environment Variables:
  GOSCOUTER_SKIP_VERSION_CHECK=1    Disable automatic update checks

Examples:
  goscouter run              # Start in foreground on http://localhost:8080
  goscouter run -d           # Start as daemon (background)
  goscouter run --debug      # Start with debug logging
  goscouter stop             # Stop the daemon
  goscouter status           # Check daemon status
  goscouter build            # Build the frontend (quiet mode)
  goscouter version          # Show version and check for updates

For more information, visit: https://github.com/nitayStain/goscouter`)
}
