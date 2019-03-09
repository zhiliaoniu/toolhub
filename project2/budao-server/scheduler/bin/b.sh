#!/bin/bash
echo $# >> ~/budao-server/scheduler/bin/a.log
if [ $# -ge 1 ];
then
    echo $1 >> ~/budao-server/scheduler/bin/a.log
fi

curl --request "POST" \
--location "http://localhost:8080/budao.MiscService/ReloadConf" \
--header "Content-Type:application/json" \
--data '{"conf_path": "/data/budao-server/scheduler/conf/budao/budao_abtest.json"}'
