package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// FIXME: This is not working in MacOS! Fix it.

func getNetworkName() (string, error) {
	if runtime.GOOS == "darwin" { // macOS
		return getNetworkNameMacOS()
	} else if runtime.GOOS == "linux" { // Linux
		return getNetworkNameLinux()
	} else if runtime.GOOS == "windows" { // Windows
		return getNetworkNameWindows()
	}

	return "Unknown", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
}

// macOS: Use airport command
func getNetworkNameMacOS() (string, error) {
	cmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-I")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, " SSID:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "Unknown", nil
}

// Linux: Use iwconfig (if available)
func getNetworkNameLinux() (string, error) {
	cmd := exec.Command("iwconfig")
	output, err := cmd.Output()
	if err != nil {
		// iwconfig might not be available, or the interface might not be Wi-Fi
		return "Unknown", nil // Or return an error if you prefer
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "ESSID:") {
			parts := strings.Split(line, "\"")
			if len(parts) > 1 {
				return parts[1], nil
			}
		}
	}

	return "Unknown", nil
}

// Windows: Use netsh (very basic example, needs refinement)
func getNetworkNameWindows() (string, error) {
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "SSID") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "Unknown", nil
}
