package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	pb "twirprpc"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"
)

func commentVideo() {
	client := pb.NewCommentServiceProtobufClient("http://localhost:8084", &http.Client{})

	var (
		req  *pb.CommentVideoRequest
		resp *pb.CommentVideoResponse
		err  error
	)

	header := &pb.Header{}
	header.UserId = "27534343609"
	req = &pb.CommentVideoRequest{}
	req.Header = header
	req.VideoId = "2"
	req.Content = "from test script"

	fmt.Printf("req:%+v\n", req)

	resp, err = client.CommentVideo(context.Background(), req)
	if err != nil {
		log.Printf("got error %v", err)
	}
	log.Printf("resp:%v\n", resp)

}

func replyComment() {
	client := pb.NewCommentServiceProtobufClient("http://localhost:8084", &http.Client{})

	var (
		req  *pb.ReplyCommentRequest
		resp *pb.ReplyCommentResponse
		err  error
	)

	header := &pb.Header{}
	header.UserId = "22688635834"
	req = &pb.ReplyCommentRequest{}
	req.Header = header
	req.Type = &pb.ReplyCommentRequest_CommentId{
		CommentId: "487473820791",
	}
	req.Content = "reply comment1"

	fmt.Printf("req:%+v\n", req)

	resp, err = client.ReplyComment(context.Background(), req)
	if err != nil {
		log.Printf("got error %v", err)
	}
	log.Printf("resp:%v\n", resp)

}

func replyReply() {
	client := pb.NewCommentServiceProtobufClient("http://localhost:8084", &http.Client{})

	var (
		req  *pb.ReplyCommentRequest
		resp *pb.ReplyCommentResponse
		err  error
	)

	header := &pb.Header{}
	header.UserId = "22688635834"
	req = &pb.ReplyCommentRequest{}
	req.Header = header
	req.Type = &pb.ReplyCommentRequest_ReplyId{
		ReplyId: "573623459033",
	}
	req.Content = "reply reply1"

	fmt.Printf("req:%+v\n", req)

	resp, err = client.ReplyComment(context.Background(), req)
	if err != nil {
		log.Printf("got error %v", err)
	}
	log.Printf("resp:%v\n", resp)

}

func getVideoCommentList() {
	client := pb.NewCommentServiceProtobufClient("http://localhost:8080", &http.Client{})

	var (
		req  *pb.GetVideoCommentListRequest
		resp *pb.GetVideoCommentListResponse
		err  error
	)

	header := &pb.Header{}
	header.UserId = "22688635834"
	req = &pb.GetVideoCommentListRequest{}
	req.Header = header
	req.VideoId = "10978572654864"

	reqStr, _ := proto.Marshal(req)
	ioutil.WriteFile("./req_get_video_comment_list.pb", reqStr, 0666)
	fmt.Printf("req:%+v\n", req)

	resp, err = client.GetVideoCommentList(context.Background(), req)
	if err != nil {
		log.Printf("got error %v", err)
	}
	log.Printf("resp:%v\n", resp)
}

func getCommentReplyListRequest() {
	client := pb.NewCommentServiceProtobufClient("http://localhost:8084", &http.Client{})

	var (
		req  *pb.GetCommentReplyListRequest
		resp *pb.GetCommentReplyListResponse
		err  error
	)

	header := &pb.Header{}
	header.UserId = "22688635834"
	req = &pb.GetCommentReplyListRequest{}
	req.Header = header
	req.CommentId = "487473820791"

	fmt.Printf("req:%+v\n", req)

	resp, err = client.GetCommentReplyList(context.Background(), req)
	if err != nil {
		log.Printf("got error %v", err)
	}
	log.Printf("resp:%v\n", resp)
}

var RedisConnPool *redis.Pool

// InitRedisConnPool function init redis connect pool
func InitRedisConnPool() {

	addr := "221.228.106.9:4019"
	maxIdle := 10000
	idleTimeout := 240 * time.Second
	connTimeout := time.Duration(1000) * time.Millisecond
	readTimeout := time.Duration(1000) * time.Millisecond
	writeTimeout := time.Duration(1000) * time.Millisecond

	RedisConnPool = &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", addr, connTimeout, readTimeout, writeTimeout)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: nil,
	}

	glog.Debug("redis connect pool init succeed")
}

func delKeyByPattern(ps []string) {
	InitRedisConnPool()
	conn := RedisConnPool.Get()

	for _, p := range ps {
		r, err := redis.Strings(conn.Do("keys", p))
		if err != nil {
			fmt.Println(err)
		}
		for _, key := range r {
			fmt.Printf("del key %s\n", key)
			_, _ = redis.Int(conn.Do("del", key))
		}
	}
}

func execFuncNum(f func(), num int) {
	i := 1
	for {
		if i > num {
			break
		}
		f()
		i++
	}
}

func main() {
	//delKeyByPattern([]string{"vcnew_*", "crnew_*", "comment_item_all"})
	//delKeyByPattern([]string{"sourcevid_*"})
	delKeyByPattern([]string{"*list*"})
	execFuncNum(commentVideo, 0)
	execFuncNum(replyComment, 0)
	execFuncNum(replyReply, 0)
	execFuncNum(getVideoCommentList, 0)
	execFuncNum(getCommentReplyListRequest, 0)
	//replyComment()
	//replyReply()
	//getVideoCommentList()
	//getCommentReplyListRequest()
}
