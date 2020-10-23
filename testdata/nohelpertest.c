#include "bpf_helpers.h"

SEC("kprobe/do_sys_open")
int open() {
    return 0;
}

char __license[] SEC("license") = "Dual BSD/GPL";
