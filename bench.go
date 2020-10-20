package ebpfbench

import "time"

type DisableFunc func() error

func EnableBPFStats() (DisableFunc, error) {
	err := writeSysctl(bpfSysctlProcfile, []byte("1"))
	if err != nil {
		return nil, err
	}
	return disableBPFStats, nil
}

func disableBPFStats() error {
	return writeSysctl(bpfSysctlProcfile, []byte("0"))
}

type BPFProgramStats struct {
	RunCount uint
	RunTime  time.Duration
}

func GetProgramStats(fd int) (*BPFProgramStats, error) {
	pi, err := bpfGetProgInfoByFD(fd)
	if err != nil {
		return nil, err
	}
	return &BPFProgramStats{
		RunCount: uint(pi.run_cnt),
		RunTime:  time.Duration(pi.run_time_ns),
	}, nil
}
