package main

// Patch engine: works on the Polish Mayday: Conflict Earth 1.1p.
// Every fix is a group of byte spans (original bytes / new bytes). The files are
// always rebuilt from the untouched original plus the selected groups, so nothing
// can be applied twice and "Restore" is exact.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Span struct {
	Off       int
	Orig, New []byte
}

type FileSpec struct {
	Name string
	Size int
	SHA  string
}

type Group struct {
	ID       string
	File     *FileSpec
	Section  string
	Label    string
	Requires []string // other fixes this one needs (applied together)
}

var (
	ExeSpec = &FileSpec{"Mayday.exe", 684032, "9f9ab3f11a25bb2241537e984e1e487eb8e97282dbeb20f6ab0d137f359ebfb4"}
	ArtSpec = &FileSpec{"mayday.art", 23185090, "08016de9343d6fd0971ff6d9ece4904d33bf7d26e1fe6502a7edd6b0d347c613"}
	Files   = []*FileSpec{ExeSpec, ArtSpec}
)

const (
	Version  = "1.0.1"
	AppTitle = "Mayday HD " + Version + " by palin"
	AppHead  = "Mayday HD " + Version + " dla Mayday: Conflict Earth PL (1.1p)"
)

// Changelog shown in the patcher window ("" starts a new heading line with "#").
var Changelog = []string{
	"#Zmiany w 1.0.1",
	"Naprawione klawisze w grze przy odświeżaniu powyżej 60 Hz: pauza (ESC), cele misji (M), opcje (O), kapitulacja, grupy bojowe (CTRL+cyfra) i tryby walki (F1–F4)",
	"#Grafika",
	"Rozdzielczość 1920×1080 - ustawiona jako domyślna (zastępuje 1024×768)",
	"Płynniejszy obraz - zgodny z odświeżaniem monitora- np. 144Hz (prędkość gry bez zmian)",
	"Płynniejsze przewijanie mapy",
	"#Minimapa w grze",
	"Lepsza przejrzystość i orientacja (jednolite tło, obrys świata, jaśniejszy teren, obwódki jednostek)",
	"Dokładna ramka widoku (kliknięcie na minimapie precyzyjnie centruje widok)",
	"#Rozgrywka",
	"Poprawione zaznaczanie jednostek (wciąż mogą być jednak problemy przy zaznaczaniu jednostek na różnych wysokościach terenu)",
	"#Terminal gry",
	"Spokojniejsze przewijanie widoku prawo-lewo",
	"#Windows 10/11",
	"Dodanie narzędzia \"cnc-ddraw\" (wymagane, by gra działała) + program „cnc-ddraw config” (w folderze z grą), gdzie można dostosować sobie mnóstwo opcji związanych z wyświetlaniem gry",
	"Naprawa potencjalnego problemu z wyciszonym dźwiękiem gry",
	"Instalacja kodeka Power!Video z płyty gry – działają filmy w intro i odprawach misji (Windows jednorazowo zapyta o zgodę administratora)",
	"#Intro",
	"Domyślnie wyłączone aby przyspieszyć uruchamianie gry (opcja poniżej do włączenia intro)",
}

// AllFixes is the full patch: every fix is always applied.
func AllFixes() map[string]bool {
	m := map[string]bool{}
	for _, g := range Groups {
		m[g.ID] = true
	}
	return m
}

// IsComplete reports whether a recognised file carries every fix of this version.
func (st *FileState) IsComplete() bool {
	if !st.Known {
		return false
	}
	for _, g := range groupsOf(st.Spec) {
		if !st.On[g.ID] {
			return false
		}
	}
	return true
}

const BackupDir = "MaydayHD_backup"

var groupSpans = map[string][]Span{}

// legacySpans: changes made by older releases that no current group produces (only Mayday.exe);
// used to recognise those files and rebuild the original from them.
var legacySpans = map[string][]Span{}

var Groups = []Group{
	{"selection", ExeSpec, "Rozgrywka", "Lepsze zaznaczanie ramką (jednostka zaznacza się, gdy ramka jej dotknie)", nil},
	{"terminal", ExeSpec, "Rozgrywka", "Spokojniejsze przewijanie terminala (stała prędkość na każdym komputerze)", nil},
	{"panel", ArtSpec, "Minimapa", "Jednolite ciemne tło minimapy", nil},
	{"outline", ExeSpec, "Minimapa", "Obrys świata na minimapie (stalowoniebieski)", nil},
	{"terrain", ExeSpec, "Minimapa", "Jaśniejszy teren na minimapie (np. mapy śnieżne)", nil},
	{"rings", ExeSpec, "Minimapa", "Czarne obwódki jednostek i budynków na minimapie", nil},
	{"frame", ExeSpec, "Minimapa", "Dokładna ramka widoku + kliknięcie centruje widok", nil},
	{"fullhd", ExeSpec, "Grafika", "Zamień rozdzielczość 1024x768 na 1920x1080 (Full HD)", nil},
	{"fps60", ExeSpec, "Grafika", "Płynny obraz: 60 klatek/s (prędkość gry bez zmian)", nil},
	{"smooth", ExeSpec, "Grafika", "Płynne przewijanie mapy (wymaga Full HD, 60 klatek/s i dokładnej ramki)", []string{"fullhd", "fps60", "frame"}},
	{"hz", ExeSpec, "Grafika", "Odświeżanie jak monitor, np. 144 Hz zamiast 60 (wymaga 60 klatek/s)", []string{"fps60"}},
	{"sound", ExeSpec, "Dźwięk", "Naprawa „brak dźwięku po starcie” (Windows 10/11)", nil},
}

func hx(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func shaHex(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// FileState describes one game file in a folder.
type FileState struct {
	Spec       *FileSpec
	Present    bool
	Known      bool            // current file = original + some of our fixes
	On         map[string]bool // fixes currently applied (only when Known)
	Original   []byte          // pristine original (rebuilt or from backup), nil if unavailable
	FromBackup bool
	Current    []byte
}

func groupsOf(spec *FileSpec) []Group {
	var r []Group
	for _, g := range Groups {
		if g.File == spec {
			r = append(r, g)
		}
	}
	return r
}

// WithRequired adds the fixes that the selected fixes depend on.
func WithRequired(want map[string]bool) map[string]bool {
	out := map[string]bool{}
	for k, v := range want {
		out[k] = v
	}
	for changed := true; changed; {
		changed = false
		for _, g := range Groups {
			if out[g.ID] {
				for _, r := range g.Requires {
					if !out[r] {
						out[r] = true
						changed = true
					}
				}
			}
		}
	}
	return out
}

func validCombo(want map[string]bool) bool {
	for _, g := range Groups {
		if want[g.ID] {
			for _, r := range g.Requires {
				if !want[r] {
					return false
				}
			}
		}
	}
	return true
}

type layer struct {
	g int
	b byte
}

func Analyze(dir string, spec *FileSpec) *FileState {
	st := &FileState{Spec: spec, On: map[string]bool{}}
	data, err := os.ReadFile(filepath.Join(dir, spec.Name))
	if err == nil {
		st.Present = true
		st.Current = data
		if len(data) == spec.Size {
			gs := groupsOf(spec)
			// every byte any fix touches: its original value and the value each fix writes
			origAt := map[int]byte{}
			layers := map[int][]layer{}
			for gi, g := range gs {
				for _, s := range groupSpans[g.ID] {
					for k := range s.New {
						origAt[s.Off+k] = s.Orig[k]
						layers[s.Off+k] = append(layers[s.Off+k], layer{gi, s.New[k]})
					}
				}
			}
			orig := append([]byte(nil), data...)
			for p, b := range origAt {
				orig[p] = b
			}
			if spec == ExeSpec {
				for _, ls := range legacySpans {
					for _, sp := range ls {
						copy(orig[sp.Off:], sp.Orig)
					}
				}
			}
			if shaHex(orig) == spec.SHA {
				for mask := 0; mask < 1<<len(gs); mask++ {
					want := map[string]bool{}
					for gi, g := range gs {
						if mask&(1<<gi) != 0 {
							want[g.ID] = true
						}
					}
					if !validCombo(want) {
						continue
					}
					match := true
					for p, ls := range layers {
						exp := origAt[p]
						for _, l := range ls { // later fixes override earlier ones
							if mask&(1<<l.g) != 0 {
								exp = l.b
							}
						}
						if data[p] != exp {
							match = false
							break
						}
					}
					if match {
						st.Known = true
						st.Original = orig
						st.On = want
						break
					}
				}
				if !st.Known {
					st.Original = orig // an older release: the original is rebuilt, no backup needed
				}
			}
		}
	}
	if !st.Known {
		st.On = map[string]bool{}
		if st.Original != nil {
			return st
		}
		if b, err := os.ReadFile(filepath.Join(dir, BackupDir, spec.Name)); err == nil && shaHex(b) == spec.SHA {
			st.Original = b
			st.FromBackup = true
		}
	}
	return st
}

func Build(orig []byte, spec *FileSpec, want map[string]bool) []byte {
	want = WithRequired(want)
	out := append([]byte(nil), orig...)
	for _, g := range groupsOf(spec) {
		if want[g.ID] {
			for _, s := range groupSpans[g.ID] {
				copy(out[s.Off:], s.New)
			}
		}
	}
	return out
}

func writeReplace(path string, data []byte) error {
	tmp := path + ".mdhd_tmp"
	if err := os.WriteFile(tmp, data, 0666); err != nil {
		return err
	}
	_ = os.Chmod(path, 0666) // clear read-only (files copied from a CD)
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// ApplyFile rebuilds one file from its original with the wanted fixes.
func ApplyFile(dir string, st *FileState, want map[string]bool) (bool, error) {
	if st.Original == nil {
		return false, fmt.Errorf("%s: nie znaleziono oryginalnej wersji pliku", st.Spec.Name)
	}
	bdir := filepath.Join(dir, BackupDir)
	bpath := filepath.Join(bdir, st.Spec.Name)
	if b, err := os.ReadFile(bpath); err != nil || shaHex(b) != st.Spec.SHA {
		if err := os.MkdirAll(bdir, 0777); err != nil {
			return false, err
		}
		if err := writeReplace(bpath, st.Original); err != nil {
			return false, fmt.Errorf("kopia zapasowa %s: %w", st.Spec.Name, err)
		}
	}
	out := Build(st.Original, st.Spec, want)
	if bytes.Equal(out, st.Current) {
		return false, nil
	}
	path := filepath.Join(dir, st.Spec.Name)
	if err := writeReplace(path, out); err != nil {
		return false, err
	}
	check, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(check, out) {
		return true, fmt.Errorf("%s: zapis nie powiódł się (plik różni się po zapisie)", st.Spec.Name)
	}
	return true, nil
}

// ---- Intro ----

const IntroName = "Intro.exe"
const IntroOff = "Intro.exe.off"

func IntroState(dir string) string {
	if _, err := os.Stat(filepath.Join(dir, IntroName)); err == nil {
		return "on"
	}
	if _, err := os.Stat(filepath.Join(dir, IntroOff)); err == nil {
		return "off"
	}
	return "missing"
}

func SetIntro(dir string, on bool) error {
	cur := IntroState(dir)
	if cur == "missing" {
		return errors.New("nie znaleziono pliku Intro.exe")
	}
	exe, off := filepath.Join(dir, IntroName), filepath.Join(dir, IntroOff)
	if on && cur == "off" {
		if err := os.Rename(off, exe); err != nil {
			return err
		}
	}
	if !on && cur == "on" {
		if err := os.Rename(exe, off); err != nil {
			return err
		}
	}
	if on {
		return setColor256(exe)
	}
	return nil
}

// merge "256COLOR" into an AppCompat "Layers" value ("~ FLAG FLAG")
func mergeLayers(cur string) string {
	f := strings.Fields(cur)
	for _, x := range f {
		if strings.EqualFold(x, "256COLOR") {
			return cur
		}
	}
	if len(f) == 0 {
		return "~ 256COLOR"
	}
	if f[0] != "~" {
		f = append([]string{"~"}, f...)
	}
	return strings.Join(append(f, "256COLOR"), " ")
}
