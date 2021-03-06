// Code generated by protoc-gen-twirp v5.4.1, DO NOT EDIT.
// source: parseurl.proto

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

// =========================
// ParseURLService Interface
// =========================

type ParseURLService interface {
	ParseURL(context.Context, *ParseURLRequest) (*ParseURLResponse, error)

	ParseExternalURL(context.Context, *ParseExternalURLRequest) (*ParseURLResponse, error)
}

// ===============================
// ParseURLService Protobuf Client
// ===============================

type parseURLServiceProtobufClient struct {
	client HTTPClient
	urls   [2]string
}

// NewParseURLServiceProtobufClient creates a Protobuf client that implements the ParseURLService interface.
// It communicates using Protobuf and can be configured with a custom HTTPClient.
func NewParseURLServiceProtobufClient(addr string, client HTTPClient) ParseURLService {
	prefix := urlBase(addr) + ParseURLServicePathPrefix
	urls := [2]string{
		prefix + "ParseURL",
		prefix + "ParseExternalURL",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &parseURLServiceProtobufClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &parseURLServiceProtobufClient{
		client: client,
		urls:   urls,
	}
}

func (c *parseURLServiceProtobufClient) ParseURL(ctx context.Context, in *ParseURLRequest) (*ParseURLResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "ParseURLService")
	ctx = ctxsetters.WithMethodName(ctx, "ParseURL")
	out := new(ParseURLResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[0], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *parseURLServiceProtobufClient) ParseExternalURL(ctx context.Context, in *ParseExternalURLRequest) (*ParseURLResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "ParseURLService")
	ctx = ctxsetters.WithMethodName(ctx, "ParseExternalURL")
	out := new(ParseURLResponse)
	err := doProtobufRequest(ctx, c.client, c.urls[1], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ===========================
// ParseURLService JSON Client
// ===========================

type parseURLServiceJSONClient struct {
	client HTTPClient
	urls   [2]string
}

// NewParseURLServiceJSONClient creates a JSON client that implements the ParseURLService interface.
// It communicates using JSON and can be configured with a custom HTTPClient.
func NewParseURLServiceJSONClient(addr string, client HTTPClient) ParseURLService {
	prefix := urlBase(addr) + ParseURLServicePathPrefix
	urls := [2]string{
		prefix + "ParseURL",
		prefix + "ParseExternalURL",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &parseURLServiceJSONClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &parseURLServiceJSONClient{
		client: client,
		urls:   urls,
	}
}

func (c *parseURLServiceJSONClient) ParseURL(ctx context.Context, in *ParseURLRequest) (*ParseURLResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "ParseURLService")
	ctx = ctxsetters.WithMethodName(ctx, "ParseURL")
	out := new(ParseURLResponse)
	err := doJSONRequest(ctx, c.client, c.urls[0], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *parseURLServiceJSONClient) ParseExternalURL(ctx context.Context, in *ParseExternalURLRequest) (*ParseURLResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "ParseURLService")
	ctx = ctxsetters.WithMethodName(ctx, "ParseExternalURL")
	out := new(ParseURLResponse)
	err := doJSONRequest(ctx, c.client, c.urls[1], in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ==============================
// ParseURLService Server Handler
// ==============================

type parseURLServiceServer struct {
	ParseURLService
	hooks *twirp.ServerHooks
}

func NewParseURLServiceServer(svc ParseURLService, hooks *twirp.ServerHooks) TwirpServer {
	return &parseURLServiceServer{
		ParseURLService: svc,
		hooks:           hooks,
	}
}

// writeError writes an HTTP response with a valid Twirp error format, and triggers hooks.
// If err is not a twirp.Error, it will get wrapped with twirp.InternalErrorWith(err)
func (s *parseURLServiceServer) writeError(ctx context.Context, resp http.ResponseWriter, err error) {
	writeError(ctx, resp, err, s.hooks)
}

// ParseURLServicePathPrefix is used for all URL paths on a twirp ParseURLService server.
// Requests are always: POST ParseURLServicePathPrefix/method
// It can be used in an HTTP mux to route twirp requests along with non-twirp requests on other routes.
const ParseURLServicePathPrefix = "/budao.ParseURLService/"

func (s *parseURLServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "ParseURLService")
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
	case "/budao.ParseURLService/ParseURL":
		s.serveParseURL(ctx, resp, req)
		return
	case "/budao.ParseURLService/ParseExternalURL":
		s.serveParseExternalURL(ctx, resp, req)
		return
	default:
		msg := fmt.Sprintf("no handler for path %q", req.URL.Path)
		err = badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, err)
		return
	}
}

func (s *parseURLServiceServer) serveParseURL(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveParseURLJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveParseURLProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *parseURLServiceServer) serveParseURLJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "ParseURL")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(ParseURLRequest)
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
	var respContent *ParseURLResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.ParseURL(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *ParseURLResponse and nil error while calling ParseURL. nil responses are not supported"))
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

func (s *parseURLServiceServer) serveParseURLProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "ParseURL")
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
	reqContent := new(ParseURLRequest)
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
	var respContent *ParseURLResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.ParseURL(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *ParseURLResponse and nil error while calling ParseURL. nil responses are not supported"))
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

func (s *parseURLServiceServer) serveParseExternalURL(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveParseExternalURLJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveParseExternalURLProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *parseURLServiceServer) serveParseExternalURLJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "ParseExternalURL")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(ParseExternalURLRequest)
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
	var respContent *ParseURLResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.ParseExternalURL(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *ParseURLResponse and nil error while calling ParseExternalURL. nil responses are not supported"))
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

func (s *parseURLServiceServer) serveParseExternalURLProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "ParseExternalURL")
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
	reqContent := new(ParseExternalURLRequest)
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
	var respContent *ParseURLResponse
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.ParseExternalURL(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *ParseURLResponse and nil error while calling ParseExternalURL. nil responses are not supported"))
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

func (s *parseURLServiceServer) ServiceDescriptor() ([]byte, int) {
	return twirpFileDescriptor4, 0
}

func (s *parseURLServiceServer) ProtocGenTwirpVersion() string {
	return "v5.4.1"
}

var twirpFileDescriptor4 = []byte{
	// 723 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xdb, 0x6e, 0xd3, 0x40,
	0x10, 0xad, 0x13, 0xe7, 0x36, 0x49, 0xdb, 0xed, 0xf6, 0x92, 0x34, 0x45, 0xa5, 0x54, 0x42, 0x8a,
	0x78, 0x08, 0x52, 0x78, 0xe1, 0xf2, 0x82, 0x43, 0x4d, 0xbc, 0xb4, 0xb5, 0x23, 0xa7, 0xa6, 0x97,
	0x97, 0xc8, 0x4e, 0xb6, 0x22, 0xe0, 0xc4, 0xc1, 0x97, 0xa6, 0xe1, 0x1f, 0xf8, 0x02, 0x5e, 0x10,
	0x5f, 0xc3, 0x3f, 0xf0, 0x0d, 0xfc, 0x03, 0xda, 0xb5, 0xd3, 0x6e, 0x53, 0x5a, 0x09, 0xf1, 0xb6,
	0x7b, 0xf6, 0xcc, 0xec, 0x39, 0x33, 0xa3, 0x81, 0xa5, 0xb1, 0xed, 0x07, 0x34, 0xf2, 0xdd, 0xfa,
	0xd8, 0xf7, 0x42, 0x0f, 0x67, 0x9c, 0xa8, 0x6f, 0x7b, 0xd5, 0x52, 0xcf, 0x1b, 0x0e, 0xbd, 0x51,
	0x0c, 0xee, 0x76, 0x60, 0xb9, 0xcd, 0x68, 0x96, 0x79, 0x60, 0xd2, 0xcf, 0x11, 0x0d, 0x42, 0xfc,
	0x18, 0xb2, 0x1f, 0xa8, 0xdd, 0xa7, 0x7e, 0x45, 0xda, 0x91, 0x6a, 0xc5, 0xc6, 0x62, 0x9d, 0x07,
	0xd6, 0x35, 0x0e, 0x9a, 0xc9, 0x23, 0xde, 0x84, 0xfc, 0xc5, 0xa0, 0x4f, 0xbd, 0xee, 0xa0, 0x5f,
	0x49, 0xed, 0x48, 0xb5, 0x82, 0x99, 0xe3, 0x77, 0xd2, 0xdf, 0xfd, 0x25, 0x41, 0x99, 0x67, 0x55,
	0x2f, 0x43, 0xea, 0x8f, 0x6c, 0xf7, 0xdf, 0xb3, 0x63, 0x90, 0x7b, 0x5e, 0x9f, 0xf2, 0xcc, 0x8b,
	0x26, 0x3f, 0x33, 0xcc, 0xf1, 0xfa, 0xd3, 0x4a, 0x9a, 0xff, 0xc6, 0xcf, 0x0c, 0x3b, 0x8f, 0x46,
	0xbd, 0x8a, 0x1c, 0x63, 0xec, 0x7c, 0x43, 0x59, 0xe6, 0x86, 0x32, 0xfc, 0x10, 0x8a, 0x7e, 0x2c,
	0xa4, 0x3b, 0x8a, 0x86, 0x95, 0x2c, 0xcf, 0x0e, 0x09, 0xa4, 0x47, 0x43, 0xfc, 0x08, 0x4a, 0x33,
	0x42, 0x38, 0x18, 0xd2, 0x4a, 0x8e, 0x33, 0x66, 0x41, 0x47, 0x83, 0x21, 0xdd, 0xfd, 0x9a, 0x01,
	0x74, 0x5d, 0xb3, 0x60, 0xec, 0x8d, 0x02, 0xca, 0x6c, 0x05, 0xa1, 0x1d, 0x46, 0xc1, 0x9c, 0xad,
	0x0e, 0x07, 0xcd, 0xe4, 0x11, 0x37, 0x40, 0x1e, 0xbb, 0xf6, 0x94, 0xdb, 0x2a, 0x36, 0x1e, 0x24,
	0xa4, 0xf9, 0x6c, 0xf5, 0xb6, 0x6b, 0x4f, 0xb5, 0x05, 0x93, 0x73, 0x59, 0x4c, 0xcf, 0xb7, 0x27,
	0xdc, 0xf6, 0x3d, 0x31, 0x6f, 0x7c, 0x7b, 0xc2, 0x62, 0x18, 0x17, 0x3f, 0x85, 0xf4, 0x84, 0x3a,
	0xbc, 0x2a, 0xc5, 0xc6, 0xd6, 0x5d, 0x21, 0xc7, 0xd4, 0xd1, 0x16, 0x4c, 0xc6, 0xac, 0x7e, 0x97,
	0x40, 0x66, 0xbf, 0xe2, 0x2d, 0x28, 0xc4, 0xc5, 0x8b, 0x7c, 0x97, 0x7b, 0x29, 0x98, 0x71, 0x35,
	0x2d, 0xdf, 0xc5, 0xaf, 0xaf, 0x9a, 0x97, 0xda, 0x49, 0xd7, 0x8a, 0x8d, 0xda, 0x7d, 0x06, 0x92,
	0x96, 0xaa, 0xa3, 0xd0, 0x9f, 0xce, 0xfa, 0x5a, 0x7d, 0x01, 0x45, 0x01, 0xc6, 0x08, 0xd2, 0x9f,
	0xe8, 0x34, 0xf9, 0x87, 0x1d, 0xf1, 0x1a, 0x64, 0x2e, 0x6c, 0x37, 0xa2, 0xc9, 0x4c, 0xc5, 0x97,
	0x97, 0xa9, 0xe7, 0x52, 0xf5, 0xb7, 0x04, 0x32, 0x33, 0xc9, 0x7a, 0x3e, 0xa2, 0x97, 0x61, 0x12,
	0xc5, 0xcf, 0x37, 0x65, 0xa7, 0xe6, 0x64, 0x0b, 0x4d, 0x1d, 0x7b, 0x41, 0xc8, 0x2b, 0x99, 0xbf,
	0x6a, 0x6a, 0xdb, 0x0b, 0x42, 0xc1, 0x99, 0x7c, 0xbf, 0x33, 0xa6, 0xe0, 0x6f, 0xce, 0x98, 0xf0,
	0xf3, 0x01, 0x75, 0xe3, 0x91, 0x2b, 0x99, 0xf1, 0xe5, 0x7f, 0xfc, 0x6e, 0x43, 0xfa, 0x98, 0x3a,
	0xb8, 0x0c, 0xb9, 0x09, 0x75, 0x84, 0x76, 0x64, 0x27, 0xd4, 0xb1, 0x7c, 0xb7, 0x99, 0x05, 0xf9,
	0x68, 0x3a, 0xa6, 0x4f, 0x7e, 0xa6, 0x60, 0xe9, 0x3d, 0xb3, 0xca, 0xb5, 0x9a, 0x91, 0x4b, 0xf1,
	0x1a, 0xa0, 0xb6, 0x62, 0x76, 0xd4, 0xae, 0x69, 0x1d, 0xa8, 0xdd, 0x13, 0xd2, 0xb2, 0x14, 0xb4,
	0x80, 0xd7, 0x61, 0x45, 0x40, 0xf7, 0x0c, 0xeb, 0x94, 0xe8, 0x48, 0xc2, 0x65, 0x58, 0x15, 0xe0,
	0x7d, 0x4b, 0x21, 0x1d, 0xcd, 0xb0, 0x50, 0x0a, 0x6f, 0x41, 0x59, 0x78, 0xd0, 0x55, 0xa2, 0x29,
	0xfa, 0x9e, 0xa5, 0xe8, 0x67, 0x04, 0xa5, 0xe7, 0x92, 0x69, 0x8a, 0xb1, 0xaf, 0xe8, 0x48, 0xc6,
	0xab, 0xb0, 0x2c, 0xc0, 0x4d, 0xa3, 0x69, 0xa0, 0x0c, 0xde, 0x00, 0x2c, 0x80, 0x87, 0x44, 0x31,
	0xda, 0x0a, 0x41, 0xd9, 0xb9, 0x1c, 0x87, 0x2a, 0x61, 0x70, 0x6e, 0x4e, 0x7d, 0xf3, 0x4c, 0x53,
	0x74, 0x94, 0x9f, 0x43, 0x8f, 0x55, 0xd2, 0x34, 0x50, 0x01, 0xaf, 0xc0, 0xa2, 0x80, 0x1e, 0x10,
	0x04, 0x78, 0x13, 0xd6, 0x05, 0xe8, 0xad, 0x4a, 0xf6, 0x88, 0xda, 0xd1, 0x2c, 0x03, 0x15, 0xf1,
	0x36, 0x54, 0x85, 0xa7, 0x53, 0x45, 0x6f, 0xb5, 0x2c, 0x45, 0x6f, 0x75, 0x34, 0xd2, 0x26, 0x3a,
	0x2a, 0x35, 0xbe, 0x49, 0xd7, 0xeb, 0xb0, 0x43, 0xfd, 0x8b, 0x41, 0x8f, 0xe2, 0x57, 0x90, 0x9f,
	0x41, 0x78, 0xe3, 0xd6, 0x54, 0xf0, 0x09, 0xaa, 0x96, 0xef, 0x98, 0x16, 0xbc, 0x9f, 0xac, 0x0a,
	0x61, 0x11, 0xe2, 0x6d, 0x91, 0x7c, 0x7b, 0x43, 0xde, 0x99, 0xac, 0x59, 0x05, 0xd4, 0xf3, 0x86,
	0xf5, 0x8f, 0x97, 0x5f, 0xea, 0x63, 0x27, 0x26, 0xb5, 0xa5, 0x1f, 0xa9, 0xf4, 0xbb, 0x93, 0x33,
	0x27, 0xcb, 0xd7, 0xf9, 0xb3, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x09, 0x3d, 0xa5, 0xf5,
	0x05, 0x00, 0x00,
}
