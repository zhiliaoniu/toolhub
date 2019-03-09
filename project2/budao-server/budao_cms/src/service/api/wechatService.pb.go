// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wechatService.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// wechat回调
type WechatRedirectRequest struct {
	Code string `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
}

func (m *WechatRedirectRequest) Reset()                    { *m = WechatRedirectRequest{} }
func (m *WechatRedirectRequest) String() string            { return proto.CompactTextString(m) }
func (*WechatRedirectRequest) ProtoMessage()               {}
func (*WechatRedirectRequest) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{0} }

func (m *WechatRedirectRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type WechatRedirectResponse struct {
	Code string      `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Msg  string      `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Data *BackStruct `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *WechatRedirectResponse) Reset()                    { *m = WechatRedirectResponse{} }
func (m *WechatRedirectResponse) String() string            { return proto.CompactTextString(m) }
func (*WechatRedirectResponse) ProtoMessage()               {}
func (*WechatRedirectResponse) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{1} }

func (m *WechatRedirectResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *WechatRedirectResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *WechatRedirectResponse) GetData() *BackStruct {
	if m != nil {
		return m.Data
	}
	return nil
}

type BackStruct struct {
	UserId string `protobuf:"bytes,1,opt,name=userId" json:"userId,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token" json:"token,omitempty"`
}

func (m *BackStruct) Reset()                    { *m = BackStruct{} }
func (m *BackStruct) String() string            { return proto.CompactTextString(m) }
func (*BackStruct) ProtoMessage()               {}
func (*BackStruct) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{2} }

func (m *BackStruct) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *BackStruct) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func init() {
	proto.RegisterType((*WechatRedirectRequest)(nil), "api.WechatRedirectRequest")
	proto.RegisterType((*WechatRedirectResponse)(nil), "api.WechatRedirectResponse")
	proto.RegisterType((*BackStruct)(nil), "api.backStruct")
}

func init() { proto.RegisterFile("wechatService.proto", fileDescriptor13) }

var fileDescriptor13 = []byte{
	// 209 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x50, 0xcb, 0x4a, 0xc5, 0x30,
	0x10, 0xa5, 0xf6, 0x7a, 0xc1, 0x11, 0x1f, 0x8c, 0x5a, 0x4a, 0xdd, 0x94, 0xba, 0x29, 0x08, 0x5d,
	0xd4, 0x9d, 0x7f, 0xd0, 0x6d, 0xba, 0x10, 0xdc, 0xa5, 0xc9, 0xa0, 0xa1, 0xd8, 0xc4, 0x64, 0xaa,
	0xbf, 0x2f, 0xa6, 0x05, 0xa9, 0x74, 0x77, 0x5e, 0x9c, 0xe4, 0x0c, 0xdc, 0x7c, 0x93, 0x7a, 0x97,
	0xdc, 0x93, 0xff, 0x32, 0x8a, 0x1a, 0xe7, 0x2d, 0x5b, 0x4c, 0xa5, 0x33, 0xd5, 0x23, 0xdc, 0xbd,
	0x44, 0x4f, 0x90, 0x36, 0x9e, 0x14, 0x0b, 0xfa, 0x9c, 0x29, 0x30, 0x22, 0x1c, 0x94, 0xd5, 0x94,
	0x27, 0x65, 0x52, 0x9f, 0x89, 0x88, 0x2b, 0x05, 0xd9, 0xff, 0x70, 0x70, 0x76, 0x0a, 0xb4, 0x97,
	0xc6, 0x6b, 0x48, 0x3f, 0xc2, 0x5b, 0x7e, 0x12, 0xa5, 0x5f, 0x88, 0x0f, 0x70, 0xd0, 0x92, 0x65,
	0x9e, 0x96, 0x49, 0x7d, 0xde, 0x5e, 0x35, 0xd2, 0x99, 0x66, 0x90, 0x6a, 0xec, 0xd9, 0xcf, 0x8a,
	0x45, 0x34, 0xab, 0x67, 0x80, 0x3f, 0x0d, 0x33, 0x38, 0xce, 0x81, 0x7c, 0xa7, 0xd7, 0xea, 0x95,
	0xe1, 0x2d, 0x9c, 0xb2, 0x1d, 0x69, 0x5a, 0xeb, 0x17, 0xd2, 0xbe, 0xc2, 0xc5, 0x66, 0x29, 0x76,
	0x70, 0xb9, 0xfd, 0x31, 0x16, 0xf1, 0xd5, 0xdd, 0xcd, 0xc5, 0xfd, 0xae, 0xb7, 0x4c, 0x1c, 0x8e,
	0xf1, 0x6a, 0x4f, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xe1, 0xd2, 0xf4, 0x33, 0x4c, 0x01, 0x00,
	0x00,
}
