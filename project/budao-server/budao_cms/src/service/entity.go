package service

import (
	"fmt"
)

type Question struct {
	Vid          string            `json:"vid"`
	Options      []*QuestionOption `json:"options"`
	Id           string            `json:"id"`
	Content      string            `json:"content"`
	OptionID     string            `json:"optionID"`
	TypeName     string            `json:"typeName"`
	QuestionType string            `json:"questionType"`
	Category     string            `json:"category"`
	CategoryName string            `json:"categoryName"`
	AnswerType   string            `json:"answerType"`
	Answer       string            `json:"answer"`
	Answers      []string          `json:"answers"`
	HardLevel    string            `json:"hardLevel"`
	Score        string            `json:"score"`
	CTime        string            `json:"cTime"`
	State        string            `json:"state"`
	UTime        string            `json:"uTime"`
	RightNum     uint              `json:"rightNum"`
	WrongNum     uint              `json:"wrongNum"`
}

type QuestionOption struct {
	Id         string `json:"id"`
	QuestionID string `json:"questionID"`
	Index      string `json:"index"`
	Content    string `json:"content"`
	IsAnswer   string `json:"isAnswer"`
	AnswerNum  string `json:"answerNum"`
}

type Video struct {
	Id            string `json:"id"`
	Vid           string `json:"vid"`
	SourceVid     string `json:"sourceVid"`
	VSourceID     string `json:"vSourceId"`
	PostVid       string `json:"postVid"`
	SourceType    string `json:"sourceType"`
	MediaName     string `json:"mediaName"`
	PraiseCount   string `json:"praiseCount"`
	FavCount      string `json:"favCount"`
	ShareCount    string `json:"shareCount"`
	CommentCount  string `json:"commentCount"`
	Topic         string `json:"topic"`
	Title         string `json:"title"`
	Duration      int    `json:"duration"`
	PlayCount     int    `json:"playCount"`
	VideoUrl      string `json:"videoUrl"`
	PlayUrl       string `json:"playUrl"`
	VideoCover    string `json:"videoCover"`
	Source        uint8  `json:"source"`
	Status        uint8  `json:"status"`
	VideoDuration int32  `json:"videoDuration"`
	VideoWidth    uint32 `json:"videoWidth"`
	VideoHeight   uint32 `json:"videoHeight"`
	TitleLen      string `json:"titleLen"`
	CommentNum    uint   `json:"commentNum"`
	FavorNum      uint   `json:"favorNum"`
	ShareNum      uint   `json:"shareNum"`
	ViewNum       uint   `json:"viewNum"`
	ParseType     int32  `json:"parseType"`
	InsertTime    string `json:"insertTime"`
}

type Topic struct {
	Id          string `json:"id"`
	TopicID     string `json:"topicID"`
	Version     string `json:"version"`
	Name        string `json:"name"`
	Pic         string `json:"pic"`
	Disable     string `json:"disable"`
	NeedLogin   string `json:"needLogin"`
	Weight      string `json:"weight"`
	Description string `json:"description"`
	Rule        string `json:"rule"`
	UserNum     string `json:"userNum"`
	VideoNum    string `json:"videoNum"`
	CTime       string `json:"cTime"`
	UTime       string `json:"uTime"`
	Hide        string `json:"hide"`
}

type Comment struct {
	VSourceID string     `json:"vSourceId"`
	Status    uint8      `json:"status"`
	Source    string     `json:"source"`
	Vid       string     `json:"vid"`
	Cid       string     `json:"cid"`
	Uid       string     `json:"uid"`
	UName     string     `json:"uName"`
	UPhoto    string     `json:"uPhoto"`
	PId       string     `json:"pId"`
	MasterId  string     `json:"masterId"`
	Content   string     `json:"content"`
	FavorNum  string     `json:"favorNum"`
	Weight    string     `json:"weight"`
	ReplyNum  string     `json:"replyNum"`
	State     string     `json:"state"`
	IsHot     byte       `json:"ishot"`
	CTime     string     `json:"cTime"`
	UTime     string     `json:"uTime"`
	Sub       []*Comment `json:"children"`
}

//用户举报
type UserReportVideo struct {
	Id     string `json:"id"`
	Uid    string `json:"uid"`
	Vid    string `json:"vid"`
	Reason string `json:"reason"`
	Status string `json:"status"`
}

//用户评论举报
type UserReportComment struct {
	Id       string `json:"id"`
	Uid      string `json:"uid"`
	Cid      string `json:"Cid"`
	DeviceId string `json:"deviceId"`
	Reason   string `json:"reason"`
	Status   string `json:"status"`
}

//规则
type Rule struct {
	Id            string `json:"id"`
	TopicId       string `json:"topicID"`
	RuleName      string `json:"ruleName"`      //规则名称
	RuleCondition string `json:"ruleCondition"` //规则条件
	RuleStatus    string `json:"ruleStatus"`    //规则状态
	TopicName     string `json:"topicName"`     //话题名称
	PostCount     string `json:"postCount"`     //规则发布次数
}

//push
type Push struct {
	Id            string `json:"id"`
	OpUid         string `json:"opUid"`         //运营uid
	PushObj       string `json:"pushObj"`       //推送对象，人群
	PushChannel   string `json:"pushChannel"`   //推送渠道,eg:客户端
	PushType      string `json:"pushType"`      //推送方式
	PushStatus    string `json:"pushStatus"`    //推送状态
	PushTitle     string `json:"pushTitle"`     //推送标题
	PushContent   string `json:"pushContent"`   //推送内容
	PushDescribe  string `json:"pushDescribe"`  //推送描述
	PushTime      string `json:"pushTime"`      //推送时间
	PushEndTime   string `json:"pushEndTime"`   //推送终止时间
	PushVideoType string `json:"pushVideoType"` //推送视频类型
	PushUrlType   string `json:"pushUrlType"`   //推送点击方式
	PushUrl       string `json:"pushUrl"`       //推送连接
	Device        string `json:"device"`        //设备名称
	Status        string `json:"status"`        //状态
	CreatedAt     string `json:"createdAt"`     //创建时间
	UpdatedAt     string `json:"updatedAt"`     //更新时间
	CommonStr     string `json:"commonStr"`
}

//视频相关统计
type StatisBizDaily struct {
	Id              string `json:"id"`
	VideoExposeNum  string `json:"vExposeNum"`   //视频播放
	VideoClictNum   string `json:"vClictNum"`    //视频点击
	VideoViewNum    string `json:"vViewNum"`     //视频观看
	VideoFavorNum   string `json:"vFavorNum"`    //视频点赞
	CommentFavorNum string `json:"commFavorNum"` //评论点赞
	CommentNum      string `json:"commNum"`      //评论发表
	TopicFollowNum  string `json:"tFollowNum"`   //话题订阅数
	CreateTime      string `json:"cTime"`
}

//用户相关统计
type StatisNewUser struct {
	Id         string `json:"id"`
	TotalNum   string `json:"totalNum"`  //当前用户总数
	ActiveNum  string `json:"activeNum"` //前一天活跃用户数
	NewNum     string `json:"newNum"`    //前一天新增用户数
	CreateTime string `json:"cTime"`
}

//banner
type Banner struct {
	Id          string `json:"id"`
	PicUrl      string `json:"picUrl"`      //图片路径
	ClickUrl    string `json:"clickUrl"`    //图片点击路径
	Position    string `json:"position"`    //图片位置
	Description string `json:"description"` //描述
	Status      string `json:"status"`      //状态
	FromTime    string `json:"fromTime"`    //banner的起始时间
	ToTime      string `json:"toTime"`      //banner的终止时间
	CreatedAt   string `json:"createdAt"`   //创建时间
	UpdatedAt   string `json:"updatedAt"`   //更新时间
}

//opRecord　运营记录
type OpRecord struct {
	UserId string //微信用户名
	Route  string //路由
}

//前端传入的分页数据
type Page struct {
	Num  int
	Size int
}

func (p Page) TurnSql() string {
	if p.Size < 0 || p.Num < 0 {
		return ""
	}
	if p.Num == 0 {
		p.Num = 0
	} else {
		p.Num -= 1
	}
	if p.Size == 0 {
		p.Size = 10
	}
	return fmt.Sprintf(` LIMIT  %v,%v`, p.Num*p.Size, p.Size)
}
