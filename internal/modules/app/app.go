package app

import (
	"os"
	"path/filepath"

	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/setting"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"github.com/ouqiang/goutil"
)

var (
	// AppDir 应用根目录
	AppDir string // 应用根目录
	// ConfDir 配置文件目录
	ConfDir string // 配置目录
	// LogDir 日志目录
	LogDir string // 日志目录
	// AppConfig 配置文件
	AppConfig string // 应用配置文件
	// Installed 应用是否已安装
	Installed bool // 应用是否安装过
	// Setting 应用配置
	Setting *setting.Setting // 应用配置
	// VersionId 版本号
	VersionId int // 版本号
	// VersionFile 版本文件
	VersionFile string // 版本号文件
)

// InitEnv 初始化
func InitEnv(versionString string) {
	logger.InitLogger()
	var err error
	AppDir, err = goutil.WorkDir()
	if err != nil {
		logger.Fatal(err)
	}
	ConfDir = filepath.Join(AppDir, "/conf")
	LogDir = filepath.Join(AppDir, "/log")
	AppConfig = filepath.Join(ConfDir, "/app.ini")
	VersionFile = filepath.Join(ConfDir, "/.version")
	createDirIfNotExists(ConfDir, LogDir)
	Installed = IsInstalled()
	VersionId = ToNumberVersion(versionString)
}

// IsInstalled 判断应用是否已安装
func IsInstalled() bool {
	_, err := os.Stat(filepath.Join(ConfDir, "/install.lock"))
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// CreateInstallLock 创建安装锁文件
func CreateInstallLock() error {
	_, err := os.Create(filepath.Join(ConfDir, "/install.lock"))
	if err != nil {
		logger.Error("创建安装锁文件conf/install.lock失败")
	}

	return err
}

// UpdateVersionFile 更新应用版本号文件
func UpdateVersionFile() {
	err := ioutil.WriteFile(VersionFile,
		[]byte(strconv.Itoa(VersionId)),
		0644,
	)

	if err != nil {
		logger.Fatal(err)
	}
}

// GetCurrentVersionId 获取应用当前版本号, 从版本号文件中读取
func GetCurrentVersionId() int {
	if !utils.FileExist(VersionFile) {
		return 0
	}

	bytes, err := ioutil.ReadFile(VersionFile)
	if err != nil {
		logger.Fatal(err)
	}

	versionId, err := strconv.Atoi(strings.TrimSpace(string(bytes)))
	if err != nil {
		logger.Fatal(err)
	}

	return versionId
}

// ToNumberVersion 把字符串版本号 a.b.c[-...] 转换为整数版本号 abc
func ToNumberVersion(versionString string) int {
	// 去掉前缀 v
	versionString = strings.TrimPrefix(versionString, "v")

	// 拆掉可能的后缀, 比如 -beta1
	if idx := strings.Index(versionString, "-"); idx != -1 {
		versionString = versionString[:idx]
	}

	// 按点分割
	segments := strings.Split(versionString, ".")

	// 补齐到3段
	for len(segments) < 3 {
		segments = append(segments, "0")
	}

	// 只保留前三段
	segments = segments[:3]

	// 拼接为 abc
	v := strings.Join(segments, "")

	versionId, err := strconv.Atoi(v)
	if err != nil {
		logger.Fatal(err)
	}
	return versionId
}

// 检测目录是否存在
func createDirIfNotExists(path ...string) {
	for _, value := range path {
		if utils.FileExist(value) {
			continue
		}
		err := os.Mkdir(value, 0755)
		if err != nil {
			logger.Fatal(fmt.Sprintf("创建目录失败:%s", err.Error()))
		}
	}
}
