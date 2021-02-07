package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

//文字列リストからNull文字を削除する関数
func remove(strings []string, search string) []string {
	result := []string{}
	for _, v := range strings {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}

//リダイレクション処理をしてコマンド実行する関数
func exe(cmd string, attr *syscall.ProcAttr) int {
	//標準入力をリダイレクション
	if strings.Contains(cmd, "<") {
		r, _ := regexp.Compile(`<\s*(\w*(\.*\w*)*)`)
		res := r.FindAllStringSubmatch(cmd, -1)

		//先頭にあるリダイレクションから処理
		for _, v := range res {
			//リダイレクション部分を削除
			cmd = strings.Replace(cmd, v[0], "", -1)
			fr, err := os.OpenFile(v[1], os.O_RDONLY, 0666)
			if err != nil {
				return -1
			}
			defer func() {
				if err := fr.Close(); err != nil {
					panic(err)
				}
			}()
			attr.Files[0] = fr.Fd()
		}
	}

	//標準出力をリダイレクション
	if strings.Contains(cmd, ">") {
		r, _ := regexp.Compile(`>\s*(\w*(\.*\w*)*)`)
		res := r.FindAllStringSubmatch(cmd, -1)

		//先頭にあるリダイレクションから処理
		for _, v := range res {
			//リダイレクション部分を削除
			cmd = strings.Replace(cmd, v[0], "", -1)
			fw, err := os.OpenFile(v[1],
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return -1
			}
			defer func() {
				if err := fw.Close(); err != nil {
					panic(err)
				}
			}()
			attr.Files[1] = fw.Fd()
			attr.Files[2] = fw.Fd()
		}
	}

	//コマンド実行処理
	cmdargv := strings.Split(cmd, " ")
	cmdargv = remove(cmdargv, "")
	cpath, err := exec.LookPath(cmdargv[0])
	if err != nil {
		fmt.Printf("%s not found in $PATH.\n", cmdargv[0])
		return -1
	}

	pid, err := syscall.ForkExec(cpath, cmdargv, attr)
	if err != nil {
		panic(err)
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		panic(err)
	}
	status, err := proc.Wait()
	if err != nil {
		panic(err)
	}
	return status.ExitCode()
}

//パイプライン処理を行う関数
func pipeline(cmdlist []string, attr *syscall.ProcAttr) int {
	//パイプラインによる分割がなかった時
	if len(cmdlist) == 1 {
		return exe(cmdlist[0], attr)
	}
	//パイプライン作成
	pin, pout, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	//先頭コマンドの標準出力をパイプラインに付け替えて実行
	attr.Files[1] = pout.Fd()
	ecode1 := exe(cmdlist[0], attr)
	pout.Close()

	//2番目のコマンドの標準入力をパイプラインに付け替えて再帰的にパイプライン処理
	attr2 := &syscall.ProcAttr{
		Files: []uintptr{pin.Fd(), 1, 2},
	}
	ecode2 := pipeline(cmdlist[1:], attr2)
	pin.Close()

	//エラーコード処理
	if ecode1 != 0 {
		return ecode1
	}
	if ecode2 != 0 {
		return ecode2
	}
	return 0
}

func main() {
	prompt := os.Args[0]
	count := 0
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s[%02d]> ", prompt, count)
		var input string
		_, err := fmt.Scan()
		if err != nil {
			log.Fatalf("failed input")
		}
		if sc.Scan() {
			input = sc.Text()
			if input == "" {
				continue
			} else if input == "bye" {
				os.Exit(0)
			}
		} else {
			fmt.Println()
			os.Exit(0)
		}

		if strings.Contains(input, "?") {
			expr := strings.Split(input, "?")
			attr := &syscall.ProcAttr{
				Files: []uintptr{0, 1, 2},
			}
			cmd := strings.Split(expr[0], "|")
			status := pipeline(cmd, attr)

			idx := 0
			if status != 0 {
				idx = 1
			}
			cmdlist := strings.Split(expr[1], ":")[idx]
			attr = &syscall.ProcAttr{
				Files: []uintptr{0, 1, 2},
			}
			cmd = strings.Split(cmdlist, "|")
			pipeline(cmd, attr)

		} else {
			attr := &syscall.ProcAttr{
				Files: []uintptr{0, 1, 2},
			}
			cmd := strings.Split(input, "|")
			pipeline(cmd, attr)
		}
		count++
	}
}
