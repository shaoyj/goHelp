package goHelp

import (
	"testing"
	"fmt"
)

func TestGetSqlConfig(t *testing.T) {
	InitApollo("127.0.0.1:18080", "80000")
	InitNamespaceEncryptAndSecret("syj.mysql", "syj.secret")
	result := GetMySqlConfigFromInstrumentation("tk-syj-5290")
	fmt.Println(result)
	fmt.Println(GetSqlTcpStyleFromInstrumentation("tk-syj-5290"))
}
