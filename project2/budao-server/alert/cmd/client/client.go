package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	pb "rpc/budao"
)

func alert(toUser string, content string) {

	client := pb.NewAlertWeChatWorkProtobufClient("http://localhost:8124", &http.Client{})

	var (
		msg *pb.AlertResponse
		err error
	)
	msg, err = client.Alert(context.Background(), &pb.AlertRequest{
		Touser:  toUser,
		Content: content,
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("ret = %v\n", msg)
}

func main() {

	toUser := flag.String("u", "", "to user")
	msg := flag.String("m", "", "alert msg")

	flag.Parse()

	fmt.Println("to user : ", *toUser)
	fmt.Println("message : ", *msg)
	alert(*toUser, *msg)
}
