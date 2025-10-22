//go:build windows

package uiwin

import (
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32           = windows.NewLazySystemDLL("user32.dll")
	procSetCursorPos = user32.NewProc("SetCursorPos")
	procGetCursorPos = user32.NewProc("GetCursorPos")
	procMouseEvent   = user32.NewProc("mouse_event")
	procKeybdEvent   = user32.NewProc("keybd_event")
	procVkKeyScanW   = user32.NewProc("VkKeyScanW")
	procSendInput    = user32.NewProc("SendInput")
)

type point struct{ X, Y int32 }

type UIWin struct {
	DelayMS int
}

func NewUIWin(delayMS int) *UIWin {
	return &UIWin{DelayMS: delayMS}
}

func (u *UIWin) sleep() {
	if u.DelayMS > 0 {
		time.Sleep(time.Duration(u.DelayMS) * time.Millisecond)
	}
}

func (u *UIWin) getCursorPos() (int, int) {
	var p point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	return int(p.X), int(p.Y)
}

func (u *UIWin) MoveSmooth(x, y int, durSec float64) {
	cx, cy := u.getCursorPos()
	steps := int(40 * durSec)
	if steps < 1 {
		procSetCursorPos.Call(uintptr(int32(x)), uintptr(int32(y)))
		return
	}
	dx, dy := x-cx, y-cy
	if steps < 1 {
		procSetCursorPos.Call(uintptr(int32(x)), uintptr(int32(y)))
		return
	}
	stepDuration := time.Duration(float64(time.Second) * durSec / float64(steps))
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		e := 1 - (1-t)*(1-t) // ease-out
		nx := cx + int(float64(dx)*e)
		ny := cy + int(float64(dy)*e)
		procSetCursorPos.Call(uintptr(int32(nx)), uintptr(int32(ny)))
		time.Sleep(stepDuration)
	}
	u.sleep()
}

const (
	MOUSEEVENTF_LEFTDOWN = 0x0002
	MOUSEEVENTF_LEFTUP   = 0x0004
	MOUSEEVENTF_WHEEL    = 0x0800
	INPUT_MOUSE          = 0
	WHEEL_DELTA          = 120
	WM_MOUSEWHEEL        = 0x020A
)

func (u *UIWin) ClickLeft() {
	procMouseEvent.Call(MOUSEEVENTF_LEFTDOWN, 0, 0, 0, 0)
	procMouseEvent.Call(MOUSEEVENTF_LEFTUP, 0, 0, 0, 0)
	u.sleep()
}

func (u *UIWin) DoubleClick() {
	u.ClickLeft()
	time.Sleep(60 * time.Millisecond)
	u.ClickLeft()
	u.sleep()
}

type MOUSEINPUT struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type INPUT struct {
	Type uint32
	_    uint32
	Mi   MOUSEINPUT
}

func (u *UIWin) ScrollNotches(notches int) {
	if notches == 0 {
		return
	}

	delta := int32(notches) * WHEEL_DELTA
	in := INPUT{
		Type: INPUT_MOUSE,
		Mi: MOUSEINPUT{
			Dx:        0,
			Dy:        0,
			MouseData: uint32(delta),
			DwFlags:   MOUSEEVENTF_WHEEL,
			Time:      0,
		},
	}
	if ret, _, _ := procSendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&in)),
		unsafe.Sizeof(in),
	); ret == 1 {
		u.sleep()
		return
	}
}

// ===== Keyboard =====
const (
	KEYEVENTF_KEYUP = 0x0002

	VK_SHIFT   = 0x10
	VK_CONTROL = 0x11
	VK_MENU    = 0x12 // ALT
	VK_LWIN    = 0x5B
	VK_RETURN  = 0x0D
	VK_LEFT    = 0x25
	VK_RIGHT   = 0x27
	VK_A       = 0x41
	VK_SPACE   = 0x20
)

func (u *UIWin) keyDown(vk byte) { procKeybdEvent.Call(uintptr(vk), 0, 0, 0) }
func (u *UIWin) keyUp(vk byte)   { procKeybdEvent.Call(uintptr(vk), 0, KEYEVENTF_KEYUP, 0) }

func (u *UIWin) Press(vk byte) {
	u.keyDown(vk)
	time.Sleep(8 * time.Millisecond)
	u.keyUp(vk)
	u.sleep()
}

func (u *UIWin) PressEnter() { u.Press(VK_RETURN) }

func (u *UIWin) CtrlA() {
	u.keyDown(VK_CONTROL)
	u.Press(VK_A)
	u.keyUp(VK_CONTROL)
	u.sleep()
}

func (u *UIWin) vkForRune(r rune) (vk byte, needShift bool, ok bool) {
	uu := utf16.Encode([]rune{r})[0]
	ret, _, _ := procVkKeyScanW.Call(uintptr(uu))
	v := uint16(ret)
	if v == 0xFFFF {
		return 0, false, false
	}
	vk = byte(v & 0x00FF)
	shiftState := byte((v & 0xFF00) >> 8)
	needShift = (shiftState & 0x01) != 0
	return vk, needShift, true
}

func (u *UIWin) TypeString(s string, perCharMS int) {
	delay := time.Duration(perCharMS) * time.Millisecond
	for _, r := range s {
		switch r {
		case '\n':
			u.PressEnter()
		case ' ':
			u.Press(VK_SPACE)
		default:
			if vk, needShift, ok := u.vkForRune(r); ok {
				if needShift {
					u.keyDown(VK_SHIFT)
					u.Press(vk)
					u.keyUp(VK_SHIFT)
				} else {
					u.Press(vk)
				}
			}
		}
		time.Sleep(delay)
	}
	u.sleep()
}

type DesktopDirection int

const (
	DesktopLeft DesktopDirection = iota
	DesktopRight
)

func (u *UIWin) SwitchDesktop(direction DesktopDirection) {
	u.keyDown(VK_CONTROL)
	u.keyDown(VK_LWIN)
	switch direction {
	case DesktopLeft:
		u.Press(VK_LEFT)
	case DesktopRight:
		u.Press(VK_RIGHT)
	}
	u.keyUp(VK_LWIN)
	u.keyUp(VK_CONTROL)
	u.sleep()
}

func (u *UIWin) Sleep(sec float64) {
	time.Sleep(time.Duration(sec * float64(time.Second)))
}
