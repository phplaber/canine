package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/phplaber/canine/pkg/scan"
)

func Output(rst []map[string]string, flag string) {
	fmt.Printf("[*] Found %d entries that are %s\n", len(rst), flag)

	for _, file := range rst {
		fmt.Printf("%10s %s %s %s %s\n", file["ftype"], file["perm"], file["owner"], file["group"], file["absPath"])
	}
}

func Scan(user string, groups string, dirpath string, SUIDFilesChan chan []map[string]string, SGIDFilesChan chan []map[string]string, writableFilesChan chan []map[string]string) {
	// SUID、SGID可执行文件和可写文件
	var SUIDFiles, SGIDFiles, writableFiles []map[string]string
	filepath.Walk(dirpath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		m := info.Mode()
		_, _, owner, group := scan.GetFileOwnership(info)
		ftype := scan.GetFileType(info)
		perm := scan.GetFilePerm(info)
		abspath := path

		// SUID 可执行文件
		if strings.Contains(ftype, "u") && (owner == "root" || owner == "system") && strings.Contains("7,5,3,1", perm[len(perm)-1:]) {
			SUIDFiles = append(SUIDFiles, map[string]string{
				"ftype":   ftype,
				"perm":    perm,
				"owner":   owner,
				"group":   group,
				"absPath": abspath,
			})
		}

		// SGID 可执行文件
		if strings.Contains(ftype, "g") && (group == "root" || group == "system") && strings.Contains("7,5,3,1", perm[len(perm)-1:]) {
			SGIDFiles = append(SGIDFiles, map[string]string{
				"ftype":   ftype,
				"perm":    perm,
				"owner":   owner,
				"group":   group,
				"absPath": abspath,
			})
		}

		// 可写文件
		// 1、文件属主为 user，且有 w 权限
		// 2、文件属组在 groups 中，且有 w 权限
		// 3、其它用户有 w 权限
		if (owner == user && m&0200 != 0) || (strings.Contains(groups, group) && m&0020 != 0) || m&0002 != 0 {
			writableFiles = append(writableFiles, map[string]string{
				"ftype":   ftype,
				"perm":    perm,
				"owner":   owner,
				"group":   group,
				"absPath": abspath,
			})
		}

		return nil
	})

	SUIDFilesChan <- SUIDFiles
	SGIDFilesChan <- SGIDFiles
	writableFilesChan <- writableFiles
}
