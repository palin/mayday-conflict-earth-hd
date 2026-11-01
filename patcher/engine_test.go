package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func load(t *testing.T, p string) []byte {
	b, err := os.ReadFile(p)
	if err != nil {
		// reference builds of the game files are not part of the repository
		t.Skipf("%s not available (game files are not distributed)", p)
	}
	return b
}

func setup(t *testing.T, exe []byte) string {
	d := t.TempDir()
	os.WriteFile(filepath.Join(d, "Mayday.exe"), exe, 0444) // read-only like a CD copy
	os.WriteFile(filepath.Join(d, "mayday.art"), load(t, "../mayday_orig.art"), 0666)
	os.WriteFile(filepath.Join(d, "Intro.exe"), []byte("x"), 0666)
	return d
}

func TestAllCombinations(t *testing.T) {
	orig := load(t, "../Mayday_orig.exe")
	d := setup(t, orig)
	gs := groupsOf(ExeSpec)
	n := 0
	for mask := 0; mask < 1<<len(gs); mask++ {
		want := map[string]bool{}
		for i, g := range gs {
			if mask&(1<<i) != 0 {
				want[g.ID] = true
			}
		}
		if !validCombo(want) {
			continue
		}
		n++
		st := Analyze(d, ExeSpec)
		if !st.Known {
			t.Fatalf("mask %d: not recognised before apply", mask)
		}
		if _, err := ApplyFile(d, st, want); err != nil {
			t.Fatal(err)
		}
		st = Analyze(d, ExeSpec)
		if !st.Known {
			t.Fatalf("mask %d: not recognised after apply", mask)
		}
		for _, g := range gs {
			if st.On[g.ID] != want[g.ID] {
				t.Fatalf("mask %d: %s detected %v, wanted %v", mask, g.ID, st.On[g.ID], want[g.ID])
			}
		}
	}
	t.Logf("%d valid combinations checked", n)
	ApplyFile(d, Analyze(d, ExeSpec), map[string]bool{})
	if !bytes.Equal(load(t, filepath.Join(d, "Mayday.exe")), orig) {
		t.Fatal("restore not exact")
	}
	// mayday.art: panel on / off / restore
	art := load(t, "../mayday_orig.art")
	for _, on := range []bool{true, false, true} {
		st := Analyze(d, ArtSpec)
		if _, err := ApplyFile(d, st, map[string]bool{"panel": on}); err != nil {
			t.Fatal(err)
		}
		if st := Analyze(d, ArtSpec); !st.Known || st.On["panel"] != on {
			t.Fatal("art state")
		}
	}
	ApplyFile(d, Analyze(d, ArtSpec), map[string]bool{})
	if !bytes.Equal(load(t, filepath.Join(d, "mayday.art")), art) {
		t.Fatal("art restore not exact")
	}
}

func TestDependencies(t *testing.T) {
	orig := load(t, "../Mayday_orig.exe")
	a := Build(orig, ExeSpec, map[string]bool{"smooth": true})
	b := Build(orig, ExeSpec, map[string]bool{"smooth": true, "fullhd": true, "fps60": true, "frame": true})
	if !bytes.Equal(a, b) {
		t.Fatal("smooth must pull in fullhd, fps60 and frame")
	}
	d := setup(t, a)
	st := Analyze(d, ExeSpec)
	if !st.Known || !st.On["smooth"] || !st.On["fullhd"] || !st.On["fps60"] || !st.On["frame"] || len(st.On) != 4 {
		t.Fatalf("detected %v", st.On)
	}
}

func TestAllEqualsV11(t *testing.T) {
	orig := load(t, "../Mayday_orig.exe")
	all := map[string]bool{}
	for _, g := range Groups {
		all[g.ID] = true
	}
	if !bytes.Equal(Build(orig, ExeSpec, all), load(t, "../Mayday_v16.exe")) {
		t.Fatal("all fixes != v16 (1.0.1)")
	}
	all["hz"] = false
	if !bytes.Equal(Build(orig, ExeSpec, all), load(t, "../Mayday_v14h.exe")) {
		t.Fatal("all but hz != v14h")
	}
	all["smooth"] = false
	if !bytes.Equal(Build(orig, ExeSpec, all), load(t, "../Mayday_v13.exe")) {
		t.Fatal("all but smooth != v13")
	}
	all["fps60"] = false
	if !bytes.Equal(Build(orig, ExeSpec, all), load(t, "../Mayday_v12.exe")) {
		t.Fatal("all but smooth/fps60 != v12")
	}
	// the user's v11 is recognised (7 fixes), v12 with all 8 exe fixes
	d := setup(t, load(t, "../Mayday_v11.exe"))
	st := Analyze(d, ExeSpec)
	if !st.Known || len(st.On) != 7 || st.On["terminal"] {
		t.Fatalf("v11 not recognised: %v %v", st.Known, st.On)
	}
	d = setup(t, load(t, "../Mayday_v14h.exe"))
	st = Analyze(d, ExeSpec)
	if !st.Known || len(st.On) != 10 {
		t.Fatalf("v11 not recognised: %v %v", st.Known, st.On)
	}
	d = setup(t, load(t, "../Mayday_v16.exe"))
	st = Analyze(d, ExeSpec)
	if !st.Known || len(st.On) != 11 || !st.IsComplete() {
		t.Fatalf("v16 not recognised as complete: %v %v", st.Known, st.On)
	}
	// release 1.0 (v15): an older version, original rebuilt without a backup, update gives v16
	d = setup(t, load(t, "../Mayday_v15.exe"))
	st = Analyze(d, ExeSpec)
	if st.Known || st.FromBackup || !bytes.Equal(st.Original, orig) {
		t.Fatalf("1.0 exe: known=%v backup=%v original ok=%v", st.Known, st.FromBackup, bytes.Equal(st.Original, orig))
	}
	if _, err := ApplyFile(d, st, AllFixes()); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(load(t, filepath.Join(d, "Mayday.exe")), load(t, "../Mayday_v16.exe")) {
		t.Fatal("update 1.0 -> 1.0.1 != v16")
	}
	if !bytes.Equal(load(t, filepath.Join(d, BackupDir, "Mayday.exe")), orig) {
		t.Fatal("backup after update is not the original")
	}
	d = setup(t, load(t, "../Mayday_v14h.exe"))
	if st = Analyze(d, ExeSpec); !st.Known || st.IsComplete() {
		t.Fatalf("v14h must be an older (incomplete) patch: %v", st.On)
	}
}

func TestUnknownUsesBackup(t *testing.T) {
	orig := load(t, "../Mayday_orig.exe")
	bad := append([]byte(nil), orig...)
	bad[0x1234] ^= 0xFF
	d := setup(t, bad)
	if st := Analyze(d, ExeSpec); st.Known || st.Original != nil {
		t.Fatal("unknown exe without backup should be refused")
	}
	os.MkdirAll(filepath.Join(d, BackupDir), 0777)
	os.WriteFile(filepath.Join(d, BackupDir, "Mayday.exe"), orig, 0666)
	st := Analyze(d, ExeSpec)
	if st.Known || !st.FromBackup {
		t.Fatal("backup not used")
	}
	if _, err := ApplyFile(d, st, map[string]bool{"fullhd": true}); err != nil {
		t.Fatal(err)
	}
	if st := Analyze(d, ExeSpec); !st.Known || !st.On["fullhd"] || len(st.On) != 1 {
		t.Fatal("rebuild from backup failed")
	}
}

func TestIntro(t *testing.T) {
	d := setup(t, load(t, "../Mayday_orig.exe"))
	if IntroState(d) != "on" {
		t.Fatal()
	}
	if SetIntro(d, false) != nil || IntroState(d) != "off" {
		t.Fatal("off")
	}
	if SetIntro(d, true) != nil || IntroState(d) != "on" {
		t.Fatal("on")
	}
	for in, out := range map[string]string{"": "~ 256COLOR", "~ RUNASADMIN": "~ RUNASADMIN 256COLOR",
		"~ 256COLOR": "~ 256COLOR", "WIN98": "~ WIN98 256COLOR"} {
		if mergeLayers(in) != out {
			t.Fatalf("merge %q -> %q", in, mergeLayers(in))
		}
	}
}

func TestDdraw(t *testing.T) {
	d := setup(t, load(t, "../Mayday_orig.exe"))
	if DdrawState(d) != "none" {
		t.Fatal("none")
	}
	if err := SetDdraw(d, true); err != nil || DdrawState(d) != "on" {
		t.Fatal("install", err)
	}
	if b, _ := os.ReadFile(filepath.Join(d, "ddraw.dll")); len(b) != 411648 {
		t.Fatal("dll size")
	}
	os.Remove(filepath.Join(d, DdrawConfig))
	SetDdraw(d, true) // already installed: only the config tool is added back
	for _, f := range []string{"ddraw.ini", DdrawConfig, "cnc-ddraw-LICENSE.txt", "Shaders/interpolation/catmull-rom-bilinear.glsl"} {
		if _, err := os.Stat(filepath.Join(d, f)); err != nil {
			t.Fatal(f)
		}
	}
	os.WriteFile(filepath.Join(d, "ddraw.ini"), []byte("mine"), 0666)
	if SetDdraw(d, false) != nil || DdrawState(d) != "off" {
		t.Fatal("off")
	}
	if SetDdraw(d, true) != nil || DdrawState(d) != "on" {
		t.Fatal("on again")
	}
	if b, _ := os.ReadFile(filepath.Join(d, "ddraw.ini")); string(b) != "mine" {
		t.Fatal("ini overwritten")
	}
}

func TestCfg(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, CfgName)
	if SetFullHDDefault(d) != nil {
		t.Fatal()
	}
	if b, _ := os.ReadFile(p); string(b) != "[Settings]\r\niXResolution=1024\r\n" {
		t.Fatalf("new: %q", b)
	}
	os.WriteFile(p, []byte("[Settings]\r\niXResolution=640\r\niCDVolume=62414\r\niFXVolume=0\r\niGameSpeed=8321\r\n"), 0666)
	SetFullHDDefault(d)
	if b, _ := os.ReadFile(p); string(b) != "[Settings]\r\niXResolution=1024\r\niCDVolume=62414\r\niFXVolume=0\r\niGameSpeed=8321\r\n" {
		t.Fatalf("update: %q", b)
	}
	os.WriteFile(p, []byte("[Settings]\r\niCDVolume=5\r\n[Other]\r\nx=1\r\n"), 0666)
	SetFullHDDefault(d)
	if b, _ := os.ReadFile(p); string(b) != "[Settings]\r\niCDVolume=5\r\niXResolution=1024\r\n[Other]\r\nx=1\r\n" {
		t.Fatalf("insert: %q", b)
	}
}

func TestCodecExtract(t *testing.T) {
	const p = "../../pv/setup_pv.exe" // the CD's PV_WIN\SETUP_PV.EXE (not part of the sources)
	if _, err := os.Stat(p); err != nil {
		t.Skip("SETUP_PV.EXE not available")
	}
	files, err := ExtractCodec(p)
	if err != nil {
		t.Fatal(err)
	}
	d := t.TempDir()
	if codecFilesPresent(d) {
		t.Fatal("empty dir has codec")
	}
	if err := WriteCodecFiles(d, files); err != nil {
		t.Fatal(err)
	}
	if !codecFilesPresent(d) {
		t.Fatal("codec not written")
	}
	g := t.TempDir()
	os.WriteFile(filepath.Join(g, "install.ini"), []byte("[Main]\r\nCDPfad=f:\r\nHDPFAD=C:\\Gry\\Mayday\r\n"), 0666)
	if cdDrive(g) != "F:" {
		t.Fatalf("cdDrive = %q", cdDrive(g))
	}
}
