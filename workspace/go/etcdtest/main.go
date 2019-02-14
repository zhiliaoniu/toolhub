package main

import (
	"github.com/coreos/etcd/clientv3"
	"log"
)

type Client struct {
	etcdClient *clientv3.Client
}

func main() {
	log.Print("begin main")
	//1.检查当前集群是不是有 leader，如果有 leader，就 watch 这个 leader，只要发现 leader 掉了，就重新开始 1。

	//2.如果没有 leader，开始 campaign，创建一个 Lessor，并且通过 etcd 的事务机制写入相关信息，如下：

	//3.成为 leader 之后，我们对定期进行保活处理

	//4.所有做完之后，监听 lessor 是否过期，以及外面是否主动退出

	log.Print("end main")
}
