#/bin/sh

etcdctl del --prefix ""

nerdctl rm -f $(nerdctl ps -aq --no-trunc)