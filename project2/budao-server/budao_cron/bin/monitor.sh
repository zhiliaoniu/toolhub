# ！！！请根据具体需要更改该脚本，此处只是默认脚本，非 100% 通用

# 该处输入内容将保存为 /data/services/$PKG_NAME-$VERSION/admin/monitor.sh

# ！！！注意！！！

# 由于系统自带 pidof 命令在 cron 中会随机 coredump ，这里的 pidof 是自定义函数，

# 仅支持 pidof <name> 的用法 ！！！！

#---------------#
# 进程数量设置  #
#---------------#

count=1

#----------------#
# 初始化日志文件 #
#----------------#

log=$INSTALL_PATH/admin/core.log

#---------------#
# 进程数量判断  #
#---------------#

pids=$(pidof $APP_NAME)

num=$(echo $pids |wc -w)

if [ $num -lt $count ] ; then

    # 进程数量异常，执行异常处理

    echo "$(date +'%F %T')| $APP_NAME num = $num [< $count] , pid=[$pids]" >> $log

    ps -lf $pids

    bash $INSTALL_PATH/admin/resolve.sh

    exit $?
else

    # 进程数量正常

    echo "current num of $APP_NAME = $num , pid=[$pids]"

    ps -lf $pids

    exit 0
fi
