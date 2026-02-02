package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	versionCheckURL     = "https://api.github.com/repos/nitayStain/goscouter/releases/latest"
	versionCheckTimeout = 3 * time.Second
	cacheDuration       = 24 * time.Hour
)

type GitHubRelease struct {
	TagName    string    `json:"tag_name"`
	Name       string    `json:"name"`
	Draft      bool      `json:"draft"`
	PreRelease bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
	Body       string    `json:"body"`
}

type VersionCache struct {
	LatestVersion string    `json:"latest_version"`
	CheckedAt     time.Time `json:"checked_at"`
}

func checkForUpdates() {
	// Skip if version check disabled via env
	if os.Getenv("GOSCOUTER_SKIP_VERSION_CHECK") == "1" {
		return
	}

	// Check cache first
	cachedVersion, cacheValid := getCachedVersion()
	var latestVersion string
	var err error

	if cacheValid && cachedVersion != "" {
		latestVersion = cachedVersion
	} else {
		// Fetch latest version (blocking)
		latestVersion, err = fetchLatestVersion()
		if err != nil {
			// Silently fail - don't bother user with network issues
			if debugMode {
				fmt.Fprintf(os.Stderr, "[DEBUG] Version check failed: %v\n", err)
			}
			return
		}

		// Cache the result
		cacheVersion(latestVersion)
	}

	// Compare versions
	if latestVersion != "" && compareVersions(Version, latestVersion) < 0 {
		promptForUpdate(latestVersion)
	}
}

func fetchLatestVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), versionCheckTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", versionCheckURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("goscouter/%s", Version))

	client := &http.Client{Timeout: versionCheckTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	// Skip drafts and pre-releases
	if release.Draft || release.PreRelease {
		return "", nil
	}

	// Clean version tag (remove 'v' prefix)
	version := strings.TrimPrefix(release.TagName, "v")
	return version, nil
}

func getCachedVersion() (string, bool) {
	cacheFile := getCacheFilePath()
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return "", false
	}

	var cache VersionCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return "", false
	}

	// Check if cache is still valid
	if time.Since(cache.CheckedAt) > cacheDuration {
		return "", false
	}

	return cache.LatestVersion, true
}

func cacheVersion(version string) {
	cache := VersionCache{
		LatestVersion: version,
		CheckedAt:     time.Now(),
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return
	}

	cacheFile := getCacheFilePath()
	os.MkdirAll(filepath.Dir(cacheFile), 0755)
	os.WriteFile(cacheFile, data, 0644)
}

func getCacheFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".goscouter", "version_cache.json")
}

func promptForUpdate(latestVersion string) {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Printf("║  ⚠️  Update Available: v%s → v%s%-*s║\n",
		Version, latestVersion,
		48-len(Version)-len(latestVersion), "")
	fmt.Println("║                                                               ║")
	fmt.Println("║  A newer version is available and may contain important       ║")
	fmt.Println("║  security fixes and improvements.                             ║")
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println()

	fmt.Print("Would you like to update now? [Y/n]: ")

	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))

	// Default to yes if just Enter is pressed
	if response == "" || response == "y" || response == "yes" {
		if err := performUpdate(); err != nil {
			fmt.Printf("\n❌ Update failed: %v\n", err)
			fmt.Println("\nManual update:")
			fmt.Println("  cd /path/to/goscouter")
			fmt.Println("  git pull")
			fmt.Println("  make install")
			fmt.Println()
			fmt.Print("Press Enter to continue with current version...")
			fmt.Scanln(&response)
			return
		}

		fmt.Println("\n✅ Update successful! Restarting...")
		time.Sleep(1 * time.Second)
		restartApp()
	} else {
		fmt.Println("\nContinuing with current version...")
		fmt.Println("To skip this prompt: export GOSCOUTER_SKIP_VERSION_CHECK=1")
		fmt.Println()
	}
}

func compareVersions(v1, v2 string) int {
	// Simple version comparison (major.minor.patch)
	// Returns: -1 if v1 < v2, 0 if equal, 1 if v1 > v2

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &p1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &p2)
		}

		if p1 < p2 {
			return -1
		}
		if p1 > p2 {
			return 1
		}
	}

	return 0
}

func performUpdate() error {
	fmt.Println("\n🔄 Updating GoScouter...")

	// Check if we're in a git repository
	if !isGitRepo() {
		return fmt.Errorf("not in a git repository - please update manually")
	}

	// Check if there are uncommitted changes
	if hasUncommittedChanges() {
		fmt.Println("⚠️  You have uncommitted changes. Stashing them...")
		if err := execCommand("git", "stash", "push", "-m", "goscouter-auto-update"); err != nil {
			return fmt.Errorf("failed to stash changes: %w", err)
		}
		defer execCommand("git", "stash", "pop")
	}

	// Pull latest changes
	fmt.Println("📥 Pulling latest changes...")
	if err := execCommand("git", "pull"); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	// Check if we're installed or running from source
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// If installed in /usr/local/bin, use make install
	if strings.Contains(execPath, "/usr/local/bin") {
		fmt.Println("🔨 Installing update...")
		if err := execCommand("make", "install"); err != nil {
			return fmt.Errorf("make install failed: %w", err)
		}
	} else {
		// Running from source, rebuild
		fmt.Println("🔨 Building update...")
		if err := execCommand("make", "build"); err != nil {
			return fmt.Errorf("make build failed: %w", err)
		}
	}

	return nil
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func hasUncommittedChanges() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func restartApp() {
	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get executable path: %v\n", err)
		os.Exit(1)
	}

	// Get current arguments
	args := os.Args

	// Execute the new version
	if err := syscall.Exec(execPath, args, os.Environ()); err != nil {
		fmt.Printf("Failed to restart: %v\n", err)
		fmt.Println("Please restart manually: goscouter run")
		os.Exit(1)
	}
}
