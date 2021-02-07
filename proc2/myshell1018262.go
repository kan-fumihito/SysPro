package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

var sc = bufio.NewScanner(os.Stdin)

func fork() uintptr {
	var pid, r2 uintptr
	if pid, r2, _ = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0); r2 == 1 {
		// fixup for macOS
		if runtime.GOOS == "darwin" && r2 == 1 {
			pid = 0
		}
	}
	return pid
}

func wait(pid uintptr) (status syscall.WaitStatus) {
	var rusage syscall.Rusage
	wpid, err := syscall.Wait4(int(pid), &status, syscall.WSTOPPED, &rusage)
	if err != nil {
		fmt.Printf("%v", err)
		syscall.Exit(1)
	}
	if status != 0 {
		fmt.Printf("Process %d existed with status(%d).\n", wpid, status)
		return
	}
	return
}

func execve(cpath string, cmd []string) error {
	return syscall.Exec(cpath, cmd[0:], os.Environ())
}

func job(cmd []string) int {

	pid := fork()
	if pid > 0 {
		if wait(pid) != 0 {
			return 1
		}
		return 0
	}

	//child
	cpath, err := exec.LookPath(cmd[0])
	if err != nil {
		fmt.Printf("%s not found in $PATH.\n", cmd[0])
		syscall.Exit(1)
	}
	err = execve(cpath, cmd)
	if err != nil {
		syscall.Exit(1)
	}

	return 0
}

func main() {
	prompt := os.Args[0]
	count := 0
	for {
		fmt.Printf("%s[%02d]> ", prompt, count)
		count++
		_, err := fmt.Scan()
		if err != nil {
			log.Fatalf("failed input")
		}
		if sc.Scan() {
			s := sc.Text()
			if s == "bye" {
				os.Exit(0)
			}
			if s == "" {
				continue
			}
			cmds := strings.Split(s, "?")
			cmds[0] = strings.TrimSpace(cmds[0])
			cmd := strings.Split(cmds[0], " ")
			n := job(cmd)
			if len(cmds) == 2 {
				next := strings.Split(cmds[1], ":")[n]
				next = strings.TrimSpace(next)
				cmd := strings.Split(next, " ")
				job(cmd)
			}
		} else {
			fmt.Println()
			os.Exit(0)
		}
	}
}
