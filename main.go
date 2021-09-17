package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker run <container> cmd args
// go run main.go run cmd 			args
// go run main.go run /bin/bash
// run = arg[1]

func main() {
	//fmt.Printf("arg[0] %v\n", os.Args[1])
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("what??")
	}
}

func run() {
	//fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Namespaces : what you can see
	// cgroups : what you can use

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("container")))

	// i need a file system (/home/rootfs) for the container.
	// to give it the libs, it needs.
	must(syscall.Chroot("/"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	must(cmd.Run())
}

func must(err error) {
	//fmt.Printf("return value =  %v\n", err)
	if err != nil {
		panic(err)
	}
}
