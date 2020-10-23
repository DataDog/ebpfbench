package ebpfbench

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

func bpf(cmd int, attr unsafe.Pointer, size uintptr) (uintptr, error) {
	var err error

	r1, _, errno := unix.Syscall(unix.SYS_BPF, uintptr(cmd), uintptr(attr), size)
	runtime.KeepAlive(attr)
	if errno != 0 {
		err = errno
	}
	return r1, err
}

type bpfObjGetInfoByFdAttr struct {
	bpf_fd   uint32
	info_len uint32
	info     unsafe.Pointer
}

type bpfProgInfo struct {
	progType                 uint32
	id                       uint32
	tag                      [unix.BPF_TAG_SIZE]byte
	jited_prog_len           uint32
	xlated_prog_len          uint32
	jited_prog_insns         unsafe.Pointer
	xlated_prog_insns        unsafe.Pointer
	load_time                uint64
	created_by_uid           uint32
	nr_map_ids               uint32
	map_ids                  unsafe.Pointer
	name                     [unix.BPF_OBJ_NAME_LEN]byte
	ifindex                  uint32
	gpl_compatible           uint32
	netns_dev                uint64
	netns_ino                uint64
	nr_jited_ksyms           uint32
	nr_jited_func_lens       uint32
	jited_ksyms              unsafe.Pointer
	jited_func_lens          unsafe.Pointer
	btf_id                   uint32
	func_info_rec_size       uint32
	func_info                unsafe.Pointer
	nr_func_info             uint32
	nr_line_info             uint32
	line_info                unsafe.Pointer
	jited_line_info          unsafe.Pointer
	nr_jited_line_info       uint32
	line_info_rec_size       uint32
	jited_line_info_rec_size uint32
	nr_prog_tags             uint32
	prog_tags                unsafe.Pointer
	run_time_ns              uint64
	run_cnt                  uint64
}

func bpfGetProgInfoByFD(fd int) (*bpfProgInfo, error) {
	pi := bpfProgInfo{}
	attr := bpfObjGetInfoByFdAttr{
		bpf_fd:   uint32(fd),
		info_len: uint32(unsafe.Sizeof(pi)),
		info:     unsafe.Pointer(&pi),
	}

	_, err := bpf(unix.BPF_OBJ_GET_INFO_BY_FD, unsafe.Pointer(&attr), unsafe.Sizeof(attr))
	if err != nil {
		return nil, fmt.Errorf("cannot get obj info by fd: %w", err)
	}
	return &pi, nil
}

type bpfEnableStatsAttr struct {
	enable_stats struct {
		statsType uint32
	}
}

func bpfEnableStats() (*wrappedFD, error) {
	attr := bpfEnableStatsAttr{}
	attr.enable_stats.statsType = unix.BPF_STATS_RUN_TIME

	fd, err := bpf(unix.BPF_ENABLE_STATS, unsafe.Pointer(&attr), unsafe.Sizeof(attr))
	if err != nil {
		return nil, fmt.Errorf("cannot enable bpf stats: %w", err)
	}
	return NewWrappedFD(int(fd)), nil
}

func cstr(s string) unsafe.Pointer {
	// zero terminate the string
	buf := make([]byte, len(s)+1)
	copy(buf, s)

	return unsafe.Pointer(&buf[0])
}

func goString(s []byte) string {
	str := string(s[:])
	if li := strings.LastIndexByte(str, 0); li > 0 {
		return str[:li]
	}
	return ""
}

func supportsBpfEnableStats() func() bool {
	var once sync.Once
	result := false

	return func() bool {
		once.Do(func() {
			fd, err := bpfEnableStats()
			if err != nil {
				return
			}
			result = true
			_ = fd.Close()
		})
		return result
	}
}
