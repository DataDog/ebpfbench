package ebpfbench

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func NewEBPFBenchmark(b *testing.B) *EBPFBenchmark {
	return &EBPFBenchmark{
		b:     b,
		progs: make(map[int]string),
	}
}

type EBPFBenchmark struct {
	b     *testing.B
	progs map[int]string
}

func (e *EBPFBenchmark) ProfileProgram(fd int, name string) {
	e.progs[fd] = name
}

func (e *EBPFBenchmark) getAllStats() (map[int]*bpfProgramStats, error) {
	res := map[int]*bpfProgramStats{}
	for fd := range e.progs {
		stats, err := getProgramStats(fd)
		if err != nil {
			return nil, err
		}
		res[fd] = stats
	}
	return res, nil
}

func (e *EBPFBenchmark) Close() {
	e.progs = make(map[int]string)
}

func (e *EBPFBenchmark) Run(fn func(*testing.B)) {
	disableFunc, err := enableBPFStats()
	if err != nil {
		e.b.Fatal(err)
	}
	defer func() { _ = disableFunc() }()

	var results map[string]*testing.BenchmarkResult
	e.b.Run("eBPF", func(b *testing.B) {
		baseline, err := e.getAllStats()
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		fn(b)
		b.StopTimer()

		post, err := e.getAllStats()
		if err != nil {
			b.Fatal(err)
		}

		// override outer variable here so we only report on the last run of results
		results = make(map[string]*testing.BenchmarkResult, len(baseline))
		for fd, base := range baseline {
			p := post[fd]
			runTime := p.RunTime - base.RunTime
			runCount := p.RunCount - base.RunCount
			name := e.progs[fd]
			if name == "" && p.Name != "" {
				name = p.Name
			}
			results[name] = &testing.BenchmarkResult{
				N: int(runCount),
				T: runTime,
			}
		}
	})
	fmt.Print(prettyPrintEBPFResults(e.b.Name(), results))
}

func prettyPrintEBPFResults(benchName string, results map[string]*testing.BenchmarkResult) string {
	maxLen := 0
	var names []string
	for name := range results {
		if len(name) > maxLen {
			maxLen = len(name)
		}
		names = append(names, name)
	}
	sort.Strings(names)
	buf := new(strings.Builder)
	for _, name := range names {
		pr := results[name]
		_, _ = fmt.Fprintf(buf, "%s/eBPF/%-*s\t%s\n", benchName, maxLen, name, pr.String())
	}
	return buf.String()
}
