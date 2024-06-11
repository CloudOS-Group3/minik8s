# Pod

### 功能使用

1. `kubectl get` 可以获取pod运行状态，可以指定namespace、name
2. `kubectl apply -f <file-path>`可以通过配置文件创建pod
3. `kubectl delete pod -n <namespace> name` 可以删除指定pod
4. pod内部容器之间支持localhost访问
5. 外部支持 pod ip+端口 访问



### 实现方式

1. pod创建
   - `apiserver` 接收请求，存入etcd，并向kafka中发送创建消息
   - `scheduler` 监听pod创建消息，进行调度，将调度节点写入pod配置文件
   - `kubelet` 监听pod创建消息（过滤掉没有调度的pod），创建pod实例
     - pull对应镜像
     - 创建pause容器，cni插件（flannel）为pause容器分配ip地址
     - 创建其他容器，配置env、command、volumn等对应参数
   - `kubelet` 用心跳机制发回pod ip
2. pod内部通讯
   - 通过flannel为pause容器分配ip地址（作为pod ip），pod ip适用于集群
   - 其余容器通过linux namespace配置ipc、utc、network，共享pause容器网络（例如：/proc/<pause_pid>/ns/ipc）
   - pod内部容器之间支持localhost访问



### CNI插件

选择flannel，通过kubenetes集群配置方式，配置了flannel插件，通过创建容器时指定 --network flannel。主要用于分配pod ip（pause容器ip）