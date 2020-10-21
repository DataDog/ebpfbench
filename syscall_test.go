package ebpfbench

import (
	"os"
	"testing"

	"github.com/cilium/ebpf"

	"github.com/cilium/ebpf/link"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go getproginfo testdata/rawtp.c

func TestGetProgInfo(t *testing.T) {
	disableFunc, err := enableBPFStats()
	if err != nil {
		t.Fatal(err)
	}
	defer disableFunc()

	specs, err := newGetproginfoSpecs()
	if err != nil {
		t.Fatal(err)
	}

	prog, err := ebpf.NewProgram(specs.ProgramOpen)
	if err != nil {
		t.Fatal(err)
	}
	defer prog.Close()

	link, err := link.AttachRawTracepoint(link.RawTracepointOptions{
		Name:    "sys_enter",
		Program: prog,
	})
	if err != nil {
		t.Fatal(err)
	}
	if link == nil {
		t.Fatal("no link")
	}
	defer link.Close()

	f, err := os.Open("/etc/os-release")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	stats, err := getProgramStats(prog.FD())
	if err != nil {
		t.Fatal(err)
	}

	if stats.RunCount == 0 {
		t.Errorf("run count should be non-zero")
	}
	if stats.RunTime == 0 {
		t.Errorf("run time should be non-zero")
	}
}
