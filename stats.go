package ebpfbench

import "time"

type disableFunc func() error

func enableBPFStats() (disableFunc, error) {
	err := writeSysctl(bpfSysctlProcfile, []byte("1"))
	if err != nil {
		return nil, err
	}
	return disableBPFStats, nil
}

func disableBPFStats() error {
	return writeSysctl(bpfSysctlProcfile, []byte("0"))
}

type bpfProgramStats struct {
	RunCount uint
	RunTime  time.Duration
}

func getProgramStats(fd int) (*bpfProgramStats, error) {
	pi, err := bpfGetProgInfoByFD(fd)
	if err != nil {
		return nil, err
	}
	return &bpfProgramStats{
		RunCount: uint(pi.run_cnt),
		RunTime:  time.Duration(pi.run_time_ns),
	}, nil
}
