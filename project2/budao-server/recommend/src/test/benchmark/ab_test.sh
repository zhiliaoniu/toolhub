#!/bin/bash

#1.服务器地址
HOST=localhost
#URL=${HOST}:8080/budao.TimeLineService/GetTimeLine
URL=103.229.151.190:8080/budao.TimeLineService/GetTimeLine

#2.请求书,客户端数
#REQ_NUM=10
REQ_NUM=20000
#CLIENT_NUM="1"
CLIENT_NUM="100 500 1000 2000"
#CLIENT_NUM="10 100 500 1000"

#3.输出log
LOG_NAME=${HOST}_result.log

#4.req
PROTOCAL="application/protobuf"
#PROTOCAL="application/json"
#REQ="req_get_video_comment_list.pb"
REQ="req_gettimeline.pb"

echo "----------------------------------------------------------------" >> ${LOG_NAME}
echo "method: "${REQ} " date: "`date` >> ${LOG_NAME}
echo "|client_num|request per second|Time per request|time per request" >> ${LOG_NAME}
for j in ${CLIENT_NUM}
do
    TMP_LOG=${HOST}.${j}.result.log
	echo -n "|$j|" >> ${LOG_NAME}
    ab -k -n ${REQ_NUM} -c $j -T ${PROTOCAL} -p ${REQ} ${URL} > ${TMP_LOG}
	grep -E "Time per request|Requests per second" ${TMP_LOG} | awk '{print $4}' | xargs | sed "s/[ ]/|/g" >> ${LOG_NAME}
	rm -f ${TMP_LOG}
	sleep 3;
done
