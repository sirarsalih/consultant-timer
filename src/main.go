package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

const (
	WS_OVERLAPPED               = 0x00000000
	WS_CAPTION                  = 0x00C00000
	WS_SYSMENU                  = 0x00080000
	WS_MINIMIZEBOX              = 0x00020000
	WS_VISIBLE                  = 0x10000000
	WS_CHILD                    = 0x40000000
	WS_TABSTOP                  = 0x00010000
	WS_BORDER                   = 0x00800000
	BS_OWNERDRAW                = 0x0000000B
	BS_AUTOCHECKBOX             = 0x00000003
	SS_CENTER                   = 0x00000001
	ES_CENTER                   = 0x00000001
	ES_READONLY                 = 0x00000800
	SW_SHOW                     = 5
	CW_USEDEFAULT       int32   = -2147483648
	WM_DESTROY                  = 0x0002
	WM_COMMAND                  = 0x0111
	WM_TIMER                    = 0x0113
	WM_SETFONT                  = 0x0030
	WM_CLOSE                    = 0x0010
	WM_DRAWITEM                 = 0x002B
	WM_SETICON                  = 0x0080
	BM_GETCHECK                 = 0x00F0
	BM_SETCHECK                 = 0x00F1
	WM_RBUTTONUP                = 0x0205
	WM_LBUTTONDBLCLK            = 0x0203
	WM_APP                      = 0x8000
	WM_TRAYICON                 = WM_APP + 1
	ICON_SMALL                  = 0
	ICON_BIG                    = 1
	COLOR_WINDOW                = 5
	ID_START                    = 1001
	ID_PAUSE                    = 1002
	ID_STOP                     = 1003
	ID_STARTUP                  = 1004
	ID_TRAY_START               = 2001
	ID_TRAY_PAUSE               = 2002
	ID_TRAY_STOP                = 2003
	ID_TRAY_SHOW                = 2004
	ID_TRAY_EXIT                = 2005
	TIMER_ID                    = 1
	DEFAULT_IDLE_MINUTE         = 5
	ODT_BUTTON                  = 4
	ODS_DISABLED                = 0x0004
	DT_CENTER                   = 0x00000001
	DT_VCENTER                  = 0x00000004
	DT_SINGLELINE               = 0x00000020
	TRANSPARENT                 = 1
	BST_UNCHECKED               = 0
	BST_CHECKED                 = 1
	REG_SZ                      = 1
	KEY_SET_VALUE               = 0x0002
	HKEY_CURRENT_USER   uintptr = 0x80000001

	NIM_ADD         = 0x00000000
	NIM_MODIFY      = 0x00000001
	NIM_DELETE      = 0x00000002
	NIF_MESSAGE     = 0x00000001
	NIF_ICON        = 0x00000002
	NIF_TIP         = 0x00000004
	MF_STRING       = 0x00000000
	MF_GRAYED       = 0x00000001
	MF_SEPARATOR    = 0x00000800
	TPM_RIGHTBUTTON = 0x0002
)

type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       syscall.Handle
}

type POINT struct{ X, Y int32 }
type RECT struct{ Left, Top, Right, Bottom int32 }
type MSG struct {
	HWnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
	Private uint32
}
type LASTINPUTINFO struct {
	CbSize uint32
	DwTime uint32
}
type DRAWITEMSTRUCT struct {
	CtlType    uint32
	CtlID      uint32
	ItemID     uint32
	ItemAction uint32
	ItemState  uint32
	HwndItem   syscall.Handle
	HDC        syscall.Handle
	RcItem     RECT
	ItemData   uintptr
}
type NOTIFYICONDATA struct {
	CbSize            uint32
	HWnd              syscall.Handle
	UID               uint32
	UFlags            uint32
	UCallbackMessage  uint32
	HIcon             syscall.Handle
	SzTip             [128]uint16
	DwState           uint32
	DwStateMask       uint32
	SzInfo            [256]uint16
	UTimeoutOrVersion uint32
	SzInfoTitle       [64]uint16
	DwInfoFlags       uint32
	GuidItem          [16]byte
	HBalloonIcon      syscall.Handle
}
type Store struct {
	Days             map[string]float64 `json:"days"`
	StartWithWindows *bool              `json:"start_with_windows,omitempty"`
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	advapi32 = syscall.NewLazyDLL("advapi32.dll")

	procRegisterClassExW    = user32.NewProc("RegisterClassExW")
	procCreateWindowExW     = user32.NewProc("CreateWindowExW")
	procDefWindowProcW      = user32.NewProc("DefWindowProcW")
	procShowWindow          = user32.NewProc("ShowWindow")
	procUpdateWindow        = user32.NewProc("UpdateWindow")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procPostQuitMessage     = user32.NewProc("PostQuitMessage")
	procSetWindowTextW      = user32.NewProc("SetWindowTextW")
	procSendMessageW        = user32.NewProc("SendMessageW")
	procEnableWindow        = user32.NewProc("EnableWindow")
	procSetTimer            = user32.NewProc("SetTimer")
	procKillTimer           = user32.NewProc("KillTimer")
	procGetLastInputInfo    = user32.NewProc("GetLastInputInfo")
	procLoadCursorW         = user32.NewProc("LoadCursorW")
	procLoadIconW           = user32.NewProc("LoadIconW")
	procFillRect            = user32.NewProc("FillRect")
	procDrawTextW           = user32.NewProc("DrawTextW")
	procCreatePopupMenu     = user32.NewProc("CreatePopupMenu")
	procAppendMenuW         = user32.NewProc("AppendMenuW")
	procTrackPopupMenu      = user32.NewProc("TrackPopupMenu")
	procDestroyMenu         = user32.NewProc("DestroyMenu")
	procGetCursorPos        = user32.NewProc("GetCursorPos")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procIsWindowVisible     = user32.NewProc("IsWindowVisible")
	procInvalidateRect      = user32.NewProc("InvalidateRect")
	procSetClassLongPtrW    = user32.NewProc("SetClassLongPtrW")

	procSetTextColor           = gdi32.NewProc("SetTextColor")
	procSetBkMode              = gdi32.NewProc("SetBkMode")
	procCreateSolidBrush       = gdi32.NewProc("CreateSolidBrush")
	procDeleteObject           = gdi32.NewProc("DeleteObject")
	procSelectObject           = gdi32.NewProc("SelectObject")
	procCreateFontW            = gdi32.NewProc("CreateFontW")
	procCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	procCreateBitmap           = gdi32.NewProc("CreateBitmap")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procEllipse                = gdi32.NewProc("Ellipse")
	procMoveToEx               = gdi32.NewProc("MoveToEx")
	procLineTo                 = gdi32.NewProc("LineTo")
	procCreatePen              = gdi32.NewProc("CreatePen")
	procCreateIconIndirect     = user32.NewProc("CreateIconIndirect")
	procDestroyIcon            = user32.NewProc("DestroyIcon")

	procGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	procGetTickCount64   = kernel32.NewProc("GetTickCount64")
	procShellNotifyIconW = shell32.NewProc("Shell_NotifyIconW")

	procRegOpenKeyExW   = advapi32.NewProc("RegOpenKeyExW")
	procRegSetValueExW  = advapi32.NewProc("RegSetValueExW")
	procRegDeleteValueW = advapi32.NewProc("RegDeleteValueW")
	procRegCloseKey     = advapi32.NewProc("RegCloseKey")

	hwndMain, hwndClock, hwndDecimal, hwndStatus syscall.Handle
	hwndStart, hwndStop, hwndStartup             syscall.Handle
	normalFont, bigFont                          syscall.Handle
	appIcon                                      syscall.Handle
	appIconOwned                                 bool
	trayData                                     NOTIFYICONDATA

	store          = Store{Days: map[string]float64{}}
	currentDate    string
	currentSeconds float64
	running        bool
	lastTick       time.Time
	sessionAdded   float64
	autoIdlePaused bool
	dataPath       string
)

type ICONINFO struct {
	FIcon    bool
	XHotspot uint32
	YHotspot uint32
	HbmMask  syscall.Handle
	HbmColor syscall.Handle
}

func utf16(s string) *uint16   { p, _ := syscall.UTF16PtrFromString(s); return p }
func loword(v uintptr) uint16  { return uint16(v & 0xffff) }
func rgb(r, g, b byte) uintptr { return uintptr(uint32(r) | uint32(g)<<8 | uint32(b)<<16) }

func createWindow(className, text string, style uint32, x, y, w, h int32, parent syscall.Handle, menu uintptr, instance syscall.Handle) syscall.Handle {
	r, _, _ := procCreateWindowExW.Call(0,
		uintptr(unsafe.Pointer(utf16(className))), uintptr(unsafe.Pointer(utf16(text))), uintptr(style),
		uintptr(x), uintptr(y), uintptr(w), uintptr(h), uintptr(parent), menu, uintptr(instance), 0)
	return syscall.Handle(r)
}
func setText(h syscall.Handle, s string) {
	procSetWindowTextW.Call(uintptr(h), uintptr(unsafe.Pointer(utf16(s))))
}
func setFont(h, f syscall.Handle) { procSendMessageW.Call(uintptr(h), WM_SETFONT, uintptr(f), 1) }
func enable(h syscall.Handle, yes bool) {
	var v uintptr
	if yes {
		v = 1
	}
	procEnableWindow.Call(uintptr(h), v)
}
func todayKey(t time.Time) string { return t.Format("2006-01-02") }

func initDataPath() {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		if d, err := os.UserConfigDir(); err == nil {
			base = d
		} else {
			base = "."
		}
	}
	dir := filepath.Join(base, "ConsultantTimer")
	_ = os.MkdirAll(dir, 0755)
	dataPath = filepath.Join(dir, "history.json")
}
func setStartWithWindows(enabled bool) bool {
	const runKey = `Software\Microsoft\Windows\CurrentVersion\Run`
	var key syscall.Handle
	r, _, _ := procRegOpenKeyExW.Call(HKEY_CURRENT_USER, uintptr(unsafe.Pointer(utf16(runKey))), 0, KEY_SET_VALUE, uintptr(unsafe.Pointer(&key)))
	if r != 0 {
		return false
	}
	defer procRegCloseKey.Call(uintptr(key))

	name := utf16("ConsultantTimer")
	if !enabled {
		r, _, _ = procRegDeleteValueW.Call(uintptr(key), uintptr(unsafe.Pointer(name)))
		// ERROR_FILE_NOT_FOUND (2) is fine: desired state is already achieved.
		return r == 0 || r == 2
	}

	exe, err := os.Executable()
	if err != nil {
		return false
	}
	exe, _ = filepath.Abs(exe)
	cmd := `"` + exe + `"`
	u := syscall.StringToUTF16(cmd)
	r, _, _ = procRegSetValueExW.Call(uintptr(key), uintptr(unsafe.Pointer(name)), 0, REG_SZ, uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)*2))
	return r == 0
}

func applyStartupPreference() {
	if store.StartWithWindows == nil {
		v := true
		store.StartWithWindows = &v
		saveStore()
	}
	setStartWithWindows(*store.StartWithWindows)
}

func setStartupCheckbox() {
	if hwndStartup == 0 || store.StartWithWindows == nil {
		return
	}
	state := uintptr(BST_UNCHECKED)
	if *store.StartWithWindows {
		state = BST_CHECKED
	}
	procSendMessageW.Call(uintptr(hwndStartup), BM_SETCHECK, state, 0)
}

func toggleStartWithWindows() {
	if hwndStartup == 0 {
		return
	}
	state, _, _ := procSendMessageW.Call(uintptr(hwndStartup), BM_GETCHECK, 0, 0)
	enabled := state == BST_CHECKED
	if setStartWithWindows(enabled) {
		store.StartWithWindows = &enabled
		saveStore()
		if enabled {
			setText(hwndStatus, "Paused - starts with Windows")
		} else {
			setText(hwndStatus, "Paused - Windows startup disabled")
		}
	} else {
		// Restore checkbox to the last saved preference on registry failure.
		setStartupCheckbox()
		setText(hwndStatus, "Could not change Windows startup setting")
	}
}

func loadStore() {
	initDataPath()
	if b, err := os.ReadFile(dataPath); err == nil {
		_ = json.Unmarshal(b, &store)
	}
	if store.Days == nil {
		store.Days = map[string]float64{}
	}
	currentDate = todayKey(time.Now())
	currentSeconds = store.Days[currentDate]
}
func saveStore() {
	if currentDate != "" {
		store.Days[currentDate] = math.Max(0, currentSeconds)
	}
	b, err := json.MarshalIndent(store, "", "  ")
	if err == nil {
		_ = os.WriteFile(dataPath, b, 0644)
	}
}
func formatHHMM(seconds float64) string {
	totalMinutes := int(math.Floor(math.Max(0, seconds) / 60.0))
	h := totalMinutes / 60
	m := totalMinutes % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}
func updateDisplay() {
	setText(hwndDecimal, fmt.Sprintf("%.2f hours", currentSeconds/3600.0))
	setText(hwndClock, formatHHMM(currentSeconds))
	updateTrayTip()
}
func idleSeconds() float64 {
	lii := LASTINPUTINFO{CbSize: uint32(unsafe.Sizeof(LASTINPUTINFO{}))}
	ok, _, _ := procGetLastInputInfo.Call(uintptr(unsafe.Pointer(&lii)))
	if ok == 0 {
		return 0
	}
	tick, _, _ := procGetTickCount64.Call()
	elapsedMs := uint32(tick) - lii.DwTime
	return float64(elapsedMs) / 1000.0
}
func ensureDate(now time.Time) {
	k := todayKey(now)
	if k == currentDate {
		return
	}
	saveStore()
	currentDate = k
	currentSeconds = store.Days[k]
	sessionAdded = 0
	lastTick = now
	updateDisplay()
}
func updateButtons() {
	// START and PAUSE share one owner-drawn button. It stays enabled and
	// changes text/color based on the current running state.
	enable(hwndStart, true)
	enable(hwndStop, running || currentSeconds > 0)
	if hwndStart != 0 {
		procInvalidateRect.Call(uintptr(hwndStart), 0, 1)
	}
}
func startTimer() {
	if running {
		return
	}
	ensureDate(time.Now())
	running = true
	autoIdlePaused = false
	lastTick = time.Now()
	sessionAdded = 0
	setText(hwndStatus, "Working")
	updateButtons()
	updateTrayTip()
}
func pauseTimer(status string) {
	if !running {
		return
	}
	tick(time.Now(), false)
	running = false
	autoIdlePaused = false
	saveStore()
	setText(hwndStatus, status)
	updateButtons()
	updateTrayTip()
}
func stopAndClear() {
	if running {
		tick(time.Now(), false)
	}
	running = false
	autoIdlePaused = false
	currentSeconds = 0
	sessionAdded = 0
	store.Days[currentDate] = 0
	saveStore()
	setText(hwndStatus, "Stopped - today cleared")
	updateDisplay()
	updateButtons()
	updateTrayTip()
}
func tick(now time.Time, checkIdle bool) {
	ensureDate(now)

	// If the timer was paused automatically due to inactivity, resume it
	// as soon as Windows reports fresh keyboard/mouse input. Manual pauses
	// never enter this state and therefore never auto-resume.
	if !running {
		if checkIdle && autoIdlePaused && idleSeconds() < 2 {
			running = true
			autoIdlePaused = false
			lastTick = now
			sessionAdded = 0
			setText(hwndStatus, "Working - resumed after idle")
			updateButtons()
			updateDisplay()
		}
		return
	}

	dt := now.Sub(lastTick).Seconds()
	if dt < 0 {
		dt = 0
	}
	if dt > 5 {
		dt = 1
	}
	currentSeconds += dt
	sessionAdded += dt
	lastTick = now

	if checkIdle {
		idle := idleSeconds()
		if idle >= DEFAULT_IDLE_MINUTE*60 {
			remove := math.Min(idle, sessionAdded)
			currentSeconds = math.Max(0, currentSeconds-remove)
			sessionAdded = math.Max(0, sessionAdded-remove)
			running = false
			autoIdlePaused = true
			saveStore()
			setText(hwndStatus, fmt.Sprintf("Paused automatically - idle %.0f min", math.Floor(idle/60)))
			updateButtons()
		}
	}
	updateDisplay()
}
func drawButton(dis *DRAWITEMSTRUCT) {
	if dis == nil || dis.CtlType != ODT_BUTTON {
		return
	}
	var fill uintptr
	var textColor uintptr = rgb(255, 255, 255)
	text := ""
	if dis.CtlID == ID_START {
		if running {
			fill = rgb(245, 196, 48)
			text = "PAUSE"
			textColor = rgb(30, 30, 30)
		} else {
			fill = rgb(34, 139, 34)
			text = "START"
		}
	} else if dis.CtlID == ID_STOP {
		fill = rgb(196, 48, 43)
		text = "STOP"
	} else {
		return
	}
	if dis.ItemState&ODS_DISABLED != 0 {
		fill = rgb(150, 150, 150)
		textColor = rgb(245, 245, 245)
	}
	brush, _, _ := procCreateSolidBrush.Call(fill)
	procFillRect.Call(uintptr(dis.HDC), uintptr(unsafe.Pointer(&dis.RcItem)), brush)
	procDeleteObject.Call(brush)
	procSetBkMode.Call(uintptr(dis.HDC), TRANSPARENT)
	procSetTextColor.Call(uintptr(dis.HDC), textColor)
	old, _, _ := procSelectObject.Call(uintptr(dis.HDC), uintptr(normalFont))
	p := utf16(text)
	procDrawTextW.Call(uintptr(dis.HDC), uintptr(unsafe.Pointer(p)), uintptr(len(text)), uintptr(unsafe.Pointer(&dis.RcItem)), DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	if old != 0 {
		procSelectObject.Call(uintptr(dis.HDC), old)
	}
}

func createClockIcon() syscall.Handle {
	const size = 32
	screenDC, _, _ := user32.NewProc("GetDC").Call(0)
	if screenDC == 0 {
		return 0
	}
	defer user32.NewProc("ReleaseDC").Call(0, screenDC)

	memDC, _, _ := procCreateCompatibleDC.Call(screenDC)
	if memDC == 0 {
		return 0
	}
	defer gdi32.NewProc("DeleteDC").Call(memDC)

	colorBmp, _, _ := procCreateCompatibleBitmap.Call(screenDC, size, size)
	if colorBmp == 0 {
		return 0
	}
	maskBmp, _, _ := procCreateBitmap.Call(size, size, 1, 1, 0)
	if maskBmp == 0 {
		procDeleteObject.Call(colorBmp)
		return 0
	}

	oldBmp, _, _ := procSelectObject.Call(memDC, colorBmp)
	white, _, _ := procCreateSolidBrush.Call(rgb(255, 255, 255))
	rc := RECT{0, 0, size, size}
	procFillRect.Call(memDC, uintptr(unsafe.Pointer(&rc)), white)
	procDeleteObject.Call(white)

	pen, _, _ := procCreatePen.Call(0, 3, rgb(40, 40, 40))
	oldPen, _, _ := procSelectObject.Call(memDC, pen)
	face, _, _ := procCreateSolidBrush.Call(rgb(245, 245, 245))
	oldBrush, _, _ := procSelectObject.Call(memDC, face)
	procEllipse.Call(memDC, 3, 3, 29, 29)
	procSelectObject.Call(memDC, oldBrush)
	procDeleteObject.Call(face)

	procMoveToEx.Call(memDC, 16, 16, 0)
	procLineTo.Call(memDC, 16, 9)
	procMoveToEx.Call(memDC, 16, 16, 0)
	procLineTo.Call(memDC, 22, 19)
	procSelectObject.Call(memDC, oldPen)
	procDeleteObject.Call(pen)
	procSelectObject.Call(memDC, oldBmp)

	info := ICONINFO{FIcon: true, HbmMask: syscall.Handle(maskBmp), HbmColor: syscall.Handle(colorBmp)}
	icon, _, _ := procCreateIconIndirect.Call(uintptr(unsafe.Pointer(&info)))
	procDeleteObject.Call(maskBmp)
	procDeleteObject.Call(colorBmp)
	return syscall.Handle(icon)
}

func copyUTF16(dst []uint16, s string) {
	u := syscall.StringToUTF16(s)
	if len(u) > len(dst) {
		u = u[:len(dst)]
	}
	copy(dst, u)
}

func addTrayIcon() {
	if hwndMain == 0 || appIcon == 0 {
		return
	}
	trayData = NOTIFYICONDATA{}
	trayData.CbSize = uint32(unsafe.Sizeof(trayData))
	trayData.HWnd = hwndMain
	trayData.UID = 1
	trayData.UFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
	trayData.UCallbackMessage = WM_TRAYICON
	trayData.HIcon = appIcon
	copyUTF16(trayData.SzTip[:], "Consultant Timer")
	procShellNotifyIconW.Call(NIM_ADD, uintptr(unsafe.Pointer(&trayData)))
	updateTrayTip()
}
func removeTrayIcon() {
	if trayData.HWnd != 0 {
		procShellNotifyIconW.Call(NIM_DELETE, uintptr(unsafe.Pointer(&trayData)))
	}
}
func updateTrayTip() {
	if trayData.HWnd == 0 {
		return
	}
	status := "Paused"
	if running {
		status = "Working"
	} else if currentSeconds <= 0 {
		status = "Stopped"
	}
	tip := fmt.Sprintf("Consultant Timer - %s - %.2f h (%s)", status, currentSeconds/3600.0, formatHHMM(currentSeconds))
	trayData.UFlags = NIF_TIP
	copyUTF16(trayData.SzTip[:], tip)
	procShellNotifyIconW.Call(NIM_MODIFY, uintptr(unsafe.Pointer(&trayData)))
}
func showTrayMenu() {
	menu, _, _ := procCreatePopupMenu.Call()
	if menu == 0 {
		return
	}
	defer procDestroyMenu.Call(menu)

	startFlags := uintptr(MF_STRING)
	pauseFlags := uintptr(MF_STRING)
	stopFlags := uintptr(MF_STRING)
	if running {
		startFlags |= MF_GRAYED
	} else {
		pauseFlags |= MF_GRAYED
	}
	if !running && currentSeconds <= 0 {
		stopFlags |= MF_GRAYED
	}
	procAppendMenuW.Call(menu, startFlags, ID_TRAY_START, uintptr(unsafe.Pointer(utf16("Start"))))
	procAppendMenuW.Call(menu, pauseFlags, ID_TRAY_PAUSE, uintptr(unsafe.Pointer(utf16("Pause"))))
	procAppendMenuW.Call(menu, stopFlags, ID_TRAY_STOP, uintptr(unsafe.Pointer(utf16("Stop / Clear today"))))
	procAppendMenuW.Call(menu, MF_SEPARATOR, 0, 0)
	procAppendMenuW.Call(menu, MF_STRING, ID_TRAY_SHOW, uintptr(unsafe.Pointer(utf16("Show"))))
	procAppendMenuW.Call(menu, MF_SEPARATOR, 0, 0)
	procAppendMenuW.Call(menu, MF_STRING, ID_TRAY_EXIT, uintptr(unsafe.Pointer(utf16("Exit"))))

	var pt POINT
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	procSetForegroundWindow.Call(uintptr(hwndMain))
	procTrackPopupMenu.Call(menu, TPM_RIGHTBUTTON, uintptr(pt.X), uintptr(pt.Y), 0, uintptr(hwndMain), 0)
}

func wndProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch loword(wParam) {
		case ID_START:
			if running {
				pauseTimer("Paused")
			} else {
				startTimer()
			}
			return 0
		case ID_TRAY_START:
			startTimer()
			return 0
		case ID_PAUSE, ID_TRAY_PAUSE:
			pauseTimer("Paused")
			return 0
		case ID_STOP, ID_TRAY_STOP:
			stopAndClear()
			return 0
		case ID_STARTUP:
			toggleStartWithWindows()
			return 0
		case ID_TRAY_SHOW:
			procShowWindow.Call(uintptr(hwndMain), SW_SHOW)
			procSetForegroundWindow.Call(uintptr(hwndMain))
			return 0
		case ID_TRAY_EXIT:
			procSendMessageW.Call(uintptr(hwndMain), WM_CLOSE, 0, 0)
			return 0
		}
	case WM_TIMER:
		if wParam == TIMER_ID {
			tick(time.Now(), true)
			return 0
		}
	case WM_DRAWITEM:
		dis := (*DRAWITEMSTRUCT)(unsafe.Pointer(lParam))
		drawButton(dis)
		return 1
	case WM_TRAYICON:
		switch uint32(lParam) {
		case WM_RBUTTONUP:
			showTrayMenu()
			return 0
		case WM_LBUTTONDBLCLK:
			procShowWindow.Call(uintptr(hwndMain), SW_SHOW)
			procSetForegroundWindow.Call(uintptr(hwndMain))
			return 0
		}
	case WM_CLOSE:
		if running {
			tick(time.Now(), false)
		}
		saveStore()
		removeTrayIcon()
	case WM_DESTROY:
		procKillTimer.Call(uintptr(hwnd), TIMER_ID)
		removeTrayIcon()
		if appIcon != 0 && appIconOwned {
			procDestroyIcon.Call(uintptr(appIcon))
		}
		appIcon = 0
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return r
}
func createFont(height, weight int32) syscall.Handle {
	r, _, _ := procCreateFontW.Call(uintptr(height), 0, 0, 0, uintptr(weight), 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(utf16("Segoe UI"))))
	return syscall.Handle(r)
}
func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	loadStore()
	applyStartupPreference()
	hInst, _, _ := procGetModuleHandleW.Call(0)
	instance := syscall.Handle(hInst)
	iconRes, _, _ := procLoadIconW.Call(uintptr(instance), 1) // embedded RT_GROUP_ICON #1
	appIcon = syscall.Handle(iconRes)
	if appIcon == 0 {
		appIcon = createClockIcon()
		appIconOwned = appIcon != 0
	}
	className := utf16("ConsultantTimerWindowV7")
	cursor, _, _ := procLoadCursorW.Call(0, 32512)
	wc := WNDCLASSEX{CbSize: uint32(unsafe.Sizeof(WNDCLASSEX{})), LpfnWndProc: syscall.NewCallback(wndProc), HInstance: instance, HIcon: appIcon, HCursor: syscall.Handle(cursor), HbrBackground: syscall.Handle(COLOR_WINDOW + 1), LpszClassName: className, HIconSm: appIcon}
	if r, _, _ := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		return
	}

	style := uint32(WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_MINIMIZEBOX | WS_VISIBLE)
	hwndMain = createWindow("ConsultantTimerWindowV7", "Consultant Timer", style, CW_USEDEFAULT, CW_USEDEFAULT, 560, 430, 0, 0, instance)
	if hwndMain == 0 {
		return
	}
	if appIcon != 0 {
		procSendMessageW.Call(uintptr(hwndMain), WM_SETICON, ICON_BIG, uintptr(appIcon))
		procSendMessageW.Call(uintptr(hwndMain), WM_SETICON, ICON_SMALL, uintptr(appIcon))
	}

	normalFont = createFont(20, 500)
	bigFont = createFont(52, 600)

	lblToday := createWindow("STATIC", "TODAY", WS_CHILD|WS_VISIBLE|SS_CENTER, 30, 20, 490, 28, hwndMain, 0, instance)
	// Decimal hours first, then hh:mm underneath.
	hwndDecimal = createWindow("EDIT", "", WS_CHILD|WS_VISIBLE|WS_BORDER|ES_CENTER|ES_READONLY, 125, 57, 300, 62, hwndMain, 0, instance)
	hwndClock = createWindow("EDIT", "", WS_CHILD|WS_VISIBLE|WS_BORDER|ES_CENTER|ES_READONLY, 175, 128, 200, 38, hwndMain, 0, instance)
	hwndStatus = createWindow("STATIC", "Stopped", WS_CHILD|WS_VISIBLE|SS_CENTER, 30, 177, 490, 28, hwndMain, 0, instance)
	hwndStart = createWindow("BUTTON", "START", WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_OWNERDRAW, 110, 215, 160, 52, hwndMain, ID_START, instance)
	hwndStop = createWindow("BUTTON", "STOP", WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_OWNERDRAW, 290, 215, 160, 52, hwndMain, ID_STOP, instance)
	hwndStartup = createWindow("BUTTON", "Start automatically when Windows starts", WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_AUTOCHECKBOX, 115, 275, 330, 32, hwndMain, ID_STARTUP, instance)
	lblIdle := createWindow("STATIC", "Idle detection: 5 minutes  |  Data stored locally", WS_CHILD|WS_VISIBLE|SS_CENTER, 30, 323, 490, 24, hwndMain, 0, instance)

	for _, h := range []syscall.Handle{lblToday, hwndClock, hwndStatus, hwndStart, hwndStop, hwndStartup, lblIdle} {
		setFont(h, normalFont)
	}
	setFont(hwndDecimal, bigFont)
	setStartupCheckbox()

	// Use a hand pointer over all action buttons.
	hand, _, _ := procLoadCursorW.Call(0, 32649) // IDC_HAND
	if hand != 0 {
		procSetClassLongPtrW.Call(uintptr(hwndStart), ^uintptr(11), hand)
	}
	updateButtons()
	updateDisplay()
	if currentSeconds > 0 {
		setText(hwndStatus, "Paused")
	}

	addTrayIcon()
	procSetTimer.Call(uintptr(hwndMain), TIMER_ID, 1000, 0)
	procShowWindow.Call(uintptr(hwndMain), SW_SHOW)
	procUpdateWindow.Call(uintptr(hwndMain))

	var msg MSG
	for {
		r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}
