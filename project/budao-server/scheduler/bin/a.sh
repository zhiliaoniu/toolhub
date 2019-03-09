#!/bin/bash
echo $# >> ~/budao-server/scheduler/bin/a.log
if [ $# -ge 1 ];
then
    echo $1 >> ~/budao-server/scheduler/bin/a.log
fi
