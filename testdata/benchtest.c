// +build ignore

char __license[] __attribute__((section("license"), used)) = "MIT";

__attribute__((section("kprobe/do_sys_open"), used))
int open() {
    return 0;
}
