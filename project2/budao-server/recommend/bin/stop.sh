# ！！！请根据具体需要更改该脚本，此处只是默认脚本，非 100% 通用

# 该处输入内容将保存为 /data/services/$PKG_NAME-$VERSION/admin/stop.sh

# 如果你是使用 root 用户停止，则无需使用 sudo ，脚本默认以 root 身份运行 ！！！！

# 如果你是使用 user_00 用户停止，请注意去掉 sudo ，否则可能导致脚本运行失败 ！！！！

# 升级时调用当前内容停止旧进程

# ！！！注意！！！

# 由于系统自带 pidof 命令在 cron 中会随机 coredump ，这里的 pidof 是自定义函数，

# 仅支持 pidof <name> 的用法 ！！！！

#----------------#
# 初始化日志文件 #
#----------------#

log=$INSTALL_PATH/admin/stop.log

true > $log

#---------------#
# 进程数量检查  #
#---------------#

pid=$(pidof $APP_NAME)

if [ -z "$pid" ] ; then
    echo "no running $APP_NAME found , already stopped"
    exit 0
fi

#---------------#
# 停止进程      #
#---------------#

for i in $pid ; do
    echo "kill $APP_NAME pid=$i [$(ps --no-headers -lf $i)]"
    kill $i
    [ $? -eq 0 ] && ( bash /data/pkg/public-scripts/func/common-cleanup.sh $i ) &
    sleep 5
done

#---------------#
# 二次确认       #
#---------------#

if [ -z "$(pidof $APP_NAME)" ] ; then
    echo "stop $APP_NAME ok, all $APP_NAME got killed"
    echo "output last 20 lines of $log"
    tail -n 20 $log
    exit 0
else 
    echo "stop $APP_NAME failed, found $APP_NAME still running . see following"
    pidof $APP_NAME | xargs -r ps -lf
    echo "output last 20 lines of $log"
    tail -n 20 $log
    exit 1
fi
