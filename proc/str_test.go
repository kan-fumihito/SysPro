package main

import (
	"strings"
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport
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
func TestSplitSuccess(t *testing.T) {
	args := SplitMultiSep("hoge -l -a ? fuga -v : echo Done", []string{":", "?"})
	if args[0] != "hoge -l -a " {
		t.Fatalf("failed test %s\n", args[0])
	}
	if args[1] != " fuga -v " {
		t.Fatalf("failed test %s\n", args[1])
	}
	if args[2] != " echo Done" {
		t.Fatalf("failed test %s\n", args[2])
	}
}

func TestSplitSpace(t *testing.T) {
	s := "ls -l   "
	cmd := strings.Split(s, " ")
	ans := []string{"ls", "-l"}
	for i := range ans {
		if len(cmd) <= i || cmd[i] != ans[i] {
			t.Fatalf("failed test %s\n", cmd)
		}
	}
}
