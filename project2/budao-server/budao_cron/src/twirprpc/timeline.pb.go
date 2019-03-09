// Code generated by protoc-gen-go. DO NOT EDIT.
// source: timeline.proto

package budao

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type GetTimeLineRequest struct {
	Header    *Header   `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Direction Direction `protobuf:"varint,2,opt,name=direction,enum=budao.Direction" json:"direction,omitempty"`
	ChannelId string    `protobuf:"bytes,3,opt,name=channel_id,json=channelId" json:"channel_id,omitempty"`
}

func (m *GetTimeLineRequest) Reset()                    { *m = GetTimeLineRequest{} }
func (m *GetTimeLineRequest) String() string            { return proto.CompactTextString(m) }
func (*GetTimeLineRequest) ProtoMessage()               {}
func (*GetTimeLineRequest) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{0} }

func (m *GetTimeLineRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetTimeLineRequest) GetDirection() Direction {
	if m != nil {
		return m.Direction
	}
	return Direction_ALL
}

func (m *GetTimeLineRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

type GetTimeLineResponse struct {
	Status     *Status     `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	ListItems  []*ListItem `protobuf:"bytes,2,rep,name=list_items,json=listItems" json:"list_items,omitempty"`
	HasMore    bool        `protobuf:"varint,3,opt,name=has_more,json=hasMore" json:"has_more,omitempty"`
	ClearCache bool        `protobuf:"varint,4,opt,name=clear_cache,json=clearCache" json:"clear_cache,omitempty"`
	Tips       string      `protobuf:"bytes,5,opt,name=tips" json:"tips,omitempty"`
}

func (m *GetTimeLineResponse) Reset()                    { *m = GetTimeLineResponse{} }
func (m *GetTimeLineResponse) String() string            { return proto.CompactTextString(m) }
func (*GetTimeLineResponse) ProtoMessage()               {}
func (*GetTimeLineResponse) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{1} }

func (m *GetTimeLineResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *GetTimeLineResponse) GetListItems() []*ListItem {
	if m != nil {
		return m.ListItems
	}
	return nil
}

func (m *GetTimeLineResponse) GetHasMore() bool {
	if m != nil {
		return m.HasMore
	}
	return false
}

func (m *GetTimeLineResponse) GetClearCache() bool {
	if m != nil {
		return m.ClearCache
	}
	return false
}

func (m *GetTimeLineResponse) GetTips() string {
	if m != nil {
		return m.Tips
	}
	return ""
}

type GetSubscribedTimeLineRequest struct {
	Header       *Header   `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Direction    Direction `protobuf:"varint,2,opt,name=direction,enum=budao.Direction" json:"direction,omitempty"`
	RefVideoId   string    `protobuf:"bytes,3,opt,name=ref_video_id,json=refVideoId" json:"ref_video_id,omitempty"`
	RefVideoTime uint64    `protobuf:"varint,4,opt,name=ref_video_time,json=refVideoTime" json:"ref_video_time,omitempty"`
}

func (m *GetSubscribedTimeLineRequest) Reset()                    { *m = GetSubscribedTimeLineRequest{} }
func (m *GetSubscribedTimeLineRequest) String() string            { return proto.CompactTextString(m) }
func (*GetSubscribedTimeLineRequest) ProtoMessage()               {}
func (*GetSubscribedTimeLineRequest) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{2} }

func (m *GetSubscribedTimeLineRequest) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetSubscribedTimeLineRequest) GetDirection() Direction {
	if m != nil {
		return m.Direction
	}
	return Direction_ALL
}

func (m *GetSubscribedTimeLineRequest) GetRefVideoId() string {
	if m != nil {
		return m.RefVideoId
	}
	return ""
}

func (m *GetSubscribedTimeLineRequest) GetRefVideoTime() uint64 {
	if m != nil {
		return m.RefVideoTime
	}
	return 0
}

type GetSubscribedTimeLineResponse struct {
	Status     *Status     `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	ListItems  []*ListItem `protobuf:"bytes,2,rep,name=list_items,json=listItems" json:"list_items,omitempty"`
	TopicNum   int32       `protobuf:"varint,3,opt,name=topic_num,json=topicNum" json:"topic_num,omitempty"`
	HasMore    bool        `protobuf:"varint,4,opt,name=has_more,json=hasMore" json:"has_more,omitempty"`
	ClearCache bool        `protobuf:"varint,5,opt,name=clear_cache,json=clearCache" json:"clear_cache,omitempty"`
	Tips       string      `protobuf:"bytes,6,opt,name=tips" json:"tips,omitempty"`
}

func (m *GetSubscribedTimeLineResponse) Reset()                    { *m = GetSubscribedTimeLineResponse{} }
func (m *GetSubscribedTimeLineResponse) String() string            { return proto.CompactTextString(m) }
func (*GetSubscribedTimeLineResponse) ProtoMessage()               {}
func (*GetSubscribedTimeLineResponse) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{3} }

func (m *GetSubscribedTimeLineResponse) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *GetSubscribedTimeLineResponse) GetListItems() []*ListItem {
	if m != nil {
		return m.ListItems
	}
	return nil
}

func (m *GetSubscribedTimeLineResponse) GetTopicNum() int32 {
	if m != nil {
		return m.TopicNum
	}
	return 0
}

func (m *GetSubscribedTimeLineResponse) GetHasMore() bool {
	if m != nil {
		return m.HasMore
	}
	return false
}

func (m *GetSubscribedTimeLineResponse) GetClearCache() bool {
	if m != nil {
		return m.ClearCache
	}
	return false
}

func (m *GetSubscribedTimeLineResponse) GetTips() string {
	if m != nil {
		return m.Tips
	}
	return ""
}

func init() {
	proto.RegisterType((*GetTimeLineRequest)(nil), "budao.GetTimeLineRequest")
	proto.RegisterType((*GetTimeLineResponse)(nil), "budao.GetTimeLineResponse")
	proto.RegisterType((*GetSubscribedTimeLineRequest)(nil), "budao.GetSubscribedTimeLineRequest")
	proto.RegisterType((*GetSubscribedTimeLineResponse)(nil), "budao.GetSubscribedTimeLineResponse")
}

func init() { proto.RegisterFile("timeline.proto", fileDescriptor10) }

var fileDescriptor10 = []byte{
	// 449 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x93, 0xdf, 0x6a, 0xd4, 0x40,
	0x14, 0xc6, 0x99, 0xfd, 0xe7, 0xe6, 0xec, 0xba, 0x2d, 0x23, 0x42, 0x1a, 0x2d, 0x86, 0xb5, 0xc2,
	0x5e, 0xe5, 0x22, 0xbe, 0x81, 0x16, 0xea, 0x4a, 0x15, 0x99, 0x15, 0x11, 0x6f, 0x42, 0x32, 0x39,
	0x25, 0x23, 0x99, 0x4c, 0x9c, 0x99, 0x14, 0xf1, 0x11, 0x7c, 0x0c, 0x5f, 0x43, 0x7d, 0x1a, 0x5f,
	0x44, 0x32, 0x9b, 0x25, 0x56, 0xbb, 0xbd, 0xb3, 0x77, 0xe1, 0x77, 0xbe, 0x43, 0xbe, 0x73, 0xbe,
	0x33, 0xb0, 0xb0, 0x42, 0x62, 0x29, 0x2a, 0x8c, 0x6a, 0xad, 0xac, 0xa2, 0xe3, 0xac, 0xc9, 0x53,
	0x15, 0xcc, 0xb9, 0x92, 0x52, 0x55, 0x5b, 0xb8, 0xfc, 0x4a, 0x80, 0x9e, 0xa1, 0x7d, 0x2b, 0x24,
	0x9e, 0x8b, 0x0a, 0x19, 0x7e, 0x6a, 0xd0, 0x58, 0xfa, 0x04, 0x26, 0x05, 0xa6, 0x39, 0x6a, 0x9f,
	0x84, 0x64, 0x35, 0x8b, 0xef, 0x46, 0xae, 0x39, 0x7a, 0xe1, 0x20, 0xeb, 0x8a, 0x34, 0x02, 0x2f,
	0x17, 0x1a, 0xb9, 0x15, 0xaa, 0xf2, 0x07, 0x21, 0x59, 0x2d, 0xe2, 0xc3, 0x4e, 0x79, 0xba, 0xe3,
	0xac, 0x97, 0xd0, 0x63, 0x00, 0x5e, 0xa4, 0x55, 0x85, 0x65, 0x22, 0x72, 0x7f, 0x18, 0x92, 0x95,
	0xc7, 0xbc, 0x8e, 0xac, 0xf3, 0xe5, 0x0f, 0x02, 0xf7, 0xae, 0x98, 0x31, 0xb5, 0xaa, 0x0c, 0xb6,
	0x6e, 0x8c, 0x4d, 0x6d, 0x63, 0xfe, 0x72, 0xb3, 0x71, 0x90, 0x75, 0x45, 0x1a, 0x01, 0x94, 0xc2,
	0xd8, 0x44, 0x58, 0x94, 0xc6, 0x1f, 0x84, 0xc3, 0xd5, 0x2c, 0x3e, 0xe8, 0xa4, 0xe7, 0xc2, 0xd8,
	0xb5, 0x45, 0xc9, 0xbc, 0xb2, 0xfb, 0x32, 0xf4, 0x08, 0xa6, 0x45, 0x6a, 0x12, 0xa9, 0x34, 0x3a,
	0x2f, 0x53, 0x76, 0xa7, 0x48, 0xcd, 0x2b, 0xa5, 0x91, 0x3e, 0x82, 0x19, 0x2f, 0x31, 0xd5, 0x09,
	0x4f, 0x79, 0x81, 0xfe, 0xc8, 0x55, 0xc1, 0xa1, 0xe7, 0x2d, 0xa1, 0x14, 0x46, 0x56, 0xd4, 0xc6,
	0x1f, 0xbb, 0x19, 0xdc, 0xf7, 0xf2, 0x27, 0x81, 0x87, 0x67, 0x68, 0x37, 0x4d, 0x66, 0xb8, 0x16,
	0x19, 0xe6, 0xb7, 0xb4, 0xd5, 0x10, 0xe6, 0x1a, 0x2f, 0x92, 0x4b, 0x91, 0xa3, 0xea, 0xf7, 0x0a,
	0x1a, 0x2f, 0xde, 0xb5, 0x68, 0x9d, 0xd3, 0x13, 0x58, 0xf4, 0x8a, 0xf6, 0x2c, 0xdc, 0x44, 0x23,
	0x36, 0xdf, 0x69, 0x5a, 0xa7, 0xcb, 0x5f, 0x04, 0x8e, 0xf7, 0xf8, 0xff, 0xbf, 0x41, 0x3c, 0x00,
	0xcf, 0xaa, 0x5a, 0xf0, 0xa4, 0x6a, 0xa4, 0x73, 0x3f, 0x66, 0x53, 0x07, 0x5e, 0x37, 0xf2, 0x4a,
	0x4a, 0xa3, 0x1b, 0x53, 0x1a, 0xef, 0x4d, 0x69, 0xd2, 0xa7, 0x14, 0x7f, 0x27, 0x70, 0xb0, 0x1b,
	0x6c, 0x83, 0xfa, 0x52, 0x70, 0xa4, 0xa7, 0x30, 0xfb, 0xe3, 0xee, 0xe8, 0x51, 0xe7, 0xf5, 0xdf,
	0x87, 0x11, 0x04, 0xd7, 0x95, 0xba, 0xed, 0x64, 0x70, 0xff, 0xda, 0xf5, 0xd1, 0xc7, 0x7d, 0xd3,
	0xde, 0xe3, 0x08, 0x4e, 0x6e, 0x16, 0x6d, 0xff, 0xf1, 0x2c, 0x80, 0x43, 0xae, 0x64, 0xf4, 0xf1,
	0xf3, 0x97, 0xa8, 0xce, 0xb6, 0x1d, 0x6f, 0xc8, 0xb7, 0xc1, 0xf0, 0xe5, 0xfb, 0x0f, 0xd9, 0xc4,
	0x3d, 0xe9, 0xa7, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc8, 0x7d, 0xa5, 0xfb, 0xf9, 0x03, 0x00,
	0x00,
}
