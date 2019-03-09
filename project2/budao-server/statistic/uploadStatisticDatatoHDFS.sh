#!/bin/bash

cd /data/budao-server/logs

dateDstFile=`date -d -1hour +%Y-%m-%d-%H`
dateGrep=`date +%Y/%m/%d`
dateHour=`date -d -1hour +%H`
dateGrep=${dateGrep}" "${dateHour}
ipaddr=`/sbin/ifconfig eth0 | grep "inet addr:" | awk '{print $2}' | cut -c 6-`
targetfilename=${ipaddr}"-statistic-server-"${dateDstFile}".log"

localfilepath="/home/lidong1/"${targetfilename}
basefilename="statistic-server.log"
sudo grep '\[I\]' -R $basefilename | grep "${dateGrep}" | awk -v FS=']' '{print $3}' | sed 's/^[ \t]*//g' > ${localfilepath}

hdfsfilepath="/user/budao/statistic-server-log/"${targetfilename}
cd /home/hadoop/bin
./hadoop fs -put ${localfilepath} ${hdfsfilepath}

cd ~
#sudo rm ${localfilepath}