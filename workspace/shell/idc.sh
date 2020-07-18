#!/bin/bash
host="change me"
while read line;
do
    #echo $line
    #ip=`echo $line| awk '{print $3}'`
    echo $line | grep -q -v "^[0-9]"
    if [ $? == 0  ];
    then
        echo $line
        continue
    fi

    #ip=`echo $line| awk '{print $1}'`
    ip=`echo $line| awk '{print $3}'`
    #echo $ip
    result=`curl -s -d "{\"systemid\":\"10022\", \"systemname\": \"dfs_fastdfs\", \"operator\": \"yangshengzhi\", \"fields\":[\"groupid\"], \"ip\":[\"${ip}\"]}" ${host}`
    `echo $result | grep -q "ok"`
    if [ $? != 0  ];
    then
        #echo "failed." $result
        continue
    fi
    result2=`echo $result|awk -F"\"" '{print $14}'`

    echo $line $result2
done < s
