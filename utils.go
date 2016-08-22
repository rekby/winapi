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

/*
return uintptr of steing for pass to winapi function and keepalive pointer to the string's buffer.
Usage:
sPtr, keepAlive := stringToUintPtr("asd")
... = winAPI(..., sPtr,...)
runtime.KeepAlive(keepAlive)

It need for GC doesn't collect converted string before winAPI function return.
 */
func stringToUintPtr(s string)(res uintptr,keepAlive *uint16){
	keepAlive = syscall.StringToUTF16Ptr(s)
	return uintptr(unsafe.Pointer(keepAlive))
}