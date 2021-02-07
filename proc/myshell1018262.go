package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var sc = bufio.NewScanner(os.Stdin)

func exe(cmd []string, attr syscall.ProcAttr) int {
	//check bye
	if cmd[0] == "bye" {
		return 1
	}

	//check cpath
	cpath, err := exec.LookPath(cmd[0])
	if err != nil {
		fmt.Printf("%s not found in $PATH.\n", cmd[0])
		return 256
	}

	//fork
	pid, err := syscall.ForkExec(cpath, cmd, &attr)
	if err != nil {
		panic(err)
	}

	//find process and wait
	proc, err := os.FindProcess(pid)
	status, err := proc.Wait()
	if err != nil {
		panic(err)
	}

	//check result
	if !status.Success() {
		fmt.Println(status.String())
		return status.ExitCode()
	}
	return status.ExitCode()
}

func main() {
	prompt := os.Args[0]
	count := 0
	for {
		fmt.Printf("%s[%d]> ", prompt, count)
		count++
		var args []string
		_, err := fmt.Scan()
		if err != nil {
			log.Fatalf("failed input")
		}
		if sc.Scan() {
			s := sc.Text()
			args = strings.Split(s, "?")
		}

		attr := syscall.ProcAttr{Files: []uintptr{0, 1, 2}}

		if args == nil {
			fmt.Println()
			return
		}
		arg := strings.TrimSpace(args[0])
		cmd := strings.Split(arg, " ")
		res := exe(cmd, attr)
		if res == 1 {
			return
		}
		if len(args) > 1 {
			nexts := strings.Split(args[1], ":")
			switch res {
			case -1:
				next := strings.TrimSpace(nexts[1])
				cmd := strings.Split(next, " ")
				exe(cmd, attr)
			case 0:
				next := strings.TrimSpace(nexts[0])
				cmd := strings.Split(next, " ")
				exe(cmd, attr)
			case 1:
				return
			}
		}
	}
}
