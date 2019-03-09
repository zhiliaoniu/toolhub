package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/sumaig/glog"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
	pb "twirprpc"
)

/**
 * 获取md5字符串
 */
func GetMd5Str(str []byte) string {
	s := string(str) + "jxz_budao"
	hmd := md5.New()
	//hmd.Write([]byte(str))
	hmd.Write([]byte(s))
	md5Str := hmd.Sum(nil)

	return hex.EncodeToString(md5Str)
}

/**
 * sha1加密串
 */
func GetSha1Str(secret string, str string) []byte {

	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(str))

	return mac.Sum(nil)
}

/**
 * 获取请求响应的初始状态
 */
func GetInitStatus() (status *pb.Status) {
	return &pb.Status{
		Code:       pb.Status_SERVER_ERR,
		Message:    "success",
		ServerTime: uint64(time.Now().Unix()),
	}
}

/**
 * sha256加密串
 */
func GetSha256Str(secret string, str string) []byte {

	key := []byte(secret)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(str))

	return mac.Sum(nil)
}

/**
 * 获取格林时间
 * Mon, 02 Jan 2006 07:04:05
 * @return string
 */
func Gmtime() string {
	local, _ := time.LoadLocation("PRC")
	timeFormat := "Mon, 02 Jan 2006 07:04:05 GMT"

	return time.Now().In(local).Format(timeFormat)
}

/**
 * t 为时间字符串，时间字符串为utc格式 eg:"2016-08-22T08:23:42Z"
 * @return string
 */
func GetTimeStr(t string) string {
	res, _ := time.Parse(time.RFC3339, t)

	return res.Format("2006-01-02 15:04:05")
}

/**
 * 下载文件
 * fileFullPath  文件路径
 */
func DownloadFile(fileFullPath string, res http.ResponseWriter) {
	glog.Info("downloadFile:" + fileFullPath)

	file, err := os.Open(fileFullPath)
	if err != nil {
		glog.Info(fmt.Sprintf("downloadFile open file err: %v", err))
		return
	}

	defer file.Close()
	fileName := path.Base(fileFullPath)
	fileName = url.QueryEscape(fileName) // 防止中文乱码
	res.Header().Set("Content-Type", "application/octet-stream")
	res.Header().Set("content-disposition", "attachment; filename=\""+fileName+"\"")
	_, err1 := io.Copy(res, file)
	if err1 != nil {
		glog.Info(fmt.Sprintf("downloadFile io.copy err: %v", err1))
		return
	}
}

//get方式获取信息
func HttpGet(url string) (string, []byte) {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return fmt.Sprintf("respose http get error %s", err), nil
	}
	buff, err := ioutil.ReadAll(response.Body)

	return "", buff
}

/**
 * AES加密
 *
 * @param  origData 加密的数据
 * @param  key      加密需要的key  key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
 */
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return crypted, nil
}

/**
 * AES解密
 *
 * @param  crypted	已加密的数据
 * @param  key      解密需要的key
 */
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)

	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
