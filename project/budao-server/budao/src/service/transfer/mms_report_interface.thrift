/**﻿# 
*# 富媒体监控系统接口（采集服务）
*# 
*# 基于：thrift-0.6.1
*#    
*#
*#
**/
namespace cpp com.yy.zbase.channel.mms.thrift.server
namespace java com.yy.process.thrift
namespace php com.yy.zbase.channel.mms.thrift.server
namespace perl com.yy.zbase.channel.mms.thrift.server
namespace py com.yy.zbase.channel.mms.thrift.server


/**
* 监控/举报信息签名认证包
*/
struct MmsSign {
	/**
	*密钥编码
	*/
	1:string appKey;
	/**
	*签名
	*/
	2:string sign;
	/**
	*签名参数名列表 
	*/
	3:list<string> signParNames;

}


/**
* 监控/举报信息上传附件
*/
struct MmsReportAttc{
	/**
	*上报类型TEXT，IMG，AUDIO，VIDEO，VIDEO_FILE，JSON
	*/
	1:string attcType;
	/**
	*上报文本，由产品接入方按照规定 json 结构上传，接入服务按照给定格式检查
	*/
	2:string attcText;
	/**
	*上报链接 
	*/
	3:string attcUrl;
	/**
	*(取消)二进制文件
	*/
	4:binary attcFile;	
}


/**
* 监控/举报信息
*/
struct MmsReport{

	/**
	*流水号（视频流系列编号 programid）
	*/
	1:string serial;
	
	/**
	*待审用户uid
	*/
	2:i64 uid;
	
	/**
	*上报时间
	*/
	3:string reportTime;	
	
	/**
	*上报主体 （备注）， 由接入双方按照约定 json 结构上传，接入服务按照给定格式检查
	*/
	4:string reportComment;
	
	/**
	*上报附件
	*/
	5:list<MmsReportAttc> attachments;
	
	/**
	*上传人uid（非必填）
	*/
	6:i64 uploadUid;
	
	/**
	*举报级别（非必填）
	*/
	7:string severity;
	
	/**
	*频道ID（非必填）
	*/
	8:i64 sid;
		
	/**
	*子频道ID（非必填）
	*/
	9:i64 ssid;
		
	/**
	*OWID（非必填）
	*/
	10:i64 owid;
	
	/**
	*人气（非必填）
	*/
	11:i64 pcu;
	
	/**
	*附带参数（非必填）
	*/
	12:string extPar;

	/**
	* 标题
	*/
	13:string title;
}


/**
* 监控/举报信息上传请求包
*/
struct MmsReportReq{
	/**
	*认证签名
	*/
	1:MmsSign mmsSign;
	/**
	*通道ID
	*/
	2:string chid;
	/**
	*应用appid
	*/
	3:string appid;
	/**
	*上报材料
	*/
	4:list<MmsReport> reports;	
}


/**
* 监控/举报信息响应记录
*/
struct MmsReportRspRec {
	/**
	*编码
	*/
	1:i32 code;
	/**
	*反馈
	*/
	2:string msg;
	/**
	*流水号 
	*/
	3:string serial;
}


/**
* 监控/举报信息响应
*/
struct MmsReportRsp {
	/**
	*编码
	*/
	1:i32 code;
	/**
	*反馈
	*/
	2:string msg;
	/**
	*流水号 
	*/
	3:list<MmsReportRspRec> mmsReportRspRecs;
}



/**
* 处罚指令执行请求包
*/
struct MmsReportCmdReq {
	/**
	*密钥编码
	*/
	1:string appKey;
	/**
	*流水号
	*/
	2:string serial;
	/**
	*指令
	*/
	3:string cmd;
	/**
	*处罚理由
	*/
	4:string reason;
	/**
	*反馈信息
	*/
	5:string msg;
	/**
	*附带参数
	*/
	6:string extPar;
	/**
	*签名
	*/
	7:string sign;
	/**
	*审核状态
	*/
	8:string status;
	
}

/**
* 处罚指令执行结果
*/
struct MmsReportCmdRsp {
	/**
	*编码
	*/
	1:i32 code;
	/**
	*反馈
	*/
	2:string msg;
}

/**
*  查询审核结果
*/
struct AuditResult {
	1:string serial;

	2:string status;

	3:string reason;
}

/**
* 基本的接口，其他接口类都需要继承此类
**/
service BaseMmsThriftServ
{
	
	/*
	* 简单的连接测试
	*/
	void ping();
	
	/*
	* 简单的连接测试
	* @param ramdonId any number
	*/
	void pingWithParam(1:i32 randomId);
}


/**
* 服务接口
**/
service MmsReportServ extends BaseMmsThriftServ
{

    /**
    *  推送监控、举报记录        由接入方实现客户端。对接到监控系统服务器 
    * @param MmsReportReq 举报信息列表
    * @return MmsReportRsp
    */
	MmsReportRsp pushReports(1:MmsReportReq mmsReportReq);
	
    /**
    *  推送处罚指令    由接入方实现服务端，供监控系统调用
    * @param MmsReportCmdReq 处罚指令
    * @return MmsReportCmdRsp 处罚指令执行结果
    */
	MmsReportCmdRsp pushReportsCmd(1:MmsReportCmdReq mmsReportCmdReq);

    /**
    *  批量查询审核记录
    * @param serialList 流水号列表
    * @return AuditResult 查询结果
    */
    list<AuditResult> batchQueryAuditStatus(1:list<string> serialList, 2:string appId);
}



