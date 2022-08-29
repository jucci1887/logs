/*
 Author: Kernel.Huang
 Mail: kernelman79@gmail.com
 Date: 3/18/21 1:01 PM
*/
package logs

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 获取当前执行程序的绝对目录路径
func GetCurrentDir() string {
	currentPath := CurrentAndAbsPath()
	return filepath.Dir(currentPath)
}

// 当前执行程序的绝对路径
func CurrentAndAbsPath() string {
	current := SetCurrentPath()
	return GetAbsPath(current)
}

// 设置当前执行程序的绝对路径
func SetCurrentPath() string {
	current := os.Args[0]
	path, err := exec.LookPath(current)
	if err != nil {
		log.Println("Set the current path error: ", err)
	}

	return path
}

// 获取当前执行程序的绝对路径
func GetAbsPath(current string) string {
	absPath, err := filepath.Abs(current)
	if err != nil {
		log.Println("Get the current absolute of path error: ", err)
	}

	return absPath
}

// 获取日志文件名
func GetLogsFilename() string {
	content := GetToml()
	return content.Zone("log").Fetch("name").ToStr()
}

// 获取日志文件内容前缀
func GetLogsPrefix() string {
	content := GetToml()
	return content.Zone("log").Fetch("prefix").ToStr()
}

// 获取日志级别, 值为OFF则关闭日志
func GetLogsLevel() string {
	content := GetToml()
	return content.Zone("log").Fetch("level").ToStr()
}

// 获取配置目录名
func GetConfigDir() string {
	return "config"
}

// 获取日志配置名
func GetConfigPath() string {
	return "logs.toml"
}

// 获取Toml配置解析服务
func GetToml() *TomlConfig {
	configDir := GetConfigDir()
	configPath := GetConfigPath()
	return Toml.NewToml(configDir, configPath)
}

// 获取日志目录
func GetLogsDir() string {
	rootPath := GetRootPath()
	content := GetToml()
	relative := content.Zone("log").Fetch("relative").ToBool()
	logDir := content.Zone("log").Fetch("dir").ToStr()

	if relative {
		return filepath.Join(rootPath, logDir, string(os.PathSeparator))
	}

	return logDir
}

// 获取路径的上个目录
func GetLastPath(currentPath string) string {
	index := strings.LastIndex(currentPath, string(os.PathSeparator))
	return currentPath[:index]
}

// 获取项目根目录
func GetRootPath() string {
	dir := GetCurrentDir()
	rootPath := GetLastPath(dir)
	return filepath.Join(rootPath, string(os.PathSeparator))
}

// Get config dir of custom
func GetCustomConfigDir(dirname string) string {
	rootPath := GetRootPath()
	return filepath.Join(rootPath, dirname, string(os.PathSeparator))
}

// Get config path
func GetCustomConfigPath(dirname string, filename string) string {
	configDir := GetCustomConfigDir(dirname)
	return filepath.Join(configDir, filename)
}
