package goHelp

import (
	"github.com/shima-park/agollo"
)

var apollo agollo.Agollo

func InitApollo(configServerURL, appID string) {

	var err error
	apollo, err = agollo.New(
		configServerURL,
		appID,
		agollo.AutoFetchOnCacheMiss(),
		agollo.FailTolerantOnBackupExists(),
	)
	if err != nil {
		panic("apollo 初始化失败")
	}
}
