package nlpservice

import (
	"context"
	"service/api"
	"net/http"
	"io/ioutil"
	"github.com/sumaig/glog"
	"github.com/golang/protobuf/proto"
	"bytes"
	"service"
	"fmt"
)

const (
	NLP_SIM_COMMENT_URL = "http://112.25.84.193:15558/nlp/comment/sim"
	NLP_SIM_COMMENT_COUNT = 50
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

//算法生成评论
func (s *Server)SimComment(ctx context.Context, req *api.SimCommentRequest) (resp *api.SimCommentResponse, err error){
	resp = &api.SimCommentResponse{
		Data:[]*api.Sentence{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	params := &api.SimSentInput{}
	params.Content = req.Content
	params.Count = NLP_SIM_COMMENT_COUNT
	sentences := getInfo(params, &api.SimSentResponse{})

	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data = sentences
	return
}

//获取信息
func getInfo(param *api.SimSentInput, respp *api.SimSentResponse)([]*api.Sentence){
	reqBodyBytes, err := proto.Marshal(param)
	if err != nil {
		glog.Info(err)
		return nil
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", NLP_SIM_COMMENT_URL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		glog.Info(err)
		return nil
	}
	req.Header.Set("Content-Type", "application/protobuf")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Info(err)
		return nil
	}

	err = proto.Unmarshal(body, respp)
	if err != nil {
		glog.Info(err)
		return nil
	}
	
	return respp.Sentences
}

