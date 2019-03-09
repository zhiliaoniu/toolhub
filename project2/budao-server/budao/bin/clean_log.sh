#!/bin/bash

# 大于15天的都删除
find /data/budao-server/logs/* -mtime +15 -delete
