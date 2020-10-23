# ebpfbench
profile eBPF programs from Go. Requires Linux kernel 5.1 or greater.

## Usage

`ebpfbench` augments the standard `testing.B` object. 

```go
func BenchmarkExample(b *testing.B) {
    eb := ebpfbench.NewEBPFBenchmark(b)
    defer eb.Close()
    
    // setup eBPF programs using cilium/ebpf, gobpf, or other libraries.
    
    fd := prog.FD()
    eb.ProfileProgram(fd, "")
    eb.Run(func(b *testing.B) {
        // exercise programs here
    })
}
```

The results per program will be output in the standard go benchmark format.
```
goos: linux
goarch: amd64
pkg: github.com/DataDog/bench-example
BenchmarkExample/eBPF-4         	                   66741	     20344 ns/op
BenchmarkExample/eBPF/kprobe/sys_bind           	       1	       568 ns/op
BenchmarkExample/eBPF/kprobe/sys_socket         	       3	      1110 ns/op
BenchmarkExample/eBPF/kprobe/tcp_cleanup_rbuf   	  266952	       295 ns/op
```
