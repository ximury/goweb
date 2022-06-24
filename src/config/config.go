package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"logger"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var configFile []byte

type Config struct {
	Log struct {
		EnableConsole     bool   `yaml:"enableConsole"`
		ConsoleLevel      string `yaml:"consoleLevel"`
		ConsoleJSONFormat bool   `yaml:"consoleJSONFormat"`
		EnableFile        bool   `yaml:"enableFile"`
		FileJSONFormat    bool   `yaml:"fileJSONFormat"`
		FileLevel         string `yaml:"fileLevel"`
		FileLocation      string `yaml:"fileLocation"`
		MaxAge            int    `yaml:"maxAge"`
		MaxSize           int    `yaml:"maxSize"`
		Compress          bool   `yaml:"compress"`
		FileExport        string `yaml:"fileExport"`
	}

	Webapi struct {
		Uri string `yaml:"uri"`
	}

	Mysqlnd struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
	}
}

func init() {
	var err error
	var configFilePath = filepath.Join(getCurrentAbPathByCaller(), "config.yaml")
	configFile, err = ioutil.ReadFile(configFilePath)
	if err != nil {
		logger.Fatalf("Read config yaml file err %v", err)
	}
}

func GetChannelConfig() (e *Config, err error) {
	err = yaml.Unmarshal(configFile, &e)
	return e, err
}

// 获取程序运行路径（go build）
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Errorf("Get current path err %v", err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	// os.Stat获取文件信息
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
