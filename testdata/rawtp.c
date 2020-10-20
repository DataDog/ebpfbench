// +build ignore

char __license[] __attribute__((section("license"), used)) = "MIT";

__attribute__((section("raw_tracepoint/sys_enter"), used)) int open() { return 0; }
