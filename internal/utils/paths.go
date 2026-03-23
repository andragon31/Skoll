package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetDefaultDataDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".skoll")
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func ResolveBinaryPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "skoll"
	}
	if res, err := filepath.EvalSymlinks(exe); err == nil {
		return res
	}
	return exe
}

func OpenCodeConfigDir() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		home, _ := os.UserHomeDir()
		appData = filepath.Join(home, "Library", "Application Support")
	}
	return filepath.Join(appData, "OpenCode")
}

func StripJSONC(data []byte) []byte {
	var out []byte
	i := 0
	for i < len(data) {
		if data[i] == '"' {
			out = append(out, data[i])
			i++
			for i < len(data) && data[i] != '"' {
				if data[i] == '\\' && i+1 < len(data) {
					out = append(out, data[i], data[i+1])
					i += 2
					continue
				}
				out = append(out, data[i])
				i++
			}
			if i < len(data) {
				out = append(out, data[i])
				i++
			}
			continue
		}
		if i+1 < len(data) && data[i] == '/' && data[i+1] == '/' {
			for i < len(data) && data[i] != '\n' {
				i++
			}
			continue
		}
		if i+1 < len(data) && data[i] == '/' && data[i+1] == '*' {
			i += 2
			for i+1 < len(data) && !(data[i] == '*' && data[i+1] == '/') {
				i++
			}
			if i+1 < len(data) {
				i += 2
			} else {
				i = len(data)
			}
			continue
		}
		out = append(out, data[i])
		i++
	}
	return out
}

func PatchSkollBINLine(src []byte, absBin string) []byte {
	const marker = `const SKOLL_BIN = process.env.SKOLL_BIN ?? "skoll"`
	replacement := fmt.Sprintf(`const SKOLL_BIN = process.env.SKOLL_BIN ?? Bun.which("skoll") ?? %q`, absBin)
	return []byte(strings.Replace(string(src), marker, replacement, 1))
}
