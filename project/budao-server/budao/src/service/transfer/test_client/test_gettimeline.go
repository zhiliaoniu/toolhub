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
	client := pb.NewTimeLineServiceProtobufClient("http://localhost:8082", &http.Client{})

	var (
		req  *pb.GetTimeLineRequest
		resp *pb.GetTimeLineResponse
		err  error
	)

	header := &pb.Header{}
	header.UserId = "123456abcdef--yangshengzhi"
	req = &pb.GetTimeLineRequest{}
	req.Header = header

	fmt.Printf("req:%+v\n", req)

	for i := 0; i < 1; i++ {
		resp, err = client.GetTimeLine(context.Background(), req)
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
