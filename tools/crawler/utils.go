package main

import (
	"fmt"
	"io"
	"os"
)

// 目录是否为空
func DirIsEmpty(path string) (bool, error) {
	dir, err := os.Open(path)
	if err == io.EOF {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	fs, err := dir.ReadDir(1)
	if err == io.EOF {
		return true, nil
	}
	return len(fs) == 0, err
}

// 目录必须为空
func DirMustEmpty(path string) {
	empty, err := DirIsEmpty(".")
	if err != nil {
		panic(err)
	}
	if !empty {
		panic(fmt.Errorf("目录必须为空"))
	}
}

// 必须创建文件夹
func MustMkdir(perm os.FileMode, names ...string) {
	for _, name := range names {
		err := os.MkdirAll(name, perm)
		if err != nil {
			panic(err)
		}
	}
}

// 必须创建文件
func MustWriteFile(name string, data []byte, perm os.FileMode) {
	err := os.WriteFile(name, data, perm)
	if err != nil {
		panic(err)
	}
}
