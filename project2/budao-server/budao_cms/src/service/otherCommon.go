package service

import (
	"bytes"
	"common"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sumaig/glog"
	"github.com/twitchtv/twirp"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

type service struct {
}

//return url
func (s *OtherServie) verfiy(req *http.Request) (string, []byte, error) {
	//todo 验证工作
	if req.Method != http.MethodPost {
		if req.Method == http.MethodOptions {
			return "", nil, fmt.Errorf(http.MethodOptions)
		}
		msg := fmt.Sprintf("unsupported method %q (only POST is allowed)", req.Method)
		err := badRouteError(msg, req.Method, req.URL.Path)
		return "", nil, err
	}

	if req.URL.Path == "/other/uploadPic" {
		return req.URL.Path, nil, nil
	}

	defer req.Body.Close()
	buf, err := ioutil.ReadAll(req.Body)
	glog.Debug(`request:{url:%s,param:%s}`, req.URL.Path, buf)
	return req.URL.Path, buf, err
}

/**
 * 返回错误信息
 */
func (s *service) writeError(ctx context.Context, resp http.ResponseWriter, err error) {
	msg := fmt.Sprintf("%s", err)
	if msg == http.MethodOptions {
		resp.Header().Add("Access-Control-Allow-Origin", "*")
		resp.Header().Add("Access-Control-Allow-Credentials", "true")
		resp.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-Automatic-Token, X-Remote-Addr")
		return
	}
	result := struct {
		Code string      `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		FAIL_CODE,
		msg,
		nil,
	}
	buff, err := json.Marshal(result)
	if err != nil {
		glog.Error("resp data struct err", err)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusBadRequest)
	resp.Write(buff)
}

// badRouteError is used when the twirp server cannot route a request
func badRouteError(msg string, method, url string) twirp.Error {
	err := twirp.NewError(twirp.BadRoute, msg)
	err = err.WithMeta("twirp_invalid_route", method+" "+url)
	return err
}

func writeJsonResp(resp http.ResponseWriter, data interface{}) {
	result := struct {
		Code string      `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		"200",
		"OK",
		data,
	}
	buff, err := json.Marshal(result)
	if err != nil {
		glog.Error("resp data struct err", err)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(buff)
}

type httpParam struct {
	Filter map[string]string `json:"filter"`
	Sort   string            `json:"sort"`
	Page
}

//过滤、分页、排序
func getSqlParam(buff []byte, filterMap map[string]string) (filterSql, sortSql, pageSql string, err error) {
	param := new(httpParam)
	err = json.Unmarshal(buff, param)
	if err != nil {
		return
	}
	//分页
	pageSql = param.TurnSql()
	//排序
	if param.Sort != "" {
		ss := strings.Split(param.Sort, "-")
		if col, ok := filterMap[ss[0]]; ok && ss[1] != "" {
			sortSql += fmt.Sprintf(`,%v %v`, strings.Split(col, `,`)[0], ss[1])
		}

		if sortSql != "" {
			sortSql = `ORDER by ` + sortSql[1:]
		}
	}

	//过滤
	for k, v := range param.Filter {
		if col, ok := filterMap[k]; ok && v != "" {
			col := strings.Split(col, `,`)
			filterSql += fmt.Sprintf(` and %v `+col[1], col[0], v)
		}
	}
	return
}

const (
	Upload_Key    = "ak_yey"
	Upload_Secret = "97cbac39d755ee2c3c63ac79944a2a477fd5faff"
	Upload_Host   = "jxzimg.bs2ul.yy.com"
)

var contentType = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
}

/**
 * 上传图片
 * 两种情况：1，图片路径 filepath不为空  2，上传图片文件 filepath为空
 *
 * @param filepath 图片路径，可以是网络图片 也可以是本地图片
 * @param filename 上传图片的名字
 * @param buff  上传图片的内容
 */
func uploadPic(filepath, filename string, buff []byte) (string, []byte) {
	var (
		bodyBuff   []byte
		fileSuffix string
	)
	if filepath != "" {
		//分为两种:本地图片路径 和 网络连接地址
		if strings.HasPrefix(strings.Trim(filepath, " "), "http") {
			//获取读片资源
			_, bodyBuff = getHttpPic(filepath)
		} else {
			//读文件
			bodyBuff, _ = ioutil.ReadFile(filepath)
		}

		filename := path.Base(filepath)
		//获取文件后缀
		fileSuffix = path.Ext(filename)
	} else {
		//获取文件后缀
		fileSuffix = path.Ext(filename)
		bodyBuff = buff
	}

	//加密文件名字
	//blen := bytes.Count(bodyBuff, nil) - 1
	//randNum := rand.Intn(blen)
	//buf := bodyBuff[:randNum]
	fname := common.GetMd5Str(bodyBuff)

	url := "http://" + Upload_Host + "/" + fname
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBuff))

	if err != nil {
		return fmt.Sprintf("newRequest error %s", err), nil
	}

	//加密过程
	expires := strconv.FormatInt(time.Now().Unix()+7200, 10)
	hmacStr := "PUT\njxzimg\n" + fname + "\n" + expires + "\n"

	hm := common.GetSha1Str(Upload_Secret, hmacStr)
	base := base64.URLEncoding.EncodeToString(hm)

	request.Header.Set("Host", Upload_Host)
	request.Header.Set("Date", common.Gmtime())
	request.Header.Set("Authorization", Upload_Key+":"+base+":"+expires)
	request.Header.Set("Content-Length", strconv.Itoa(len(bodyBuff)))
	//根据图片类型不同
	if v, ok := contentType[fileSuffix]; ok {
		request.Header.Set("Content-Type", v) //"image/jpeg"  text/plain
	} else {
		return fmt.Sprintf("content-type no, fileSuffix is %s", fileSuffix), nil
	}

	respose, err := client.Do(request)

	if respose.StatusCode != 200 {
		return fmt.Sprintf("respose statusCode not 200"), nil
	}

	if respose.Body != nil {
		defer respose.Body.Close()
		body, err := ioutil.ReadAll(respose.Body)
		glog.Info(respose.Status)
		if err != nil {
			fmt.Sprintf("respose error %s", err)
		}

		return url, body
	}

	return url, nil
}
