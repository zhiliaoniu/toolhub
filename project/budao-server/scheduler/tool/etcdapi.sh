# 增加etcd成员
# see https://coreos.com/etcd/docs/latest/v2/members_api.html
curl http://127.0.0.1:21987/v2/members -XPOST -H "Content-Type: application/json" -d '{"peerURLs": ["http://127.0.0.1:23989"]}'