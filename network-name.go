package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func getNetworkName() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return getNetworkNameMacOS()
	case "linux":
		return getNetworkNameLinux()
	case "windows":
		return getNetworkNameWindows()
	default:
		return "Unknown", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// macOS: find the Wi-Fi interface via networksetup, then query its SSID.
func getNetworkNameMacOS() (string, error) {
	iface, err := findWifiInterfaceMacOS()
	if err != nil {
		return "", err
	}
	if iface == "" {
		return "Unknown", nil
	}

	out, err := exec.Command("/usr/sbin/networksetup", "-getairportnetwork", iface).Output()
	if err != nil {
		return "", err
	}

	// Output: "Current Wi-Fi Network: <SSID>" or
	//         "You are not associated with an AirPort network.\n"
	line := strings.TrimSpace(string(out))
	ssid, found := strings.CutPrefix(line, "Current Wi-Fi Network: ")
	if !found {
		return "Unknown", nil
	}
	return ssid, nil
}

// findWifiInterfaceMacOS returns the BSD device name of the Wi-Fi adapter
// (e.g. "en0") by parsing `networksetup -listallhardwareports`.
func findWifiInterfaceMacOS() (string, error) {
	out, err := exec.Command("/usr/sbin/networksetup", "-listallhardwareports").Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if strings.Contains(line, "Wi-Fi") || strings.Contains(line, "AirPort") {
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(lines[j], "Device: ") {
					return strings.TrimPrefix(lines[j], "Device: "), nil
				}
				if lines[j] == "" {
					break
				}
			}
		}
	}
	return "", nil
}

// Linux: use iwconfig if available.
func getNetworkNameLinux() (string, error) {
	out, err := exec.Command("iwconfig").Output()
	if err != nil {
		return "Unknown", nil
	}

	for line := range strings.SplitSeq(string(out), "\n") {
		if strings.Contains(line, "ESSID:") {
			parts := strings.Split(line, "\"")
			if len(parts) > 1 {
				return parts[1], nil
			}
		}
	}
	return "Unknown", nil
}

// Windows: use netsh to query the active Wi-Fi SSID.
func getNetworkNameWindows() (string, error) {
	out, err := exec.Command("netsh", "wlan", "show", "interfaces").Output()
	if err != nil {
		return "", err
	}

	for line := range strings.SplitSeq(string(out), "\n") {
		// Match " SSID : value" but not "BSSID"
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "SSID") && !strings.HasPrefix(trimmed, "BSSID") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "Unknown", nil
}
