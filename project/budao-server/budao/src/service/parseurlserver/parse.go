package parseurlserver

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	pb "twirprpc"

	simplejson "github.com/bitly/go-simplejson"
	uuid "github.com/satori/go.uuid"
	"github.com/sumaig/glog"
)

func parseKuaishou(body string) (playURL string, header map[string]string, next bool) {
	hasFuncJump := strings.Contains(body, "function jump()")
	if hasFuncJump == false {
		a := strings.Split(body, "video src=\"")
		b := strings.Split(a[1], "\" autoplay")
		glog.Debug(b[0])
		playURL = b[0]

		header = make(map[string]string, 0)
		header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.62 Safari/537.36"

		next = false

		glog.Debug("kuaishou playurl:%s", playURL)
	} else {
		glog.Debug("^^^^^^^^^^^^^^^kuaishou need to parse again")
		a := strings.Split(body, "'cookie' : ")
		b := strings.Split(a[1], "'uri' :")
		c := strings.TrimSpace(b[0])
		c = strings.TrimSuffix(c, ",")
		cookie := strings.Trim(c, "\"")
		cookie = "ksbv_sign_javascript=(" + cookie + ")"

		m := strings.Split(body, "'uri' : ")
		n := strings.Split(m[1], "}")
		s := strings.TrimSpace(n[0])
		s = strings.TrimSuffix(s, ",")
		uri := strings.Trim(s, "\"")

		hasksjssig := strings.Contains(uri, "ksjs_sig")
		if hasksjssig == false {
			timeStamp := time.Now().UnixNano() / 1000
			temp := "?ksjs_sig=" + strconv.FormatUint(uint64(timeStamp), 10)
			uri = uri + temp
		}
		playURL = uri
		glog.Debug("kuaishou uri:%s", playURL)

		header = make(map[string]string, 0)
		header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.62 Safari/537.36"
		header["cookie"] = cookie

		next = true
	}

	return
}

func parseWeibo(body string) (playURL string, header map[string]string) {
	header = make(map[string]string, 0)
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.62 Safari/537.36"
	header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	header["Accept-Language"] = "zh-CN,zh;q=0.9"

	a := strings.Split(body, "\"media_info\":")
	b := strings.Split(a[1], "}")
	c := strings.TrimSpace(b[0])
	e := strings.TrimLeft(c, "{")
	m := strings.TrimSpace(e)

	stream1 := strings.Split(m, "\"stream_url_hd\":")
	stream1[0] = strings.TrimSpace(stream1[0])
	stream1[0] = strings.TrimSuffix(stream1[0], "\",")
	stream := strings.Split(stream1[0], "\":")
	stream[1] = strings.TrimSpace(stream[1])
	streamURL := strings.TrimPrefix(stream[1], "\"")

	stream1[1] = strings.TrimSpace(stream1[1])
	exist := strings.Contains(stream1[1], "\"duration\"")
	var streamHD string
	if exist == true {
		target := strings.Split(stream1[1], "\"duration\"")
		target[0] = strings.TrimSpace(target[0])
		target[0] = strings.TrimSuffix(target[0], ",")
		if target[0] != "null" {
			target[0] = strings.TrimSuffix(target[0], "\"")
			target[0] = strings.TrimPrefix(target[0], "\"")
			streamHD = target[0]
		}
	} else {
		target := strings.Trim(stream1[1], "\"")
		if target != "null" {
			streamHD = target
		}
	}

	if streamHD != "" {
		playURL = streamHD
	} else {
		playURL = streamURL
	}

	glog.Debug("playurl:%s", playURL)

	return
}

func parseBzhan(body string) (videoURL string, headers map[string]string) {
	tempA := strings.Split(body, "(")
	tempB := strings.Split(tempA[1], ")")
	target := tempB[0]

	js, err := simplejson.NewJson([]byte(target))
	if err != nil {
		glog.Error("simplejson.NewJson failed. err:%v", err)
		return "", headers
	}

	videoArr, err := js.Get("durl").Array()
	for i := range videoArr {
		video := js.Get("durl").GetIndex(i)
		videoURL = video.Get("url").MustString()
		if videoURL != "" {
			break
		}
	}

	headers = make(map[string]string, 0)
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.62 Safari/537.36"
	headers["Referer"] = "http://www.bilibili.com"
	timeStamp := time.Now().UnixNano() / 1000
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

	return videoURL, headers
}

func parseLiVideoBody(body string, header *pb.Header) (videoURL string, headers map[string]string) {
	headers = make(map[string]string, 0)
	headers["X-Channel-Code"] = "official"
	headers["X-Platform-Type"] = "1"
	headers["X-Serial-Num"] = strconv.FormatUint(uint64(time.Now().Unix()), 10)
	headers["Accept-Language"] = "zh-Hans-CN;q=1"
	headers["Accept"] = "*/*"
	headers["User-Agent"] = "LiVideoIOS / 4.2.4(iPhone;iOS11.3.1;Scale / 3.00)"
	headers["Cookie"] = generateCookie()
	headers["X-Client-Agent"] = header.GetDeviceInfo().GetOsVersion()
	headers["X-Client-Version"] = header.GetDeviceInfo().GetAppVersion()
	headers["X-Platform-Version"] = "11.3.1"

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		glog.Error("simplejson.NewJson failed. err:%v", err)
		return "", headers
	}

	resultMsg, err := js.Get("resultMsg").String()
	if resultMsg != "success" {
		glog.Error("parseLiVideoBody client spider failed.")
		return "", headers
	}

	tempMap := make(map[string]string, 0)
	videoArr, err := js.Get("content").Get("videos").Array()
	for i := range videoArr {
		video := js.Get("content").Get("videos").GetIndex(i)
		tempMap[video.Get("tag").MustString()] = video.Get("url").MustString()
	}

	if value, ok := tempMap["fhd"]; ok {
		return value, headers
	}
	if value, ok := tempMap["hd"]; ok {
		return value, headers
	}
	if value, ok := tempMap["ld"]; ok {
		return value, headers
	}
	if value, ok := tempMap["sd"]; ok {
		return value, headers
	}

	return "", headers
}

func parseHaokanVideoBody(body string) (videoURL string, headers map[string]string) {
	glog.Debug("parseHaokanVideoBody, body:", body)

	headers = make(map[string]string, 0)
	headers["Accept"] = "*/*"
	headers["User-Agent"] = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E302 haokan/3.8.0.11 (Baidu; P2 11.3.1)/1.3.11_3,01enohP/381d/26026263E1FE591501C1F0FDDF43B0EB44F0D13EAOMDHANFGQL/1"
	headers["Accept-Language"] = "zh-Hans-CN;q=1"
	headers["X-TurboNet-Info"] = "2.8.1763.15"

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		glog.Error("simplejson.NewJson failed. err:%v", err)
		return "", headers
	}

	status := js.Get("video/detail").Get("status").MustInt()
	if status != 0 {
		glog.Error("parseHaokanVideoBody client spider failed.")
		return "", headers
	}

	sc := js.Get("video/detail").Get("data").Get("video_list").Get("sc").MustString()
	if sc != "" {
		return sc, headers
	}

	hd := js.Get("video/detail").Get("data").Get("video_list").Get("hd").MustString()
	if hd != "" {
		return hd, headers
	}

	sd := js.Get("video/detail").Get("data").Get("video_list").Get("sd").MustString()
	if sd != "" {
		return sd, headers
	}

	return "", headers
}

func parseBoboVideoBody(body string) (videoURL string, headers map[string]string) {
	glog.Debug("bobo*************")
	glog.Debug(body)

	headers = make(map[string]string, 0)
	headers["User-Agent"] = "Dalvik/2.1.0 (Linux; U; Android 7.0; Redmi Note 4X MIUI/V9.5.1.0.NCFCNFA)"

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		glog.Error("simplejson.NewJson failed. err:%v", err)
		return "", headers
	}

	msg := js.Get("msg").MustString()
	if msg != "ok" {
		glog.Error("parseBoboVideoBody client spider failed.")
		return "", headers
	}

	var url string
	var url2 string
	videoArr, err := js.Get("data").Get("videoUrl").Array()
	for i := range videoArr {
		video := js.Get("data").Get("videoUrl").GetIndex(i)
		url = video.Get("url").MustString()
		url2 = video.Get("url2").MustString()
	}

	glog.Debug("url:", url)
	glog.Debug("url2:", url2)
	if url != "" {
		return url, headers
	}
	if url2 != "" {
		return url2, headers
	}

	return "", headers
}

func parseMiaopaiVideoBody(body string) (videoURL string, headers map[string]string) {
	headers = make(map[string]string, 0)
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

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		glog.Error("simplejson.NewJson failed. err:%v", err)
		return "", headers
	}

	status, err := js.Get("status").Int()
	if status != 200 {
		glog.Error("parseMiaopaiVideoBody client spider failed.")
		return "", headers
	}

	result := js.Get("result").GetIndex(0)
	scheme := result.Get("scheme").MustString()
	host := result.Get("host").MustString()
	path := result.Get("path").MustString()
	videoURL = scheme + host + path
	glog.Debug("lid")
	glog.Debug(videoURL)

	return videoURL, headers
}
