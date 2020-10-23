package ebpfbench

import (
	"io/ioutil"
	"os"
	"testing"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpf -cflags $CFLAGS nohelpertest testdata/nohelpertest.c -- -O2
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpf -cflags $CFLAGS helpertest testdata/helpertest.c -- -O2

func openBench(b *testing.B) {
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
}

func BenchmarkNoHelper(b *testing.B) {
	// setup ebpf benchmark
	eb := NewEBPFBenchmark(b)
	defer eb.Close()

	// setup ebpf kprobe
	specs, err := newNohelpertestSpecs()
	if err != nil {
		b.Fatal(err)
	}
	objs, err := specs.Load(nil)
	if err != nil {
		b.Fatal(err)
	}
	defer objs.Close()

	efd, err := kprobeAttach(false, "do_sys_open", objs.ProgramOpen.FD())
	if err != nil {
		b.Fatal(err)
	}
	defer kprobeDetach(efd)

	// register probe with benchmark and run
	eb.ProfileProgram(objs.ProgramOpen.FD(), "")
	eb.Run(openBench)
}

func BenchmarkHelper(b *testing.B) {
	// setup ebpf benchmark
	eb := NewEBPFBenchmark(b)
	defer eb.Close()

	// setup ebpf kprobe
	specs, err := newHelpertestSpecs()
	if err != nil {
		b.Fatal(err)
	}
	objs, err := specs.Load(nil)
	if err != nil {
		b.Fatal(err)
	}
	defer objs.Close()

	efd, err := kprobeAttach(false, "do_sys_open", objs.ProgramOpen.FD())
	if err != nil {
		b.Fatal(err)
	}
	defer kprobeDetach(efd)

	// register probe with benchmark and run
	eb.ProfileProgram(objs.ProgramOpen.FD(), "")
	eb.Run(openBench)
}
