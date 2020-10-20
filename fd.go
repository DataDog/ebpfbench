package ebpfbench

import (
	"fmt"
	"runtime"

	"golang.org/x/sys/unix"
)

type wrappedFD struct {
	sysfd int
}

func NewWrappedFD(raw int) *wrappedFD {
	f := &wrappedFD{sysfd: raw}
	runtime.SetFinalizer(f, (*wrappedFD).Close)
	return f
}

func (f *wrappedFD) Raw() (int, error) {
	if f.sysfd < 0 {
		return 0, fmt.Errorf("use of closed wrappedFD")
	}
	return f.sysfd, nil
}

func (f *wrappedFD) Close() error {
	if f.sysfd < 0 {
		return nil
	}
	val := f.sysfd
	f.sysfd = -1
	runtime.SetFinalizer(f, nil)
	return unix.Close(val)
}
