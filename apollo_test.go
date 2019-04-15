package goHelp

import (
	"fmt"
	"github.com/shima-park/agollo"
	"testing"
)

func TestInit(t *testing.T) {
	agollo.Init(
		"ip:8088",
		"appId",
		func(options *agollo.Options) {
			//options.DefaultNamespace="application"

		},
	)

	// 获取默认配置中cluster=default namespace=application key=Name的值
	fmt.Println("Name:", agollo.Get("Name"))

	// 获取默认配置中cluster=default namespace=application key=Name的值，提供默认值返回
	fmt.Println("gray_scale_merchant:", agollo.Get("gray_scale_merchant", agollo.WithDefault("YourDefaultValue")))
	syj := agollo.Get("gray_scale_merchant")
	fmt.Println(len(syj))
	if syj == "" {
		fmt.Println("true")
	}

	// 获取默认配置中cluster=default namespace=Test.Namespace key=Name的值，提供默认值返回
	fmt.Println("YourConfigKey2:", agollo.Get("YourConfigKey2", agollo.WithDefault("YourDefaultValue"), agollo.WithNamespace("YourNamespace")))

	// 获取namespace下的所有配置项
	fmt.Println("Configuration of the namespace:", agollo.GetNameSpace("application"))

	// TEST.Namespace1是非预加载的namespace
	// agollo初始化是带上agollo.AutoFetchOnCacheMiss()可选项的话
	// 陪到非预加载的namespace，会去apollo缓存接口获取配置
	// 未配置的话会返回空或者传入的默认值选项
	fmt.Println(agollo.Get("Name", agollo.WithDefault("foo"), agollo.WithNamespace("TEST.Namespace1")))

	// 如果想监听并同步服务器配置变化，启动apollo长轮训
	// 返回一个期间发生错误的error channel,按照需要去处理
	errorCh := agollo.Start()

	// 监听apollo配置更改事件
	// 返回namespace和其变化前后的配置,以及可能出现的error
	watchCh := agollo.Watch()

	go func() {
		for {
			select {
			case err := <-errorCh:
				fmt.Println("Error:", err)
			case update := <-watchCh:
				fmt.Println("Apollo Update:", update)
			}
		}
	}()

	select {}

}
