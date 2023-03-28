package utils

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/zly-app/zapp/logger"
	"github.com/zlyuancn/zstr"
	"go.uber.org/zap"
)

// 必须读取内嵌文件数据
func MustReadEmbedFile(fs embed.FS, file string) []byte {
	data, err := fs.ReadFile(file)
	if err != nil {
		logger.Log.Fatal("读取内嵌文件数据失败", zap.String("file", file), zap.Error(err))
	}
	return data
}

// 必须读取文件数据
func MustReadFile(file string) []byte {
	data, err := os.ReadFile(file)
	if err != nil {
		logger.Log.Fatal("读取文件数据失败", zap.String("file", file), zap.Error(err))
	}
	return data
}

// 目录是否为空
func DirIsEmpty(path string) (bool, error) {
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	if dir != nil {
		defer dir.Close()
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
		logger.Log.Fatal("检查目录为空失败", zap.String("path", path), zap.Error(err))
	}
	if !empty {
		logger.Log.Fatal("目录必须为空", zap.String("path", path))
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
		logger.Log.Fatal("创建文件夹失败", zap.String("path", name), zap.Error(err))
	}
}

// 必须创建文件夹或者文件夹为空
func MustCreateDirOrDirIsEmpty(name string, perm ...os.FileMode) {
	var p os.FileMode = 0666
	if len(perm) > 0 {
		p = perm[0]
	}

	dir, err := os.Open(name)
	if dir != nil {
		defer dir.Close()
	}
	if err == nil { // 文件夹已存在
		_, err := dir.ReadDir(1)
		if err == io.EOF {
			return
		}
		if err != nil {
			logger.Log.Fatal("判断为空文件夹失败", zap.String("path", name), zap.Error(err))
		}
		logger.Log.Fatal("文件夹不为空", zap.String("path", name))
	}
	if os.IsNotExist(err) { // 文件夹不存在
		err = os.MkdirAll(name, p)
		if err == nil {
			return
		}
		logger.Log.Fatal("创建文件夹失败", zap.String("path", name), zap.Error(err))
	}
	logger.Log.Fatal("创建文件夹失败", zap.String("path", name), zap.Error(err))
}

// 必须创建文件夹并且是创建
func MustMkdirAndIsCreate(name string, perm ...os.FileMode) {
	var p os.FileMode = 0666
	if len(perm) > 0 {
		p = perm[0]
	}
	dir, err := os.Open(name)
	if dir != nil {
		defer dir.Close()
	}
	if err == nil {
		logger.Log.Fatal("创建文件夹失败, 文件夹已存在", zap.String("path", name))
	}
	if !os.IsNotExist(err) {
		logger.Log.Fatal("创建文件夹失败", zap.String("path", name), zap.Error(err))
	}

	err = os.MkdirAll(name, p)
	if err != nil {
		logger.Log.Fatal("创建文件夹失败", zap.String("path", name), zap.Error(err))
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
		logger.Log.Fatal("创建文件失败", zap.Error(err))
	}
}

// 检查是否存在路径并指定为文件夹或文件, 读取失败或者权限不足会fatal
func CheckHasPath(path string, isDir bool) bool {
	d, err := os.Open(path)
	if d != nil {
		defer d.Close()
	}
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		logger.Log.Fatal("读取失败", zap.String("path", path), zap.Error(err))
	}
	s, err := d.Stat()
	if err != nil {
		logger.Log.Fatal("读取失败", zap.String("path", path), zap.Error(err))
	}
	return s.IsDir() == isDir
}

// 必须存在路径
func MustHasPath(path string, isDir bool) {
	d, err := os.Open(path)
	if d != nil {
		defer d.Close()
	}
	if err != nil {
		if os.IsNotExist(err) {
			logger.Log.Fatal("path不存在", zap.String("path", path))
		}
		logger.Log.Fatal("读取失败", zap.String("path", path), zap.Error(err))
	}
	s, err := d.Stat()
	if err != nil {
		logger.Log.Fatal("读取失败", zap.String("path", path), zap.Error(err))
	}
	if isDir && !s.IsDir() {
		logger.Log.Fatal("path不是一个目录", zap.String("path", path))
	}
	if !isDir && s.IsDir() {
		logger.Log.Fatal("path不是一个文件", zap.String("path", path))
	}
}

// 必须不存在某个路径
func MustNoExistPath(path string) {
	dir, err := os.Open(path)
	if dir != nil {
		defer dir.Close()
	}
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logger.Log.Fatal("读取失败", zap.String("path", path), zap.Error(err))
	}
	logger.Log.Fatal("路径已存在", zap.String("path", path))
}

// 必须在项目中, 依次向上查找并进入项目根目录, 返回项目名
func MustEnterProject() string {
	wd := MustGetWorkdir()
	for {
		pn := searchProjectName(wd)
		if pn != "" {
			// 进入工程目录
			if err := os.Chdir(wd); err != nil {
				logger.Log.Fatal("进入工程目录失败", zap.String("path", wd), zap.Error(err))
			}
			return pn
		}

		wd2 := filepath.Dir(wd)
		if wd2 == wd { // 已经是根目录了
			logger.Log.Fatal("无法读取项目名")
		}
		wd = wd2
	}
}

// 查找项目名
func searchProjectName(path string) string {
	file := MustDirJoin(path, ".crawler")
	if !CheckHasPath(file, false) {
		return ""
	}

	data, err := os.ReadFile(file)
	if err != nil {
		logger.Log.Fatal("读取项目文件失败", zap.String("file", file), zap.Error(err))
	}
	var projectName string
	for _, suf := range []string{"\n", "\r\n", "\r", ""} {
		projectName = ExtractMiddleText(string(data), "project=", suf, "", false)
		if projectName != "" {
			return projectName
		}
	}
	if projectName == "" {
		logger.Log.Fatal("读取项目文件失败", zap.String("file", file), zap.Error(fmt.Errorf("无法获取项目名")))
	}
	return projectName
}

/*
*提取中间文本

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
		logger.Log.Fatal("无法获取当前目录", zap.Error(err))
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
		logger.Log.Fatal("获取绝对路径失败", zap.String("path", path2), zap.Error(err))
	}
	return path
}

// 构建模板参数
func MakeTemplateArgs(projectName, env string) map[string]interface{} {
	numCpu := runtime.NumCPU()
	if numCpu < 1 {
		numCpu = 1
	}
	return map[string]interface{}{
		"project_name": projectName,      // 项目名
		"project_dir":  MustGetWorkdir(), // 项目路径
		"env":          env,              // 环境名
		"date":         time.Now().Format("2006-01-02"),
		"time":         time.Now().Format("15:04:05"),
		"date_time":    time.Now().Format("2006-01-02 15:04:05"),
		"num_cpu":      numCpu,
	}
}

// 渲染模板
func RenderTemplate(template string, args map[string]interface{}) string {
	out := template
	for k, v := range args {
		out = strings.ReplaceAll(out, "{@"+k+"}", zstr.GetString(v))
	}
	return out
}

// 必须获取指定路径的文件夹名
func MustGetDirName(path string) string {
	return filepath.Base(path)
}
