package goHelp

import (
	"strings"
	"github.com/shima-park/agollo"
	"encoding/hex"
	"crypto/aes"
	"crypto/sha1"
	"bytes"
	"errors"
)

var (
	namespaceEncrypt, namespaceSecret string
	DecryptError                      = errors.New("decrypt error")
	InvalidConfig                     = errors.New("decode error, check config please")
)

type MySQLConfig struct {
	Url      string `json:"url"`
	Database string `json:"data_base"`
	Params   string `json:"params"`
	Username string `json:"user_name"`
	Password string `json:"password"`
}

//
func InitNamespaceEncryptAndSecret(encrypt, secret string) {
	namespaceEncrypt, namespaceSecret = encrypt, secret
}

//获取 pSplit 和 uSplit
func getPSplitAndUSplit(url string) (string, *[]string, *[]string) {
	if namespaceEncrypt == "" || namespaceSecret == "" {
		panic("should exec func InitNamespaceEncryptAndSecret() ---> init namespaceEncrypt,namespaceSecret")
	}

	params := "useUnicode=yes&characterEncoding=UTF-8&zeroDateTimeBehavior=convertToNull"
	if strings.Contains(url, "?") {
		us := strings.Split(url, "?")
		url = us[0]
		params += "&" + us[1]
	}

	plain := decrypt(url)
	if plain == "" {
		panic("decode error, check config please")
		return "", nil, nil
	}

	pSplit := strings.Split(plain, "|")
	if len(pSplit) != 3 {
		panic("decode error, check config please. len(pSplit) != 3")
	}
	uSplit := strings.Split(pSplit[0], "/")

	return params, &pSplit, &uSplit;
}

//
func GetMySqlConfigFromInstrumentation(url string) *MySQLConfig {
	params, pSplit, uSplit := getPSplitAndUSplit(url)
	return &MySQLConfig{
		Url:      (*uSplit)[0],
		Database: (*uSplit)[1],
		Params:   params,
		Username: (*pSplit)[1],
		Password: (*pSplit)[2],
	}
}

//返回结果格式: username:password@Wc7++@tcp(127.0.0.1:3306)/database?params
func GetSqlTcpStyleFromInstrumentation(url string) string {
	params, pSplit, uSplit := getPSplitAndUSplit(url)
	var buffer bytes.Buffer
	buffer.WriteString((*pSplit)[1])
	buffer.WriteString(":")
	buffer.WriteString((*pSplit)[2])
	buffer.WriteString("@tcp(")
	buffer.WriteString((*uSplit)[0])
	buffer.WriteString(")/")
	buffer.WriteString((*uSplit)[1])
	buffer.WriteString("?")
	buffer.WriteString(params)
	return buffer.String()
}

//返回结果格式: username:password@Wc7++@tcp(127.0.0.1:3306)/database
func GetSqlTcpStyleFromInstrumentationNoParams(url string) string {
	_, pSplit, uSplit := getPSplitAndUSplit(url)
	var buffer bytes.Buffer
	buffer.WriteString((*pSplit)[1])
	buffer.WriteString(":")
	buffer.WriteString((*pSplit)[2])
	buffer.WriteString("@tcp(")
	buffer.WriteString((*uSplit)[0])
	buffer.WriteString(")/")
	buffer.WriteString((*uSplit)[1])
	return buffer.String()
}

func decrypt(key string) string {
	c := getCipherText(key)
	secretKey := getSecretKey()
	ciphertext, _ := hex.DecodeString(c)
	plain := aesEcbDecrypt(ciphertext, sha1PRNG(secretKey))
	return string(plain)
}

func getCipherText(key string) string {
	return GoHelpApollo.Get(key,
		agollo.WithDefault(""),
		agollo.WithNamespace(namespaceEncrypt))
}

func getSecretKey() string {
	return GoHelpApollo.Get("database.decode.key",
		agollo.WithDefault(""),
		agollo.WithNamespace(namespaceSecret))
}

func sha1PRNG(seed string) []byte {
	h := sha1.New()
	h.Write([]byte(seed))
	bs := h.Sum(nil)

	h = sha1.New()
	h.Write([]byte(bs))
	bs = h.Sum(nil)
	return bs[:16]
}

// JDK 默认AES为AES/ECB/PKCS5Padding
// 解密时NoPadding也可以
func aesEcbDecrypt(ciphertext, key []byte) (content []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	blockMode := NewECBDecrypter(block)
	content = make([]byte, len(ciphertext))
	blockMode.CryptBlocks(content, ciphertext)
	content = PKCS5NoPadding(content)
	return
}

func PKCS5NoPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
