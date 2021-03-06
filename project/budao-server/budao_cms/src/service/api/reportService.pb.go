// Code generated by protoc-gen-go. DO NOT EDIT.
// source: reportService.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// app统计每日视频相关数-导出execl
type ReportVideoExeclResponse struct {
	Code string                `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Msg  string                `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Data *ReportVideoExeclBack `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ReportVideoExeclResponse) Reset()                    { *m = ReportVideoExeclResponse{} }
func (m *ReportVideoExeclResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportVideoExeclResponse) ProtoMessage()               {}
func (*ReportVideoExeclResponse) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{0} }

func (m *ReportVideoExeclResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ReportVideoExeclResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ReportVideoExeclResponse) GetData() *ReportVideoExeclBack {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReportVideoExeclBack struct {
	FileUrl string `protobuf:"bytes,1,opt,name=fileUrl" json:"fileUrl,omitempty"`
}

func (m *ReportVideoExeclBack) Reset()                    { *m = ReportVideoExeclBack{} }
func (m *ReportVideoExeclBack) String() string            { return proto.CompactTextString(m) }
func (*ReportVideoExeclBack) ProtoMessage()               {}
func (*ReportVideoExeclBack) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{1} }

func (m *ReportVideoExeclBack) GetFileUrl() string {
	if m != nil {
		return m.FileUrl
	}
	return ""
}

// app统计每日用户相关数-导出execl
type ReportUserExeclResponse struct {
	Code string               `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Msg  string               `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Data *ReportUserExeclBack `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ReportUserExeclResponse) Reset()                    { *m = ReportUserExeclResponse{} }
func (m *ReportUserExeclResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportUserExeclResponse) ProtoMessage()               {}
func (*ReportUserExeclResponse) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{2} }

func (m *ReportUserExeclResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ReportUserExeclResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ReportUserExeclResponse) GetData() *ReportUserExeclBack {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReportUserExeclBack struct {
	FileUrl string `protobuf:"bytes,1,opt,name=fileUrl" json:"fileUrl,omitempty"`
}

func (m *ReportUserExeclBack) Reset()                    { *m = ReportUserExeclBack{} }
func (m *ReportUserExeclBack) String() string            { return proto.CompactTextString(m) }
func (*ReportUserExeclBack) ProtoMessage()               {}
func (*ReportUserExeclBack) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{3} }

func (m *ReportUserExeclBack) GetFileUrl() string {
	if m != nil {
		return m.FileUrl
	}
	return ""
}

// 每日视频数据统计
type ReportVideoStatisticsResponse struct {
	Code string                     `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Msg  string                     `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Data *ReportVideoStatisticsBack `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ReportVideoStatisticsResponse) Reset()                    { *m = ReportVideoStatisticsResponse{} }
func (m *ReportVideoStatisticsResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportVideoStatisticsResponse) ProtoMessage()               {}
func (*ReportVideoStatisticsResponse) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{4} }

func (m *ReportVideoStatisticsResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ReportVideoStatisticsResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ReportVideoStatisticsResponse) GetData() *ReportVideoStatisticsBack {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReportVideoStatisticsBack struct {
	Data  []*StatisBizDaily `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
	Count uint64            `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
}

func (m *ReportVideoStatisticsBack) Reset()                    { *m = ReportVideoStatisticsBack{} }
func (m *ReportVideoStatisticsBack) String() string            { return proto.CompactTextString(m) }
func (*ReportVideoStatisticsBack) ProtoMessage()               {}
func (*ReportVideoStatisticsBack) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{5} }

func (m *ReportVideoStatisticsBack) GetData() []*StatisBizDaily {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *ReportVideoStatisticsBack) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type StatisBizDaily struct {
	Id              string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	VideoExposeNum  string `protobuf:"bytes,2,opt,name=videoExposeNum" json:"videoExposeNum,omitempty"`
	VideoClictNum   string `protobuf:"bytes,3,opt,name=videoClictNum" json:"videoClictNum,omitempty"`
	VideoViewNum    string `protobuf:"bytes,4,opt,name=videoViewNum" json:"videoViewNum,omitempty"`
	VideoFavorNum   string `protobuf:"bytes,5,opt,name=videoFavorNum" json:"videoFavorNum,omitempty"`
	CommentFavorNum string `protobuf:"bytes,6,opt,name=commentFavorNum" json:"commentFavorNum,omitempty"`
	CommentNum      string `protobuf:"bytes,7,opt,name=commentNum" json:"commentNum,omitempty"`
	TopicFollowNum  string `protobuf:"bytes,8,opt,name=topicFollowNum" json:"topicFollowNum,omitempty"`
	CreateTime      string `protobuf:"bytes,9,opt,name=createTime" json:"createTime,omitempty"`
}

func (m *StatisBizDaily) Reset()                    { *m = StatisBizDaily{} }
func (m *StatisBizDaily) String() string            { return proto.CompactTextString(m) }
func (*StatisBizDaily) ProtoMessage()               {}
func (*StatisBizDaily) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{6} }

func (m *StatisBizDaily) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *StatisBizDaily) GetVideoExposeNum() string {
	if m != nil {
		return m.VideoExposeNum
	}
	return ""
}

func (m *StatisBizDaily) GetVideoClictNum() string {
	if m != nil {
		return m.VideoClictNum
	}
	return ""
}

func (m *StatisBizDaily) GetVideoViewNum() string {
	if m != nil {
		return m.VideoViewNum
	}
	return ""
}

func (m *StatisBizDaily) GetVideoFavorNum() string {
	if m != nil {
		return m.VideoFavorNum
	}
	return ""
}

func (m *StatisBizDaily) GetCommentFavorNum() string {
	if m != nil {
		return m.CommentFavorNum
	}
	return ""
}

func (m *StatisBizDaily) GetCommentNum() string {
	if m != nil {
		return m.CommentNum
	}
	return ""
}

func (m *StatisBizDaily) GetTopicFollowNum() string {
	if m != nil {
		return m.TopicFollowNum
	}
	return ""
}

func (m *StatisBizDaily) GetCreateTime() string {
	if m != nil {
		return m.CreateTime
	}
	return ""
}

// 统计每日用户数
type ReportUserStatisticsResponse struct {
	Code string                    `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Msg  string                    `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Data *ReportUserStatisticsBack `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ReportUserStatisticsResponse) Reset()                    { *m = ReportUserStatisticsResponse{} }
func (m *ReportUserStatisticsResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportUserStatisticsResponse) ProtoMessage()               {}
func (*ReportUserStatisticsResponse) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{7} }

func (m *ReportUserStatisticsResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ReportUserStatisticsResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ReportUserStatisticsResponse) GetData() *ReportUserStatisticsBack {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReportUserStatisticsBack struct {
	Data  []*StatisUser `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
	Count uint64        `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
}

func (m *ReportUserStatisticsBack) Reset()                    { *m = ReportUserStatisticsBack{} }
func (m *ReportUserStatisticsBack) String() string            { return proto.CompactTextString(m) }
func (*ReportUserStatisticsBack) ProtoMessage()               {}
func (*ReportUserStatisticsBack) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{8} }

func (m *ReportUserStatisticsBack) GetData() []*StatisUser {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *ReportUserStatisticsBack) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type StatisUser struct {
	Id         string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	TotalNum   string `protobuf:"bytes,2,opt,name=totalNum" json:"totalNum,omitempty"`
	ActiveNum  string `protobuf:"bytes,3,opt,name=activeNum" json:"activeNum,omitempty"`
	NewNum     string `protobuf:"bytes,4,opt,name=newNum" json:"newNum,omitempty"`
	CreateTime string `protobuf:"bytes,5,opt,name=createTime" json:"createTime,omitempty"`
}

func (m *StatisUser) Reset()                    { *m = StatisUser{} }
func (m *StatisUser) String() string            { return proto.CompactTextString(m) }
func (*StatisUser) ProtoMessage()               {}
func (*StatisUser) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{9} }

func (m *StatisUser) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *StatisUser) GetTotalNum() string {
	if m != nil {
		return m.TotalNum
	}
	return ""
}

func (m *StatisUser) GetActiveNum() string {
	if m != nil {
		return m.ActiveNum
	}
	return ""
}

func (m *StatisUser) GetNewNum() string {
	if m != nil {
		return m.NewNum
	}
	return ""
}

func (m *StatisUser) GetCreateTime() string {
	if m != nil {
		return m.CreateTime
	}
	return ""
}

// 内容数据  话题数据
type ReportContentTopicDataResponse struct {
	Code string                      `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Msg  string                      `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Data *ReportContentTopicDataBack `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ReportContentTopicDataResponse) Reset()                    { *m = ReportContentTopicDataResponse{} }
func (m *ReportContentTopicDataResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportContentTopicDataResponse) ProtoMessage()               {}
func (*ReportContentTopicDataResponse) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{10} }

func (m *ReportContentTopicDataResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ReportContentTopicDataResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ReportContentTopicDataResponse) GetData() *ReportContentTopicDataBack {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReportContentTopicDataBack struct {
	Data    []*TopicOther `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
	Count   uint64        `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
	TOCount uint64        `protobuf:"varint,3,opt,name=tOCount" json:"tOCount,omitempty"`
	TCCount uint64        `protobuf:"varint,4,opt,name=tCCount" json:"tCCount,omitempty"`
}

func (m *ReportContentTopicDataBack) Reset()                    { *m = ReportContentTopicDataBack{} }
func (m *ReportContentTopicDataBack) String() string            { return proto.CompactTextString(m) }
func (*ReportContentTopicDataBack) ProtoMessage()               {}
func (*ReportContentTopicDataBack) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{11} }

func (m *ReportContentTopicDataBack) GetData() []*TopicOther {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *ReportContentTopicDataBack) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *ReportContentTopicDataBack) GetTOCount() uint64 {
	if m != nil {
		return m.TOCount
	}
	return 0
}

func (m *ReportContentTopicDataBack) GetTCCount() uint64 {
	if m != nil {
		return m.TCCount
	}
	return 0
}

type TopicOther struct {
	TopicName  string `protobuf:"bytes,1,opt,name=topicName" json:"topicName,omitempty"`
	VideoCount string `protobuf:"bytes,2,opt,name=videoCount" json:"videoCount,omitempty"`
	QueCount   string `protobuf:"bytes,3,opt,name=queCount" json:"queCount,omitempty"`
	UserCount  string `protobuf:"bytes,4,opt,name=userCount" json:"userCount,omitempty"`
}

func (m *TopicOther) Reset()                    { *m = TopicOther{} }
func (m *TopicOther) String() string            { return proto.CompactTextString(m) }
func (*TopicOther) ProtoMessage()               {}
func (*TopicOther) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{12} }

func (m *TopicOther) GetTopicName() string {
	if m != nil {
		return m.TopicName
	}
	return ""
}

func (m *TopicOther) GetVideoCount() string {
	if m != nil {
		return m.VideoCount
	}
	return ""
}

func (m *TopicOther) GetQueCount() string {
	if m != nil {
		return m.QueCount
	}
	return ""
}

func (m *TopicOther) GetUserCount() string {
	if m != nil {
		return m.UserCount
	}
	return ""
}

// 内容数据  发布数据
type ReportContentPostDataResponse struct {
	Code string                     `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Msg  string                     `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Data *ReportContentPostDataBack `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ReportContentPostDataResponse) Reset()                    { *m = ReportContentPostDataResponse{} }
func (m *ReportContentPostDataResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportContentPostDataResponse) ProtoMessage()               {}
func (*ReportContentPostDataResponse) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{13} }

func (m *ReportContentPostDataResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ReportContentPostDataResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ReportContentPostDataResponse) GetData() *ReportContentPostDataBack {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReportContentPostDataBack struct {
	PostVideoCount    uint64 `protobuf:"varint,1,opt,name=postVideoCount" json:"postVideoCount,omitempty"`
	PostVideoQueCount uint64 `protobuf:"varint,2,opt,name=postVideoQueCount" json:"postVideoQueCount,omitempty"`
}

func (m *ReportContentPostDataBack) Reset()                    { *m = ReportContentPostDataBack{} }
func (m *ReportContentPostDataBack) String() string            { return proto.CompactTextString(m) }
func (*ReportContentPostDataBack) ProtoMessage()               {}
func (*ReportContentPostDataBack) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{14} }

func (m *ReportContentPostDataBack) GetPostVideoCount() uint64 {
	if m != nil {
		return m.PostVideoCount
	}
	return 0
}

func (m *ReportContentPostDataBack) GetPostVideoQueCount() uint64 {
	if m != nil {
		return m.PostVideoQueCount
	}
	return 0
}

// pushVV列表
type ReportPushVVListRequest struct {
	Date      string `protobuf:"bytes,1,opt,name=date" json:"date,omitempty"`
	PhotoType string `protobuf:"bytes,2,opt,name=photoType" json:"photoType,omitempty"`
}

func (m *ReportPushVVListRequest) Reset()                    { *m = ReportPushVVListRequest{} }
func (m *ReportPushVVListRequest) String() string            { return proto.CompactTextString(m) }
func (*ReportPushVVListRequest) ProtoMessage()               {}
func (*ReportPushVVListRequest) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{15} }

func (m *ReportPushVVListRequest) GetDate() string {
	if m != nil {
		return m.Date
	}
	return ""
}

func (m *ReportPushVVListRequest) GetPhotoType() string {
	if m != nil {
		return m.PhotoType
	}
	return ""
}

type ReportPushVVListResponse struct {
	Code string                  `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Msg  string                  `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Data []*ReportPushVVListBack `protobuf:"bytes,3,rep,name=data" json:"data,omitempty"`
}

func (m *ReportPushVVListResponse) Reset()                    { *m = ReportPushVVListResponse{} }
func (m *ReportPushVVListResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportPushVVListResponse) ProtoMessage()               {}
func (*ReportPushVVListResponse) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{16} }

func (m *ReportPushVVListResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ReportPushVVListResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ReportPushVVListResponse) GetData() []*ReportPushVVListBack {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReportPushVVListBack struct {
	Vid      string `protobuf:"bytes,1,opt,name=vid" json:"vid,omitempty"`
	SumVV    string `protobuf:"bytes,2,opt,name=sumVV" json:"sumVV,omitempty"`
	Title    string `protobuf:"bytes,3,opt,name=title" json:"title,omitempty"`
	Duration string `protobuf:"bytes,4,opt,name=duration" json:"duration,omitempty"`
}

func (m *ReportPushVVListBack) Reset()                    { *m = ReportPushVVListBack{} }
func (m *ReportPushVVListBack) String() string            { return proto.CompactTextString(m) }
func (*ReportPushVVListBack) ProtoMessage()               {}
func (*ReportPushVVListBack) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{17} }

func (m *ReportPushVVListBack) GetVid() string {
	if m != nil {
		return m.Vid
	}
	return ""
}

func (m *ReportPushVVListBack) GetSumVV() string {
	if m != nil {
		return m.SumVV
	}
	return ""
}

func (m *ReportPushVVListBack) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *ReportPushVVListBack) GetDuration() string {
	if m != nil {
		return m.Duration
	}
	return ""
}

// pushVV数
type ReportPushVVRequest struct {
	STime     string `protobuf:"bytes,1,opt,name=sTime" json:"sTime,omitempty"`
	ETime     string `protobuf:"bytes,2,opt,name=eTime" json:"eTime,omitempty"`
	PhotoType string `protobuf:"bytes,3,opt,name=photoType" json:"photoType,omitempty"`
}

func (m *ReportPushVVRequest) Reset()                    { *m = ReportPushVVRequest{} }
func (m *ReportPushVVRequest) String() string            { return proto.CompactTextString(m) }
func (*ReportPushVVRequest) ProtoMessage()               {}
func (*ReportPushVVRequest) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{18} }

func (m *ReportPushVVRequest) GetSTime() string {
	if m != nil {
		return m.STime
	}
	return ""
}

func (m *ReportPushVVRequest) GetETime() string {
	if m != nil {
		return m.ETime
	}
	return ""
}

func (m *ReportPushVVRequest) GetPhotoType() string {
	if m != nil {
		return m.PhotoType
	}
	return ""
}

type ReportPushVVResponse struct {
	Code string              `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Msg  string              `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Data []*ReportPushVVBack `protobuf:"bytes,3,rep,name=data" json:"data,omitempty"`
}

func (m *ReportPushVVResponse) Reset()                    { *m = ReportPushVVResponse{} }
func (m *ReportPushVVResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportPushVVResponse) ProtoMessage()               {}
func (*ReportPushVVResponse) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{19} }

func (m *ReportPushVVResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ReportPushVVResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ReportPushVVResponse) GetData() []*ReportPushVVBack {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReportPushVVBack struct {
	Date  string `protobuf:"bytes,1,opt,name=date" json:"date,omitempty"`
	SumVV string `protobuf:"bytes,2,opt,name=sumVV" json:"sumVV,omitempty"`
}

func (m *ReportPushVVBack) Reset()                    { *m = ReportPushVVBack{} }
func (m *ReportPushVVBack) String() string            { return proto.CompactTextString(m) }
func (*ReportPushVVBack) ProtoMessage()               {}
func (*ReportPushVVBack) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{20} }

func (m *ReportPushVVBack) GetDate() string {
	if m != nil {
		return m.Date
	}
	return ""
}

func (m *ReportPushVVBack) GetSumVV() string {
	if m != nil {
		return m.SumVV
	}
	return ""
}

// 运营当日操作数
type ReportOpTodayResponse struct {
	Code string        `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Msg  string        `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Data *ReportOpBack `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ReportOpTodayResponse) Reset()                    { *m = ReportOpTodayResponse{} }
func (m *ReportOpTodayResponse) String() string            { return proto.CompactTextString(m) }
func (*ReportOpTodayResponse) ProtoMessage()               {}
func (*ReportOpTodayResponse) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{21} }

func (m *ReportOpTodayResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ReportOpTodayResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ReportOpTodayResponse) GetData() *ReportOpBack {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReportOpBack struct {
	AllPush     uint32 `protobuf:"varint,1,opt,name=allPush" json:"allPush,omitempty"`
	HasQuestion uint32 `protobuf:"varint,2,opt,name=hasQuestion" json:"hasQuestion,omitempty"`
}

func (m *ReportOpBack) Reset()                    { *m = ReportOpBack{} }
func (m *ReportOpBack) String() string            { return proto.CompactTextString(m) }
func (*ReportOpBack) ProtoMessage()               {}
func (*ReportOpBack) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{22} }

func (m *ReportOpBack) GetAllPush() uint32 {
	if m != nil {
		return m.AllPush
	}
	return 0
}

func (m *ReportOpBack) GetHasQuestion() uint32 {
	if m != nil {
		return m.HasQuestion
	}
	return 0
}

func init() {
	proto.RegisterType((*ReportVideoExeclResponse)(nil), "api.ReportVideoExeclResponse")
	proto.RegisterType((*ReportVideoExeclBack)(nil), "api.ReportVideoExeclBack")
	proto.RegisterType((*ReportUserExeclResponse)(nil), "api.ReportUserExeclResponse")
	proto.RegisterType((*ReportUserExeclBack)(nil), "api.ReportUserExeclBack")
	proto.RegisterType((*ReportVideoStatisticsResponse)(nil), "api.ReportVideoStatisticsResponse")
	proto.RegisterType((*ReportVideoStatisticsBack)(nil), "api.ReportVideoStatisticsBack")
	proto.RegisterType((*StatisBizDaily)(nil), "api.StatisBizDaily")
	proto.RegisterType((*ReportUserStatisticsResponse)(nil), "api.ReportUserStatisticsResponse")
	proto.RegisterType((*ReportUserStatisticsBack)(nil), "api.ReportUserStatisticsBack")
	proto.RegisterType((*StatisUser)(nil), "api.StatisUser")
	proto.RegisterType((*ReportContentTopicDataResponse)(nil), "api.ReportContentTopicDataResponse")
	proto.RegisterType((*ReportContentTopicDataBack)(nil), "api.ReportContentTopicDataBack")
	proto.RegisterType((*TopicOther)(nil), "api.TopicOther")
	proto.RegisterType((*ReportContentPostDataResponse)(nil), "api.ReportContentPostDataResponse")
	proto.RegisterType((*ReportContentPostDataBack)(nil), "api.ReportContentPostDataBack")
	proto.RegisterType((*ReportPushVVListRequest)(nil), "api.ReportPushVVListRequest")
	proto.RegisterType((*ReportPushVVListResponse)(nil), "api.ReportPushVVListResponse")
	proto.RegisterType((*ReportPushVVListBack)(nil), "api.ReportPushVVListBack")
	proto.RegisterType((*ReportPushVVRequest)(nil), "api.ReportPushVVRequest")
	proto.RegisterType((*ReportPushVVResponse)(nil), "api.ReportPushVVResponse")
	proto.RegisterType((*ReportPushVVBack)(nil), "api.ReportPushVVBack")
	proto.RegisterType((*ReportOpTodayResponse)(nil), "api.ReportOpTodayResponse")
	proto.RegisterType((*ReportOpBack)(nil), "api.ReportOpBack")
}

func init() { proto.RegisterFile("reportService.proto", fileDescriptor7) }

var fileDescriptor7 = []byte{
	// 995 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x57, 0x5d, 0x73, 0xdb, 0x44,
	0x14, 0x1d, 0xd9, 0x4e, 0x5a, 0xdf, 0x38, 0x69, 0xba, 0x49, 0x8a, 0xe3, 0x49, 0x42, 0x50, 0xf9,
	0x08, 0x33, 0x25, 0x40, 0xfa, 0xc4, 0x0c, 0x4f, 0x75, 0xe9, 0x03, 0x30, 0x71, 0xa2, 0x3a, 0x7e,
	0xe0, 0x85, 0x59, 0xa4, 0xa5, 0xd9, 0x41, 0xf6, 0x2a, 0xab, 0x95, 0x83, 0xe1, 0x95, 0x37, 0x66,
	0xf8, 0x61, 0xfc, 0x00, 0x7e, 0x0f, 0xb3, 0x1f, 0x92, 0x56, 0xf2, 0xca, 0xd3, 0xf8, 0x4d, 0xf7,
	0xdc, 0xeb, 0xa3, 0x73, 0xcf, 0xde, 0xdd, 0x95, 0x61, 0x8f, 0x93, 0x84, 0x71, 0xf1, 0x96, 0xf0,
	0x39, 0x0d, 0xc9, 0x79, 0xc2, 0x99, 0x60, 0xa8, 0x8d, 0x13, 0x3a, 0xe8, 0x85, 0x6c, 0x3a, 0x65,
	0x33, 0x0d, 0xf9, 0x0c, 0xfa, 0x81, 0xaa, 0x9c, 0xd0, 0x88, 0xb0, 0xef, 0x7e, 0x27, 0x61, 0x1c,
	0x90, 0x34, 0x61, 0xb3, 0x94, 0x20, 0x04, 0x9d, 0x90, 0x45, 0xa4, 0xef, 0x9d, 0x7a, 0x67, 0xdd,
	0x40, 0x3d, 0xa3, 0x5d, 0x68, 0x4f, 0xd3, 0x77, 0xfd, 0x96, 0x82, 0xe4, 0x23, 0xfa, 0x02, 0x3a,
	0x11, 0x16, 0xb8, 0xdf, 0x3e, 0xf5, 0xce, 0xb6, 0x2e, 0x0e, 0xcf, 0x71, 0x42, 0xcf, 0xeb, 0x94,
	0xaf, 0x70, 0xf8, 0x5b, 0xa0, 0xca, 0xfc, 0xaf, 0x60, 0xdf, 0x95, 0x45, 0x7d, 0x78, 0xf4, 0x2b,
	0x8d, 0xc9, 0x0d, 0x8f, 0xcd, 0xfb, 0xf2, 0xd0, 0x9f, 0xc2, 0x07, 0xfa, 0x17, 0x37, 0x29, 0xe1,
	0xeb, 0x28, 0x7c, 0x51, 0x51, 0xd8, 0xb7, 0x14, 0x16, 0x8c, 0x96, 0xc0, 0x2f, 0x61, 0xcf, 0x91,
	0x5c, 0xa1, 0x6f, 0x01, 0xc7, 0x56, 0x47, 0x6f, 0x05, 0x16, 0x34, 0x15, 0x34, 0x4c, 0x1f, 0xa8,
	0xf2, 0xa2, 0xa2, 0xf2, 0xa4, 0xee, 0x63, 0xc9, 0x6b, 0x69, 0xfd, 0x09, 0x0e, 0x1b, 0x4b, 0xd0,
	0x67, 0x86, 0xd0, 0x3b, 0x6d, 0x9f, 0x6d, 0x5d, 0xec, 0x29, 0x42, 0x5d, 0xf2, 0x8a, 0xfe, 0xf1,
	0x1a, 0xd3, 0x78, 0xa1, 0x59, 0xd0, 0x3e, 0x6c, 0x84, 0x2c, 0x9b, 0x09, 0xa5, 0xa6, 0x13, 0xe8,
	0xc0, 0xff, 0xb7, 0x05, 0x3b, 0xd5, 0x72, 0xb4, 0x03, 0x2d, 0x1a, 0x99, 0x36, 0x5a, 0x34, 0x42,
	0x9f, 0xc2, 0xce, 0x5c, 0xaf, 0x62, 0xc2, 0x52, 0x72, 0x99, 0x4d, 0x4d, 0x3f, 0x35, 0x14, 0x7d,
	0x0c, 0xdb, 0x0a, 0x19, 0xc6, 0x34, 0x14, 0xb2, 0xac, 0xad, 0xca, 0xaa, 0x20, 0xf2, 0xa1, 0xa7,
	0x80, 0x09, 0x25, 0xf7, 0xb2, 0xa8, 0xa3, 0x8a, 0x2a, 0x58, 0xc1, 0xf4, 0x06, 0xcf, 0x19, 0x97,
	0x45, 0x1b, 0x16, 0x53, 0x0e, 0xa2, 0x33, 0x78, 0x22, 0x87, 0x9c, 0xcc, 0x44, 0x51, 0xb7, 0xa9,
	0xea, 0xea, 0x30, 0x3a, 0x01, 0x30, 0x90, 0x2c, 0x7a, 0xa4, 0x8a, 0x2c, 0x44, 0x76, 0x28, 0x58,
	0x42, 0xc3, 0x37, 0x2c, 0x8e, 0x99, 0x52, 0xf5, 0x58, 0x77, 0x58, 0x45, 0x15, 0x0f, 0x27, 0x58,
	0x90, 0x31, 0x9d, 0x92, 0x7e, 0xd7, 0xf0, 0x14, 0x88, 0x7f, 0x0f, 0x47, 0xe5, 0x50, 0xad, 0x3d,
	0x22, 0x5f, 0x57, 0x46, 0xe4, 0xb8, 0x36, 0xc8, 0xce, 0x09, 0xb9, 0xc9, 0xf7, 0xf7, 0x72, 0x05,
	0x7a, 0x5e, 0x19, 0x90, 0x27, 0xd6, 0x80, 0xc8, 0xe2, 0x95, 0xc3, 0xf1, 0x8f, 0x07, 0x50, 0x96,
	0x2e, 0x0d, 0xc6, 0x00, 0x1e, 0x0b, 0x26, 0x70, 0x5c, 0x8e, 0x44, 0x11, 0xa3, 0x23, 0xe8, 0xe2,
	0x50, 0xd0, 0x39, 0x29, 0x07, 0xa1, 0x04, 0xd0, 0x33, 0xd8, 0x9c, 0xd9, 0xcb, 0x6f, 0xa2, 0x9a,
	0xc1, 0x1b, 0x4b, 0x06, 0xff, 0x09, 0x27, 0xba, 0xcf, 0x21, 0x9b, 0x09, 0x32, 0x13, 0x63, 0xb9,
	0x3e, 0xaf, 0xb1, 0xc0, 0x0f, 0xb4, 0xf8, 0x65, 0xc5, 0xe2, 0x0f, 0x2d, 0x8b, 0xeb, 0xc4, 0x96,
	0xc9, 0x7f, 0x7b, 0x30, 0x68, 0x2e, 0x72, 0xfa, 0xac, 0x2a, 0x46, 0xe2, 0x76, 0xb5, 0xcf, 0xf2,
	0xd4, 0x11, 0xa3, 0xa1, 0xc2, 0xdb, 0x0a, 0xcf, 0x43, 0x95, 0x19, 0xea, 0x4c, 0xc7, 0x64, 0x74,
	0xe8, 0xff, 0xe5, 0x01, 0x94, 0xf4, 0xd2, 0x6f, 0x35, 0xac, 0x97, 0x78, 0x9a, 0x37, 0x5f, 0x02,
	0xd2, 0x57, 0xbd, 0x0b, 0x8b, 0x77, 0x77, 0x03, 0x0b, 0x91, 0x2b, 0x79, 0x97, 0x91, 0x52, 0x41,
	0x37, 0x28, 0x62, 0xc9, 0x9c, 0xa5, 0x84, 0x97, 0x22, 0xba, 0x41, 0x09, 0x94, 0xc7, 0xa2, 0xf1,
	0xe4, 0x8a, 0xa5, 0x62, 0x8d, 0x05, 0x69, 0x3e, 0x16, 0x6b, 0xbc, 0xd6, 0x7a, 0xdc, 0xe5, 0xc7,
	0xa2, 0xa3, 0x44, 0x6e, 0xe9, 0x84, 0xa5, 0xfa, 0xc4, 0xd4, 0xd2, 0x3d, 0xe5, 0x5f, 0x0d, 0x45,
	0x2f, 0xe0, 0x69, 0x81, 0x5c, 0xe7, 0x16, 0xe8, 0xc5, 0x59, 0x4e, 0xf8, 0x3f, 0xe4, 0x97, 0xd4,
	0x55, 0x96, 0xde, 0x4e, 0x26, 0x3f, 0xd2, 0x54, 0x04, 0xe4, 0x2e, 0x23, 0xa9, 0x90, 0x7d, 0x46,
	0x58, 0x14, 0x7d, 0xca, 0x67, 0x69, 0x5d, 0x72, 0xcb, 0x04, 0x1b, 0x2f, 0x12, 0x62, 0xba, 0x2d,
	0x81, 0xf2, 0x52, 0xb6, 0xc9, 0xd6, 0xbc, 0x94, 0xdb, 0xb5, 0x4b, 0xb9, 0xa4, 0xb4, 0x0c, 0x4b,
	0xf2, 0x4b, 0xb9, 0x9a, 0x95, 0xc4, 0xf3, 0x62, 0x63, 0xcb, 0x47, 0x39, 0xa6, 0x69, 0x36, 0x9d,
	0x4c, 0xcc, 0xcb, 0x74, 0x20, 0x51, 0x41, 0x45, 0x4c, 0xcc, 0x88, 0xe8, 0x40, 0xce, 0x4e, 0x94,
	0x71, 0x2c, 0x28, 0x9b, 0x99, 0xf1, 0x28, 0x62, 0xff, 0xe7, 0xfc, 0x96, 0xd5, 0x6f, 0xcc, 0xbd,
	0x92, 0xf4, 0x6a, 0x87, 0x7b, 0x86, 0x5e, 0x06, 0x12, 0xd5, 0xfb, 0xde, 0xbc, 0x54, 0x05, 0x55,
	0x0f, 0xdb, 0x75, 0x0f, 0xdf, 0x55, 0x5b, 0x7a, 0xa0, 0x7f, 0x9f, 0x57, 0xfc, 0x3b, 0x58, 0xf2,
	0xcf, 0xf2, 0xee, 0x5b, 0xd8, 0xad, 0x67, 0x9c, 0x4b, 0xee, 0x74, 0xce, 0x8f, 0xe0, 0x40, 0xff,
	0x7a, 0x94, 0x8c, 0x59, 0x84, 0x17, 0x0f, 0xd4, 0xf9, 0x49, 0x65, 0x77, 0x3c, 0xb5, 0x74, 0x8e,
	0x12, 0x4b, 0xe3, 0xf7, 0xd0, 0xb3, 0x51, 0x79, 0x78, 0xe0, 0x38, 0x96, 0x82, 0x15, 0xff, 0x76,
	0x90, 0x87, 0xe8, 0x14, 0xb6, 0x6e, 0x71, 0x7a, 0x2d, 0x17, 0x43, 0x2e, 0x5b, 0x4b, 0x65, 0x6d,
	0xe8, 0xe2, 0xbf, 0x0d, 0xd8, 0x0e, 0xec, 0x8f, 0x4b, 0xf4, 0x0d, 0x6c, 0x73, 0xbb, 0x07, 0xd4,
	0x53, 0x3a, 0x2e, 0xb3, 0x38, 0x0e, 0xc8, 0xdd, 0x60, 0x50, 0x51, 0x55, 0xed, 0x72, 0x08, 0x3d,
	0x6e, 0x99, 0x87, 0xfa, 0x4b, 0x4e, 0x9b, 0xc9, 0x18, 0x1c, 0x3a, 0x32, 0x86, 0x64, 0x04, 0xbb,
	0xbc, 0x36, 0xbd, 0xe8, 0xc8, 0x39, 0xf2, 0x39, 0xd9, 0x71, 0x43, 0xd6, 0x10, 0x06, 0x70, 0xc0,
	0x5d, 0xe7, 0x07, 0xd2, 0x83, 0x70, 0x9d, 0x11, 0xbe, 0xb0, 0xe9, 0xfc, 0xe6, 0x53, 0xa9, 0xe0,
	0x1c, 0xc3, 0x33, 0xee, 0xbc, 0x22, 0x9a, 0x48, 0x9f, 0xaf, 0xb8, 0x7b, 0x0a, 0xd6, 0x2b, 0xd8,
	0xe7, 0x8e, 0xeb, 0xbd, 0x89, 0xf3, 0xa3, 0xc6, 0x4f, 0x86, 0xe5, 0xde, 0x6b, 0x9f, 0x94, 0xef,
	0xd3, 0x7b, 0xd3, 0x07, 0xf0, 0x15, 0x1c, 0xba, 0x54, 0xaa, 0x8f, 0xeb, 0x26, 0xde, 0x23, 0xd7,
	0x67, 0xba, 0xa5, 0x72, 0xe0, 0x54, 0xb9, 0x92, 0xf2, 0xd8, 0xf9, 0xdf, 0x24, 0xe7, 0xfc, 0x65,
	0x53, 0xfd, 0x23, 0x7a, 0xf9, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x19, 0x9c, 0x88, 0x42, 0x3b,
	0x0d, 0x00, 0x00,
}
