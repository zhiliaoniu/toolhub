# ！！！请根据具体需要更改该脚本，此处只是默认脚本，非 100% 通用

# 该处输入内容将保存为 /data/services/$PKG_NAME-$VERSION/admin/resolve.sh

# ！！！注意！！！

# 由于系统自带 pidof 命令在 cron 中会随机 coredump ，这里的 pidof 是自定义函数，

# 仅支持 pidof <name> 的用法 ！！！！

#----------------#
# 初始化日志文件 #
#----------------#

log=$INSTALL_PATH/admin/resolve.log

true > $log

#--------------#
# 执行启动脚本 #
#--------------#

bash $INSTALL_PATH/admin/start.sh &>$log
