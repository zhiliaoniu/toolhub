// Code generated by protoc-gen-twirp v5.4.1, DO NOT EDIT.
// source: like.proto

package budao

import bytes "bytes"
import strings "strings"
import context "context"
import fmt "fmt"
import ioutil "io/ioutil"
import http "net/http"

import jsonpb "github.com/golang/protobuf/jsonpb"
import proto "github.com/golang/protobuf/proto"
import twirp "github.com/twitchtv/twirp"
import ctxsetters "github.com/twitchtv/twirp/ctxsetters"
import cityhash "github.com/zhenjl/cityhash"
import glog "github.com/sumaig/glog"
import strconv "strconv"

// =====================
// LikeService Interface
// =====================

type LikeService interface {
	LikeVideo(context.Context, *LikeVideoRequest) (*LikeVideoResponse, error)

	LikeComment(context.Context, *LikeCommentRequest) (*LikeCommentResponse, error)
}

// ===========================
// LikeService Protobuf Client
// ===========================

type likeServiceProtobufClient struct {
	client HTTPClient
	urls   [2]string
}

// NewLikeServiceProtobufClient creates a Protobuf client that implements the LikeService interface.
// It communicates using Protobuf and can be configured with a custom HTTPClient.
func NewLikeServiceProtobufClient(addr string, client HTTPClient) LikeService {
	prefix := urlBase(addr) + LikeServicePathPrefix
	urls := [2]string{
		prefix + "LikeVideo",
		prefix + "LikeComment",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &likeServiceProtobufClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &likeServiceProtobufClient{
		client: client,
		urls:   urls,
	}
}

func (c *likeServiceProtobufClient) LikeVideo(ctx context.Context, in *LikeVideoRequest) (*LikeVideoResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "LikeService")
	ctx = ctxsetters.WithMethodName(ctx, "LikeVideo")
	out := new(LikeVideoResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[0], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceProtobufClient) LikeComment(ctx context.Context, in *LikeCommentRequest) (*LikeCommentResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "LikeService")
	ctx = ctxsetters.WithMethodName(ctx, "LikeComment")
	out := new(LikeCommentResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[1], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// =======================
// LikeService JSON Client
// =======================

type likeServiceJSONClient struct {
	client HTTPClient
	urls   [2]string
}

// NewLikeServiceJSONClient creates a JSON client that implements the LikeService interface.
// It communicates using JSON and can be configured with a custom HTTPClient.
func NewLikeServiceJSONClient(addr string, client HTTPClient) LikeService {
	prefix := urlBase(addr) + LikeServicePathPrefix
	urls := [2]string{
		prefix + "LikeVideo",
		prefix + "LikeComment",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &likeServiceJSONClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &likeServiceJSONClient{
		client: client,
		urls:   urls,
	}
}

func (c *likeServiceJSONClient) LikeVideo(ctx context.Context, in *LikeVideoRequest) (*LikeVideoResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "LikeService")
	ctx = ctxsetters.WithMethodName(ctx, "LikeVideo")
	out := new(LikeVideoResponse)
	err := doJSONRequest(ctx, c.client, c.urls[0], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceJSONClient) LikeComment(ctx context.Context, in *LikeCommentRequest) (*LikeCommentResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "LikeService")
	ctx = ctxsetters.WithMethodName(ctx, "LikeComment")
	out := new(LikeCommentResponse)
	err := doJSONRequest(ctx, c.client, c.urls[1], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ==========================
// LikeService Server Handler
// ==========================

type likeServiceServer struct {
	LikeService
	hooks *twirp.ServerHooks
}

func NewLikeServiceServer(svc LikeService, hooks *twirp.ServerHooks) TwirpServer {
	return &likeServiceServer{
		LikeService: svc,
		hooks:       hooks,
	}
}

// writeError writes an HTTP response with a valid Twirp error format, and triggers hooks.
// If err is not a twirp.Error, it will get wrapped with twirp.InternalErrorWith(err)
func (s *likeServiceServer) writeError(ctx context.Context, resp http.ResponseWriter, err error) {
	writeError(ctx, resp, err, s.hooks)
}

// LikeServicePathPrefix is used for all URL paths on a twirp LikeService server.
// Requests are always: POST LikeServicePathPrefix/method
// It can be used in an HTTP mux to route twirp requests along with non-twirp requests on other routes.
const LikeServicePathPrefix = "/budao.LikeService/"

func (s *likeServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "LikeService")
	ctx = ctxsetters.WithRequestIp(ctx, req.Header.Get("Remote_addr"))
	ctx = ctxsetters.WithSIG(ctx, req.Header.Get("sig"))
	ctx = ctxsetters.WithResponseWriter(ctx, resp)

	var err error
	ctx, err = callRequestReceived(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	if req.Method != "POST" {
		msg := fmt.Sprintf("unsupported method %q (only POST is allowed)", req.Method)
		err = badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, err)
		return
	}

	switch req.URL.Path {
	case "/budao.LikeService/LikeVideo":
		s.serveLikeVideo(ctx, resp, req)
		return
	case "/budao.LikeService/LikeComment":
		s.serveLikeComment(ctx, resp, req)
		return
	default:
		msg := fmt.Sprintf("no handler for path %q", req.URL.Path)
		err = badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, err)
		return
	}
}

func (s *likeServiceServer) serveLikeVideo(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveLikeVideoJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveLikeVideoProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *likeServiceServer) serveLikeVideoJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "LikeVideo")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(LikeVideoRequest)
	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err = unmarshaler.Unmarshal(req.Body, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request json")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// check client signature
	sig, _ := twirp.SIG(ctx)
	if sig != "" {
		buf, _ := proto.Marshal(reqContent)
		seed := uint64(307976497148328517)
		ch := cityhash.CityHash64WithSeed(buf, uint32(len(buf)), seed)
		sigNum, _ := strconv.ParseUint(sig, 10, 64)
		if sigNum != ch {
			fmt.Printf("req has bad_signature. req:%v, correct_sig:%d\n", req, ch)
		} else {
			fmt.Printf("req has correct_signature. req:%v\n", req)
		}
	} else {
		fmt.Printf("req has no_signature req:%v\n", req)
	}

	// Call service method
	var respContent *LikeVideoResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.LikeVideo(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *LikeVideoResponse and nil error while calling LikeVideo. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: true}
	if err = marshaler.Marshal(&buf, respContent); err != nil {
		err = wrapErr(err, "failed to marshal json response")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)

	respBytes := buf.Bytes()
	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *likeServiceServer) serveLikeVideoProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "LikeVideo")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		err = wrapErr(err, "failed to read request body")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}
	reqContent := new(LikeVideoRequest)
	if err = proto.Unmarshal(buf, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request proto")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// check client signature
	sig, _ := twirp.SIG(ctx)
	if sig != "" {
		buf, _ := proto.Marshal(reqContent)
		seed := uint64(307976497148328517)
		ch := cityhash.CityHash64WithSeed(buf, uint32(len(buf)), seed)
		sigNum, _ := strconv.ParseUint(sig, 10, 64)
		if sigNum != ch {
			glog.Error("req has bad_signature. req:%v, correct_sig:%d", req, ch)
		} else {
			glog.Debug("req has correct_signature. req:%v", req)
		}
	} else {
		glog.Error("req has no_signature req:%v", req)
	}

	// Call service method
	var respContent *LikeVideoResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.LikeVideo(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *LikeVideoResponse and nil error while calling LikeVideo. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	respBytes, err := proto.Marshal(respContent)
	if err != nil {
		err = wrapErr(err, "failed to marshal proto response")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/protobuf")
	resp.WriteHeader(http.StatusOK)
	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *likeServiceServer) serveLikeComment(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveLikeCommentJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveLikeCommentProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *likeServiceServer) serveLikeCommentJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "LikeComment")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(LikeCommentRequest)
	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err = unmarshaler.Unmarshal(req.Body, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request json")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// check client signature
	sig, _ := twirp.SIG(ctx)
	if sig != "" {
		buf, _ := proto.Marshal(reqContent)
		seed := uint64(307976497148328517)
		ch := cityhash.CityHash64WithSeed(buf, uint32(len(buf)), seed)
		sigNum, _ := strconv.ParseUint(sig, 10, 64)
		if sigNum != ch {
			fmt.Printf("req has bad_signature. req:%v, correct_sig:%d\n", req, ch)
		} else {
			fmt.Printf("req has correct_signature. req:%v\n", req)
		}
	} else {
		fmt.Printf("req has no_signature req:%v\n", req)
	}

	// Call service method
	var respContent *LikeCommentResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.LikeComment(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *LikeCommentResponse and nil error while calling LikeComment. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: true}
	if err = marshaler.Marshal(&buf, respContent); err != nil {
		err = wrapErr(err, "failed to marshal json response")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)

	respBytes := buf.Bytes()
	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *likeServiceServer) serveLikeCommentProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "LikeComment")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		err = wrapErr(err, "failed to read request body")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}
	reqContent := new(LikeCommentRequest)
	if err = proto.Unmarshal(buf, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request proto")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// check client signature
	sig, _ := twirp.SIG(ctx)
	if sig != "" {
		buf, _ := proto.Marshal(reqContent)
		seed := uint64(307976497148328517)
		ch := cityhash.CityHash64WithSeed(buf, uint32(len(buf)), seed)
		sigNum, _ := strconv.ParseUint(sig, 10, 64)
		if sigNum != ch {
			glog.Error("req has bad_signature. req:%v, correct_sig:%d", req, ch)
		} else {
			glog.Debug("req has correct_signature. req:%v", req)
		}
	} else {
		glog.Error("req has no_signature req:%v", req)
	}

	// Call service method
	var respContent *LikeCommentResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.LikeComment(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *LikeCommentResponse and nil error while calling LikeComment. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	respBytes, err := proto.Marshal(respContent)
	if err != nil {
		err = wrapErr(err, "failed to marshal proto response")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/protobuf")
	resp.WriteHeader(http.StatusOK)
	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *likeServiceServer) ServiceDescriptor() ([]byte, int) {
	return twirpFileDescriptor2, 0
}

func (s *likeServiceServer) ProtocGenTwirpVersion() string {
	return "v5.4.1"
}

var twirpFileDescriptor2 = []byte{
	// 336 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x52, 0x4f, 0x4b, 0xfb, 0x40,
	0x14, 0xfc, 0x6d, 0xfb, 0x6b, 0x6c, 0x5f, 0xad, 0xb4, 0xab, 0x60, 0x1a, 0x10, 0x4a, 0x40, 0xa8,
	0x1e, 0x72, 0xa8, 0x77, 0xc1, 0x7f, 0x60, 0xb5, 0x48, 0xd9, 0xa2, 0x88, 0x17, 0x49, 0x77, 0x17,
	0x5c, 0x6b, 0xb2, 0x31, 0xd9, 0x14, 0xf1, 0xe2, 0xc9, 0x4f, 0xe0, 0x37, 0xf0, 0x93, 0x4a, 0x76,
	0x17, 0x9b, 0xaa, 0x47, 0xbd, 0xe5, 0xcd, 0x4c, 0xe6, 0xcd, 0xbc, 0x04, 0xe0, 0x41, 0xcc, 0x78,
	0x90, 0xa4, 0x52, 0x49, 0x5c, 0x9b, 0xe6, 0x2c, 0x94, 0xde, 0x2a, 0x95, 0x51, 0x24, 0x63, 0x03,
	0xfa, 0x2f, 0xd0, 0x1e, 0x89, 0x19, 0xbf, 0x12, 0x8c, 0x4b, 0xc2, 0x1f, 0x73, 0x9e, 0x29, 0xbc,
	0x0d, 0xce, 0x1d, 0x0f, 0x19, 0x4f, 0x5d, 0xd4, 0x43, 0xfd, 0xe6, 0xa0, 0x15, 0xe8, 0x37, 0x83,
	0x53, 0x0d, 0x12, 0x4b, 0xe2, 0x1d, 0x70, 0x42, 0xaa, 0x84, 0x8c, 0xdd, 0x4a, 0x0f, 0xf5, 0xd7,
	0x06, 0x1d, 0x2b, 0x2b, 0xfc, 0x0e, 0x34, 0x41, 0xac, 0x00, 0x77, 0xa1, 0x3e, 0x2f, 0x36, 0xdc,
	0x0a, 0xe6, 0x56, 0x7b, 0xa8, 0xdf, 0x20, 0x2b, 0x7a, 0x1e, 0x32, 0x7f, 0x0c, 0x9d, 0x52, 0x80,
	0x2c, 0x91, 0x71, 0xc6, 0x8b, 0x04, 0x99, 0x0a, 0x55, 0x9e, 0x7d, 0x49, 0x30, 0xd1, 0x20, 0xb1,
	0x24, 0xde, 0x80, 0x1a, 0x95, 0x79, 0xac, 0x74, 0x80, 0x16, 0x31, 0x83, 0xff, 0x8a, 0x00, 0x17,
	0x96, 0x47, 0x32, 0x8a, 0x78, 0xac, 0xfe, 0xae, 0xd5, 0x16, 0x00, 0x35, 0x3b, 0x16, 0xbd, 0x1a,
	0x16, 0x19, 0x32, 0x9f, 0xc0, 0xfa, 0x52, 0x8c, 0x5f, 0xe8, 0xb6, 0xeb, 0x03, 0x2c, 0x82, 0xe0,
	0x3a, 0xfc, 0x1f, 0x0d, 0xcf, 0x4f, 0xda, 0xff, 0x30, 0x80, 0x73, 0x79, 0xa1, 0x9f, 0xd1, 0xe0,
	0x0d, 0x41, 0xb3, 0x10, 0x4d, 0x78, 0x3a, 0x17, 0x94, 0xe3, 0x7d, 0x68, 0x7c, 0x5e, 0x18, 0x6f,
	0x96, 0xea, 0x94, 0x3f, 0xba, 0xe7, 0x7e, 0x27, 0x6c, 0xe0, 0x63, 0x63, 0x67, 0x7b, 0xe0, 0x6e,
	0x49, 0xb8, 0x7c, 0x62, 0xcf, 0xfb, 0x89, 0x32, 0x2e, 0x87, 0x1e, 0xb4, 0xa9, 0x8c, 0x82, 0xfb,
	0xa7, 0xe7, 0x20, 0x99, 0x1a, 0xdd, 0x18, 0xbd, 0x57, 0xaa, 0x67, 0xd7, 0x37, 0x53, 0x47, 0xff,
	0x8b, 0x7b, 0x1f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xdb, 0xc9, 0x94, 0x3c, 0xae, 0x02, 0x00, 0x00,
}
