// Code generated by protoc-gen-go. DO NOT EDIT.
// source: misc.proto

package budao

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ReportVideoRequest struct {
	Header  *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	VideoId string  `protobuf:"bytes,2,opt,name=video_id,json=videoId" json:"video_id,omitempty"`
	Reason  string  `protobuf:"bytes,3,opt,name=reason" json:"reason,omitempty"`
}

func (m *ReportVideoRequest) Reset()                    { *m = ReportVideoRequest{} }
func (m *ReportVideoRequest) String() string            { return proto.CompactTextString(m) }
func (*ReportVideoRequest) ProtoMessage()               {}
func (*ReportVideoRequest) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

func (m *ReportVideoRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ReportVideoRequest) GetVideoId() string {
	if m != nil {
		return m.VideoId
	}
	return ""
}

func (m *ReportVideoRequest) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type ReportVideoResponse struct {
	Status *Status `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *ReportVideoResponse) Reset()                    { *m = ReportVideoResponse{} }
func (m *ReportVideoResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportVideoResponse) ProtoMessage()               {}
func (*ReportVideoResponse) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{1} }

func (m *ReportVideoResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

type ReportCommentRequest struct {
	Header    *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	CommentId string  `protobuf:"bytes,2,opt,name=comment_id,json=commentId" json:"comment_id,omitempty"`
	Reason    string  `protobuf:"bytes,3,opt,name=reason" json:"reason,omitempty"`
}

func (m *ReportCommentRequest) Reset()                    { *m = ReportCommentRequest{} }
func (m *ReportCommentRequest) String() string            { return proto.CompactTextString(m) }
func (*ReportCommentRequest) ProtoMessage()               {}
func (*ReportCommentRequest) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{2} }

func (m *ReportCommentRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ReportCommentRequest) GetCommentId() string {
	if m != nil {
		return m.CommentId
	}
	return ""
}

func (m *ReportCommentRequest) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type ReportCommentResponse struct {
	Status *Status `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *ReportCommentResponse) Reset()                    { *m = ReportCommentResponse{} }
func (m *ReportCommentResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportCommentResponse) ProtoMessage()               {}
func (*ReportCommentResponse) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{3} }

func (m *ReportCommentResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

type UserFeedbackRequest struct {
	Header   *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Feedback string  `protobuf:"bytes,2,opt,name=feedback" json:"feedback,omitempty"`
	Contact  string  `protobuf:"bytes,3,opt,name=contact" json:"contact,omitempty"`
}

func (m *UserFeedbackRequest) Reset()                    { *m = UserFeedbackRequest{} }
func (m *UserFeedbackRequest) String() string            { return proto.CompactTextString(m) }
func (*UserFeedbackRequest) ProtoMessage()               {}
func (*UserFeedbackRequest) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{4} }

func (m *UserFeedbackRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *UserFeedbackRequest) GetFeedback() string {
	if m != nil {
		return m.Feedback
	}
	return ""
}

func (m *UserFeedbackRequest) GetContact() string {
	if m != nil {
		return m.Contact
	}
	return ""
}

type UserFeedbackResponse struct {
	Status *Status `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *UserFeedbackResponse) Reset()                    { *m = UserFeedbackResponse{} }
func (m *UserFeedbackResponse) String() string            { return proto.CompactTextString(m) }
func (*UserFeedbackResponse) ProtoMessage()               {}
func (*UserFeedbackResponse) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{5} }

func (m *UserFeedbackResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

type PostDeviceTokenRequest struct {
	Header      *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	DeviceToken string  `protobuf:"bytes,2,opt,name=device_token,json=deviceToken" json:"device_token,omitempty"`
}

func (m *PostDeviceTokenRequest) Reset()                    { *m = PostDeviceTokenRequest{} }
func (m *PostDeviceTokenRequest) String() string            { return proto.CompactTextString(m) }
func (*PostDeviceTokenRequest) ProtoMessage()               {}
func (*PostDeviceTokenRequest) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{6} }

func (m *PostDeviceTokenRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *PostDeviceTokenRequest) GetDeviceToken() string {
	if m != nil {
		return m.DeviceToken
	}
	return ""
}

type PostDeviceTokenResponse struct {
	Status *Status `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *PostDeviceTokenResponse) Reset()                    { *m = PostDeviceTokenResponse{} }
func (m *PostDeviceTokenResponse) String() string            { return proto.CompactTextString(m) }
func (*PostDeviceTokenResponse) ProtoMessage()               {}
func (*PostDeviceTokenResponse) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{7} }

func (m *PostDeviceTokenResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

type GetRemoteConfigRequest struct {
	Header        *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	SubmoduleHash string  `protobuf:"bytes,2,opt,name=submodule_hash,json=submoduleHash" json:"submodule_hash,omitempty"`
}

func (m *GetRemoteConfigRequest) Reset()                    { *m = GetRemoteConfigRequest{} }
func (m *GetRemoteConfigRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRemoteConfigRequest) ProtoMessage()               {}
func (*GetRemoteConfigRequest) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{8} }

func (m *GetRemoteConfigRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetRemoteConfigRequest) GetSubmoduleHash() string {
	if m != nil {
		return m.SubmoduleHash
	}
	return ""
}

type GetRemoteConfigResponse struct {
	Status *Status `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *GetRemoteConfigResponse) Reset()                    { *m = GetRemoteConfigResponse{} }
func (m *GetRemoteConfigResponse) String() string            { return proto.CompactTextString(m) }
func (*GetRemoteConfigResponse) ProtoMessage()               {}
func (*GetRemoteConfigResponse) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{9} }

func (m *GetRemoteConfigResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

type GetVideoInfoRequest struct {
	Header  *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	VideoId string  `protobuf:"bytes,2,opt,name=video_id,json=videoId" json:"video_id,omitempty"`
}

func (m *GetVideoInfoRequest) Reset()                    { *m = GetVideoInfoRequest{} }
func (m *GetVideoInfoRequest) String() string            { return proto.CompactTextString(m) }
func (*GetVideoInfoRequest) ProtoMessage()               {}
func (*GetVideoInfoRequest) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{10} }

func (m *GetVideoInfoRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetVideoInfoRequest) GetVideoId() string {
	if m != nil {
		return m.VideoId
	}
	return ""
}

type GetVideoInfoResponse struct {
	Status   *Status   `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	ListItem *ListItem `protobuf:"bytes,2,opt,name=list_item,json=listItem" json:"list_item,omitempty"`
}

func (m *GetVideoInfoResponse) Reset()                    { *m = GetVideoInfoResponse{} }
func (m *GetVideoInfoResponse) String() string            { return proto.CompactTextString(m) }
func (*GetVideoInfoResponse) ProtoMessage()               {}
func (*GetVideoInfoResponse) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{11} }

func (m *GetVideoInfoResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *GetVideoInfoResponse) GetListItem() *ListItem {
	if m != nil {
		return m.ListItem
	}
	return nil
}

func init() {
	proto.RegisterType((*ReportVideoRequest)(nil), "budao.ReportVideoRequest")
	proto.RegisterType((*ReportVideoResponse)(nil), "budao.ReportVideoResponse")
	proto.RegisterType((*ReportCommentRequest)(nil), "budao.ReportCommentRequest")
	proto.RegisterType((*ReportCommentResponse)(nil), "budao.ReportCommentResponse")
	proto.RegisterType((*UserFeedbackRequest)(nil), "budao.UserFeedbackRequest")
	proto.RegisterType((*UserFeedbackResponse)(nil), "budao.UserFeedbackResponse")
	proto.RegisterType((*PostDeviceTokenRequest)(nil), "budao.PostDeviceTokenRequest")
	proto.RegisterType((*PostDeviceTokenResponse)(nil), "budao.PostDeviceTokenResponse")
	proto.RegisterType((*GetRemoteConfigRequest)(nil), "budao.GetRemoteConfigRequest")
	proto.RegisterType((*GetRemoteConfigResponse)(nil), "budao.GetRemoteConfigResponse")
	proto.RegisterType((*GetVideoInfoRequest)(nil), "budao.GetVideoInfoRequest")
	proto.RegisterType((*GetVideoInfoResponse)(nil), "budao.GetVideoInfoResponse")
}

func init() { proto.RegisterFile("misc.proto", fileDescriptor4) }

var fileDescriptor4 = []byte{
	// 521 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x95, 0xdd, 0x6e, 0xd3, 0x40,
	0x10, 0x85, 0x95, 0x46, 0xa4, 0xc9, 0x38, 0xa1, 0x68, 0x13, 0x42, 0xea, 0x52, 0x54, 0x2c, 0x55,
	0xea, 0x05, 0xf2, 0x45, 0xb8, 0x05, 0x84, 0x68, 0x45, 0x9a, 0x0a, 0xa4, 0xc8, 0xe5, 0x4f, 0xdc,
	0x44, 0xfe, 0x99, 0x90, 0x25, 0xd9, 0xdd, 0xe0, 0x5d, 0x57, 0x88, 0xc7, 0xe1, 0x9d, 0x78, 0x1f,
	0x64, 0x7b, 0xeb, 0xc4, 0x89, 0x8b, 0xb4, 0x12, 0x77, 0xd9, 0x99, 0x9d, 0x73, 0x3e, 0xcb, 0x67,
	0x62, 0x00, 0x46, 0x65, 0xe8, 0xae, 0x62, 0xa1, 0x04, 0xb9, 0x17, 0x24, 0x91, 0x2f, 0xec, 0x76,
	0x28, 0x18, 0x13, 0x3c, 0x2f, 0x3a, 0x1c, 0x88, 0x87, 0x2b, 0x11, 0xab, 0x4f, 0x34, 0x42, 0xe1,
	0xe1, 0x8f, 0x04, 0xa5, 0x22, 0xa7, 0xd0, 0x98, 0xa3, 0x1f, 0x61, 0x3c, 0xa8, 0x9d, 0xd4, 0xce,
	0xac, 0x61, 0xc7, 0xcd, 0x66, 0xdd, 0xcb, 0xac, 0xe8, 0xe9, 0x26, 0x39, 0x84, 0xe6, 0x4d, 0x3a,
	0x36, 0xa5, 0xd1, 0x60, 0xef, 0xa4, 0x76, 0xd6, 0xf2, 0xf6, 0xb3, 0xf3, 0x38, 0x22, 0x7d, 0x68,
	0xc4, 0xe8, 0x4b, 0xc1, 0x07, 0xf5, 0xac, 0xa1, 0x4f, 0xce, 0x0b, 0xe8, 0x96, 0xfc, 0xe4, 0x4a,
	0x70, 0x89, 0xa9, 0xa1, 0x54, 0xbe, 0x4a, 0xe4, 0x96, 0xe1, 0x75, 0x56, 0xf4, 0x74, 0xd3, 0x51,
	0xd0, 0xcb, 0xa7, 0xcf, 0x05, 0x63, 0xc8, 0x95, 0x21, 0xef, 0x31, 0x40, 0x98, 0x0f, 0xae, 0x89,
	0x5b, 0xba, 0xf2, 0x0f, 0xe6, 0x57, 0xf0, 0x70, 0xcb, 0xd5, 0x8c, 0x3a, 0x86, 0xee, 0x47, 0x89,
	0xf1, 0x5b, 0xc4, 0x28, 0xf0, 0xc3, 0x85, 0x21, 0xb4, 0x0d, 0xcd, 0x99, 0x9e, 0xd4, 0xc8, 0xc5,
	0x99, 0x0c, 0x60, 0x3f, 0x14, 0x5c, 0xf9, 0xa1, 0xd2, 0xc8, 0xb7, 0x47, 0xe7, 0x25, 0xf4, 0xca,
	0x9e, 0x66, 0xc8, 0x01, 0xf4, 0x27, 0x42, 0xaa, 0x0b, 0xbc, 0xa1, 0x21, 0x7e, 0x10, 0x0b, 0xe4,
	0x86, 0xd4, 0x4f, 0xa1, 0x1d, 0x65, 0xc3, 0x53, 0x95, 0x4e, 0x6b, 0x72, 0x2b, 0x5a, 0x0b, 0x3a,
	0xaf, 0xe1, 0xd1, 0x8e, 0x87, 0x19, 0xe5, 0x0c, 0xfa, 0x23, 0x54, 0x1e, 0x32, 0xa1, 0xf0, 0x5c,
	0xf0, 0x19, 0xfd, 0x66, 0x48, 0x79, 0x0a, 0xf7, 0x65, 0x12, 0x30, 0x11, 0x25, 0x4b, 0x9c, 0xce,
	0x7d, 0x39, 0xd7, 0x9c, 0x9d, 0xa2, 0x7a, 0xe9, 0xcb, 0x79, 0x4a, 0xba, 0xe3, 0x63, 0x46, 0xfa,
	0x19, 0xba, 0x23, 0xcc, 0x33, 0x3f, 0xe6, 0xb3, 0xff, 0xb7, 0x67, 0xce, 0x02, 0x7a, 0x65, 0x61,
	0x23, 0x2e, 0xf2, 0x0c, 0x5a, 0x4b, 0x2a, 0xd5, 0x94, 0x2a, 0x64, 0x99, 0xb4, 0x35, 0x3c, 0xd0,
	0x37, 0xdf, 0x51, 0xa9, 0xc6, 0x0a, 0x99, 0xd7, 0x5c, 0xea, 0x5f, 0xc3, 0x3f, 0x75, 0xb0, 0xde,
	0x53, 0x19, 0x5e, 0x63, 0x9c, 0xbe, 0x33, 0x72, 0x01, 0xd6, 0xc6, 0x32, 0x93, 0x43, 0x3d, 0xb9,
	0xfb, 0x87, 0x62, 0xdb, 0x55, 0x2d, 0x8d, 0x7a, 0x05, 0x9d, 0xd2, 0x7a, 0x91, 0xa3, 0xd2, 0xe5,
	0xf2, 0xaa, 0xdb, 0x8f, 0xab, 0x9b, 0x5a, 0x6b, 0x04, 0xed, 0xcd, 0xd8, 0x93, 0x5b, 0xdf, 0x8a,
	0xfd, 0xb3, 0x8f, 0x2a, 0x7b, 0x5a, 0x68, 0x02, 0x07, 0x5b, 0xe1, 0x24, 0xc7, 0xfa, 0x7e, 0xf5,
	0x62, 0xd8, 0x4f, 0xee, 0x6a, 0xaf, 0x15, 0xb7, 0x42, 0x54, 0x28, 0x56, 0x87, 0xb8, 0x50, 0xbc,
	0x2b, 0x7b, 0x23, 0x68, 0x6f, 0xbe, 0xfb, 0xe2, 0x61, 0x2b, 0x92, 0x56, 0x3c, 0x6c, 0x55, 0x58,
	0xde, 0xd8, 0xf0, 0x20, 0x14, 0xcc, 0xfd, 0xfe, 0xf3, 0x97, 0xbb, 0x0a, 0xf2, 0x8b, 0x93, 0xda,
	0xef, 0xbd, 0xfa, 0xd5, 0x97, 0xaf, 0x41, 0x23, 0xfb, 0x4e, 0x3c, 0xff, 0x1b, 0x00, 0x00, 0xff,
	0xff, 0x3f, 0xe5, 0xa7, 0x13, 0x4a, 0x06, 0x00, 0x00,
}
