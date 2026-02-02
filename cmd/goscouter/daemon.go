package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	pidFile = ".goscouter.pid"
	logFile = ".goscouter.log"
)

func getPIDFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, pidFile)
}

func getLogFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, logFile)
}

func startDaemon() error {
	// Check if already running
	if isRunning() {
		pid, _ := readPID()
		return fmt.Errorf("goscouter is already running (PID: %d)", pid)
	}

	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Create log file
	logPath := getLogFilePath()
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	// Start process in background
	cmd := exec.Command(execPath, "run", "--daemon")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	// Save PID before releasing
	pid := cmd.Process.Pid

	// Write PID file
	if err := writePID(pid); err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	// Detach from parent
	cmd.Process.Release()

	fmt.Printf("✓ GoScouter started in daemon mode (PID: %d)\n", pid)
	fmt.Printf("  Logs: %s\n", logPath)
	fmt.Printf("  Web Interface: http://localhost:8080\n")

	return nil
}

func stopDaemon() error {
	pid, err := readPID()
	if err != nil {
		return fmt.Errorf("goscouter is not running")
	}

	// Find process
	process, err := os.FindProcess(pid)
	if err != nil {
		removePID()
		return fmt.Errorf("goscouter is not running")
	}

	// Send SIGTERM
	if err := process.Signal(syscall.SIGTERM); err != nil {
		// Process might be already dead
		removePID()
		return fmt.Errorf("goscouter is not running")
	}

	// Remove PID file
	removePID()

	fmt.Printf("✓ GoScouter stopped (PID: %d)\n", pid)
	return nil
}

func statusDaemon() error {
	if !isRunning() {
		fmt.Println("GoScouter is not running")
		return nil
	}

	pid, _ := readPID()
	fmt.Printf("GoScouter is running (PID: %d)\n", pid)
	fmt.Printf("  Logs: %s\n", getLogFilePath())
	fmt.Printf("  Web Interface: http://localhost:8080\n")

	return nil
}

func isRunning() bool {
	pid, err := readPID()
	if err != nil {
		return false
	}

	// Check if process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		removePID()
		return false
	}

	// Send signal 0 to check if process is alive
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		removePID()
		return false
	}

	return true
}

func readPID() (int, error) {
	data, err := os.ReadFile(getPIDFilePath())
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return pid, nil
}

func writePID(pid int) error {
	pidPath := getPIDFilePath()
	return os.WriteFile(pidPath, []byte(fmt.Sprintf("%d", pid)), 0644)
}

func removePID() {
	os.Remove(getPIDFilePath())
}
