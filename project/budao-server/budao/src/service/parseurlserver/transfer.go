package parseurlserver

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/sumaig/glog"

	pb "twirprpc"
)

func transferBoboVideoURL(videoURL string) (headers map[string]string, postBody []byte) {
	headers = make(map[string]string, 0)
	field := make(map[string]string, 0)

	field["_aKey"] = "ANDROID"
	field["_dId"] = "Redmi+Note+4X"
	u5, _ := uuid.NewV4()
	u5Slice := strings.Split(u5.String(), "-")
	devid := strings.ToUpper(strings.Join(u5Slice, ""))
	field["_devid"] = devid[0:32]
	u6, _ := uuid.NewV4()
	u6Slice := strings.Split(u6.String(), "-")
	u6ID := strings.Join(u6Slice, "")
	imei := u6ID[:16]
	field["_imei"] = imei
	field["_lang"] = "zh_CN"
	field["_nId"] = "1"
	field["_pName"] = "tv.yixia.bobo"
	field["_pcId"] = "xiaomi_market"
	field["_pgLoad"] = "0"
	u2, _ := uuid.NewV4()
	requestID := u2.String()
	field["_reqId"] = requestID
	field["_reqNum"] = "0"
	field["_rt"] = "0"
	field["_t"] = strconv.FormatUint(uint64(time.Now().Unix()), 10)
	field["_uId"] = "0"
	u1, _ := uuid.NewV4()
	u1Slice := strings.Split(u1.String(), "-")
	udid := strings.ToUpper(strings.Join(u1Slice, ""))
	field["_udid"] = udid
	field["_vApp"] = "8403"
	field["_vName"] = "2.8.6"
	field["_vOs"] = "7.0"
	field["country"] = "CN"
	target := strings.Split(videoURL, "channel/")
	field["videoId"] = target[1]

	var fieldKeySlice []string
	for key := range field {
		fieldKeySlice = append(fieldKeySlice, key)
	}
	sort.Strings(fieldKeySlice)
	var readyString []string
	for i := 0; i < len(field); i++ {
		readyString = append(readyString, fieldKeySlice[i], field[fieldKeySlice[i]])
	}
	readyString = append(readyString, "Cc$nceR6qGg5^Pdv%4@C")
	fieldStr := strings.Join(readyString, "")
	glog.Debug("bobo, fieldStr:%s", fieldStr)
	h := md5.New()
	h.Write([]byte(fieldStr))
	fieldByte := h.Sum(nil)
	temp := hex.EncodeToString(fieldByte)
	temp = temp[2:22]
	glog.Debug("bobo, sign:%s", temp)
	field["_sign"] = temp

	glog.Debug("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	data := url.Values{}
	for key, value := range field {
		glog.Debug("key:%v", key)
		glog.Debug("value:%v", value)
		data.Add(key, value)
	}
	glog.Debug("data:%v", data)
	postBody = []byte(data.Encode())
	glog.Debug("##################################")
	glog.Debug("postbody:%v", postBody)
	glog.Debug("bodyString:%v", string(postBody))

	headers["User-Agent"] = "Dalvik/2.1.0 (Linux; U; Android 7.0; Redmi Note 4X MIUI/V9.5.1.0.NCFCNFA)"
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	return headers, postBody
}

func transferLiVideoURL(videoURL string, header *pb.Header) (headers map[string]string, field []byte) {
	headers = make(map[string]string, 0)

	strArr := strings.Split(videoURL, "_")
	data := url.Values{}
	data.Add("contId", strArr[1])
	field = []byte(data.Encode())

	headers["X-Channel-Code"] = "official"
	headers["X-Platform-Type"] = "1"
	headers["X-Serial-Num"] = strconv.FormatUint(uint64(time.Now().Unix()), 10)
	headers["Accept-Language"] = "zh-Hans-CN;q=1"
	headers["Accept"] = "*/*"
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["User-Agent"] = "LiVideoIOS / 4.2.4(iPhone;iOS11.3.1;Scale / 3.00)"
	headers["Cookie"] = generateCookie()
	headers["X-Client-Agent"] = header.GetDeviceInfo().GetOsVersion()
	headers["X-Client-Version"] = header.GetDeviceInfo().GetAppVersion()
	headers["X-Platform-Version"] = "11.3.1"

	return headers, field
}

func transferKuaishouVideoURL() (headers map[string]string) {
	header := make(map[string]string, 0)
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.62 Safari/537.36"
	header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	header["Accept-Language"] = "zh-CN,zh;q=0.9"

	return header
}

func transferWeiboVideoURL() (headers map[string]string) {
	header := make(map[string]string, 0)
	header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	header["Accept-Language"] = "zh-CN,zh;q=0.9"
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.62 Safari/537.36"

	return header
}

func transferHaokanVideoURL(videoURL string) (headers map[string]string, field []byte, target string) {
	glog.Debug("start to transferHaokanVideoURL")
	headers = make(map[string]string, 0)

	strArr := strings.Split(videoURL, "_")
	tempStr := strings.TrimSuffix(strArr[1], "\"}")
	detail := fmt.Sprintf("method=get&vid=%s", tempStr)

	data := url.Values{}
	data.Add("video/detail", detail)
	field = []byte(data.Encode())

	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Accept"] = "*/*"
	headers["User-Agent"] = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E302 haokan/3.8.0.11 (Baidu; P2 11.3.1)/1.3.11_3,01enohP/381d/26026263E1FE591501C1F0FDDF43B0EB44F0D13EAOMDHANFGQL/1"
	headers["Accept-Language"] = "zh-Hans-CN;q=1"
	headers["X-TurboNet-Info"] = "2.8.1763.15"

	u5, _ := uuid.NewV4()
	u5Slice := strings.Split(u5.String(), "-")
	imei := strings.ToUpper(strings.Join(u5Slice, ""))
	u6, _ := uuid.NewV4()
	u6Slice := strings.Split(u6.String(), "-")
	u6ID := strings.ToUpper(strings.Join(u6Slice, ""))
	u7, _ := uuid.NewV4()
	u7Slice := strings.Split(u7.String(), "-")
	u7ID := strings.ToUpper(strings.Join(u7Slice, ""))
	cuid := u6ID + u7ID[:19]
	u8, _ := uuid.NewV4()
	u8Slice := strings.Split(u8.String(), "-")
	hid := strings.ToUpper(strings.Join(u8Slice, ""))

	target = fmt.Sprintf("https://sv.baidu.com/haokan/api?imei=%s&cuid=%s&os=ios&osbranch=i0&apiv=3.8.0.11&appv=1&version=3.8.0.11&hid=%s", imei, cuid, hid)
	glog.Debug("end to transferHaokanVideoURL")

	return headers, field, target
}

func transferBzhanVideoURL(videoURL string) (headers map[string]string, target string) {
	tempA := strings.Split(videoURL, "av")
	tempB := strings.Split(tempA[1], "/")
	aid := tempB[0]

	glog.Debug("transferBzhanVideoURL***************************")
	glog.Debug("aid:%s", aid)
	headers = make(map[string]string, 0)
	field := make(map[string]string, 0)

	field["callback"] = "callbackfunction"
	field["aid"] = aid
	field["page"] = "1"
	field["platform"] = "html5"
	field["quality"] = "1"
	field["vtype"] = "mp4"
	field["type"] = "jsonp"

	timeStamp := time.Now().UnixNano() / 1000
	tempStr := "bilibili_" + strconv.FormatUint(uint64(timeStamp), 10)
	h := md5.New()
	h.Write([]byte(tempStr))
	tempByte := h.Sum(nil)
	token := hex.EncodeToString(tempByte)
	field["token"] = token
	field["is_preview"] = "1"

	target = ""
	for key, value := range field {
		target = target + key + "=" + value + "&"
	}
	target = strings.TrimSuffix(target, "&")
	glog.Debug("target:%s", target)

	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.62 Safari/537.36"
	headers["Referer"] = "http://www.bilibili.com"
	purlToken := "bilibili_" + strconv.FormatUint(uint64(timeStamp), 10)

	u1, _ := uuid.NewV4()
	u1Slice := strings.Split(u1.String(), "-")
	sign := strings.Join(u1Slice, "")
	rand.Seed(time.Now().Unix())
	uid := strconv.FormatUint(uint64(20000+rand.Intn(9999)), 10)
	buvid3 := sign + uid + "infoc"
	cookie := fmt.Sprintf("purl_token=%s; buvid3=%s", purlToken, buvid3)
	headers["cookie"] = cookie
	glog.Debug(headers)

	return headers, target
}

func transferMiaopaiVideoURL() (headers map[string]string, target string) {
	glog.Debug("1###############################")
	headers = make(map[string]string, 0)
	field := make(map[string]string, 0)

	rand.Seed(time.Now().Unix())
	u1, _ := uuid.NewV4()
	u1Slice := strings.Split(u1.String(), "-")
	sign := strings.Join(u1Slice, "")
	field["sign"] = sign
	field["time_stamp"] = strconv.FormatUint(uint64(time.Now().Unix()), 10)
	field["mpflag"] = "16"
	field["refer_pg"] = "118"
	field["pg"] = "118"
	field["abld"] = strconv.FormatUint(uint64(78+rand.Intn(25)), 10)
	field["appName"] = "美拍"
	field["brand"] = "iphone"
	car := []string{"中国移动", "中国联通", "中国电信"}
	n := rand.Intn(2)
	field["carrier"] = car[n]
	field["cpu"] = "CPU_TYPE_ARM64"
	field["density"] = ""
	u0, _ := uuid.NewV4()
	field["devId"] = u0.String()
	field["dpi"] = "750x1624"
	field["facturer"] = "iphone"

	u2, _ := uuid.NewV4()
	u2Slice := strings.Split(u2.String(), "-")
	ikg_udid := strings.Join(u2Slice, "")
	field["idfa"] = u2.String()
	field["imei"] = ""
	field["ikg_udid"] = ikg_udid
	field["model"] = "Unknown iPhone"
	field["net"] = "1"
	field["os"] = "ios"
	field["pName"] = "com.yixia.iphone"
	field["partnerId"] = "1"
	field["pcId"] = "AppStore"
	field["plat"] = "iphone"
	field["platformId"] = "1"
	field["resolution"] = "750x1624"
	field["sys_version"] = "iOS11.3.1"
	field["timestamp"] = strconv.FormatUint(uint64(time.Now().Unix()), 10)

	u3, _ := uuid.NewV4()
	u3Slice := strings.Split(u3.String(), "-")
	udid := strings.Join(u3Slice, "")
	field["udid"] = udid
	rint9 := rand.Int31() >> 23
	field["unique_id"] = udid + strconv.FormatUint(uint64(rint9), 10)
	field["userId"] = ""
	field["vApp"] = "6670"
	field["vName"] = "6.7.66"
	field["vOs"] = "11.3.1"
	field["vend"] = "miaopai"
	field["version"] = "6.7.66"

	target = ""
	for key, value := range field {
		target = target + key + "=" + value + "&"
	}
	target = strings.TrimSuffix(target, "&")
	glog.Debug("2###############################")

	u6, _ := uuid.NewV4()
	u6slice := strings.Split(u6.String(), "-")
	kg_udid := strings.Join(u6slice, "")
	u7, _ := uuid.NewV4()
	u7Slice := strings.Split(u7.String(), "-")
	sessionId := strings.Join(u7Slice, "")
	u8, _ := uuid.NewV4()
	u8Slice := strings.Split(u8.String(), "-")
	newudid := strings.Join(u8Slice, "")
	aliyungf_tc := randstr(12) + "+" + randstr(19)
	cookie := fmt.Sprintf("kg_udid=%s; sessionId=%s; udid=%s; aliyungf_tc=%s", kg_udid, sessionId, newudid, aliyungf_tc)
	headers["Cookie"] = cookie

	headers["vend"] = "miaopai"
	headers["Accept"] = "*/*"
	headers["User-Agent"] = "MiaoPai/6.7.66 (iPhone; iOS 11.3.1; Scale/3.00)"
	u5, _ := uuid.NewV4()
	u5Slice := strings.Split(u5.String(), "-")
	newsign := strings.Join(u5Slice, "")
	headers["sign"] = newsign
	headers["Accept-Language"] = "zh-Hans"
	headers["appid"] = "442"
	glog.Debug("3###############################")
	return headers, target
}

func generateCookie() string {
	u1, _ := uuid.NewV4()
	pearUUID := strings.ToUpper(u1.String())

	u2, _ := uuid.NewV4()
	u2Slice := strings.Split(u2.String(), "-")
	jsessionid := strings.ToUpper(strings.Join(u2Slice, ""))

	cookie := fmt.Sprintf("JSESSIONID=%s;PEAR_UUID=%s", jsessionid, pearUUID)

	return cookie
}

func extractCharacter(str string) string {
	i := strings.LastIndex(str, "show/")
	j := strings.LastIndex(str, ".htm")
	target := str[i+5 : j]

	return target
}

func randstr(l int) string {
	var inibyte []byte
	var result bytes.Buffer
	for i := 48; i < 123; i++ {
		switch {
		case i < 58:
			inibyte = append(inibyte, byte(i))
		case i >= 65 && i < 91:
			inibyte = append(inibyte, byte(i))
		case i >= 97 && i < 123:
			inibyte = append(inibyte, byte(i))
		}

	}
	var temp byte
	for i := 0; i < l; {
		if inibyte[randInt(0, len(inibyte))] != temp {
			temp = inibyte[randInt(0, len(inibyte))]
			result.WriteByte(temp)
			i++
		}

	}
	return result.String()
}

func randInt(min int, max int) byte {
	rand.Seed(time.Now().UnixNano())
	return byte(min + rand.Intn(max-min))
}
