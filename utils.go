package goHelp

import (
	log "github.com/sirupsen/logrus"
	"github.com/tinyhubs/properties"
	"os"
	"unicode"
)

var (
	env      string
	hostname string
)

//Get env
func GetEnv() string {
	if env != "" {
		return env
	}

	//获取环境变量
	env = "dev"
	envFilePath := os.Getenv("ENV_FILE_PATH")
	file, err := os.OpenFile(envFilePath, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		env = "dev"
	} else {
		doc, err := properties.Load(file)
		if nil != err {
			log.Error("加载失败")
		} else {
			env = doc.String("env")
			hostname = doc.String("hostname")
		}
	}

	return env
}

//Get hostName
func GetHostname() string {
	if hostname == "" {
		GetEnv()
		if hostname == "" {
			hostname = "unKnow_hostname"
		}
	}
	return hostname
}

//Get project dir
func GetProjectDir(projectPath string) string {
	//获取配置文件
	return os.Getenv("GOPATH") + projectPath
}

//校验是否是中文
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}
