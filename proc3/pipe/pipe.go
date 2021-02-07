package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// split the args before and after the pipe separator
	var args1 []string
	var args2 []string
	args := []string{"pipe", "ls", "l", "|", "ls", "k"}
	for i, v := range args {
		if v == "|" {
			args1 = make([]string, i-1)
			copy(args1, args[1:i])
			args2 = args[i+1:]
		}
	}

	// make a pipe
	pin, pout, err := os.Pipe()
	// execute the 1st command
	cmdpath1, err := exec.LookPath(args1[0])
	if err != nil {
		log.Fatalf("%s not found in $PATH.", args1[0])
	}
	attr1 := syscall.ProcAttr{
		Files: []uintptr{0, pout.Fd(), 2}}
	_, err = syscall.ForkExec(cmdpath1, args1, &attr1)
	if err != nil {
		panic(err)
	}
	pout.Close()
	// execute the 2nd command
	cmdpath2, err := exec.LookPath(args2[0])
	if err != nil {
		log.Fatalf("%s not found in $PATH.", args2[0])
	}
	attr2 := syscall.ProcAttr{
		Files: []uintptr{pin.Fd(), 1, 2}}
	pid, err := syscall.ForkExec(cmdpath2, args2, &attr2)
	if err != nil {
		panic(err)
	}
	pin.Close()
	// wait for the 2nd command to complete
	proc, err := os.FindProcess(pid)
	status, err := proc.Wait()
	if err != nil {
		panic(err)
	}
	if !status.Success() {
		fmt.Println("FAILED", status.String())
		fmt.Println(status.String())
	}

}
