# ！！！请根据具体需要更改该脚本，此处只是默认脚本，非 100% 通用

# 该处输入内容将保存为 /data/services/$PKG_NAME-$VERSION/admin/restart.sh

# 如果你是使用 root 用户重启，则无需使用 sudo ，脚本默认以 root 身份运行 ！！！！

# 如果你是使用 user_00 用户重启，请注意去掉 sudo ，否则可能导致脚本运行失败 ！！！！

# ！！！注意！！！

# 由于系统自带 pidof 命令在 cron 中会随机 coredump ，这里的 pidof 是自定义函数，

# 仅支持 pidof <name> 的用法 ！！！！

#---------------#
# 先停进程      #   
#---------------#
bash $INSTALL_PATH/admin/stop.sh  || exit 1

sleep 5

#---------------#
# 再起进程      #   
#---------------#

bash $INSTALL_PATH/admin/start.sh || exit 1
