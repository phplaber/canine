package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

func usage() {
	fmt.Println(
		`canine v0.1
	A tool for find andriod attack surface of file system
	Usage: canine -u [user] -g [groups] dirpath1 dirpath2 ...`)

	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	owner := flag.String("u", "", "username, e.g. shell")
	groups := flag.String("g", "", "groupname(s), e.g. shell,log,sdcard_rw")

	flag.Usage = usage
	flag.Parse()

	dirpaths := flag.Args()

	fmt.Println("[*] Scanning...")
	var wg sync.WaitGroup
	for _, dirpath := range dirpaths {
		wg.Add(1)
		// 避免在每个协程闭包中重复利用相同的 dirpath 值
		dirpath := dirpath
		go func() {
			defer wg.Done()
			Scan(*owner, *groups, dirpath)
		}()
	}
	wg.Wait()
}
