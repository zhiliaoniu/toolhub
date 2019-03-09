// Code generated by protoc-gen-go. DO NOT EDIT.
// source: comment.proto

/*
Package budao is a generated protocol buffer package.

It is generated from these files:
	comment.proto
	common.proto
	config.proto
	like.proto
	misc.proto
	parseurl.proto
	question.proto
	recommend.proto
	report.proto
	share.proto
	timeline.proto
	topic.proto
	transfer.proto
	user.proto

It has these top-level messages:
	CommentVideoRequest
	CommentVideoResponse
	ReplyCommentRequest
	ReplyCommentResponse
	GetVideoCommentListRequest
	GetVideoCommentListResponse
	GetCommentReplyListRequest
	GetCommentReplyListResponse
	DeviceIOSInfo
	DeviceAndroidInfo
	DeviceInfo
	EnvironmentInfo
	Header
	Status
	ShareItem
	TopicItem
	VideoItem
	OptionItem
	QuestionItem
	UserItem
	CommentItem
	ReplyItem
	ArticleItem
	RecommendItem
	SwitchItem
	ListItem
	ChannelItem
	TabItem
	GetChannelListRequest
	GetChannelListResponse
	GetTabListRequest
	GetTabListResponse
	LikeVideoRequest
	LikeVideoResponse
	LikeCommentRequest
	LikeCommentResponse
	ReportVideoRequest
	ReportVideoResponse
	ReportCommentRequest
	ReportCommentResponse
	UserFeedbackRequest
	UserFeedbackResponse
	PostDeviceTokenRequest
	PostDeviceTokenResponse
	GetRemoteConfigRequest
	GetRemoteConfigResponse
	GetVideoInfoRequest
	GetVideoInfoResponse
	ParseURLRequest
	ParseExternalURLRequest
	ParseURLResponse
	AnswerQuestionRequest
	AnswerQuestionResponse
	RecommendVideoItem
	GetRecommendVideoRequest
	GetRecommendVideoResponse
	ReloadIndexRequest
	ReloadIndexResponse
	VideoBase
	VideoClick
	VideoProgress
	VideoExposure
	BannerClick
	ReportStatisDataRequest
	ReportStatisDataResponse
	ReportBatchDataRequest
	ReportBatchDataResponse
	ShareVideoBottomPageRequest
	ShareVideoBottomPageResponse
	ShareTopicPageRequest
	ShareTopicPageResponse
	GetTimeLineRequest
	GetTimeLineResponse
	GetSubscribedTimeLineRequest
	GetSubscribedTimeLineResponse
	BannerItem
	GetTopicListRequest
	GetTopicListResponse
	GetTopicVideoListRequest
	GetTopicVideoListResponse
	SubscribeTopicRequest
	SubscribeTopicResponse
	GetSubscribedTopicListRequest
	GetSubscribedTopicListResponse
	PostVideo
	PostVideoResult
	PostVideosRequest
	PostVideosResponse
	AuditContentRequest
	AuditContentResponse
	LoginRequest
	LoginResponse
	UserInfo
*/
package budao

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CommentVideoRequest struct {
	Header  *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	VideoId string  `protobuf:"bytes,2,opt,name=video_id,json=videoId" json:"video_id,omitempty"`
	Content string  `protobuf:"bytes,3,opt,name=content" json:"content,omitempty"`
}

func (m *CommentVideoRequest) Reset()                    { *m = CommentVideoRequest{} }
func (m *CommentVideoRequest) String() string            { return proto.CompactTextString(m) }
func (*CommentVideoRequest) ProtoMessage()               {}
func (*CommentVideoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *CommentVideoRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *CommentVideoRequest) GetVideoId() string {
	if m != nil {
		return m.VideoId
	}
	return ""
}

func (m *CommentVideoRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type CommentVideoResponse struct {
	Status      *Status      `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	CommentItem *CommentItem `protobuf:"bytes,2,opt,name=comment_item,json=commentItem" json:"comment_item,omitempty"`
}

func (m *CommentVideoResponse) Reset()                    { *m = CommentVideoResponse{} }
func (m *CommentVideoResponse) String() string            { return proto.CompactTextString(m) }
func (*CommentVideoResponse) ProtoMessage()               {}
func (*CommentVideoResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CommentVideoResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *CommentVideoResponse) GetCommentItem() *CommentItem {
	if m != nil {
		return m.CommentItem
	}
	return nil
}

type ReplyCommentRequest struct {
	Header *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	// Types that are valid to be assigned to Type:
	//	*ReplyCommentRequest_CommentId
	//	*ReplyCommentRequest_ReplyId
	Type    isReplyCommentRequest_Type `protobuf_oneof:"Type"`
	Content string                     `protobuf:"bytes,4,opt,name=content" json:"content,omitempty"`
}

func (m *ReplyCommentRequest) Reset()                    { *m = ReplyCommentRequest{} }
func (m *ReplyCommentRequest) String() string            { return proto.CompactTextString(m) }
func (*ReplyCommentRequest) ProtoMessage()               {}
func (*ReplyCommentRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type isReplyCommentRequest_Type interface {
	isReplyCommentRequest_Type()
}

type ReplyCommentRequest_CommentId struct {
	CommentId string `protobuf:"bytes,2,opt,name=comment_id,json=commentId,oneof"`
}
type ReplyCommentRequest_ReplyId struct {
	ReplyId string `protobuf:"bytes,3,opt,name=reply_id,json=replyId,oneof"`
}

func (*ReplyCommentRequest_CommentId) isReplyCommentRequest_Type() {}
func (*ReplyCommentRequest_ReplyId) isReplyCommentRequest_Type()   {}

func (m *ReplyCommentRequest) GetType() isReplyCommentRequest_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *ReplyCommentRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ReplyCommentRequest) GetCommentId() string {
	if x, ok := m.GetType().(*ReplyCommentRequest_CommentId); ok {
		return x.CommentId
	}
	return ""
}

func (m *ReplyCommentRequest) GetReplyId() string {
	if x, ok := m.GetType().(*ReplyCommentRequest_ReplyId); ok {
		return x.ReplyId
	}
	return ""
}

func (m *ReplyCommentRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ReplyCommentRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ReplyCommentRequest_OneofMarshaler, _ReplyCommentRequest_OneofUnmarshaler, _ReplyCommentRequest_OneofSizer, []interface{}{
		(*ReplyCommentRequest_CommentId)(nil),
		(*ReplyCommentRequest_ReplyId)(nil),
	}
}

func _ReplyCommentRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ReplyCommentRequest)
	// Type
	switch x := m.Type.(type) {
	case *ReplyCommentRequest_CommentId:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.CommentId)
	case *ReplyCommentRequest_ReplyId:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.ReplyId)
	case nil:
	default:
		return fmt.Errorf("ReplyCommentRequest.Type has unexpected type %T", x)
	}
	return nil
}

func _ReplyCommentRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ReplyCommentRequest)
	switch tag {
	case 2: // Type.comment_id
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Type = &ReplyCommentRequest_CommentId{x}
		return true, err
	case 3: // Type.reply_id
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Type = &ReplyCommentRequest_ReplyId{x}
		return true, err
	default:
		return false, nil
	}
}

func _ReplyCommentRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ReplyCommentRequest)
	// Type
	switch x := m.Type.(type) {
	case *ReplyCommentRequest_CommentId:
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.CommentId)))
		n += len(x.CommentId)
	case *ReplyCommentRequest_ReplyId:
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.ReplyId)))
		n += len(x.ReplyId)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ReplyCommentResponse struct {
	Status    *Status    `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	ReplyItem *ReplyItem `protobuf:"bytes,2,opt,name=reply_item,json=replyItem" json:"reply_item,omitempty"`
}

func (m *ReplyCommentResponse) Reset()                    { *m = ReplyCommentResponse{} }
func (m *ReplyCommentResponse) String() string            { return proto.CompactTextString(m) }
func (*ReplyCommentResponse) ProtoMessage()               {}
func (*ReplyCommentResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ReplyCommentResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *ReplyCommentResponse) GetReplyItem() *ReplyItem {
	if m != nil {
		return m.ReplyItem
	}
	return nil
}

type GetVideoCommentListRequest struct {
	Header          *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	VideoId         string  `protobuf:"bytes,2,opt,name=video_id,json=videoId" json:"video_id,omitempty"`
	LastCommentId   string  `protobuf:"bytes,3,opt,name=last_comment_id,json=lastCommentId" json:"last_comment_id,omitempty"`
	LastCommentTime uint64  `protobuf:"varint,4,opt,name=last_comment_time,json=lastCommentTime" json:"last_comment_time,omitempty"`
}

func (m *GetVideoCommentListRequest) Reset()                    { *m = GetVideoCommentListRequest{} }
func (m *GetVideoCommentListRequest) String() string            { return proto.CompactTextString(m) }
func (*GetVideoCommentListRequest) ProtoMessage()               {}
func (*GetVideoCommentListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GetVideoCommentListRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetVideoCommentListRequest) GetVideoId() string {
	if m != nil {
		return m.VideoId
	}
	return ""
}

func (m *GetVideoCommentListRequest) GetLastCommentId() string {
	if m != nil {
		return m.LastCommentId
	}
	return ""
}

func (m *GetVideoCommentListRequest) GetLastCommentTime() uint64 {
	if m != nil {
		return m.LastCommentTime
	}
	return 0
}

type GetVideoCommentListResponse struct {
	Status          *Status        `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	HotCommentItems []*CommentItem `protobuf:"bytes,2,rep,name=hot_comment_items,json=hotCommentItems" json:"hot_comment_items,omitempty"`
	CommentItems    []*CommentItem `protobuf:"bytes,3,rep,name=comment_items,json=commentItems" json:"comment_items,omitempty"`
	HasMore         bool           `protobuf:"varint,4,opt,name=has_more,json=hasMore" json:"has_more,omitempty"`
}

func (m *GetVideoCommentListResponse) Reset()                    { *m = GetVideoCommentListResponse{} }
func (m *GetVideoCommentListResponse) String() string            { return proto.CompactTextString(m) }
func (*GetVideoCommentListResponse) ProtoMessage()               {}
func (*GetVideoCommentListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GetVideoCommentListResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *GetVideoCommentListResponse) GetHotCommentItems() []*CommentItem {
	if m != nil {
		return m.HotCommentItems
	}
	return nil
}

func (m *GetVideoCommentListResponse) GetCommentItems() []*CommentItem {
	if m != nil {
		return m.CommentItems
	}
	return nil
}

func (m *GetVideoCommentListResponse) GetHasMore() bool {
	if m != nil {
		return m.HasMore
	}
	return false
}

type GetCommentReplyListRequest struct {
	Header        *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	CommentId     string  `protobuf:"bytes,2,opt,name=comment_id,json=commentId" json:"comment_id,omitempty"`
	LastReplyId   string  `protobuf:"bytes,3,opt,name=last_reply_id,json=lastReplyId" json:"last_reply_id,omitempty"`
	LastReplyTime uint64  `protobuf:"varint,4,opt,name=last_reply_time,json=lastReplyTime" json:"last_reply_time,omitempty"`
}

func (m *GetCommentReplyListRequest) Reset()                    { *m = GetCommentReplyListRequest{} }
func (m *GetCommentReplyListRequest) String() string            { return proto.CompactTextString(m) }
func (*GetCommentReplyListRequest) ProtoMessage()               {}
func (*GetCommentReplyListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *GetCommentReplyListRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetCommentReplyListRequest) GetCommentId() string {
	if m != nil {
		return m.CommentId
	}
	return ""
}

func (m *GetCommentReplyListRequest) GetLastReplyId() string {
	if m != nil {
		return m.LastReplyId
	}
	return ""
}

func (m *GetCommentReplyListRequest) GetLastReplyTime() uint64 {
	if m != nil {
		return m.LastReplyTime
	}
	return 0
}

type GetCommentReplyListResponse struct {
	Status      *Status      `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	CommentItem *CommentItem `protobuf:"bytes,2,opt,name=comment_item,json=commentItem" json:"comment_item,omitempty"`
	ReplyItems  []*ReplyItem `protobuf:"bytes,3,rep,name=reply_items,json=replyItems" json:"reply_items,omitempty"`
	HasMore     bool         `protobuf:"varint,4,opt,name=has_more,json=hasMore" json:"has_more,omitempty"`
}

func (m *GetCommentReplyListResponse) Reset()                    { *m = GetCommentReplyListResponse{} }
func (m *GetCommentReplyListResponse) String() string            { return proto.CompactTextString(m) }
func (*GetCommentReplyListResponse) ProtoMessage()               {}
func (*GetCommentReplyListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *GetCommentReplyListResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *GetCommentReplyListResponse) GetCommentItem() *CommentItem {
	if m != nil {
		return m.CommentItem
	}
	return nil
}

func (m *GetCommentReplyListResponse) GetReplyItems() []*ReplyItem {
	if m != nil {
		return m.ReplyItems
	}
	return nil
}

func (m *GetCommentReplyListResponse) GetHasMore() bool {
	if m != nil {
		return m.HasMore
	}
	return false
}

func init() {
	proto.RegisterType((*CommentVideoRequest)(nil), "budao.CommentVideoRequest")
	proto.RegisterType((*CommentVideoResponse)(nil), "budao.CommentVideoResponse")
	proto.RegisterType((*ReplyCommentRequest)(nil), "budao.ReplyCommentRequest")
	proto.RegisterType((*ReplyCommentResponse)(nil), "budao.ReplyCommentResponse")
	proto.RegisterType((*GetVideoCommentListRequest)(nil), "budao.GetVideoCommentListRequest")
	proto.RegisterType((*GetVideoCommentListResponse)(nil), "budao.GetVideoCommentListResponse")
	proto.RegisterType((*GetCommentReplyListRequest)(nil), "budao.GetCommentReplyListRequest")
	proto.RegisterType((*GetCommentReplyListResponse)(nil), "budao.GetCommentReplyListResponse")
}

func init() { proto.RegisterFile("comment.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 571 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0x26, 0x6d, 0xe9, 0xcf, 0x69, 0x4b, 0x37, 0x77, 0x17, 0x5d, 0x2a, 0x44, 0x89, 0x04, 0x9a,
	0xb8, 0x08, 0xa2, 0x08, 0x71, 0xc7, 0xc5, 0x76, 0xd1, 0x0d, 0x81, 0x84, 0xbc, 0x09, 0x21, 0x84,
	0x14, 0xa5, 0x89, 0xa5, 0x04, 0x35, 0x71, 0x16, 0xbb, 0x13, 0xe5, 0x59, 0xb8, 0xe2, 0x01, 0xb8,
	0xe0, 0x2d, 0x78, 0x03, 0xae, 0x78, 0x16, 0x14, 0xdb, 0x4d, 0x9d, 0x35, 0x99, 0x54, 0x84, 0x76,
	0x17, 0x9f, 0xf3, 0x9d, 0x73, 0xec, 0xef, 0xfb, 0xec, 0x40, 0xdf, 0xa3, 0x51, 0x44, 0x62, 0x6e,
	0x27, 0x29, 0xe5, 0x14, 0xdd, 0x9d, 0x2f, 0x7d, 0x97, 0x9a, 0xbd, 0x2c, 0x4a, 0x63, 0x19, 0xb4,
	0x2e, 0x61, 0x78, 0x22, 0x51, 0xef, 0x43, 0x9f, 0x50, 0x4c, 0x2e, 0x97, 0x84, 0x71, 0xf4, 0x08,
	0x9a, 0x01, 0x71, 0x7d, 0x92, 0x8e, 0x8c, 0x89, 0x71, 0xd4, 0x9d, 0xf6, 0x6d, 0x51, 0x6c, 0x9f,
	0x8a, 0x20, 0x56, 0x49, 0x74, 0x08, 0xed, 0xab, 0xac, 0xcc, 0x09, 0xfd, 0x51, 0x6d, 0x62, 0x1c,
	0x75, 0x70, 0x4b, 0xac, 0xcf, 0x7c, 0x34, 0x82, 0x96, 0x47, 0x63, 0x4e, 0x62, 0x3e, 0xaa, 0xcb,
	0x8c, 0x5a, 0x5a, 0x1c, 0x0e, 0x8a, 0x23, 0x59, 0x42, 0x63, 0x46, 0xb2, 0x99, 0x8c, 0xbb, 0x7c,
	0xc9, 0xae, 0xcd, 0x3c, 0x17, 0x41, 0xac, 0x92, 0xe8, 0x05, 0xf4, 0xd4, 0xb9, 0x9c, 0x90, 0x93,
	0x48, 0xcc, 0xed, 0x4e, 0x91, 0x02, 0xab, 0xce, 0x67, 0x9c, 0x44, 0xb8, 0xeb, 0x6d, 0x16, 0xd6,
	0x37, 0x03, 0x86, 0x98, 0x24, 0x8b, 0x95, 0x42, 0xec, 0x78, 0xd2, 0x07, 0x00, 0xf9, 0x54, 0x75,
	0xd6, 0xd3, 0x3b, 0xb8, 0xb3, 0x9e, 0xe0, 0xa3, 0x31, 0xb4, 0xd3, 0xac, 0x7d, 0x96, 0xae, 0xab,
	0x74, 0x4b, 0x44, 0x8a, 0x64, 0x34, 0x0a, 0x64, 0x1c, 0x37, 0xa1, 0x71, 0xb1, 0x4a, 0x88, 0x15,
	0xc3, 0x41, 0x71, 0x77, 0xbb, 0x91, 0xf2, 0x14, 0x40, 0x4d, 0xdf, 0x50, 0xb2, 0xa7, 0xa0, 0xa2,
	0xaf, 0x20, 0xa4, 0x93, 0xae, 0x3f, 0xad, 0x9f, 0x06, 0x98, 0x33, 0x22, 0x15, 0x50, 0x33, 0xdf,
	0x84, 0x8c, 0xff, 0x3f, 0xfd, 0x1f, 0xc3, 0x60, 0xe1, 0x32, 0xee, 0x68, 0xac, 0x49, 0x1f, 0xf4,
	0xb3, 0xf0, 0x49, 0xce, 0xdb, 0x13, 0xd8, 0x2f, 0xe0, 0x78, 0x18, 0x11, 0x41, 0x52, 0x03, 0x0f,
	0x34, 0xe4, 0x45, 0x18, 0x11, 0xeb, 0xb7, 0x01, 0xe3, 0xd2, 0x4d, 0xef, 0x46, 0xd6, 0x2b, 0xd8,
	0x0f, 0xa8, 0xb6, 0x33, 0x4e, 0x22, 0x36, 0xaa, 0x4d, 0xea, 0x15, 0x36, 0x1a, 0x04, 0x94, 0x6b,
	0x6b, 0x86, 0x5e, 0xe6, 0x37, 0x4b, 0xd5, 0xd6, 0x2b, 0x6b, 0x7b, 0x9e, 0x5e, 0x78, 0x08, 0xed,
	0xc0, 0x65, 0x4e, 0x44, 0x53, 0x79, 0xc4, 0x36, 0x6e, 0x05, 0x2e, 0x7b, 0x4b, 0x53, 0x62, 0xfd,
	0x90, 0x7a, 0xe4, 0xf2, 0x27, 0x8b, 0xd5, 0x3f, 0xe8, 0x71, 0x7f, 0xdb, 0xa5, 0xba, 0x47, 0x2d,
	0x10, 0xe4, 0x3b, 0x45, 0xa3, 0xe2, 0x6e, 0x16, 0xc4, 0xca, 0xaa, 0x6b, 0xdd, 0x24, 0x46, 0x53,
	0xa3, 0x9f, 0xa3, 0x84, 0x16, 0xbf, 0xa4, 0x16, 0xdb, 0x1b, 0xbe, 0x8d, 0xdb, 0x8c, 0x9e, 0x41,
	0x77, 0xe3, 0xf7, 0xb5, 0x00, 0xdb, 0x86, 0x87, 0xdc, 0xf0, 0x37, 0x91, 0x3f, 0xfd, 0x53, 0x83,
	0x7b, 0x6a, 0xd4, 0x39, 0x49, 0xaf, 0x42, 0x8f, 0xa0, 0x19, 0xf4, 0xf4, 0x47, 0x0a, 0x99, 0xc5,
	0x1d, 0xe9, 0x8f, 0xa5, 0x39, 0x2e, 0xcd, 0x29, 0x1e, 0x66, 0xd0, 0xd3, 0x2f, 0x76, 0xde, 0xa8,
	0xe4, 0x2d, 0xca, 0x1b, 0x95, 0xbe, 0x04, 0x9f, 0x60, 0x58, 0xe2, 0x7d, 0xf4, 0x50, 0xd5, 0x54,
	0x5f, 0x66, 0xd3, 0xba, 0x09, 0x52, 0xe8, 0x7e, 0x5d, 0x4d, 0xbd, 0x7b, 0x85, 0x35, 0xf5, 0xee,
	0x55, 0x66, 0x38, 0x36, 0x61, 0xcf, 0xa3, 0x91, 0xfd, 0xf9, 0xcb, 0x57, 0x3b, 0x99, 0x4b, 0xfc,
	0x3b, 0xe3, 0x7b, 0xad, 0xfe, 0xfa, 0xc3, 0xc7, 0x79, 0x53, 0xfc, 0x88, 0x9e, 0xff, 0x0d, 0x00,
	0x00, 0xff, 0xff, 0x1b, 0xd5, 0xa7, 0xc9, 0xae, 0x06, 0x00, 0x00,
}
