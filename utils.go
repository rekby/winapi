package winapi

import (
	"unsafe"
	"syscall"
)

func ptrToBool(v uintptr)bool{
	if v == 0 {
		return false
	}
	return true
}

func stringToUintPtr(s string)uintptr{
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}