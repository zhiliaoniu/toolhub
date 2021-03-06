// Code generated by protoc-gen-twirp v5.4.1, DO NOT EDIT.
// source: topic.proto

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

// ======================
// TopicService Interface
// ======================

type TopicService interface {
	GetTopicList(context.Context, *GetTopicListRequest) (*GetTopicListResponse, error)

	GetTopicVideoList(context.Context, *GetTopicVideoListRequest) (*GetTopicVideoListResponse, error)

	SubscribeTopic(context.Context, *SubscribeTopicRequest) (*SubscribeTopicResponse, error)

	GetSubscribedTopicList(context.Context, *GetSubscribedTopicListRequest) (*GetSubscribedTopicListResponse, error)
}

// ============================
// TopicService Protobuf Client
// ============================

type topicServiceProtobufClient struct {
	client HTTPClient
	urls   [4]string
}

// NewTopicServiceProtobufClient creates a Protobuf client that implements the TopicService interface.
// It communicates using Protobuf and can be configured with a custom HTTPClient.
func NewTopicServiceProtobufClient(addr string, client HTTPClient) TopicService {
	prefix := urlBase(addr) + TopicServicePathPrefix
	urls := [4]string{
		prefix + "GetTopicList",
		prefix + "GetTopicVideoList",
		prefix + "SubscribeTopic",
		prefix + "GetSubscribedTopicList",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &topicServiceProtobufClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &topicServiceProtobufClient{
		client: client,
		urls:   urls,
	}
}

func (c *topicServiceProtobufClient) GetTopicList(ctx context.Context, in *GetTopicListRequest) (*GetTopicListResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "TopicService")
	ctx = ctxsetters.WithMethodName(ctx, "GetTopicList")
	out := new(GetTopicListResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[0], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topicServiceProtobufClient) GetTopicVideoList(ctx context.Context, in *GetTopicVideoListRequest) (*GetTopicVideoListResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "TopicService")
	ctx = ctxsetters.WithMethodName(ctx, "GetTopicVideoList")
	out := new(GetTopicVideoListResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[1], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topicServiceProtobufClient) SubscribeTopic(ctx context.Context, in *SubscribeTopicRequest) (*SubscribeTopicResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "TopicService")
	ctx = ctxsetters.WithMethodName(ctx, "SubscribeTopic")
	out := new(SubscribeTopicResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[2], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topicServiceProtobufClient) GetSubscribedTopicList(ctx context.Context, in *GetSubscribedTopicListRequest) (*GetSubscribedTopicListResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "TopicService")
	ctx = ctxsetters.WithMethodName(ctx, "GetSubscribedTopicList")
	out := new(GetSubscribedTopicListResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[3], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ========================
// TopicService JSON Client
// ========================

type topicServiceJSONClient struct {
	client HTTPClient
	urls   [4]string
}

// NewTopicServiceJSONClient creates a JSON client that implements the TopicService interface.
// It communicates using JSON and can be configured with a custom HTTPClient.
func NewTopicServiceJSONClient(addr string, client HTTPClient) TopicService {
	prefix := urlBase(addr) + TopicServicePathPrefix
	urls := [4]string{
		prefix + "GetTopicList",
		prefix + "GetTopicVideoList",
		prefix + "SubscribeTopic",
		prefix + "GetSubscribedTopicList",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &topicServiceJSONClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &topicServiceJSONClient{
		client: client,
		urls:   urls,
	}
}

func (c *topicServiceJSONClient) GetTopicList(ctx context.Context, in *GetTopicListRequest) (*GetTopicListResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "TopicService")
	ctx = ctxsetters.WithMethodName(ctx, "GetTopicList")
	out := new(GetTopicListResponse)
	err := doJSONRequest(ctx, c.client, c.urls[0], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topicServiceJSONClient) GetTopicVideoList(ctx context.Context, in *GetTopicVideoListRequest) (*GetTopicVideoListResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "TopicService")
	ctx = ctxsetters.WithMethodName(ctx, "GetTopicVideoList")
	out := new(GetTopicVideoListResponse)
	err := doJSONRequest(ctx, c.client, c.urls[1], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topicServiceJSONClient) SubscribeTopic(ctx context.Context, in *SubscribeTopicRequest) (*SubscribeTopicResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "TopicService")
	ctx = ctxsetters.WithMethodName(ctx, "SubscribeTopic")
	out := new(SubscribeTopicResponse)
	err := doJSONRequest(ctx, c.client, c.urls[2], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topicServiceJSONClient) GetSubscribedTopicList(ctx context.Context, in *GetSubscribedTopicListRequest) (*GetSubscribedTopicListResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "TopicService")
	ctx = ctxsetters.WithMethodName(ctx, "GetSubscribedTopicList")
	out := new(GetSubscribedTopicListResponse)
	err := doJSONRequest(ctx, c.client, c.urls[3], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ===========================
// TopicService Server Handler
// ===========================

type topicServiceServer struct {
	TopicService
	hooks *twirp.ServerHooks
}

func NewTopicServiceServer(svc TopicService, hooks *twirp.ServerHooks) TwirpServer {
	return &topicServiceServer{
		TopicService: svc,
		hooks:        hooks,
	}
}

// writeError writes an HTTP response with a valid Twirp error format, and triggers hooks.
// If err is not a twirp.Error, it will get wrapped with twirp.InternalErrorWith(err)
func (s *topicServiceServer) writeError(ctx context.Context, resp http.ResponseWriter, err error) {
	writeError(ctx, resp, err, s.hooks)
}

// TopicServicePathPrefix is used for all URL paths on a twirp TopicService server.
// Requests are always: POST TopicServicePathPrefix/method
// It can be used in an HTTP mux to route twirp requests along with non-twirp requests on other routes.
const TopicServicePathPrefix = "/budao.TopicService/"

func (s *topicServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "TopicService")
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
	case "/budao.TopicService/GetTopicList":
		s.serveGetTopicList(ctx, resp, req)
		return
	case "/budao.TopicService/GetTopicVideoList":
		s.serveGetTopicVideoList(ctx, resp, req)
		return
	case "/budao.TopicService/SubscribeTopic":
		s.serveSubscribeTopic(ctx, resp, req)
		return
	case "/budao.TopicService/GetSubscribedTopicList":
		s.serveGetSubscribedTopicList(ctx, resp, req)
		return
	default:
		msg := fmt.Sprintf("no handler for path %q", req.URL.Path)
		err = badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, err)
		return
	}
}

func (s *topicServiceServer) serveGetTopicList(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveGetTopicListJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveGetTopicListProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *topicServiceServer) serveGetTopicListJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "GetTopicList")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(GetTopicListRequest)
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
	var respContent *GetTopicListResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.GetTopicList(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *GetTopicListResponse and nil error while calling GetTopicList. nil responses are not supported"))
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

func (s *topicServiceServer) serveGetTopicListProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "GetTopicList")
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
	reqContent := new(GetTopicListRequest)
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
	var respContent *GetTopicListResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.GetTopicList(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *GetTopicListResponse and nil error while calling GetTopicList. nil responses are not supported"))
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

func (s *topicServiceServer) serveGetTopicVideoList(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveGetTopicVideoListJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveGetTopicVideoListProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *topicServiceServer) serveGetTopicVideoListJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "GetTopicVideoList")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(GetTopicVideoListRequest)
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
	var respContent *GetTopicVideoListResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.GetTopicVideoList(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *GetTopicVideoListResponse and nil error while calling GetTopicVideoList. nil responses are not supported"))
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

func (s *topicServiceServer) serveGetTopicVideoListProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "GetTopicVideoList")
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
	reqContent := new(GetTopicVideoListRequest)
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
	var respContent *GetTopicVideoListResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.GetTopicVideoList(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *GetTopicVideoListResponse and nil error while calling GetTopicVideoList. nil responses are not supported"))
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

func (s *topicServiceServer) serveSubscribeTopic(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveSubscribeTopicJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveSubscribeTopicProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *topicServiceServer) serveSubscribeTopicJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "SubscribeTopic")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(SubscribeTopicRequest)
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
	var respContent *SubscribeTopicResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.SubscribeTopic(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *SubscribeTopicResponse and nil error while calling SubscribeTopic. nil responses are not supported"))
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

func (s *topicServiceServer) serveSubscribeTopicProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "SubscribeTopic")
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
	reqContent := new(SubscribeTopicRequest)
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
	var respContent *SubscribeTopicResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.SubscribeTopic(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *SubscribeTopicResponse and nil error while calling SubscribeTopic. nil responses are not supported"))
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

func (s *topicServiceServer) serveGetSubscribedTopicList(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveGetSubscribedTopicListJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveGetSubscribedTopicListProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *topicServiceServer) serveGetSubscribedTopicListJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "GetSubscribedTopicList")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(GetSubscribedTopicListRequest)
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
	var respContent *GetSubscribedTopicListResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.GetSubscribedTopicList(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *GetSubscribedTopicListResponse and nil error while calling GetSubscribedTopicList. nil responses are not supported"))
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

func (s *topicServiceServer) serveGetSubscribedTopicListProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "GetSubscribedTopicList")
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
	reqContent := new(GetSubscribedTopicListRequest)
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
	var respContent *GetSubscribedTopicListResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.GetSubscribedTopicList(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *GetSubscribedTopicListResponse and nil error while calling GetSubscribedTopicList. nil responses are not supported"))
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

func (s *topicServiceServer) ServiceDescriptor() ([]byte, int) {
	return twirpFileDescriptor10, 0
}

func (s *topicServiceServer) ProtocGenTwirpVersion() string {
	return "v5.4.1"
}

var twirpFileDescriptor10 = []byte{
	// 682 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0x5d, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0x71, 0x9a, 0x9f, 0x49, 0xd2, 0xb4, 0x4b, 0x29, 0xa9, 0x4b, 0xa1, 0xb2, 0xa8, 0x54,
	0xf5, 0xc1, 0xa8, 0x81, 0x1e, 0x80, 0x20, 0x54, 0x82, 0x28, 0xaa, 0x9c, 0xa6, 0x42, 0xbc, 0x04,
	0xff, 0xac, 0x94, 0x05, 0xff, 0xe1, 0x5d, 0x57, 0x88, 0x1b, 0xc0, 0x1b, 0x57, 0xe0, 0x99, 0x23,
	0x20, 0x24, 0x0e, 0xc3, 0x3d, 0x90, 0xc7, 0xeb, 0xa4, 0x71, 0x12, 0xa1, 0x54, 0xe2, 0x6d, 0x3d,
	0xdf, 0xec, 0xcc, 0x37, 0xdf, 0xcc, 0xac, 0xa1, 0x21, 0xc2, 0x88, 0x39, 0x46, 0x14, 0x87, 0x22,
	0x24, 0x6b, 0x76, 0xe2, 0x5a, 0xa1, 0xd6, 0x74, 0x42, 0xdf, 0x0f, 0x83, 0xcc, 0xa8, 0x7f, 0x55,
	0x00, 0x7a, 0x56, 0x10, 0xd0, 0xb8, 0x2f, 0xa8, 0x4f, 0x76, 0xa1, 0x6e, 0xe3, 0xd7, 0x88, 0xb9,
	0x1d, 0x65, 0x5f, 0x39, 0xac, 0x9b, 0xb5, 0xcc, 0xd0, 0x77, 0xc9, 0x5d, 0xa8, 0x46, 0xcc, 0x19,
	0x25, 0xb1, 0xd7, 0x29, 0x21, 0x54, 0x89, 0x98, 0x33, 0x8c, 0xbd, 0x1c, 0xf0, 0xdd, 0x93, 0x8e,
	0x3a, 0x01, 0xce, 0xdc, 0x13, 0xb2, 0x03, 0x35, 0x8f, 0x05, 0x1f, 0xf0, 0x4a, 0x19, 0x91, 0x6a,
	0xfa, 0x9d, 0xde, 0x21, 0x50, 0x76, 0x29, 0x77, 0x3a, 0x6b, 0x68, 0xc6, 0xb3, 0xfe, 0x0e, 0x6e,
	0x9f, 0x52, 0x71, 0x91, 0x72, 0x7e, 0xc5, 0xb8, 0x30, 0xe9, 0xc7, 0x84, 0x72, 0x41, 0x0e, 0xa0,
	0x32, 0xa6, 0x96, 0x4b, 0x63, 0x64, 0xd4, 0xe8, 0xb6, 0x0c, 0xac, 0xc4, 0x78, 0x81, 0x46, 0x53,
	0x82, 0x44, 0x87, 0x96, 0x67, 0x71, 0x31, 0xc2, 0x9a, 0x53, 0xfe, 0x19, 0xc9, 0x46, 0x6a, 0xc4,
	0x98, 0x7d, 0x57, 0xff, 0xad, 0xc0, 0xd6, 0x6c, 0x0a, 0x1e, 0x85, 0x01, 0xa7, 0x69, 0x0e, 0x2e,
	0x2c, 0x91, 0xf0, 0x42, 0x8e, 0x01, 0x1a, 0x4d, 0x09, 0x92, 0x63, 0x29, 0xe9, 0x88, 0x09, 0xea,
	0xf3, 0x4e, 0x69, 0x5f, 0x3d, 0x6c, 0x74, 0x37, 0xa4, 0x6f, 0x96, 0x44, 0x50, 0xdf, 0x04, 0x91,
	0x1f, 0x79, 0xaa, 0xc1, 0xd8, 0xe2, 0x23, 0x3f, 0x8c, 0x29, 0xaa, 0x53, 0x33, 0xab, 0x63, 0x8b,
	0x9f, 0x85, 0x31, 0x25, 0x4f, 0xa0, 0x99, 0xab, 0x8d, 0xe1, 0xca, 0x18, 0x6e, 0x53, 0x86, 0x9b,
	0xb6, 0xc5, 0x6c, 0xd8, 0x93, 0x33, 0xd7, 0x7f, 0x28, 0xd0, 0xc9, 0x6b, 0xb8, 0x64, 0x2e, 0x0d,
	0x6f, 0xa0, 0xd5, 0x0e, 0xd4, 0x0a, 0x32, 0x55, 0x33, 0xca, 0xee, 0x44, 0xc6, 0xab, 0x34, 0x74,
	0x8a, 0xab, 0x53, 0x19, 0x31, 0x5d, 0xdf, 0x25, 0x47, 0xb0, 0x89, 0x3e, 0x51, 0x62, 0x7b, 0x8c,
	0x8f, 0x47, 0x82, 0xf9, 0x14, 0x1b, 0x5c, 0x36, 0xdb, 0x29, 0x70, 0x9e, 0xd9, 0x2f, 0x98, 0x4f,
	0xf5, 0x5f, 0x0a, 0xec, 0x2c, 0xa0, 0xbb, 0x9a, 0xee, 0x8f, 0x00, 0xa6, 0xba, 0x23, 0xe3, 0x45,
	0xb2, 0xd7, 0x27, 0xb2, 0x13, 0x03, 0xc0, 0x63, 0x5c, 0x48, 0x61, 0x55, 0x14, 0xb6, 0x2d, 0x2f,
	0xa4, 0x04, 0x32, 0x7f, 0x4f, 0x9e, 0x66, 0xbb, 0x54, 0x9e, 0xe9, 0x92, 0xfe, 0x45, 0x81, 0x3b,
	0x83, 0xc4, 0xe6, 0x4e, 0xcc, 0x6c, 0x8a, 0xc9, 0x56, 0x14, 0xdb, 0x80, 0x8a, 0xe5, 0x08, 0x16,
	0x06, 0x48, 0x7c, 0xbd, 0xbb, 0x9d, 0xd7, 0x98, 0x07, 0x7d, 0x8a, 0xa8, 0x29, 0xbd, 0x66, 0x9a,
	0xa3, 0xce, 0x34, 0x47, 0x1f, 0xc2, 0x76, 0x91, 0xca, 0x6a, 0x42, 0x6e, 0xc1, 0x9a, 0x13, 0x26,
	0x81, 0x40, 0x2a, 0x2d, 0x33, 0xfb, 0xd0, 0xbf, 0x29, 0xb0, 0x77, 0x4a, 0xc5, 0x24, 0xb4, 0x7b,
	0xd3, 0x1d, 0x3c, 0x80, 0x75, 0x9e, 0x07, 0xc9, 0xa6, 0xa2, 0x84, 0x53, 0xd1, 0x9a, 0x58, 0xd3,
	0x99, 0x98, 0x5f, 0x55, 0x75, 0x7e, 0x55, 0x7f, 0x2a, 0x70, 0x7f, 0x19, 0xa7, 0xff, 0xbe, 0xb4,
	0xf3, 0x75, 0xa8, 0x8b, 0xea, 0x58, 0x3e, 0x35, 0x47, 0xc7, 0xd0, 0x2e, 0xf4, 0x97, 0xb4, 0xa0,
	0x3e, 0x18, 0xf6, 0x06, 0xcf, 0xcc, 0x7e, 0xef, 0xf9, 0xc6, 0x2d, 0xd2, 0x86, 0xc6, 0xf0, 0xf5,
	0xd4, 0xa0, 0x74, 0xff, 0x94, 0xa0, 0x89, 0x74, 0x06, 0x34, 0xbe, 0x62, 0x0e, 0x25, 0xa7, 0xd0,
	0xbc, 0xfe, 0x58, 0x11, 0x4d, 0x72, 0x5e, 0xf0, 0x48, 0x6a, 0xbb, 0x0b, 0x31, 0x29, 0xd4, 0x25,
	0x6c, 0xce, 0xad, 0x20, 0x79, 0x50, 0xb8, 0x51, 0x7c, 0x4b, 0xb4, 0xfd, 0xe5, 0x0e, 0x32, 0xee,
	0x19, 0xac, 0xcf, 0x8e, 0x23, 0xb9, 0x57, 0x9c, 0xed, 0xeb, 0x0b, 0xa3, 0xed, 0x2d, 0x41, 0x65,
	0x38, 0x0a, 0xdb, 0x8b, 0x3b, 0x4e, 0x1e, 0x4e, 0xa9, 0x2c, 0x1f, 0x52, 0xed, 0xe0, 0x1f, 0x5e,
	0x59, 0x9a, 0x9e, 0x06, 0x1b, 0x4e, 0xe8, 0x1b, 0xef, 0x3f, 0x7d, 0x36, 0x22, 0x3b, 0xbb, 0x72,
	0xae, 0x7c, 0x2f, 0xa9, 0x2f, 0xdf, 0xbc, 0xb5, 0x2b, 0xf8, 0x5b, 0x7c, 0xfc, 0x37, 0x00, 0x00,
	0xff, 0xff, 0x42, 0x73, 0xd7, 0x72, 0x3a, 0x07, 0x00, 0x00,
}
