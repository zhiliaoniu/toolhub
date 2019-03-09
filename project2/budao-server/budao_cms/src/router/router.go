package router

import (
	"github.com/twitchtv/twirp"
	"net/http"
	"service"
	"service/api"
	"service/bannerservice"
	"service/blacklistservice"
	"service/commonservice"
	"service/keyword"
	"service/pushservice"
	"service/reportservice"
	"service/spiderservice"
	"service/topicservice"
	"service/userservice"
	"service/videoservice"
	"service/wechatservice"
	"twirphook"
	"service/nlpservice"
)

var (
	Mux *http.ServeMux // With read-write lock
)

// Registermulroutes function register multiple route pattern.
func Registermulroutes() {
	//hook := statsd.NewStatsdServerHooks(common.GetLogStater())
	hook := twirphook.NewHookServerHooks() //common.NewHookServerHooks()
	Mux = http.NewServeMux()

	Mux.Handle("/other/", service.NewOtherServie())

	userservice := api.NewUserServiceServer(userservice.GetServer(), hook)
	Mux.Handle(api.UserServicePathPrefix, twirphook.WithXAutoToken(userservice))

	videoHandel := api.NewVideoServiceServer(videoservice.GetServer(), hook)
	Mux.Handle(api.VideoServicePathPrefix, twirphook.WithXAutoToken(videoHandel))

	//blacklistHandel := api.NewBlacklistServiceServer(blacklistservice.GetServer(), hook)
	//Mux.Handle(api.BlacklistServicePathPrefix, blacklistHandel)
	blacklistHandel := api.NewBlacklistServiceServer(blacklistservice.GetServer(), hook)
	Mux.Handle(api.BlacklistServicePathPrefix, twirphook.WithXAutoToken(blacklistHandel))

	spider := spiderservice.GetServer()
	spiderHandle := api.NewSpiderServiceServer(spider, hook)
	Mux.Handle(api.SpiderServicePathPrefix, twirphook.WithXAutoToken(spiderHandle))

	commonhandle := api.NewCommonServiceServer(commonservice.GetServer(), hook)
	Mux.Handle(api.CommonServicePathPrefix, twirphook.WithXAutoToken(commonhandle))

	bannerervice := api.NewBannerServiceServer(bannerservice.GetServer(), hook)
	Mux.Handle(api.BannerServicePathPrefix, twirphook.WithXAutoToken(bannerervice))

	pushservice := api.NewPushServiceServer(pushservice.GetServer(), hook)
	Mux.Handle(api.PushServicePathPrefix, twirphook.WithXAutoToken(pushservice))

	topicservice := api.NewTopicServiceServer(topicservice.GetServer(), hook)
	Mux.Handle(api.TopicServicePathPrefix, twirphook.WithXAutoToken(topicservice))

	wechatservice := api.NewWechatServiceServer(wechatservice.GetServer(), &twirp.ServerHooks{})
	Mux.Handle(api.WechatServicePathPrefix, wechatservice)

	reportservice := api.NewReportServiceServer(reportservice.GetServer(), hook)
	Mux.Handle(api.ReportServicePathPrefix, twirphook.WithXAutoToken(reportservice))

	keywordservice := api.NewKeywordServiceServer(keyword.GetServer(), hook)
	Mux.Handle(api.KeywordServicePathPrefix, twirphook.WithXAutoToken(keywordservice))

	nlpservice := api.NewNlpServiceServer(nlpservice.GetServer(), hook)
	Mux.Handle(api.NlpServicePathPrefix, twirphook.WithXAutoToken(nlpservice))

}

// GetMux function get global Mux variables.
func GetMux() *http.ServeMux {
	return Mux
}
