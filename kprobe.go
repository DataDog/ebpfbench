package ebpfbench

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func cstr(s string) unsafe.Pointer {
	// zero terminate the string
	buf := make([]byte, len(s)+1)
	copy(buf, s)

	return unsafe.Pointer(&buf[0])
}

func kprobeAttach(retprobe bool, name string, progFd int) (int, error) {
	var err error
	attr := unix.PerfEventAttr{}
	attr.Type, err = kprobePerfType()
	if err != nil {
		return 0, fmt.Errorf("unable to determine kprobe perf type: %w", err)
	}

	if retprobe {
		bit, err := kretprobeBit()
		if err != nil {
			return 0, fmt.Errorf("unable to determine kretprobe bit: %w", err)
		}
		attr.Config |= 1 << bit
	}

	attr.Size = uint32(unsafe.Sizeof(attr))
	attr.Ext1 = uint64(uintptr(cstr(name)))
	attr.Ext2 = 0

	efd, err := unix.PerfEventOpen(&attr, -1, 0, -1, unix.PERF_FLAG_FD_CLOEXEC)
	if efd < 0 || err != nil {
		return 0, fmt.Errorf("perf_event_open error: %w", err)
	}
	if _, _, err := unix.Syscall(unix.SYS_IOCTL, uintptr(efd), unix.PERF_EVENT_IOC_SET_BPF, uintptr(progFd)); err != 0 {
		return 0, fmt.Errorf("error attaching bpf program to perf event: %w", err)
	}
	if _, _, err := unix.Syscall(unix.SYS_IOCTL, uintptr(efd), unix.PERF_EVENT_IOC_ENABLE, 0); err != 0 {
		return 0, fmt.Errorf("error enabling perf event: %w", err)
	}
	return efd, nil
}

func kprobeDetach(efd int) error {
	if _, _, err := unix.Syscall(unix.SYS_IOCTL, uintptr(efd), unix.PERF_EVENT_IOC_DISABLE, 0); err != 0 {
		return fmt.Errorf("error disabling perf event: %w", err)
	}
	return unix.Close(efd)
}

func kprobePerfType() (uint32, error) {
	f, err := os.Open("/sys/bus/event_source/devices/kprobe/type")
	if err != nil {
		return 0, err
	}
	defer func() { _ = f.Close() }()

	var kt int
	_, err = fmt.Fscanf(f, "%d\n", &kt)
	return uint32(kt), err
}

func kretprobeBit() (uint32, error) {
	f, err := os.Open("/sys/bus/event_source/devices/kprobe/format/retprobe")
	if err != nil {
		return 0, err
	}
	defer func() { _ = f.Close() }()

	var kt int
	_, err = fmt.Fscanf(f, "config:%d\n", &kt)
	return uint32(kt), err
}
