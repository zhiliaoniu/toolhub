package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	pb "twirprpc"

	"github.com/twitchtv/twirp"
)

func main() {
	client := pb.NewTransferProtobufClient("http://localhost:8084", &http.Client{})

	var (
		req  *pb.AuditCommentRequest
		resp *pb.AuditCommentResponse
		err  error
	)

	req = &pb.AuditCommentRequest{}
	req.CommentId = 123
	//req.Content = "我爱北京天安门，法轮功，习近平，金瓶梅，苍井空"
	req.Content = "从前有座山，山上有座庙，i love china"

	for i := 0; i < 1; i++ {
		resp, err = client.AuditComment(context.Background(), req)
		if err != nil {
			if twerr, ok := err.(twirp.Error); ok {
				if twerr.Meta("retryable") != "" {
					// Log the error and go again.
					log.Printf("got error %q, retrying", twerr)
					continue
				}
			}
			// This was some fatal error!
			log.Fatal(err)
		}
	}
	fmt.Printf("resp:[ %+v ]\n", resp)
}
