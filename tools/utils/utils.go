package utils

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 必须读取内嵌文件数据
func MustReadEmbedFile(fs embed.FS, file string) []byte {
	data, err := fs.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return data
}

// 必须读取文件数据
func MustReadFile(file string) []byte {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return data
}

// 目录是否为空
func DirIsEmpty(path string) (bool, error) {
	dir, err := os.Open(path)
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
	empty, err := DirIsEmpty(path)
	if err != nil {
		panic(err)
	}
	if !empty {
		panic(fmt.Errorf("目录必须为空"))
	}
}

// 必须创建文件夹
func MustMkdir(name string, perm ...os.FileMode) {
	var p os.FileMode = 0666
	if len(perm) > 0 {
		p = perm[0]
	}
	err := os.MkdirAll(name, p)
	if err != nil {
		panic(err)
	}
}

// 必须创建文件夹并且是创建
func MustMkdirAndIsCreate(name string, perm ...os.FileMode) {
	var p os.FileMode = 0666
	if len(perm) > 0 {
		p = perm[0]
	}
	_, err := os.Open(name)
	if err == nil {
		panic(fmt.Errorf("文件夹'%s'已存在", name))
	}
	if !os.IsNotExist(err) {
		panic(err)
	}

	err = os.MkdirAll(name, p)
	if err != nil {
		panic(err)
	}
}

// 必须创建文件
func MustWriteFile(name string, data []byte, perm ...os.FileMode) {
	var p os.FileMode = 0666
	if len(perm) > 0 {
		p = perm[0]
	}
	err := os.WriteFile(name, data, p)
	if err != nil {
		panic(err)
	}
}

// 检查是否存在路径并指定为文件夹或文件, 读取失败或者权限不足会panic
func CheckHasPath(path string, isDir bool) bool {
	d, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	s, err := d.Stat()
	if err != nil {
		panic(fmt.Errorf("'%s'读取失败"))
	}
	return s.IsDir() == isDir
}

// 必须存在路径
func MustHasPath(path string, isDir bool) {
	d, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Errorf("'%s'不存在", path))
		}
		panic(err)
	}
	s, err := d.Stat()
	if err != nil {
		panic(fmt.Errorf("'%s'读取失败"))
	}
	if isDir && !s.IsDir() {
		panic(fmt.Errorf("'%s'不是一个目录", path))
	}
	if !isDir && s.IsDir() {
		panic(fmt.Errorf("'%s'不是一个文件", path))
	}
}

// 必须不存在某个路径
func MustNoExistPath(path string) {
	_, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	panic(fmt.Errorf("'%s'已存在", path))
}

// 必须在项目目录中并获取项目名
func MustGetProjectName() string {
	if !CheckHasPath("go.mod", false) ||
		!CheckHasPath("component", true) ||
		!CheckHasPath("configs", true) ||
		!CheckHasPath("spiders", true) {
		_, _ = os.Stderr.WriteString("必须在项目中\n")
		os.Exit(1)
	}

	data, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}
	projectName := ExtractMiddleText(string(data), "module ", "\n", "", false)
	projectName = strings.TrimSuffix(projectName, "\r")
	if projectName == "" {
		panic(errors.New("无法读取项目名"))
	}
	return projectName
}

// 必须在项目目录中
func MustInProjectDir() {
	_ = MustGetProjectName()
}

/**提取中间文本
  s 原始文本
  pre 提取数据的前面的数据, 如果为空则从最开头提取
  suf 提取数据的后面的数据, 如果为空则提取到结尾
  def 找不到时返回的默认数据
  greedy 贪婪的, 默认从开头开始查找suf, 如果是贪婪的则从结尾开始查找suf
*/
func ExtractMiddleText(s, pre, suf, def string, greedy bool) string {
	var start int  // 开始位置
	if pre != "" { // 需要查找开始数据
		k := strings.Index(s, pre)
		if k == -1 {
			return def
		}
		start = k + len(pre)
	}

	if suf == "" {
		return s[start:]
	}

	// 结束位置
	var end int
	if greedy { // 贪婪的从结尾开始查找suf
		end = strings.LastIndex(s[start:], suf)
	} else {
		end = strings.Index(s[start:], suf)
	}
	if end == -1 {
		return def
	}
	end += start // 添加偏移
	return s[start:end]
}

// 必须获取工作目录
func MustGetWorkdir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// 必须连接目录
func MustDirJoin(path1, path2 string) string {
	if !filepath.IsAbs(path2) {
		path2 = filepath.Join(path1, path2)
	}
	path, err := filepath.Abs(path2)
	if err != nil {
		panic(err)
	}
	return path
}

// 构建模板参数
func MakeTemplateArgs(projectName string) map[string]interface{} {
	numCpu := runtime.NumCPU()
	if numCpu < 1 {
		numCpu = 1
	}
	return map[string]interface{}{
		"project_name": projectName,      // 项目名
		"project_dir":  MustGetWorkdir(), // 项目路径
		"date":         time.Now().Format("2006-01-02"),
		"time":         time.Now().Format("15:04:05"),
		"date_time":    time.Now().Format("2006-01-02 15:04:05"),
		"num_cpu":      numCpu,
	}
}
