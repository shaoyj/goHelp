package goHelp

import "github.com/shima-park/agollo"

//Init apollo
func InitApollo(configServerURL, appID string) {
	agollo.Init(
		configServerURL,
		appID,
		func(options *agollo.Options) {
			options.FailTolerantOnBackupExists = true
		},
	)
}
