package goHelp

import (
	"github.com/shima-park/agollo"
)

var GoHelpApollo agollo.Agollo

//初始化
func InitApollo(configServerURL, appID string) {

	var err error
	GoHelpApollo, err = agollo.New(
		configServerURL,
		appID,
		agollo.AutoFetchOnCacheMiss(),
		agollo.FailTolerantOnBackupExists(),
	)
	if err != nil {
		panic("apollo 初始化失败")
	}
}

//初始化，不生成临时文件
func InitApolloNoFile(configServerURL, appID string) {

	var err error
	GoHelpApollo, err = agollo.New(
		configServerURL,
		appID,
		agollo.AutoFetchOnCacheMiss(),
		agollo.BackupFile("/dev/null"),
	)
	if err != nil {
		panic("apollo 初始化失败")
	}
}
