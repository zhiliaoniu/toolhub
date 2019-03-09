package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	pb "twirprpc"

	"github.com/twitchtv/twirp"
)

func main() {
	client := pb.NewTransferProtobufClient("http://localhost:8080", &http.Client{})

	var (
		req  *pb.PostVideosRequest
		resp *pb.PostVideosResponse
		err  error
	)

	req = &pb.PostVideosRequest{}
	tVideos := make([]*pb.PostVideo, 0)
	postVideo := &pb.PostVideo{}
	randSouceVid := strconv.FormatUint(rand.New(rand.NewSource(time.Now().UnixNano())).Uint64(), 10)
	postVideo.SourceVid = randSouceVid
	postVideo.Pic = "http://v1.dwstatic.com/bd/201709//09/pic/3ad7a91aa4f1b259d46640a041a00000.jpg?w=240&h=320"
	postVideo.VideoUrl = "http://v1.dwstatic.com/bd/201709/09/video/3ad7a91aa4f1b259d46640a041a00000.mp4"
	postVideo.Title = "兄弟问你个事"
	postVideo.Width = 240
	postVideo.Height = 320
	postVideo.Duration = 7
	postVideo.EVideoSource = pb.EVideoSource_E_VideoSource_CRAWL_BZHAN
	tVideos = append(tVideos, postVideo)
	req.TVideos = tVideos

	for i := 0; i < 1; i++ {
		resp, err = client.PostVideos(context.Background(), req)
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
	fmt.Printf("resp:%+v\n", resp)
}
