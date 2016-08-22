package winapi

import (
	"syscall"
	"unsafe"
	"unicode/utf16"
)

var (
	libKernel32         = syscall.NewLazyDLL("kernel32.dll")
	procFindFirstVolume = libKernel32.NewProc("FindFirstVolumeW")
	procFindNextVolume  = libKernel32.NewProc("FindNextVolumeW")
	procFindVolumeClose = libKernel32.NewProc("FindVolumeClose")
	procGetVolumePathNamesForVolumeName = libKernel32.NewProc("GetVolumePathNamesForVolumeNameW")
)

func FindFirstVolume() (string, HANDLE, syscall.Errno) {
	var bufLen = 50 // "\\?\Volume{8e1acfaf-0000-0000-0000-100000000000}\" + zero byte
	for {
		buf := make([]uint16, bufLen)
		handleUintPtr, _, err := procFindFirstVolume.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(bufLen))
		if err.(syscall.Errno) == ERROR_FILENAME_EXCED_RANGE {
			bufLen *= 2
			continue
		}
		return syscall.UTF16ToString(buf), HANDLE(handleUintPtr), err.(syscall.Errno)
	}
}

func FindNextVolume(handle HANDLE) (string, bool, syscall.Errno) {
	var bufLen = 50 // "\\?\Volume{8e1acfaf-0000-0000-0000-100000000000}\" + zero byte
	for {
		buf := make([]uint16, bufLen)
		resUintPtr, _, err := procFindNextVolume.Call(uintptr(handle), uintptr(unsafe.Pointer(&buf[0])), uintptr(bufLen))
		if err.(syscall.Errno) == ERROR_FILENAME_EXCED_RANGE {
			bufLen *= 2
			continue
		}

		return syscall.UTF16ToString(buf), ptrToBool(resUintPtr), err.(syscall.Errno)
	}
}

func FindVolumeClose(handle HANDLE) (bool, syscall.Errno) {
	resUintPtr, _, err := procFindVolumeClose.Call(uintptr(handle))
	return resUintPtr != 0, err.(syscall.Errno)
}

/*
volumeName - A volume GUID path for the volume. A volume GUID path is of the form "\\?\Volume{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}\"
             Same form as return FindFirstVolume, FindNextVolume.
*/
func GetVolumePathNamesForVolumeName(volumeName string) ([]string, bool, syscall.Errno) {
	var bufLen = 256
	var returnBufLen uint32
	for {
		buf := make([]uint16, bufLen)
		resUintPtr, _, err := procGetVolumePathNamesForVolumeName.Call(
			stringToUintPtr(volumeName),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(bufLen),
			uintptr(unsafe.Pointer(&returnBufLen)))
		errno := err.(syscall.Errno)
		if !ptrToBool(resUintPtr) {
			if errno == ERROR_MORE_DATA {
				bufLen = int(returnBufLen)
				continue
			} else {
				return nil, false, errno
			}
		}

		res := []string{}
		for index := uint32(0); index < returnBufLen; {
			end := index
			for buf[end] != 0 && end < returnBufLen {
				end++
			}
			if end == index {
				break
			}
			s := buf[index:end]
			res = append(res, string(utf16.Decode(s)))
			index = end + 1
		}
		return res, true, errno
	}
}