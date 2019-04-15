package goHelp

import (
	"fmt"
	"strings"
	"testing"
)

func TestDataFormat(t *testing.T) {
	fmt.Println(123)
	NewZipkinTracerFromKafka(strings.Split("10.46.67.89:9092", ","), false,
		"", "gohelp")

	CommitZeus("test_gohelp")
}

func TestEquals(t *testing.T) {
	flag := strings.EqualFold(GetEnv(), "FAT")
	fmt.Println(flag)
}
