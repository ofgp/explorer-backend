/*
作者：雲凌禹
链接：https://www.jianshu.com/p/9c1c8958b279
*/

package chainapi

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"

	"github.com/astaxie/beego"
)

var AuthTokenIv = "rst@123456--java"

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(AuthTokenIv))
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, []byte(AuthTokenIv))
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func DoAesEncrypt(data map[string]interface{}, keyStr string) (string, error) {
	key := []byte(keyStr)
	params, err := json.Marshal(data)
	if err != nil {
		beego.Error(err)
		return "", err
	}
	result, err := AesEncrypt(params, key)
	if err != nil {
		beego.Error(err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

func DoAesDecrypt(data string, keyStr string) ([]byte, error) {
	key := []byte(keyStr)
	params, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	result, err := AesDecrypt(params, key)
	return result, err
}
