// Code generated by protoc-gen-twirp v5.4.1, DO NOT EDIT.
// source: keywordService.proto

package api

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

// ========================
// KeywordService Interface
// ========================

type KeywordService interface {
	QueryVideoTitle(context.Context, *QueryListRequest) (*VideoListResponse, error)

	ReplaceVideoTitle(context.Context, *OptionVideoTitleReq) (*CommonResponse, error)

	QueryCommentContent(context.Context, *QueryCommentContentReq) (*CommentListResponse, error)

	OptionCommentContent(context.Context, *OptionCommentContentReq) (*CommonResponse, error)
}

// ==============================
// KeywordService Protobuf Client
// ==============================

type keywordServiceProtobufClient struct {
	client HTTPClient
	urls   [4]string
}

// NewKeywordServiceProtobufClient creates a Protobuf client that implements the KeywordService interface.
// It communicates using Protobuf and can be configured with a custom HTTPClient.
func NewKeywordServiceProtobufClient(addr string, client HTTPClient) KeywordService {
	prefix := urlBase(addr) + KeywordServicePathPrefix
	urls := [4]string{
		prefix + "QueryVideoTitle",
		prefix + "ReplaceVideoTitle",
		prefix + "QueryCommentContent",
		prefix + "OptionCommentContent",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &keywordServiceProtobufClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &keywordServiceProtobufClient{
		client: client,
		urls:   urls,
	}
}

func (c *keywordServiceProtobufClient) QueryVideoTitle(ctx context.Context, in *QueryListRequest) (*VideoListResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "api")
	ctx = ctxsetters.WithServiceName(ctx, "KeywordService")
	ctx = ctxsetters.WithMethodName(ctx, "QueryVideoTitle")
	out := new(VideoListResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[0], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keywordServiceProtobufClient) ReplaceVideoTitle(ctx context.Context, in *OptionVideoTitleReq) (*CommonResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "api")
	ctx = ctxsetters.WithServiceName(ctx, "KeywordService")
	ctx = ctxsetters.WithMethodName(ctx, "ReplaceVideoTitle")
	out := new(CommonResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[1], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keywordServiceProtobufClient) QueryCommentContent(ctx context.Context, in *QueryCommentContentReq) (*CommentListResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "api")
	ctx = ctxsetters.WithServiceName(ctx, "KeywordService")
	ctx = ctxsetters.WithMethodName(ctx, "QueryCommentContent")
	out := new(CommentListResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[2], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keywordServiceProtobufClient) OptionCommentContent(ctx context.Context, in *OptionCommentContentReq) (*CommonResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "api")
	ctx = ctxsetters.WithServiceName(ctx, "KeywordService")
	ctx = ctxsetters.WithMethodName(ctx, "OptionCommentContent")
	out := new(CommonResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[3], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ==========================
// KeywordService JSON Client
// ==========================

type keywordServiceJSONClient struct {
	client HTTPClient
	urls   [4]string
}

// NewKeywordServiceJSONClient creates a JSON client that implements the KeywordService interface.
// It communicates using JSON and can be configured with a custom HTTPClient.
func NewKeywordServiceJSONClient(addr string, client HTTPClient) KeywordService {
	prefix := urlBase(addr) + KeywordServicePathPrefix
	urls := [4]string{
		prefix + "QueryVideoTitle",
		prefix + "ReplaceVideoTitle",
		prefix + "QueryCommentContent",
		prefix + "OptionCommentContent",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &keywordServiceJSONClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &keywordServiceJSONClient{
		client: client,
		urls:   urls,
	}
}

func (c *keywordServiceJSONClient) QueryVideoTitle(ctx context.Context, in *QueryListRequest) (*VideoListResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "api")
	ctx = ctxsetters.WithServiceName(ctx, "KeywordService")
	ctx = ctxsetters.WithMethodName(ctx, "QueryVideoTitle")
	out := new(VideoListResponse)
	err := doJSONRequest(ctx, c.client, c.urls[0], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keywordServiceJSONClient) ReplaceVideoTitle(ctx context.Context, in *OptionVideoTitleReq) (*CommonResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "api")
	ctx = ctxsetters.WithServiceName(ctx, "KeywordService")
	ctx = ctxsetters.WithMethodName(ctx, "ReplaceVideoTitle")
	out := new(CommonResponse)
	err := doJSONRequest(ctx, c.client, c.urls[1], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keywordServiceJSONClient) QueryCommentContent(ctx context.Context, in *QueryCommentContentReq) (*CommentListResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "api")
	ctx = ctxsetters.WithServiceName(ctx, "KeywordService")
	ctx = ctxsetters.WithMethodName(ctx, "QueryCommentContent")
	out := new(CommentListResponse)
	err := doJSONRequest(ctx, c.client, c.urls[2], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keywordServiceJSONClient) OptionCommentContent(ctx context.Context, in *OptionCommentContentReq) (*CommonResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "api")
	ctx = ctxsetters.WithServiceName(ctx, "KeywordService")
	ctx = ctxsetters.WithMethodName(ctx, "OptionCommentContent")
	out := new(CommonResponse)
	err := doJSONRequest(ctx, c.client, c.urls[3], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// =============================
// KeywordService Server Handler
// =============================

type keywordServiceServer struct {
	KeywordService
	hooks *twirp.ServerHooks
}

func NewKeywordServiceServer(svc KeywordService, hooks *twirp.ServerHooks) TwirpServer {
	return &keywordServiceServer{
		KeywordService: svc,
		hooks:          hooks,
	}
}

// writeError writes an HTTP response with a valid Twirp error format, and triggers hooks.
// If err is not a twirp.Error, it will get wrapped with twirp.InternalErrorWith(err)
func (s *keywordServiceServer) writeError(ctx context.Context, resp http.ResponseWriter, err error) {
	writeError(ctx, resp, err, s.hooks)
}

// KeywordServicePathPrefix is used for all URL paths on a twirp KeywordService server.
// Requests are always: POST KeywordServicePathPrefix/method
// It can be used in an HTTP mux to route twirp requests along with non-twirp requests on other routes.
const KeywordServicePathPrefix = "/api.KeywordService/"

func (s *keywordServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = ctxsetters.WithPackageName(ctx, "api")
	ctx = ctxsetters.WithServiceName(ctx, "KeywordService")
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
		if req.Method == "OPTIONS" {
			resp.Header().Add("Access-Control-Allow-Origin", "*")
			resp.Header().Add("Access-Control-Allow-Credentials", "true")
			resp.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-Automatic-Token,X-Remote-Addr")
			return
		}
		msg := fmt.Sprintf("unsupported method %q (only POST is allowed)", req.Method)
		err = badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, err)
		return
	}

	switch req.URL.Path {
	case "/api.KeywordService/QueryVideoTitle":
		s.serveQueryVideoTitle(ctx, resp, req)
		return
	case "/api.KeywordService/ReplaceVideoTitle":
		s.serveReplaceVideoTitle(ctx, resp, req)
		return
	case "/api.KeywordService/QueryCommentContent":
		s.serveQueryCommentContent(ctx, resp, req)
		return
	case "/api.KeywordService/OptionCommentContent":
		s.serveOptionCommentContent(ctx, resp, req)
		return
	default:
		msg := fmt.Sprintf("no handler for path %q", req.URL.Path)
		err = badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, err)
		return
	}
}

func (s *keywordServiceServer) serveQueryVideoTitle(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveQueryVideoTitleJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveQueryVideoTitleProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *keywordServiceServer) serveQueryVideoTitleJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Access-Control-Allow-Origin", "*")
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "QueryVideoTitle")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(QueryListRequest)
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
	var respContent *VideoListResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.QueryVideoTitle(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *VideoListResponse and nil error while calling QueryVideoTitle. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: false}
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

func (s *keywordServiceServer) serveQueryVideoTitleProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "QueryVideoTitle")
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
	reqContent := new(QueryListRequest)
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
	var respContent *VideoListResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.QueryVideoTitle(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *VideoListResponse and nil error while calling QueryVideoTitle. nil responses are not supported"))
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

func (s *keywordServiceServer) serveReplaceVideoTitle(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveReplaceVideoTitleJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveReplaceVideoTitleProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *keywordServiceServer) serveReplaceVideoTitleJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Access-Control-Allow-Origin", "*")
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "ReplaceVideoTitle")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(OptionVideoTitleReq)
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
	var respContent *CommonResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.ReplaceVideoTitle(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *CommonResponse and nil error while calling ReplaceVideoTitle. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: false}
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

func (s *keywordServiceServer) serveReplaceVideoTitleProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "ReplaceVideoTitle")
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
	reqContent := new(OptionVideoTitleReq)
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
	var respContent *CommonResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.ReplaceVideoTitle(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *CommonResponse and nil error while calling ReplaceVideoTitle. nil responses are not supported"))
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

func (s *keywordServiceServer) serveQueryCommentContent(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveQueryCommentContentJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveQueryCommentContentProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *keywordServiceServer) serveQueryCommentContentJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Access-Control-Allow-Origin", "*")
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "QueryCommentContent")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(QueryCommentContentReq)
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
	var respContent *CommentListResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.QueryCommentContent(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *CommentListResponse and nil error while calling QueryCommentContent. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: false}
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

func (s *keywordServiceServer) serveQueryCommentContentProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "QueryCommentContent")
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
	reqContent := new(QueryCommentContentReq)
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
	var respContent *CommentListResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.QueryCommentContent(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *CommentListResponse and nil error while calling QueryCommentContent. nil responses are not supported"))
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

func (s *keywordServiceServer) serveOptionCommentContent(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveOptionCommentContentJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveOptionCommentContentProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *keywordServiceServer) serveOptionCommentContentJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Access-Control-Allow-Origin", "*")
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "OptionCommentContent")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(OptionCommentContentReq)
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
	var respContent *CommonResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.OptionCommentContent(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *CommonResponse and nil error while calling OptionCommentContent. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: false}
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

func (s *keywordServiceServer) serveOptionCommentContentProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "OptionCommentContent")
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
	reqContent := new(OptionCommentContentReq)
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
	var respContent *CommonResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.OptionCommentContent(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *CommonResponse and nil error while calling OptionCommentContent. nil responses are not supported"))
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

func (s *keywordServiceServer) ServiceDescriptor() ([]byte, int) {
	return twirpFileDescriptor3, 0
}

func (s *keywordServiceServer) ProtocGenTwirpVersion() string {
	return "v5.4.1"
}

var twirpFileDescriptor3 = []byte{
	// 301 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcf, 0x4e, 0x83, 0x40,
	0x10, 0xc6, 0xd3, 0x56, 0xab, 0x9d, 0x18, 0xff, 0x2c, 0x15, 0x57, 0xf4, 0xd0, 0x70, 0xea, 0x89,
	0x83, 0xbe, 0x80, 0x91, 0x93, 0xda, 0xc4, 0x88, 0x46, 0xbd, 0x19, 0xa4, 0x63, 0xb2, 0xb1, 0xec,
	0xac, 0xb0, 0x48, 0xea, 0x4b, 0xf9, 0x8a, 0x86, 0xa1, 0xf6, 0x4f, 0xca, 0x45, 0x6f, 0x3b, 0xdf,
	0x37, 0xf9, 0xcd, 0xc7, 0x0c, 0xd0, 0x7f, 0xc7, 0x69, 0x49, 0xd9, 0xf8, 0x1e, 0xb3, 0x4f, 0x95,
	0x60, 0x60, 0x32, 0xb2, 0x24, 0x3a, 0xb1, 0x51, 0xde, 0x4e, 0x42, 0x69, 0x4a, 0xba, 0x96, 0xfc,
	0x67, 0x70, 0xef, 0x0a, 0xcc, 0xa6, 0x21, 0xa5, 0x29, 0x6a, 0x1b, 0x92, 0xb6, 0xa8, 0x6d, 0x84,
	0x1f, 0x42, 0xc2, 0xd6, 0x0c, 0x22, 0x5b, 0x83, 0xd6, 0xb0, 0x17, 0xfd, 0x96, 0x62, 0x1f, 0x3a,
	0xba, 0x48, 0x65, 0x7b, 0xd0, 0x1a, 0x6e, 0x46, 0xd5, 0x53, 0x08, 0xd8, 0xc8, 0xd5, 0x17, 0xca,
	0x0e, 0x4b, 0xfc, 0xf6, 0xdf, 0xe0, 0xe8, 0xd6, 0x58, 0x45, 0xfa, 0x2f, 0xe8, 0x63, 0xd8, 0xd6,
	0x58, 0xbe, 0xb0, 0xd5, 0xae, 0x2d, 0x8d, 0xe5, 0x53, 0x65, 0xb9, 0xd0, 0x25, 0xe6, 0xf1, 0x94,
	0x5e, 0x34, 0xab, 0xfc, 0x6b, 0x70, 0xea, 0x39, 0x8f, 0x6a, 0x8c, 0xf4, 0xa0, 0xec, 0x04, 0xff,
	0x3b, 0xe3, 0xec, 0xbb, 0x0d, 0xbb, 0x37, 0x2b, 0x9b, 0x13, 0x17, 0xb0, 0xc7, 0x0b, 0x5a, 0xd0,
	0xc5, 0x61, 0x10, 0x1b, 0x15, 0xb0, 0x3a, 0x52, 0x79, 0xf5, 0x45, 0x05, 0xe6, 0xd6, 0x73, 0x59,
	0xe6, 0xbe, 0x5a, 0xce, 0x0d, 0xe9, 0x1c, 0xc5, 0x25, 0x1c, 0x44, 0x68, 0x26, 0x71, 0x82, 0x4b,
	0x0c, 0xc9, 0xcd, 0x0d, 0xc1, 0x3d, 0x87, 0x9d, 0x90, 0x8f, 0x34, 0x67, 0x8c, 0xc0, 0x69, 0x38,
	0x93, 0x38, 0x59, 0x24, 0x59, 0xdb, 0xb2, 0x27, 0xe7, 0x20, 0xd4, 0x76, 0x25, 0xd1, 0x15, 0xf4,
	0x9b, 0x4e, 0x23, 0x4e, 0x97, 0x42, 0xad, 0xf3, 0x9a, 0x82, 0xbd, 0x76, 0xf9, 0x37, 0x3a, 0xff,
	0x09, 0x00, 0x00, 0xff, 0xff, 0x8c, 0x36, 0x4d, 0xb4, 0x71, 0x02, 0x00, 0x00,
}
