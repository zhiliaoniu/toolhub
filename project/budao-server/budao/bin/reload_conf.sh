#!/bin/bash
echo $# >> ./reload_conf.log
if [ $# -ge 1 ];
then
    echo $1 >> ./reload_conf.log
fi

curl --request "POST" \
--location "http://localhost:8080/budao.MiscService/ReloadConf" \
--header "Content-Type:application/json" \
--data '{"conf_path": "/data/budao-server/scheduler/conf/budao/budao_abtest.json"}'
