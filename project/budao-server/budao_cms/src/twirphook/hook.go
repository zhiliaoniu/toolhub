package twirphook

import (
	"bytes"
	"common"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/sumaig/glog"
	"github.com/twitchtv/twirp"
	"io/ioutil"
	"net/http"
	"service"
	"strconv"
	"strings"
	"time"
)

const (
	VERFITY_STR = "budao.ops"
	AES_KEY     = "sfe023f_9fd&fwfl"
)

var (
	userId string
	ok     bool
)

func NewHookServerHooks() *twirp.ServerHooks {
	hooks := &twirp.ServerHooks{}

	//在客户端请求的时候触发
	hooks.RequestReceived = func(ctx context.Context) (context.Context, error) {
		var userIdstr string
		method := ctx.Value("METHOD")
		if method != nil {
			method = method.(string)
		}
		userId := ctx.Value("USERID")
		if userId != nil {
			userIdstr = userId.(string)
		}
		if method != http.MethodOptions {
			if !common.GetConfig().Condition {
				if len(userIdstr) <= 0 {
					return ctx, twirp.NewError(twirp.Unauthenticated, "授权失败")
				}
			}
		}

		return ctx, nil
	}

	hooks.ResponseSent = func(ctx context.Context) {
		//记录运营操作记录
		method := ctx.Value("METHOD")
		if method != nil {
			method = method.(string)
		}
		code, _ := twirp.StatusCode(ctx)
		if method != http.MethodOptions {
			service.RerecordOpLog(code, ctx)
		}
	}

	return hooks
}

/**
 * 获取header头
 */
func WithXAutoToken(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("X-Automatic-Token")
		ip := r.Header.Get("X-Remote-Addr")

		if r.Method != http.MethodOptions {
			if !common.GetConfig().Condition {
				//验证token
				ok, userId = verfiyToken(token)
				if !ok {
					//解决跨域问题
					w.Header().Add("Access-Control-Allow-Origin", "*")
				}
			}
		}

		ctx = context.WithValue(ctx, "USERID", userId)
		ctx = context.WithValue(ctx, "METHOD", r.Method)
		ctx = context.WithValue(ctx, "X-Automatic-Token", token)
		ctx = context.WithValue(ctx, "X-Remote-Addr", ip)

		//go 读两次body
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close() //  must close
		//请求信息
		ctx = context.WithValue(ctx, "PARAMS", string(bodyBytes))
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		ctx = context.WithValue(ctx, "ROUTE", r.URL.Path)
		fmt.Println(r.Method)
		r = r.WithContext(ctx)
		base.ServeHTTP(w, r)
	})
}

//4otEkYAcGiWvWqT1BhaGv+nyF1IpmtgYwLI6XP6JqNAGZP1LpM7cUDAPu9e+LmLaZFWjhq/ohyyZvktgfIxVdg==
/**
 * 1.获取header里面的token值
 * 2.将token进行解密
 * 3.验证是否合法
 */
func verfiyToken(token string) (isture bool, userId string) {
	//token := req.Header.Get("X-Automatic-Token")
	if len(token) <= 0 {
		return false, ""
	}

	t, err1 := base64.StdEncoding.DecodeString(token)
	if err1 != nil {
		glog.Info(err1)
		return false, ""
	}

	//解密
	origData, err := common.AesDecrypt(t, []byte(AES_KEY))
	if err != nil {
		glog.Info(err)
		return false, ""
	}
	strs := strings.Split(string(origData), "_")
	ver := common.GetMd5Str([]byte(VERFITY_STR))
	//是否过期
	nowTime := time.Now().Unix()
	expireTime, _ := strconv.ParseInt(strs[2], 10, 64)
	if strings.EqualFold(strs[0], ver) && nowTime < expireTime {
		//相等合法
		return true, strs[1]
	}

	return false, ""
}
