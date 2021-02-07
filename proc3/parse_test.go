package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func SplitMultiSep(s string, sep []string) []string {
	var ret []string
	ret = strings.Split(s, sep[0])
	if len(sep) > 1 {
		ret2 := []string{}
		for _, r := range ret {
			ret2 = append(ret2, SplitMultiSep(r, sep[1:])...)
		}
		ret = ret2
	}
	return ret
}

/*
func TestPipeParse(t *testing.T) {
	input := "  ls < out.txt"
	s := strings.Split(input, " ")
	//s = remove(s, "")
	fmt.Println(s)
	for _, v := range s {
		fmt.Printf("'%s'\n", v)
	}
}
*/
func TestReg(t *testing.T) {
	cmd := "wc < a.txt"
	r, _ := regexp.Compile(`<\s*(\w*(\.*\w*)*)`)
	res := r.FindAllStringSubmatch(cmd, -1)
	fmt.Println(res)
	//先頭にあるリダイレクションから処理
	for _, v := range res {
		//リダイレクション部分を削除
		cmd = strings.Replace(cmd, v[0], "", -1)
		fmt.Println(cmd)
	}
}
