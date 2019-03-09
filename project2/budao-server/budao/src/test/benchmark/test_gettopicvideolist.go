package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pb "twirprpc"

	"github.com/golang/protobuf/proto"
)

func main() {
	client := pb.NewTimeLineServiceProtobufClient("http://localhost:8080", &http.Client{})

	req := &pb.GetTimeLineRequest{
		Header: &pb.Header{
			UserId: "49515001920",
			DeviceInfo: &pb.DeviceInfo{
				DeviceId: "c0137b852d295d17",
			},
		},
	}
	reqStr, _ := proto.Marshal(req)
	ioutil.WriteFile("./req_gettimeline.pb", reqStr, 0666)

	fmt.Printf("req:%+v\n", req)

	resp, err := client.GetTimeLine(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("resp:[ %+v ]\n", resp)
}
