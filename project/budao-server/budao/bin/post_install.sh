# 该处输入内容将保存为 /data/services/$PKG_NAME-$VERSION/admin/post-install.sh, tar包解压之后执行。

# 注意！安装后脚本【固定】是使用 【root】 身份运行的 ！！！

# 这是因为大部分安装后操作需要 root 权限，例如创建目录，创建符号连接，安装 lib 库文件等。

# 如果你希望用其他身份运行，记得使用 su - <user> -c "<cmd>" 的方式 ！！！

# 如果 <cmd> 部分比较长，建议将 <cmd> 部分保存为一个脚本文件，放到打包的 bin/ 目录下，例如 test.sh

# 并使用诸如 su - <user> -c "bash $INSTALL_PATH/bin/test.sh" 方式调用 ！！！


log="$INSTALL_PATH/admin/post_install.log"
echo "post_install begin"
echo `date` >> $log
echo "pwd:`pwd`" >> $log 
echo "-----begin post-install.sh" >> $log
cd $INSTALL_PATH
echo `tree ./` >> $log
echo "-----end post-install.sh" >> $log
