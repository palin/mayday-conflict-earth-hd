//go:build windows

package main

// Win32 user interface built only on the Go standard library (syscall).

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	comctl32 = syscall.NewLazyDLL("comctl32.dll")
	advapi32 = syscall.NewLazyDLL("advapi32.dll")

	pRegisterClassExW   = user32.NewProc("RegisterClassExW")
	pCreateWindowExW    = user32.NewProc("CreateWindowExW")
	pDefWindowProcW     = user32.NewProc("DefWindowProcW")
	pGetMessageW        = user32.NewProc("GetMessageW")
	pTranslateMessage   = user32.NewProc("TranslateMessage")
	pDispatchMessageW   = user32.NewProc("DispatchMessageW")
	pIsDialogMessageW   = user32.NewProc("IsDialogMessageW")
	pPostQuitMessage    = user32.NewProc("PostQuitMessage")
	pSendMessageW       = user32.NewProc("SendMessageW")
	pMessageBoxW        = user32.NewProc("MessageBoxW")
	pSetWindowTextW     = user32.NewProc("SetWindowTextW")
	pEnableWindow       = user32.NewProc("EnableWindow")
	pLoadCursorW        = user32.NewProc("LoadCursorW")
	pSetCursor          = user32.NewProc("SetCursor")
	pShowWindow         = user32.NewProc("ShowWindow")
	pUpdateWindow       = user32.NewProc("UpdateWindow")
	pAdjustWindowRectEx = user32.NewProc("AdjustWindowRectEx")
	pGetSystemMetrics   = user32.NewProc("GetSystemMetrics")
	pSetProcessDPIAware = user32.NewProc("SetProcessDPIAware")
	pGetDC              = user32.NewProc("GetDC")
	pReleaseDC          = user32.NewProc("ReleaseDC")
	pDestroyWindow      = user32.NewProc("DestroyWindow")
	pSetWindowPos       = user32.NewProc("SetWindowPos")
	pGetDeviceCaps      = gdi32.NewProc("GetDeviceCaps")
	pSelectObject       = gdi32.NewProc("SelectObject")
	pDrawTextW          = user32.NewProc("DrawTextW")
	pCreateFontW        = gdi32.NewProc("CreateFontW")
	pGetModuleHandleW   = kernel32.NewProc("GetModuleHandleW")
	pCreateActCtxW      = kernel32.NewProc("CreateActCtxW")
	pActivateActCtx     = kernel32.NewProc("ActivateActCtx")
	pSHBrowseForFolderW = shell32.NewProc("SHBrowseForFolderW")
	pSHGetPathFromIDLW  = shell32.NewProc("SHGetPathFromIDListW")
	pCoInitialize       = ole32.NewProc("CoInitialize")
	pCoTaskMemFree      = ole32.NewProc("CoTaskMemFree")
	pShellExecuteExW    = shell32.NewProc("ShellExecuteExW")
	pWaitForSingleObj   = kernel32.NewProc("WaitForSingleObject")
	pGetExitCodeProcess = kernel32.NewProc("GetExitCodeProcess")
	pCloseHandle        = kernel32.NewProc("CloseHandle")
	pInitCommonCtrlsEx  = comctl32.NewProc("InitCommonControlsEx")
	pRegCreateKeyExW    = advapi32.NewProc("RegCreateKeyExW")
	pRegSetValueExW     = advapi32.NewProc("RegSetValueExW")
)

const (
	wsOverlapped     = 0x00000000
	wsCaption        = 0x00C00000
	wsSysMenu        = 0x00080000
	wsMinimizeBox    = 0x00020000
	wsVisible        = 0x10000000
	wsChild          = 0x40000000
	wsTabStop        = 0x00010000
	wsExControlPar   = 0x00010000
	wsExClientEdge   = 0x00000200
	bsPushButton     = 0x0
	bsDefPushButton  = 0x1
	bsAutoCheckbox   = 0x3
	bsGroupBox       = 0x7
	esAutoHScroll    = 0x80
	esReadOnly       = 0x800
	wmDestroy        = 0x0002
	wmClose          = 0x0010
	wmSetFont        = 0x0030
	wmCommand        = 0x0111
	bmGetCheck       = 0x00F0
	bmSetCheck       = 0x00F1
	mbOK             = 0x0
	mbYesNo          = 0x4
	mbIconError      = 0x10
	mbIconWarning    = 0x30
	mbIconInfo       = 0x40
	idYes            = 6
	colorBtnFace     = 15
	logPixelsY       = 90
	hkeyCurrentUser  = 0x80000001
	hkeyLocalMachine = 0x80000002
	keyRead          = 0x20019
	keySetValue      = 0x0002
	keyWow64_32      = 0x0200
	regSZ            = 1
)

type wndClassEx struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

type msgT struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	ptX     int32
	ptY     int32
	private uint32
}

type actCtx struct {
	cbSize                 uint32
	dwFlags                uint32
	lpSource               *uint16
	wProcessorArchitecture uint16
	wLangId                uint16
	lpAssemblyDirectory    *uint16
	lpResourceName         *uint16
	lpApplicationName      *uint16
	hModule                uintptr
}

type browseInfo struct {
	hwndOwner      uintptr
	pidlRoot       uintptr
	pszDisplayName *uint16
	lpszTitle      *uint16
	ulFlags        uint32
	lpfn           uintptr
	lParam         uintptr
	iImage         int32
}

type rect struct{ l, t, r, b int32 }

func u16(s string) *uint16 { p, _ := syscall.UTF16PtrFromString(s); return p }

// ---- registry ----

func setColor256(path string) error {
	var key syscall.Handle
	sub := u16(`Software\Microsoft\Windows NT\CurrentVersion\AppCompatFlags\Layers`)
	r, _, _ := pRegCreateKeyExW.Call(hkeyCurrentUser, uintptr(unsafe.Pointer(sub)), 0, 0, 0,
		keySetValue|keyRead, 0, uintptr(unsafe.Pointer(&key)), 0)
	if r != 0 {
		return fmt.Errorf("rejestr: błąd %d", r)
	}
	defer syscall.RegCloseKey(key)
	cur := regString(key, path)
	val := mergeLayers(cur)
	if val == cur {
		return nil
	}
	data, _ := syscall.UTF16FromString(val)
	r, _, _ = pRegSetValueExW.Call(uintptr(key), uintptr(unsafe.Pointer(u16(path))), 0, regSZ,
		uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)*2))
	if r != 0 {
		return fmt.Errorf("rejestr: błąd %d", r)
	}
	return nil
}

func regString(key syscall.Handle, name string) string {
	var typ, n uint32
	if syscall.RegQueryValueEx(key, u16(name), nil, &typ, nil, &n) != nil || n == 0 {
		return ""
	}
	buf := make([]uint16, n/2+1)
	if syscall.RegQueryValueEx(key, u16(name), nil, &typ, (*byte)(unsafe.Pointer(&buf[0])), &n) != nil {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

func regRead(root uintptr, sub, name string) string {
	var key syscall.Handle
	if syscall.RegOpenKeyEx(syscall.Handle(root), u16(sub), 0, keyRead|keyWow64_32, &key) != nil {
		return ""
	}
	defer syscall.RegCloseKey(key)
	return regString(key, name)
}

const drivers32 = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Drivers32`

func sysWow64() string {
	root := os.Getenv("SystemRoot")
	if root == "" {
		root = `C:\Windows`
	}
	return filepath.Join(root, "SysWOW64")
}

// the video codec works for the game (a 32-bit program) when its files are in SysWOW64
// and it is registered in the 32-bit view of Drivers32
func codecInstalled() bool {
	return regRead(hkeyLocalMachine, drivers32, CodecReg) != "" && codecFilesPresent(sysWow64())
}

// runCodecHelper runs in the elevated copy of the patcher ("--codec <SETUP_PV.EXE>").
func runCodecHelper(setup string) int {
	files, err := ExtractCodec(setup)
	if err != nil {
		return 2
	}
	if err := WriteCodecFiles(sysWow64(), files); err != nil {
		return 3
	}
	var key syscall.Handle
	r, _, _ := pRegCreateKeyExW.Call(hkeyLocalMachine, uintptr(unsafe.Pointer(u16(drivers32))), 0, 0, 0,
		keySetValue|keyWow64_32, 0, uintptr(unsafe.Pointer(&key)), 0)
	if r != 0 {
		return 4
	}
	defer syscall.RegCloseKey(key)
	data, _ := syscall.UTF16FromString("pvez32.drv")
	r, _, _ = pRegSetValueExW.Call(uintptr(key), uintptr(unsafe.Pointer(u16(CodecReg))), 0, regSZ,
		uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)*2))
	if r != 0 {
		return 4
	}
	return 0
}

type shellExecInfo struct {
	cbSize       uint32
	fMask        uint32
	hwnd         uintptr
	lpVerb       *uint16
	lpFile       *uint16
	lpParameters *uint16
	lpDirectory  *uint16
	nShow        int32
	hInstApp     uintptr
	lpIDList     uintptr
	lpClass      *uint16
	hkeyClass    uintptr
	dwHotKey     uint32
	hIcon        uintptr
	hProcess     uintptr
}

// installCodec starts an elevated copy of the patcher (one Windows administrator prompt)
// that extracts the codec from the CD and installs it.
func installCodec(setup string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	sei := shellExecInfo{fMask: 0x40 | 0x100, hwnd: mainWnd, lpVerb: u16("runas"), lpFile: u16(exe),
		lpParameters: u16(`--codec "` + setup + `"`), nShow: 0}
	sei.cbSize = uint32(unsafe.Sizeof(sei))
	if r, _, e := pShellExecuteExW.Call(uintptr(unsafe.Pointer(&sei))); r == 0 {
		if e == syscall.Errno(1223) {
			return errors.New("instalacja kodeka anulowana (brak zgody administratora)")
		}
		return e
	}
	defer pCloseHandle.Call(sei.hProcess)
	pWaitForSingleObj.Call(sei.hProcess, 0xFFFFFFFF)
	var code uint32
	pGetExitCodeProcess.Call(sei.hProcess, uintptr(unsafe.Pointer(&code)))
	switch code {
	case 0:
		return nil
	case 2:
		return errors.New("nie udało się wypakować kodeka z płyty (nieznana wersja SETUP_PV.EXE)")
	case 3:
		return errors.New("nie udało się zapisać plików kodeka w " + sysWow64())
	case 4:
		return errors.New("nie udało się zarejestrować kodeka w rejestrze")
	}
	return fmt.Errorf("błąd %d", code)
}

// game folder from the Wise uninstaller entry ("...\UNWISE.EXE ...\INSTALL.LOG")
func installedDir() string {
	s := regRead(hkeyLocalMachine, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Mayday 1.1p`, "UninstallString")
	i := strings.LastIndex(strings.ToUpper(s), "INSTALL.LOG")
	if i < 0 {
		return ""
	}
	start := strings.LastIndex(s[:i], " ")
	if j := strings.Index(strings.ToUpper(s), "UNWISE.EXE"); j >= 0 && j < i {
		start = j + len("UNWISE.EXE")
	}
	return strings.Trim(strings.TrimSpace(s[start+1:i]), `"\`)
}

// ---- UI ----

var (
	hInst      uintptr
	mainWnd    uintptr
	editDir    uintptr
	statusLbl  uintptr
	statusHead uintptr
	btnApply   uintptr
	btnRestore uintptr
	introChk   uintptr
	font       uintptr
	fontBold   uintptr
	scale      = 1.0
	curDir     string
	installed  bool // this version is fully installed: the main button only saves the intro setting
)

const (
	idBrowse  = 10
	idApply   = 11
	idRestore = 12
	idClose   = 13
	idIntro   = 14
)

func sc(v int) uintptr { return uintptr(int(float64(v)*scale + 0.5)) }

func create(ex uint32, class, text string, style uint32, x, y, w, h int, id uintptr) uintptr {
	hw, _, _ := pCreateWindowExW.Call(uintptr(ex), uintptr(unsafe.Pointer(u16(class))),
		uintptr(unsafe.Pointer(u16(text))), uintptr(style|wsChild|wsVisible),
		sc(x), sc(y), sc(w), sc(h), mainWnd, id, hInst, 0)
	pSendMessageW.Call(hw, wmSetFont, font, 1)
	return hw
}

func setText(hw uintptr, s string) { pSetWindowTextW.Call(hw, uintptr(unsafe.Pointer(u16(s)))) }
func enable(hw uintptr, on bool) {
	v := uintptr(0)
	if on {
		v = 1
	}
	pEnableWindow.Call(hw, v)
}
func isChecked(hw uintptr) bool { r, _, _ := pSendMessageW.Call(hw, bmGetCheck, 0, 0); return r == 1 }
func setCheck(hw uintptr, on bool) {
	v := uintptr(0)
	if on {
		v = 1
	}
	pSendMessageW.Call(hw, bmSetCheck, v, 0)
}

func msgBox(text string, flags uintptr) uintptr {
	r, _, _ := pMessageBoxW.Call(mainWnd, uintptr(unsafe.Pointer(u16(text))),
		uintptr(unsafe.Pointer(u16(AppTitle))), flags)
	return r
}

func waitCursor() {
	c, _, _ := pLoadCursorW.Call(0, 32514)
	pSetCursor.Call(c)
}

// textHeight measures word-wrapped text (in unscaled units) for a given width.
func textHeight(text string, f uintptr, w int) int {
	dc, _, _ := pGetDC.Call(mainWnd)
	old, _, _ := pSelectObject.Call(dc, f)
	r := rect{0, 0, int32(sc(w)), 0}
	pDrawTextW.Call(dc, uintptr(unsafe.Pointer(u16(text))), ^uintptr(0), uintptr(unsafe.Pointer(&r)),
		0x400|0x10|0x800) // DT_CALCRECT | DT_WORDBREAK | DT_NOPREFIX
	pSelectObject.Call(dc, old)
	pReleaseDC.Call(mainWnd, dc)
	return int(float64(r.b)/scale+0.999) + 1
}

// buildUI creates all controls and returns the client height used.
func buildUI() int {
	const W = 600
	title := create(0, "STATIC", AppHead, 0x80, 14, 12, W-28, 22, 0)
	pSendMessageW.Call(title, wmSetFont, fontBold, 1)
	create(0, "STATIC", "Folder gry:", 0, 14, 42, 200, 18, 0)
	editDir = create(wsExClientEdge, "EDIT", "", esAutoHScroll|esReadOnly|wsTabStop, 14, 61, W-28-110, 24, 0)
	create(0, "BUTTON", "Wybierz…", bsPushButton|wsTabStop, W-14-100, 60, 100, 26, idBrowse)
	statusHead = create(0, "STATIC", "", 0x80, 14, 94, W-28, 20, 0)
	pSendMessageW.Call(statusHead, wmSetFont, fontBold, 1)
	statusLbl = create(0, "STATIC", "", 0x80, 14, 114, W-28, 54, 0)

	// changelog: headings in bold, items as wrapped bullet lines
	y := 172
	box := create(0, "BUTTON", "Zawartość patcha "+Version+":", bsGroupBox, 10, y, W-20, 10, 0)
	yy := y + 22
	for _, line := range Changelog {
		if strings.HasPrefix(line, "#") {
			if yy > y+22 {
				yy += 4
			}
			hw := create(0, "STATIC", line[1:], 0x80, 24, yy, W-48, 18, 0)
			pSendMessageW.Call(hw, wmSetFont, fontBold, 1)
			yy += 19
			continue
		}
		create(0, "STATIC", "•", 0x80, 26, yy, 10, 18, 0)
		h := textHeight(line, font, W-62)
		create(0, "STATIC", line, 0x80, 38, yy, W-62, h, 0)
		yy += h + 1
	}
	pSetWindowPos.Call(box, 0, sc(10), sc(y), sc(W-20), sc(yy+8-y), 0x0004)
	y = yy + 16
	create(0, "BUTTON", "Intro", bsGroupBox, 10, y, W-20, 48, 0)
	introChk = create(0, "BUTTON", "Odtwarzaj intro przy starcie (ustawia tryb 256 kolorów dla Intro.exe)",
		bsAutoCheckbox|wsTabStop, 24, y+20, W-48, 22, idIntro)
	y += 48 + 12
	btnApply = create(0, "BUTTON", "Zainstaluj patch", bsDefPushButton|wsTabStop, W-14-90-8-140-8-170, y, 170, 30, idApply)
	btnRestore = create(0, "BUTTON", "Przywróć oryginał", bsPushButton|wsTabStop, W-14-90-8-140, y, 140, 30, idRestore)
	create(0, "BUTTON", "Zamknij", bsPushButton|wsTabStop, W-14-90, y, 90, 30, idClose)
	return y + 30 + 14
}

const winW = 600

func refresh() {
	setText(editDir, curDir)
	head, ok, canRestore, update := "", curDir != "", false, false
	var lines []string
	installed = false
	if curDir == "" {
		head = "Wybierz folder gry"
		lines = append(lines, "Wskaż folder, w którym zainstalowana jest gra (z plikiem Mayday.exe).")
	} else {
		original, complete := true, true
		for _, f := range Files {
			st := Analyze(curDir, f)
			switch {
			case !st.Present && st.Original == nil:
				lines = append(lines, "✖ "+f.Name+": brak pliku – zainstaluj grę ponownie z CD.")
				ok = false
			case st.Known:
				if len(st.On) > 0 {
					original = false
					canRestore = true
				}
				if !st.IsComplete() {
					complete = false
				}
			case st.Original != nil: // unknown content, original in the backup folder
				original, complete, canRestore = false, false, true
			default:
				lines = append(lines, "✖ "+f.Name+": nieznana wersja pliku – patch działa tylko z polską wersją 1.1p. Zainstaluj grę ponownie z CD.")
				ok = false
			}
		}
		ddraw := DdrawState(curDir)
		switch {
		case !ok:
			head = "Nie można zainstalować patcha"
		case original:
			head = "✓ Polska wersja 1.1p – oryginalne pliki."
			lines = append([]string{"Patch nie jest zainstalowany."}, lines...)
		case complete && ddraw == "on":
			head = "✓ Mayday HD " + Version + " jest zainstalowany."
			lines = append([]string{"Oryginalne pliki zapisane w " + BackupDir + "."}, lines...)
			installed = true
		default:
			head = "⚠ Zainstalowana starsza wersja patcha."
			lines = append([]string{"Kliknij „Zaktualizuj do " + Version + "”. Oryginalne pliki są w " + BackupDir + "."}, lines...)
			update = true
		}
	}
	intro := "missing"
	if curDir != "" {
		intro = IntroState(curDir)
	}
	if ok && intro == "missing" {
		lines = append(lines, "Intro.exe: nie znaleziono.")
	}
	codec := codecInstalled()
	if ok && !codec {
		if FindSetupPV(curDir) != "" {
			lines = append(lines, "Kodek filmów (intro, odprawy misji) zostanie zainstalowany z płyty gry.")
		} else {
			lines = append(lines, "⚠ Brak kodeka filmów, a płyta gry nie jest zamontowana – zamontuj ją, aby działały filmy.")
		}
	}
	setText(statusHead, head)
	setText(statusLbl, strings.Join(lines, "\r\n"))
	// intro: off by default until the patch is installed, then the real state
	setCheck(introChk, (installed || update) && intro == "on")
	enable(introChk, ok && intro != "missing")
	switch {
	case installed && !codec:
		setText(btnApply, "Zainstaluj kodek filmów")
	case installed:
		setText(btnApply, "Zapisz ustawienie intro")
	case update:
		setText(btnApply, "Zaktualizuj do "+Version)
	default:
		setText(btnApply, "Zainstaluj patch")
	}
	enable(btnApply, ok)
	enable(btnRestore, ok && canRestore)
}

func plural(n int) string {
	switch {
	case n == 1:
		return "poprawka"
	case n%10 >= 2 && n%10 <= 4 && (n%100 < 12 || n%100 > 14):
		return "poprawki"
	}
	return "poprawek"
}

func doApply(restore bool) {
	if curDir == "" {
		return
	}
	want := AllFixes()
	if restore {
		if msgBox("Przywrócić oryginalne pliki Mayday.exe i mayday.art?\r\n(Patch Mayday HD zostanie usunięty; cnc-ddraw zostaje, bo bez niego gra nie działa na Windows 10/11.)", mbYesNo|mbIconWarning) != idYes {
			return
		}
		want = map[string]bool{}
	}
	wasInstalled := installed
	waitCursor()
	changed := false
	for _, f := range Files {
		st := Analyze(curDir, f)
		ch, err := ApplyFile(curDir, st, want)
		if err != nil {
			hint := ""
			if errors.Is(err, syscall.Errno(32)) || errors.Is(err, syscall.Errno(33)) { // sharing/lock violation
				hint = "\r\n\r\nPlik jest używany – zamknij grę i spróbuj ponownie."
			} else if errors.Is(err, fs.ErrPermission) {
				hint = "\r\n\r\nBrak uprawnień do zapisu w folderze gry. Uruchom patcher jako administrator (prawy przycisk → Uruchom jako administrator)."
			}
			msgBox("Błąd: "+err.Error()+hint, mbOK|mbIconError)
			refresh()
			return
		}
		changed = changed || ch
	}
	var notes []string
	if !restore {
		if changed {
			// only when the game files change: later the player's own choice in the game menu stays
			if err := SetFullHDDefault(curDir); err != nil {
				notes = append(notes, "mayday.cfg: błąd – "+err.Error())
			}
		}
		if err := SetDdraw(curDir, true); err != nil {
			notes = append(notes, "cnc-ddraw: błąd – "+err.Error())
		}
	}
	if !restore && IntroState(curDir) != "missing" {
		on := isChecked(introChk)
		if on != (IntroState(curDir) == "on") || on {
			if err := SetIntro(curDir, on); err != nil {
				notes = append(notes, "Intro: błąd – "+err.Error())
			}
		}
	}
	codecMsg := ""
	if !restore && !codecInstalled() {
		if setup := FindSetupPV(curDir); setup == "" {
			codecMsg = "Kodek filmów: nie znaleziono płyty gry (PV_WIN\\SETUP_PV.EXE). Zamontuj obraz płyty i uruchom patcher ponownie, inaczej filmy (intro, odprawy misji) nie będą działać."
		} else if err := installCodec(setup); err != nil {
			codecMsg = "Kodek filmów: " + err.Error() + ". Uruchom patcher ponownie, aby spróbować jeszcze raz."
		} else {
			codecMsg = "Kodek filmów (Power!Video) zainstalowany z płyty gry."
		}
	}
	refresh()
	var msg string
	switch {
	case restore:
		msg = "Przywrócono oryginalne pliki gry."
	case wasInstalled && codecMsg != "":
		msg = "Zapisano ustawienia."
	case wasInstalled:
		msg = "Zapisano ustawienie intro."
	default:
		msg = "Mayday HD " + Version + " został zainstalowany.\r\n\r\nGra uruchomi się w 1920x1080. Oryginalne pliki są w folderze " + BackupDir + "."
	}
	if !restore {
		if isChecked(introChk) {
			msg += "\r\nIntro: włączone (tryb 256 kolorów)."
		} else {
			msg += "\r\nIntro: wyłączone."
		}
	}
	if codecMsg != "" {
		msg += "\r\n" + codecMsg
	}
	if len(notes) > 0 {
		msg += "\r\n\r\n" + strings.Join(notes, "\r\n")
	}
	if !restore && !wasInstalled {
		msg += "\r\n\r\nJeśli gra mimo to startuje bez dźwięku: uruchom grę, ustaw Mayday na 100 w Mikserze głośności Windows i wyjdź z gry przez menu. Wystarczy raz."
	}
	msgBox(msg, mbOK|mbIconInfo)
}

func browse() {
	buf := make([]uint16, 260)
	bi := browseInfo{hwndOwner: mainWnd, pszDisplayName: &buf[0],
		lpszTitle: u16("Wybierz folder gry Mayday (z plikiem Mayday.exe):"), ulFlags: 0x1 | 0x40}
	pidl, _, _ := pSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl == 0 {
		return
	}
	defer pCoTaskMemFree.Call(pidl)
	path := make([]uint16, 1024)
	if r, _, _ := pSHGetPathFromIDLW.Call(pidl, uintptr(unsafe.Pointer(&path[0]))); r == 0 {
		return
	}
	dir := syscall.UTF16ToString(path)
	if _, err := os.Stat(filepath.Join(dir, ExeSpec.Name)); err != nil {
		if _, err2 := os.Stat(filepath.Join(dir, "Mayday", ExeSpec.Name)); err2 == nil {
			dir = filepath.Join(dir, "Mayday")
		} else {
			msgBox("W tym folderze nie ma pliku Mayday.exe.", mbOK|mbIconWarning)
			return
		}
	}
	curDir = dir
	refresh()
}

func wndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmCommand:
		switch wParam & 0xFFFF {
		case idBrowse:
			browse()
		case idApply:
			doApply(false)
		case idRestore:
			doApply(true)
		case idClose, 2: // 2 = IDCANCEL (Esc)
			pDestroyWindow.Call(hwnd)
		}
		return 0
	case wmClose:
		pDestroyWindow.Call(hwnd)
		return 0
	case wmDestroy:
		pPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

const manifest = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
<dependency><dependentAssembly><assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls"
 version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/></dependentAssembly></dependency>
</assembly>`

func enableVisualStyles() {
	f, err := os.CreateTemp("", "mdhd*.manifest")
	if err != nil {
		return
	}
	f.WriteString(manifest)
	f.Close()
	defer os.Remove(f.Name())
	ac := actCtx{lpSource: u16(f.Name())}
	ac.cbSize = uint32(unsafe.Sizeof(ac))
	h, _, _ := pCreateActCtxW.Call(uintptr(unsafe.Pointer(&ac)))
	if h != ^uintptr(0) {
		var cookie uintptr
		pActivateActCtx.Call(h, uintptr(unsafe.Pointer(&cookie)))
	}
	icc := [2]uint32{8, 0x4000 | 0xFF}
	pInitCommonCtrlsEx.Call(uintptr(unsafe.Pointer(&icc)))
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--codec" {
		os.Exit(runCodecHelper(os.Args[2]))
	}
	runtime.LockOSThread()
	pSetProcessDPIAware.Call()
	enableVisualStyles()
	pCoInitialize.Call(0)
	hInst, _, _ = pGetModuleHandleW.Call(0)

	dc, _, _ := pGetDC.Call(0)
	dpi, _, _ := pGetDeviceCaps.Call(dc, logPixelsY)
	pReleaseDC.Call(0, dc)
	if dpi > 0 {
		scale = float64(dpi) / 96.0
	}
	fh := uintptr(-int(9*float64(dpi)/72.0 + 0.5))
	font, _, _ = pCreateFontW.Call(fh, 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(u16("Segoe UI"))))
	fontBold, _, _ = pCreateFontW.Call(fh, 0, 0, 0, 700, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(u16("Segoe UI"))))

	cls := u16("MaydayHDPatcher")
	cur, _, _ := pLoadCursorW.Call(0, 32512)
	wc := wndClassEx{lpfnWndProc: syscall.NewCallback(wndProc), hInstance: hInst, hCursor: cur,
		hbrBackground: colorBtnFace + 1, lpszClassName: cls}
	wc.cbSize = uint32(unsafe.Sizeof(wc))
	pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	style := uintptr(wsOverlapped | wsCaption | wsSysMenu | wsMinimizeBox)
	mainWnd, _, _ = pCreateWindowExW.Call(wsExControlPar, uintptr(unsafe.Pointer(cls)),
		uintptr(unsafe.Pointer(u16(AppTitle))), style, 0, 0, 100, 100, 0, 0, hInst, 0)
	clientH := buildUI()
	r := rect{0, 0, int32(sc(winW)), int32(sc(clientH))}
	pAdjustWindowRectEx.Call(uintptr(unsafe.Pointer(&r)), style, 0, wsExControlPar)
	w, h := uintptr(r.r-r.l), uintptr(r.b-r.t)
	sw, _, _ := pGetSystemMetrics.Call(0)
	shh, _, _ := pGetSystemMetrics.Call(1)
	pSetWindowPos.Call(mainWnd, 0, (sw-w)/2, (shh-h)/2, w, h, 0x0004) // SWP_NOZORDER

	// starting folder: command line, folder of the patcher, or the installed game
	for _, cand := range []string{argDir(), exeDir(), installedDir()} {
		if cand != "" {
			if _, err := os.Stat(filepath.Join(cand, ExeSpec.Name)); err == nil {
				curDir = cand
				break
			}
		}
	}
	refresh()
	pShowWindow.Call(mainWnd, 1)
	pUpdateWindow.Call(mainWnd)

	var m msgT
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if r == 0 || int32(r) == -1 {
			break
		}
		if d, _, _ := pIsDialogMessageW.Call(mainWnd, uintptr(unsafe.Pointer(&m))); d != 0 {
			continue
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func argDir() string {
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "--") {
		return os.Args[1]
	}
	return ""
}

func exeDir() string {
	p, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(p)
}
