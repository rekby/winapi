package winapi

import "syscall"

const (
	ERROR_NO_MORE_FILES        syscall.Errno = 18
	ERROR_FILENAME_EXCED_RANGE               = 206
	ERROR_MORE_DATA                          = 234
)
