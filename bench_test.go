package ebpfbench

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cilium/ebpf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go benchtest testdata/benchtest.c

func BenchmarkTest(b *testing.B) {
	// setup ebpf benchmark
	eb := NewEBPFBenchmark(b)
	defer eb.Close()

	// setup ebpf kprobe
	specs, err := newBenchtestSpecs()
	if err != nil {
		b.Fatal(err)
	}
	prog, err := ebpf.NewProgram(specs.ProgramOpen)
	if err != nil {
		b.Fatal(err)
	}
	defer prog.Close()
	efd, err := kprobeAttach(false, "do_sys_open", prog.FD())
	if err != nil {
		b.Fatal(err)
	}
	defer kprobeDetach(efd)

	// register probe with benchmark and run
	eb.ProfileProgram(prog.FD(), "kprobe/do_sys_open")
	eb.Run(func(b *testing.B) {
		// open b.N temp files
		for i := 0; i < b.N; i++ {
			f, err := ioutil.TempFile(os.TempDir(), "ebpf-benchtest-*")
			if err != nil {
				b.Fatal(err)
			}
			_, err = f.Write([]byte{1})
			if err != nil {
				b.Fatal(err)
			}
			fn := f.Name()
			_ = f.Close()
			_ = os.Remove(fn)
		}
	})
}
