package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"LightControl/src/Config"
	"LightControl/src/Dates"
	"LightControl/src/Hue/Colors"
)

var (
	user32                        = syscall.NewLazyDLL("user32.dll")
	kernel32                      = syscall.NewLazyDLL("kernel32.dll")
	procEnumWindows               = user32.NewProc("EnumWindows")
	procGetWindowTextW            = user32.NewProc("GetWindowTextW")
	procGetWindowThreadProcessId  = user32.NewProc("GetWindowThreadProcessId")
	procOpenProcess               = kernel32.NewProc("OpenProcess")
	procQueryFullProcessImageName = kernel32.NewProc("QueryFullProcessImageNameW")
	procCloseHandle               = kernel32.NewProc("CloseHandle")
)

var enumProcessFunc uintptr
var getWindowTitlesFunc func() []string
var initializedProcessesFunc = false

const (
	ProcessQueryInformation = 0x0400
	ProcessVmRead           = 0x0010
)

func EnumWindows(enumFunc uintptr, lParam uintptr) bool {
	ret, _, _ := procEnumWindows.Call(enumFunc, lParam)
	return ret != 0
}

func initializeEnumProcessFunc() {
	var titles []string
	enumProcessFunc = syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		_, windowPID := GetWindowThreadProcessId(hwnd)
		title := GetWindowText(hwnd)
		name := GetProcessName(windowPID)
		if title != "" {
			titles = append(titles, strings.ToLower(fmt.Sprintf("%s %s", title, name)))
		}
		return 1 // Continue enumeration
	})

	getWindowTitlesFunc = func() []string {
		titles = make([]string, 0)
		EnumWindows(enumProcessFunc, 0)
		return titles
	}
	initializedProcessesFunc = true
}

func GetWindowText(hwnd syscall.Handle) string {
	buf := make([]uint16, 256) // Buffer for window title (256 characters max)
	ret, _, _ := procGetWindowTextW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if ret == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

func GetWindowThreadProcessId(hwnd syscall.Handle) (threadID, processID uint32) {
	_, _, _ = procGetWindowThreadProcessId.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&processID)),
	)
	return
}

func GetProcessName(pid uint32) string {
	// Open the process with required access rights
	hProcess, _, _ := procOpenProcess.Call(
		uintptr(ProcessQueryInformation|ProcessVmRead),
		0,
		uintptr(pid),
	)
	if hProcess == 0 {
		return ""
	}
	defer procCloseHandle.Call(hProcess)

	// Query the process image name
	buf := make([]uint16, syscall.MAX_PATH)
	size := uint32(len(buf))
	ret, _, _ := procQueryFullProcessImageName.Call(
		hProcess,
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if ret == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

func FindWindowNames() []string {
	if !initializedProcessesFunc {
		initializeEnumProcessFunc()
	}
	return getWindowTitlesFunc()
}

func GetProcessColor(config *Config.Config) (Colors.Color, error) {
	processNames := FindWindowNames()

	resultPriority := math.MaxInt
	var resultColor Colors.Color

	now := time.Now()
	dayTime := Dates.NewDayTime(now.Hour(), now.Minute(), now.Second())

	for _, processName := range processNames {
		i := 0
		if found, color := getTimedProgramColor(processName, config, dayTime, &i, resultPriority); found {
			resultColor = color
			resultPriority = i
		}
		if found, color := getProgramColor(processName, config, &i, resultPriority); found {
			resultColor = color
			resultPriority = i
		}
	}

	if resultPriority == math.MaxInt {
		return nil, errors.New("no matching process found")
	}

	return resultColor, nil
}

func getTimedProgramColor(processName string, config *Config.Config, dayTime Dates.DayTime, i *int, maxPriority int) (bool, Colors.Color) {
	for _, timedProgram := range config.TimedPrograms {
		if *i >= maxPriority {
			return false, nil
		}
		if strings.Contains(processName, timedProgram.Name) && timedProgram.Span.Contains(dayTime) {
			// Check if color is valid
			_, err := timedProgram.Color.GetColor()
			if err != nil {
				continue
			}

			return true, timedProgram.Color
		}
		*i++
	}
	return false, nil
}

func getProgramColor(processName string, config *Config.Config, i *int, maxPriority int) (bool, Colors.Color) {
	for _, program := range config.Programs {
		if *i >= maxPriority {
			return false, nil
		}
		if strings.Contains(processName, program.Name) {
			// Check if color is valid
			_, err := program.Color.GetColor()
			if err != nil {
				continue
			}

			return true, program.Color
		}
		*i++
	}
	return false, nil
}
