# Service & kubeproxy
### service结构定义

通过yaml文件创建

定义selector规则，端口映射，ClusterIp（用户定义的如果被占用会重新分配），对外访问的nodePort

### ClusterIp生成

ip范围：10.96.0.0/16，etcd中存储ClusterIp使用情况（map ClusterIp --> ServiceName）。

预先设置dummy网卡，并将service Cluster ip加入网络

```shell
ip L a minik8s0 type dummy

ip addr add 10.96.0.2/32 dev

echo 1 > /proc/sys/net/ipv4/vs/conntrack
```



### 流量控制

使用 IPVS 控制流量转发。负载均衡选择RoundRobin。

1. 创建service ：添加service（ClusterIp：port）到 IPVS 中。在添加之前，会检查service是否已存在于 IPVS 中，如果已存在则跳过。
   等效指令：`ipvsadm -A -t <ClusterIP>:<Port> -s rr`
2. 根据service对应的endpoints，创建路由转发规则
   `ipvsadm -a -t <ClusterIP>:<Port> -r <PodIP>:<PodPort> -m`
3. 删除service
   `ipvsadm -D -t <ClusterIP>:<Port>`
4. 删除路由转发规则
   `ipvsadm -d -t <ClusterIP>:<Port> -r <PodIP>:<PodPort>`
5. NodePort配置：
   1. 将主机端口访问转发到service ClusterIP：Port
   2. 实现指令：`iptables -t nat -A PREROUTING -p tcp --dport <NodePort> -j DNAT --to-destination <ClusterIP>:<Port>`
6. 相关组件：

- `endpoint_controller`：维护endpoints的动态更新
- `kubeproxy`：监听service创建、删除等请求