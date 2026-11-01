package main

// mayday.cfg is a Windows INI file written by the game ([Settings] iXResolution=640|800|1024).
// With the Full HD fix, 1024 means 1920x1080.

import (
	"os"
	"path/filepath"
	"strings"
)

const CfgName = "mayday.cfg"

// SetCfgValue sets key=value in [section], keeping every other line as it is.
func SetCfgValue(dir, section, key, value string) error {
	path := filepath.Join(dir, CfgName)
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	inSec, secFound, done := false, false, false
	insertAt := -1
	for i, l := range lines {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
			if inSec && !done {
				insertAt = i
			}
			inSec = strings.EqualFold(strings.Trim(t, "[]"), section)
			if inSec {
				secFound = true
			}
			continue
		}
		if inSec && !done {
			if k, _, ok := strings.Cut(t, "="); ok && strings.EqualFold(strings.TrimSpace(k), key) {
				lines[i] = key + "=" + value
				done = true
			}
		}
	}
	switch {
	case done:
	case secFound && insertAt >= 0:
		lines = append(lines[:insertAt], append([]string{key + "=" + value}, lines[insertAt:]...)...)
	case secFound: // section is the last one
		lines = append(lines, key+"="+value)
	default:
		lines = append(lines, "["+section+"]", key+"="+value)
	}
	out := strings.Join(lines, "\r\n") + "\r\n"
	if string(data) == out {
		return nil
	}
	return writeReplace(path, []byte(out))
}

func SetFullHDDefault(dir string) error { return SetCfgValue(dir, "Settings", "iXResolution", "1024") }
