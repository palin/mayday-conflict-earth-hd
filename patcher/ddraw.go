package main

// cnc-ddraw (MIT licence, github.com/FunkyFr3sh/cnc-ddraw): a DirectDraw replacement.
// Without it Windows 10/11 refuses the game's 256-colour screen mode
// ("Direct Draw Init Failed (80004001)").

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

//go:embed ddraw
var ddrawFiles embed.FS

const DdrawDll = "ddraw.dll"
const DdrawOff = "ddraw.dll.off"

// DdrawState: "on" (ddraw.dll present), "off" (renamed by us), "none" (not installed)
func DdrawState(dir string) string {
	if _, err := os.Stat(filepath.Join(dir, DdrawDll)); err == nil {
		return "on"
	}
	if _, err := os.Stat(filepath.Join(dir, DdrawOff)); err == nil {
		return "off"
	}
	return "none"
}

const DdrawConfig = "cnc-ddraw config.exe"

func SetDdraw(dir string, on bool) error {
	if err := setDdraw(dir, on); err != nil || !on {
		return err
	}
	// the settings tool is added also to folders that already had cnc-ddraw
	dst := filepath.Join(dir, DdrawConfig)
	if _, err := os.Stat(dst); err == nil {
		return nil
	}
	data, err := ddrawFiles.ReadFile("ddraw/" + DdrawConfig)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0666)
}

func setDdraw(dir string, on bool) error {
	cur := DdrawState(dir)
	dll, off := filepath.Join(dir, DdrawDll), filepath.Join(dir, DdrawOff)
	switch {
	case !on && cur == "on":
		return os.Rename(dll, off)
	case on && cur == "off":
		return os.Rename(off, dll)
	case on && cur == "none":
		return fs.WalkDir(ddrawFiles, "ddraw", func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			rel, _ := filepath.Rel("ddraw", filepath.FromSlash(p))
			dst := filepath.Join(dir, rel)
			if path.Base(p) == "ddraw.ini" {
				if _, err := os.Stat(dst); err == nil {
					return nil // keep the player's own settings
				}
			}
			data, err := ddrawFiles.ReadFile(p)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(dst), 0777); err != nil {
				return err
			}
			return os.WriteFile(dst, data, 0666)
		})
	}
	return nil
}
