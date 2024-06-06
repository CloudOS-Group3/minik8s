#/bin/sh

etcdctl del --prefix "" --endpoints=192.168.3.8:2379

nerdctl rm -f $(nerdctl ps -aq --no-trunc)

ipvsadm -C