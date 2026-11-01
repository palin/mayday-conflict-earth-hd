package main

// Power!Video (PVEZ) video codec for the intro and the mission briefings.
// The codec files are taken from the player's own game CD (PV_WIN\SETUP_PV.EXE,
// a Wise installer that stores them as raw deflate streams), never shipped with the patcher.

import (
	"bufio"
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	SetupPVRel  = `PV_WIN\SETUP_PV.EXE`
	SetupPVSize = 3732620
	SetupPVSHA  = "85c34db0be2cbcf7f8653abd20ccac98f3d52f9111723f6f587b053df146f72e"
	CodecReg    = "vidc.pvez"
)

type codecFile struct {
	Name string
	Off  int // start of the deflate stream inside SETUP_PV.EXE
	Size int
	SHA  string
}

var CodecFiles = []codecFile{
	{"pvez32.drv", 0x17149a, 105472, "15a79bacb5444a43bc53f751218f9b0594e8f1a164e3224e563c6254afef9de4"},
	{"Pvezw32d.dll", 0x17d12f, 247808, "00fa9357556e9389821607b13427cc1078e1d17541fb0c08f2c6ad8836935ef0"},
}

// ExtractCodec reads the codec files out of the CD's SETUP_PV.EXE.
func ExtractCodec(setupPath string) (map[string][]byte, error) {
	d, err := os.ReadFile(setupPath)
	if err != nil {
		return nil, err
	}
	if len(d) != SetupPVSize || shaHex(d) != SetupPVSHA {
		return nil, fmt.Errorf("%s: nieznana wersja pliku", setupPath)
	}
	out := map[string][]byte{}
	for _, f := range CodecFiles {
		b, err := io.ReadAll(io.LimitReader(flate.NewReader(bytes.NewReader(d[f.Off:])), int64(f.Size)+1))
		if err != nil || len(b) != f.Size || shaHex(b) != f.SHA {
			return nil, fmt.Errorf("%s: nie udało się wypakować %s", setupPath, f.Name)
		}
		out[f.Name] = b
	}
	return out, nil
}

// cdDrive reads the CD drive letter the game was installed from (install.ini, CDPfad=F:).
func cdDrive(gameDir string) string {
	f, err := os.Open(filepath.Join(gameDir, "install.ini"))
	if err != nil {
		return ""
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		k, v, ok := strings.Cut(strings.TrimSpace(sc.Text()), "=")
		if ok && strings.EqualFold(strings.TrimSpace(k), "CDPfad") {
			v = strings.TrimSpace(v)
			if len(v) >= 2 && v[1] == ':' {
				return strings.ToUpper(v[:1]) + ":"
			}
		}
	}
	return ""
}

// FindSetupPV looks for SETUP_PV.EXE on the game CD: first the drive from install.ini,
// then every drive letter.
func FindSetupPV(gameDir string) string {
	var drives []string
	if d := cdDrive(gameDir); d != "" {
		drives = append(drives, d)
	}
	for c := 'D'; c <= 'Z'; c++ {
		drives = append(drives, string(c)+":")
	}
	for _, d := range drives {
		p := d + `\` + SetupPVRel
		if st, err := os.Stat(p); err == nil && st.Size() == SetupPVSize {
			return p
		}
	}
	return ""
}

// WriteCodecFiles copies the codec files into the 32-bit system folder (SysWOW64).
func WriteCodecFiles(sysDir string, files map[string][]byte) error {
	for _, f := range CodecFiles {
		p := filepath.Join(sysDir, f.Name)
		if cur, err := os.ReadFile(p); err == nil && bytes.Equal(cur, files[f.Name]) {
			continue
		}
		if err := os.WriteFile(p, files[f.Name], 0644); err != nil {
			return err
		}
	}
	return nil
}

// codecFilesPresent: both files in sysDir with the expected content.
func codecFilesPresent(sysDir string) bool {
	for _, f := range CodecFiles {
		st, err := os.Stat(filepath.Join(sysDir, f.Name))
		if err != nil || st.Size() <= 0 {
			return false
		}
	}
	return true
}
