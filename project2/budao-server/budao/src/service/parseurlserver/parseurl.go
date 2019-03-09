package parseurlserver

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/sumaig/glog"
	"github.com/twitchtv/twirp"

	"base"
	"common"
	"db"
	pb "twirprpc"
)

// Server identify for parseurl RPC
type Server struct {
	logStater base.LogStater
}

// GetServer return server of parseurl service
func GetServer() *Server {
	server := &Server{}
	server.initServer()

	return server
}

func (s *Server) initServer() {
	s.logStater = base.GetLogStater()
}

// ParseURL function resolve the url of the playable video for the client.
func (s *Server) ParseURL(ctx context.Context, req *pb.ParseURLRequest) (resp *pb.ParseURLResponse, err error) {
	clientIp, _ := twirp.RequestIp(ctx)
	glog.Debug("clientIp:%v, req:%v", clientIp, req)
	resp = &pb.ParseURLResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			s.logStater.IncInternalError("ParseURL_ERR", "twirp.status_codes.ParseURL.500", 1, 1)
			err = nil
		}
	}()

	vid, _ := strconv.ParseUint(req.GetVideoId(), 10, 64)
	videoTableName, err := db.GetTableName("video_", vid)
	sqlString := fmt.Sprintf("select parse_type, videourl from %s where vid = %v", videoTableName, vid)
	var parseType int32
	var videoURL string
	parseType = -1
	tempRow, err := db.QueryRow(common.BUDAODB, sqlString)
	err = tempRow.Scan(&parseType, &videoURL)
	if err != nil {
		glog.Error("ParseURL query video failed. err:%v", err)
		return
	}
	glog.Debug("ParseURL videoSource:%d, videourl:%s", parseType, videoURL)

	var play *pb.ParseURLResponse_Play
	var craw *pb.ParseURLResponse_Craw
	var web *pb.ParseURLResponse_Web
	pType := pb.VideoParseRule(parseType)
	if pType == pb.VideoParseRule_PARSE_RULE_DOUYIN ||
		pType == pb.VideoParseRule_PARSE_RULE_NEIHANDUANZI {
		play = &pb.ParseURLResponse_Play{}
	} else if pType == pb.VideoParseRule_PARSE_RULE_HAOKAN ||
		pType == pb.VideoParseRule_PARSE_RULE_BOBO ||
		pType == pb.VideoParseRule_PARSE_RULE_LI ||
		pType == pb.VideoParseRule_PARSE_RULE_BZHAN ||
		pType == pb.VideoParseRule_PARSE_RULE_KUAISHOU ||
		pType == pb.VideoParseRule_PARSE_RULE_MIAOPAI ||
		pType == pb.VideoParseRule_PARSE_RULE_WEIBO {
		glog.Debug("test^^^^^^^^^^^^^^^^^^^^^^^^^")
		craw = &pb.ParseURLResponse_Craw{}
	} else if pType == pb.VideoParseRule_PARSE_RULE_MEIPAI ||
		pType == pb.VideoParseRule_PARSE_RULE_FEIDIESHUO ||
		pType == pb.VideoParseRule_PARSE_RULE_YANGGUANGSHIPIN ||
		pType == pb.VideoParseRule_PARSE_RULE_XIGUA {
		web = &pb.ParseURLResponse_Web{}
	} else {
		resp.Status.Code = pb.Status_OK
		resp.Status.Message = "Video type does not exist."
		web = &pb.ParseURLResponse_Web{}
		web.WebUrl = videoURL
		resp.Type = &pb.ParseURLResponse_Web_{
			Web: web,
		}
		return
	}

	switch pType {
	case pb.VideoParseRule_PARSE_RULE_DOUYIN:
		play.VideoUrl = videoURL
		resp.Type = &pb.ParseURLResponse_Play_{
			Play: play,
		}
	case pb.VideoParseRule_PARSE_RULE_NEIHANDUANZI:
		play.VideoUrl = videoURL
		resp.Type = &pb.ParseURLResponse_Play_{
			Play: play,
		}
	case pb.VideoParseRule_PARSE_RULE_KUAISHOU:
		header := transferKuaishouVideoURL()
		craw.Next = "kuaishou_next_0"
		craw.RequestPost = false
		craw.Header = header
		craw.VideoUrl = videoURL
		resp.Type = &pb.ParseURLResponse_Craw_{
			Craw: craw,
		}
	case pb.VideoParseRule_PARSE_RULE_WEIBO:
		header := transferWeiboVideoURL()
		craw.Next = "weibo_next_0"
		craw.RequestPost = false
		craw.Header = header
		craw.VideoUrl = videoURL
		resp.Type = &pb.ParseURLResponse_Craw_{
			Craw: craw,
		}
	case pb.VideoParseRule_PARSE_RULE_HAOKAN:
		glog.Debug("test^^^^^^^^^^^^^^^^^^^^^^^^^transferHaokanVideoURL")
		headers, field, target := transferHaokanVideoURL(videoURL)
		glog.Debug("field:%v", field)
		glog.Debug("parseurl haokan, target:", target)
		craw.Next = "haokan_next_0"
		craw.RequestPost = true
		craw.Header = headers
		craw.Field = field
		craw.VideoUrl = target
		resp.Type = &pb.ParseURLResponse_Craw_{
			Craw: craw,
		}
	case pb.VideoParseRule_PARSE_RULE_BOBO:
		craw.Next = "bobo_next_0"
		craw.RequestPost = true
		headers, field := transferBoboVideoURL(videoURL)
		craw.Header = headers
		craw.Field = field
		craw.VideoUrl = "https://api.bbobo.com/v1/video/play.json"
		resp.Type = &pb.ParseURLResponse_Craw_{
			Craw: craw,
		}
	case pb.VideoParseRule_PARSE_RULE_MEIPAI:
		web.WebUrl = videoURL
		resp.Type = &pb.ParseURLResponse_Web_{
			Web: web,
		}
	case pb.VideoParseRule_PARSE_RULE_BZHAN:
		craw.Next = "bzhan_next_0"
		craw.RequestPost = false
		headers, target := transferBzhanVideoURL(videoURL)
		glog.Debug("^^^^^^^^^^^^^^^^^")
		glog.Debug(target)
		craw.Header = headers
		glog.Debug(headers)
		craw.VideoUrl = "https://api.bilibili.com/playurl" + "?" + target
		resp.Type = &pb.ParseURLResponse_Craw_{
			Craw: craw,
		}
	case pb.VideoParseRule_PARSE_RULE_FEIDIESHUO:
		web.WebUrl = videoURL
		resp.Type = &pb.ParseURLResponse_Web_{
			Web: web,
		}
	case pb.VideoParseRule_PARSE_RULE_YANGGUANGSHIPIN:
		web.WebUrl = videoURL
		resp.Type = &pb.ParseURLResponse_Web_{
			Web: web,
		}
	case pb.VideoParseRule_PARSE_RULE_XIGUA:
		web.WebUrl = videoURL
		resp.Type = &pb.ParseURLResponse_Web_{
			Web: web,
		}
	case pb.VideoParseRule_PARSE_RULE_LI:
		craw.Next = "li_next_0"
		craw.RequestPost = true
		headers, field := transferLiVideoURL(videoURL, req.GetHeader())
		craw.Header = headers
		craw.Field = field
		craw.VideoUrl = "http://app.pearvideo.com/clt/jsp/v4/content.jsp"
		resp.Type = &pb.ParseURLResponse_Craw_{
			Craw: craw,
		}
	case pb.VideoParseRule_PARSE_RULE_MIAOPAI:
		craw.Next = "miaopai_next_0"
		craw.RequestPost = false

		videoURL = fmt.Sprintf("http://gslb.miaopai.com/stream/%s.json", extractCharacter(videoURL))
		headers, target := transferMiaopaiVideoURL()
		tempurl := videoURL + "?" + target
		urlCode, _ := url.Parse(tempurl)
		query := urlCode.Query().Encode()
		craw.VideoUrl = videoURL + "?" + query
		craw.Header = headers
		resp.Type = &pb.ParseURLResponse_Craw_{
			Craw: craw,
		}
	}

	resp.Status.Code = pb.Status_OK
	glog.Debug(resp)

	return
}

// ParseExternalURL function resolve the URL of the outer chain video for the client.
func (s *Server) ParseExternalURL(ctx context.Context, req *pb.ParseExternalURLRequest) (resp *pb.ParseURLResponse, err error) {
	clientIp, _ := twirp.RequestIp(ctx)
	glog.Debug("clientIp:%v, req:%v", clientIp, req)
	resp = &pb.ParseURLResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			s.logStater.IncInternalError("ParseExternalURL_ERR", "twirp.status_codes.ParseExternalURL.500", 1, 1)
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	funcTag := req.GetFunc()
	body := req.GetBody()
	var videoURL string
	var header map[string]string
	var next bool
	//recover panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				glog.Error("parse external url panic. req:%v", req)
				err = errors.New("parse external url panic.")
				return
				//panic(r)
			}
		}()
		videoURL, header, next = GetVideoAndHeader(body, funcTag, req.GetHeader())
	}()
	if err != nil {
		return
	}

	if videoURL == "" {
		// parse failed
		web := &pb.ParseURLResponse_Web{}
		vid, _ := strconv.ParseUint(req.GetVideoId(), 10, 64)
		videoTableName, err := db.GetTableName("video_", vid)
		sqlString := fmt.Sprintf("select videourl from %s where vid = %v", videoTableName, vid)
		tempRow, err := db.QueryRow(common.BUDAODB, sqlString)
		err = tempRow.Scan(&videoURL)
		if err != nil {
			glog.Error("ParseExternalURL query video failed. err:%v", err)
		}
		web.WebUrl = videoURL
		resp.Type = &pb.ParseURLResponse_Web_{
			Web: web,
		}
		err = errors.New("ParseExternalURL parse body failed")
	} else {
		// parse success
		if next == true {
			glog.Debug("^^^^^^^^^^^^^ready to craw kuaishou again")
			craw := &pb.ParseURLResponse_Craw{}
			craw.Next = "kuaishou_next_0"
			craw.RequestPost = false
			craw.VideoUrl = videoURL
			craw.Header = header
			resp.Type = &pb.ParseURLResponse_Craw_{
				Craw: craw,
			}
		} else {
			play := &pb.ParseURLResponse_Play{}
			play.VideoUrl = videoURL
			play.Header = header
			resp.Type = &pb.ParseURLResponse_Play_{
				Play: play,
			}
		}
	}

	resp.Status.Code = pb.Status_OK
	glog.Debug(resp)

	return
}

func GetVideoAndHeader(body, funcTag string, reqHeader *pb.Header) (videoURL string, header map[string]string, next bool) {
	next = false
	switch funcTag {
	case "haokan_next_0":
		videoURL, header = parseHaokanVideoBody(body)
	case "kuaishou_next_0":
		videoURL, header, next = parseKuaishou(body)
	case "weibo_next_0":
		videoURL, header = parseWeibo(body)
	case "li_next_0":
		videoURL, header = parseLiVideoBody(body, reqHeader)
	case "miaopai_next_0":
		videoURL, header = parseMiaopaiVideoBody(body)
	case "bobo_next_0":
		videoURL, header = parseBoboVideoBody(body)
	case "bzhan_next_0":
		videoURL, header = parseBzhan(body)
		glog.Debug("^^^^^^^^^^^^^^^^^")
		glog.Debug(videoURL)
	}
	return
}
