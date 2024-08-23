# minik8s
minik8s是一个简单的容器编排工具，其架构仿照kubernetes的基本架构进行实现。项目的开发语言是go语言。

## 项目总体架构
开发语言：go语言。选择go语言的原因是，kubernete官方的实现也是基于go语言实现的，此外，k8s用到的大部分组件，如etcd，containerd等，在go语言中有很好的支持，所以本项目选择使用go语言。
项目架构：参考官方k8s的实现，我们也实现了大部分k8s的组件和功能，具体而言，有如下组件。
- 客户端（kubectl）：是用户与minik8s系统进行交互的入口，可以通过kubectl执行各种类型的执行，对集群的资源进行管理
- 控制平面
  - apiserver：是控制平面的核心，提供了一系列Restful的接口，供其他组件调用
  - etcd：是minik8s的数据库，集群中的所有资源都会以json字符串的形式存在etcd中
  - scheduler：负责将创建pod的请求按照一定的算法，均衡地发送到不同的工作节点上
  - controller：负责动态监控集群中各种功能的更新，维护集群状态
  - kafka：是控制平面与工作节点交互的入口，所有发送到工作节点的消息都通过kafka异步地发送
- 工作节点
  - kubeproxy：监控endpoint的变化，维护路由转发规则
  - kubelet：控制某一工作节点上所有的pod的状态，管理pod的生命周期


## 软件栈
在实现minik8s系统的过程中，我们用到的开源软件栈列举如下。
- github.com/IBM/sarama：kafka的go语言实现
- github.com/containerd/containerd：容器运行时环境
- github.com/fatih/color：命令行彩色输出
- github.com/hashicorp/consul/api：Cousul服务器相关库，用于Prometheus服务发现
- github.com/moby/ipvs：Linux 内核提供的一种负载均衡技术
- github.com/olekukonko/tablewriter：命令行打印表格
- github.com/prometheus/client_golang：Prometheus的go语言客户端
- github.com/spf13/cobra：go语言命令行工具实现框架
- go.etcd.io/etcd/client/v3：go语言etcd客户端支持
- gopkg.in/yaml.v3：yaml文件解析工具

## 分工和贡献
1. 卢天宇
- kubectl命令行编写
- 部分API对象结构定义
- 部分apiserver路由规则及处理函数实现
- deployment抽象定义与实现
- hpa抽象定义与实现
- serverless的scale-to-zero逻辑实现
- 部分辅助功能实现（例如命令行彩色表格打印，log格式化输出等）
- 项目验收答辩文档编写
- 答辩视频录制
2. 杜心敏
- pod抽象实现，设置容器参数，pod内部localhost访问，pod间通信（cni）
- 实现service抽象，对pod动态选择，ipvs流量转发，集群内clusterIp访问，外部nodePort访问
- service多机部署
- function抽象实现，镜像构建
- workflow实现，分支规则判断
- severless并发测试与分支测试
- 相关资源增删改查、前端展示
- 文档编写、视频录制
3. 陈英昊
- node抽象和控制器
- kubelet心跳机制
- scheduler调度器
- DNS功能
- etcd数据库接口的包装，包括增删改查和监听
- kafka消息队列接口的包装，并通过包装好的发布/订阅模型实现监听机制的模板
- serverless架构设计，job抽象的实现和控制器，事件触发器和http触发器的设计
- 视频录制和文档编写


## 项目管理
### 项目分支
- feature/*分支：各个具体功能对应的分支，如feature/controller，feature/apiserver
- dev(elopment)分支：开发分支，当功能分支完成一个功能的开发时，合并到dev分支
- master分支：项目的发行分支，当项目通过所有测试并具有相对稳定的功能时，才会合并到master分支
### CI/CD与测试
本项目采用的是github的持续集成和部署工具，具体而言，我们在项目的.github/workflow目录下，编写了测试文件，每当进行push或pull request的时候，都会借助go原生支持的测试功能go test ./...对项目中所有以_test.go结尾的文件进行测试，当且仅当所有测试全部通过之后才把新的提交合并到dev或master分支。
### 功能开发流程
- 迭代式开发：分为三次迭代，每次进行进度检验、功能测试、任务分配、计划制定与修改。
- 合作方式：基本保证每天4小时成员集中在图书馆自习室开发新功能，加强合作沟通。
- 功能开发的过程：需求分析->设计API对象->运行逻辑代码->前端展示代码->单元测试->多机部署->集成测试
## 项目组件详解
### kubectl
kubectl是使用minik8s的用户对集群内的资源进行操作和管理的工具，本项目基于cobra开源库实现了所有功能。具体而言，kubectl支持四种操作：
- kubectl get：用于获取集群中各种资源的信息并以表格的形式打印在命令行中。
- kubectl apply：用于在集群中创建一个新的资源。
- kubectl delete：用于删除集群中的某一资源。
- kubectl http：用于触发serverless服务中的函数
### apiserver
apiserver是控制平面的核心，提供一系列Restful接口，控制平面的所有组件和kubectl均通过http请求与apiserver交互，同时apiserver会管理集群中资源的变化，通过kafka发布消息，从而与组件进行交互
### etcd
etcd是minik8s数据存储的中心，存放所有需要存储的对象（比如Pod，Service等），所有存储的对象都放在/registry开头的路径下，通过apiserver可以访问/修改资源，同时apiserver也监听etcd中资源的变化
同时，在serverless中，etcd开放了/trigger开头的路径作为事件触发器，用户通过在/trigger/<namespace>/<name>中存放数据可以触发function
### kafka
kafka是minik8s的消息队列，apiserver通过kafka发布集群中资源的变化，组件通过订阅这些变化来执行一定的操作，保证集群正常运行
### scheduler
scheduler负责调度集群当中的Pod，监听集群中Pod和Node的变化
- 当监听到Pod变化时，scheduler检查该Pod是否被调度。若Pod未被调度，scheduler会根据round robin策略调度到已经Ready的Node上
- 当监听到Node变化时，刷新在scheduler中存有的Node以保证能正确完成调度
### controller
minik8s中共有8个controller：
1. Deployment Controller：
  - Deployment controller每隔一段固定（30秒）的时间，根据已经创建的deployment资源，来自动增删pod副本
2. DNS Controller：
  - DNS Controller负责集群中DNS规则的修改，监听集群中DNS和Service的变化
  - 当监听到DNS变化时，DNS Controller会写入nginx配置文件和/etc/hosts文件，修改DNS规则
  - 当监听到Service变化时，DNS Controller会根据Service的变化重新修改DNS规则
3. Endpoint Controller：
  - 监听pod创建、删除、更新（过滤掉label不变的情况），修改对应的endpoints
  - 监听service创建、删除、更新（过滤掉label不变的情况），修改对应的endpoints
  - pod、service变化的消息格式包含新对象、旧对象，需要修改新abel、旧label对应的两个endpoints
  - etcd中存储了LabelIndex结构，实现双向索引：由label唯一确定；存储label对应的services，pods（namespace + name）
  - Endpoint结构：
    - Endpoint作为service的status，表示service需要管理的所有pod。一个service对应一组endpoints
    - EndPoint里面保存对应的pod ip和所有container暴露的ports，以及对应的service port，也就是一个pod会对应一个endpoint
4. HPA Controller：
  - Hpa controller每隔一段固定的时间（5秒）去一次检查各个已经创建的hpa资源是否处于可更新的状态，例如，一个Hpa资源要求每30秒更新一次，那么hpa controller将会在第六次检查的时候更新该hpa的状态
  - 在更新hpa的状态时，系统会自动获取该hpa对应的所有pod的资源利用情况，例如cpu使用率和memory占用率，并根据该数据与hpa资源定义的扩缩容标准值，在最大和最小副本数量之间，自动扩缩容
5. Job Controller：
  - Job Controller负责serverless中Job的运行，监听集群中Job和Pod的变化
  - 当监听到Job变化时，若该Job仍未运行，Job Controller会根据当前Pod是否创建，运行该Job或放入等待队列中
  - 当监听到Pod变化时，Job Controller会检查当前等待队列，运行在该Pod上的Job
6. Node Controller：
  - Node Controller负责集群中Node生命周期的检查，监听集群中Node的变化，同时周期性检查Node的存活状况
  - 当监听到Node的变化时，刷新controller中存有的Node信息
  - Node Controller每30秒会启动一次检查，当有Node超过2分钟未更新自身的消息（发送心跳），则认为该Node已经失去连接，通知apiserver逐出该Node，重新部署在该Node上的Pod
7. Serverless Controller：
  - 监听Job的创建，一旦收到消息，新创建一个Job资源，并为其自动分配一个可用的pod，如果没用则自动创建
  - 会定期检查所有用于serverless服务的pod的空闲时间，如果空闲时间超过最大空闲时间限制，那么serverless controller则会自动将该pod删除
8. Workflow Controller：管理workflow调用运行状态
  - 监听上一级函数完成情况，获取返回结果
  - 根据返回结果进行rule判断，选择下一个执行的函数
  - 组装传入参数，包装http trigger，触发下一个函数执行
  - workflow全部执行完后，写入result发回给apiserver
### kubeproxy
kubeproxy运行在每个节点上，监听service创建、删除等请求，调用ipvs层更新路由转发规则。
### kubelet
kubelet负责集群Node上Pod和DNS的管理，监听集群中Pod和DNS的变化，并周期性检查Pod的状态，给apiserver发送心跳
- 当监听到Pod变化时，kubelet判断是否要创建/删除该Pod，并创建/删除该Pod，更新Node上Pod的信息
- 当监听到DNS变化时，kubelet修改/etc/hosts文件，保证DNS正常运行
- kubelet每30秒会检查该Node上的所有Pod状态，更新Pod的信息，并发送给apiserver
## 所有实现的功能
### 多机minik8s
1. Node相关使用：
  - 通过kubectl apply -f <filename>可以为集群加入Node
  - 通过kubectl get node可以获取集群中所有Node
  - 通过kubectl delete node <name>可以删除集群中的某个Node
2. Node结构定义：包含Node的名字和IP，以及该Node上PodIP的范围，通过yaml文件创建的Node状态为Unknown，需要启动对应的kubelet才能将状态变为Ready。
3. Node启动：启动kubelet即为启动Node，kubelet使用本机hostname作为该Node的名字，启动kubelet时会从apiserver获取该Node存储在etcd中的信息
4. Node生命维持：kubelet每隔30秒会给apiserver发送心跳，心跳中包含在该Node上所有Pod的最新信息（比如获取到的CPU/内存占用率），同时刷新Node的心跳时间，并且将Node的状态置为Ready
5. 多机调度：scheduler会记录下所有当前集群中的Node信息，主要关心Node的状态，当一个Pod需要调度时，状态不为Ready的所有Node都不会被调度，其他Node以round robin的策略进行调度
6. Node死亡：当一个Node超过2分钟没有给apiserver发送心跳时，Node Controller认定该Node死亡，此时Node Controller会给apiserver发出删除该Node的请求，并且将该Node上所有Pod的nodename清空，scheduler会重新调度这些Pod
### Pod
#### 功能使用
1. kubectl get 可以获取pod运行状态，可以指定namespace、name
2. kubectl apply -f <file-path>可以通过配置文件创建pod
3. kubectl delete pod -n <namespace> name 可以删除指定pod
4. pod内部容器之间支持localhost访问
5. 外部支持 pod ip+端口 访问
#### 实现方式
1. pod创建
  - apiserver 接收请求，存入etcd，并向kafka中发送创建消息
  - scheduler 监听pod创建消息，进行调度，将调度节点写入pod配置文件
  - kubelet 监听pod创建消息（过滤掉没有调度的pod），创建pod实例
    - pull对应镜像
    - 创建pause容器，cni插件（flannel）为pause容器分配ip地址
    - 创建其他容器，配置env、command、volumn等对应参数
  - kubelet 用心跳机制发回pod ip
2. pod内部通讯
  - 通过flannel为pause容器分配ip地址（作为pod ip），pod ip适用于集群
  - 其余容器通过linux namespace配置ipc、utc、network，共享pause容器网络（例如：/proc/<pause_pid>/ns/ipc）
  - pod内部容器之间支持localhost访问
#### CNI插件
选择flannel，通过kubenetes集群配置方式，配置了flannel插件，通过创建容器时指定 --network flannel。主要用于分配pod ip（pause容器ip）
### Service
1. service结构定义（通过yaml文件创建）：定义selector规则，端口映射，ClusterIp（用户定义的如果被占用会重新分配），对外访问的nodePort
2. ClusterIp生成：10.96.0.0/16，etcd中存储ClusterIp使用情况（map ClusterIp --> ServiceName）。
预先设置dummy网卡，并将service Cluster ip加入网络
ip L a minik8s0 type dummy
ip addr add 10.96.0.2/32 dev
echo 1 > /proc/sys/net/ipv4/vs/conntrack
3. 流量控制：使用 IPVS 控制流量转发。负载均衡选择RoundRobin。
  1. 创建service ：添加service（ClusterIp：port）到 IPVS 中。在添加之前，会检查service是否已存在于 IPVS 中，如果已存在则跳过。
  等效指令：ipvsadm -A -t <ClusterIP>:<Port> -s rr
  2. 根据service对应的endpoints，创建路由转发规则
  ipvsadm -a -t <ClusterIP>:<Port> -r <PodIP>:<PodPort> -m
  3. 删除service
  ipvsadm -D -t <ClusterIP>:<Port>
  4. 删除路由转发规则
  ipvsadm -d -t <ClusterIP>:<Port> -r <PodIP>:<PodPort>
4. NodePort配置：
  - 将主机端口访问转发到service ClusterIP：Port
  - 实现指令：iptables -t nat -A PREROUTING -p tcp --dport <NodePort> -j DNAT --to-destination <ClusterIP>:<Port>
5. 相关组件：
  - endpoint_controller：维护endpoints的动态更新
  - kubeproxy：监听service创建、删除等请求
### Deployment
1. Deployment的使用
  - kubectl apply -f <filename>可以创建一个deployment资源
  - kubectl get deployment [filename]可以获取对应filename的deployment，或者全部deployment的信息
  - kubectl delete deployment <filename>可以删除一个deployment资源
2. deployment文件中最主要的两个字段是：replicas和template，前者规定了该deployment需要集群中维持多少个pod副本，而后者则规定了每个pod的详细信息。
3. 每隔一段时间，deployment controller会向apiserver发送请求，获取当前集群中所有的deployment资源和pod资源，然后遍历每一个deployment，计算当前有多少pod属于这个deployment，然后如果发现实际pod数量与deployment规定的副本数量不同，则进行相应的增删操作
### HPA
1. Hpa的使用
  - kubectl apply -f <filename>可以创建一个hpa资源
  - kubectl get hpa [filename]可以获取对应filename的hpa，或者全部hpa的信息
  - kubectl delete hpa <filename>可以删除一个hpa资源
2. hpa的原理与deployment类似，会定期从apiserver获取当前集群中所有的hpa和pod资源，并遍历所有hpa，查找属于某一个hpa的所有pod
3. 对于找到的这些pod，首先获取他们的系统资源使用率（CPU和memory），如果他们超过了hpa文件规定扩缩容标准，则进行扩容，反之则进行缩容。例如，某一hpa的总体资源利用率为30%，而hpa文件规定的扩缩容标准为15%，则需要扩容两倍
4. Hpa controller还能自定义hpa的扩缩容间隔，这是通过在每次更新时，检查某一hpa的倒计时是否完成来完成的
### DNS
1. DNS相关使用：
  - 通过kubectl apply -f <filename>可以为集群加入一个DNS规则
  - 通过kubectl get dns可以获取集群中所有DNS规则
  - 通过kubectl delete dns <name>可以删除集群中的某个DNS规则
2. DNS结构定义：DNS中主要包含host和paths，host即为域名部分，paths即为子路径部分，每个DNS只能定义一个域名，但是可以有多个子路径，其中每个子路径对应一个Service的一个端口
3. 宿主机上DNS的实现：DNS分为域名部分和子路径部分的实现，其中域名部分通过修改/etc/hosts文件实现，子路径部分通过nginx的proxy_pass实现。
  - 对于域名部分，当集群中有一条新的DNS规则时，修改所有机器上的/etc/hosts，使得该DNS的域名被解析为192.168.3.8（集群中Master的IP），这样可以使得所有前往该域名的请求都经过Master节点上的nginx，从而得到转发
  - 对于子路径部分，当集群中有一条新的DNS规则时，apiserver查找对应Service的ClusterIP，并通过kafka将所有子路径对应的ClusterIP和Port发给DNS Controller，DNS Controller在/etc/nginx/conf.d文件夹中给每一个host都创建了一个nginx配置文件，配置文件内容如下：
server {
    listen 80;
    server_name <host>;
    location /<path1> {
        proxy_pass http://<ClusterIP>:<Port>;
    }
    location /<path2> {
        ...
    }
    ...
}
    修改完配置文件后，DNS Controller重启nginx，使得修改的配置文件生效。通过域名部分和子路径部分的设置，就可以实现集群中所有主机的DNS功能
4. Pod内部DNS的实现：在启动每个容器时，会将对应宿主机的/etc/hosts文件通过bind mount的方式绑定到容器中，因此容器中所有对于域名的请求也会转发到Master节点上的nginx并进行转发，从而实现了Pod内部的DNS功能
### 容错
1. etcd持久化：集群中所有资源均在etcd中持久化，在崩溃后数据仍然存储在etcd中，重启后可以拿取数据恢复，因此在重启后仍可以获取Pod和Service的信息
2. kubelet的容错：当控制面崩溃后，kubelet发送心跳会失败，但是kubelet仍然不会终止运行，而是继续管理Pod；当kubelet崩溃（未超过2分钟，集群仍认为该Node存活）重启时，kubelet会从apiserver中拿取所有需要的数据，继续管理Pod
3. controller的容错：当controller崩溃重启时，会从apiserver中拿取其所管理的资源，继续进行管理
### Serverless
1. Serverless相关使用：
  - 通过kubectl apply -f <filename>添加一个function，workflow或事件触发器
  - 通过kubectl get function查看集群中的function
  - 通过kubectl get workflow查看集群中的workflow
  - 通过kubectl get job查看function的调用结果/workflow的中间结果
  - 通过kubectl get trigger result查看workflow的最终结果
  - 通过kubectl delete job <name>删除一个job
  - 通过kubectl delete function <name>删除一个function
  - 通过kubectl delete workflow <name>删除一个workflow
  - 通过kubectl delete trigger <function_name>删除一个事件触发器
  - 通过kubectl http触发一个function/workflow
  - 通过给apiserver发送POST请求触发一个function/workflow，路径为/api/v1/namespaces/<function_namespace>/functions(workflow)/<function_name>/run
  - 通过往etcd中存储数据触发一个function，路径/trigger/<function_namespace>/<function_name>
2. Serverless架构图：
[图片]
3. Function抽象：
  - 用户可以通过单个python文件定义function，python文件名无要求，但是主函数名必须为main；
  - 用户需要指定参数和返回值的数量、类型，也可指定名字（用于workflow之间传递指定参数）；
  - 用户可以指定函数invoke方式：http或event trigger；
  - 用户需在对应文件夹下提供requirement.txt
4. Trigger抽象：无论通过http触发还是事件触发器触发，apiserver都会统一生成一条Trigger消息，Trigger消息包含本次调用的uuid，以及被调用的函数和参数，这条消息通过kafka发给Serverless Controller
5. Job抽象：Serverless Controller接收到Trigger消息后，生成一个Job，Job沿用Trigger消息的uuid，函数名称和参数，但是除此之外，Job还记录了当前调用处于的状态（Created，Running和Ended），调用的结果以及该次调用在哪个Pod上执行
6. Job对象的生成：Serverless Controller接收到Trigger消息后，检查被调用函数对应的Serverless实例（Pod）：
  - 若当前有空闲的Pod，Serverless Controller直接将该Pod分配给该Job
  - 若没有空闲的Pod，Serverless Controller创建一个新的Pod分配给该Job，这个Pod会在Job被发回给apiserver的时候被真正创建
  - Serverless Controller创建的Job处于Created状态
7. Serverless实例的创建：apiserver接收到从Serverless Controller创建的Job后，检查其中的Pod是否已被创建，若Pod仍未创建，表明这是一个新的Pod，按照Pod创建流程创建
8. Job对象的执行：Job Controller监听Job的变化，当一个Job被创建后，检查其中Pod是否已经被创建：
  - 若Pod没被创建，这个Job会进入等待队列当中，直到Job Controller监听到对应Pod被创建的消息，拿到Pod的IP，才会给Pod IP发http请求，从而调用对应函数
  - 若Pod已经被创建，直接给Pod发http请求调用函数
  - 调用函数的http请求会立刻返回，此时Job Controller把Job标记为Running，并发往apiserver
9. 容器内镜像的生成：
  1. 将所有需要的文件复制到工作路径，安装运行环境(pip install --no-cache-dir -r requirement.txt）
  2. 在容器中写入一个server.py，在容器运行时监听8080端口，并且执行用户函数，将结果发送至集群的apiserver
  3. 将镜像push到minik8s所维护的镜像仓库
该镜像仓库是运行在master  IP 下的 5050 端口上的一个 Docker registry
10. 容器内函数的执行：
  - 通过http、event trigger触发，根据函数配置文件中的参数（数量、类型）定义，组装参数
  - 通过json格式传入容器8080端口，并触发函数运行
  - 结果直接发给apiserver
11. 函数执行结果的返回：apiserver收到容器返回的结果后，取出对应的Job并将结果存入其中，将Job标记为Ended，代表一次函数调用的结束
12. Serverless实例的销毁（scale-to-0）：serverless controller会每隔一段固定时间检查controller所控制的所有pod的空闲时间，如果有pod处于空闲状态，并且空闲时间超过了预先设置的最大上限，则删除该pod
13. Workflow：
  - 每个节点的函数用namespace + name指定
  - 规则匹配采用switch-case模式，匹配满足express的第一个，否则执行default
  - expression需要参数名variable、判断关系opt（＞、=等）、值Value、数据类型type，支持多个expression的&&逻辑
  - workflow具体运行触发管理，由workflow_controller控制

#### Severless应用选择
[图片]
设计了一个分支较多的应用，在主机起了一个程序模拟内存数据库。
并发测试的时候（20并发），数据库容量设置了10，同时覆盖了不同分支选择。
并发测试配置如下：
[图片]
[图片]
### PV/PVC
1. 本系统的持久化存储是通过NFS实现的
2. 启动一台主机作为NFS的服务器，用于客户端将他们的目录挂载到服务器
3. PV有两种创建方式
  - 静态创建：用户自主通过yaml文件创建一个持久化卷资源
  - 动态方式：用户在创建pvc时发现系统中没有符合要求的pv资源，则系统会自动生成一个满足pvc要求的pv
4. pv被创建时，会在NFS服务器上生成一个挂载目录；创建pvc时，系统将符合要求的pv与该pvc进行绑定；创建pod时，通过制定一个pvc来与pvc进行绑定，同时在容器内部生成一个挂载点，与NFS服务器的挂载目录对应
5. 由于NFS服务器始终保持运行，即使pod崩溃，那么再次启动时，如果将新pod绑定到相同的pvc，那么之前的数据依然存在
### Prometheus
1. 集群节点资源的监控：在每个节点上启动node_exporter，暴露9100端口，从<NodeIP>:9100/metrics中就可以拿到node_exporter的信息，node_exporter包含CPU、内存等节点的信息
2. Pod中用户自定义资源的监控：在Pod配置文件中添加字段标明需要对Prometheus暴露的端口（例如2112），从<PodIP>:<Port>/metrics中可以获取用户自定义资源的信息
3. Prometheus服务发现：服务发现采用consul作为注册中心，位于Master节点上的8500端口。当节点启动或Pod创建时，在consul上进行服务的注册；当节点删除或Pod删除时，对应在consul上注销相应的服务，节点对应的服务名为node-exporter-<nodename>，而用户自定义Pod则为user-pod-<namespace>-<name>-<container_index>-<port_index>，对应的IP同上
4. Grafana配置：对外界安装Grafana的机器暴露Master的9000端口，在Grafana中抓取Master的9000端口的所有数据并进行筛选，选出node_exporter中的CPU和内存的信息和用户的自定义信息
### GPU
#### GPU任务提交流程：
1. kubectl 通过配置文件创建GPU-job
配置文件主要参数有：job-name, ntasks-per-node 等超算平台需要的参数，sourcePath 代码和编译文件存放路径。
2. apiServer 存入GPU-job，并发送创建pod消息
3. kubelet 创建pod实例
创建容器的时候，将配置文件中的参数，通过环境变量形式传入。
将所需要的cuda程序代码和编译文件，通过volumn映射传入容器，
创建container所使用的镜像中包含gpu_server.py，主要作用有：
  - 通过os.getenv 获取相关参数（环境变量全大写、下划线）
  - 通过参数编写作业脚本.slurm 
  - 连接超算平台，传送对应文件
  - 加载运行环境，编译程序，提交作业
  - 轮询查看是否完成（查看指定路径是否有结果文件输出）
  - 完成后给apiServer发消息（job-name和status），结果文件通过volumn映射传回主机
4. apiServer 收到job完成消息，修改job状态

#### CUDA程序说明
示例程序位于minik8s/testdata/Gpu
程序主要有以下几个部分：
1. CUDA 核函数
利用当前线程的行和列索引，编写运算逻辑
global void matrixMulKernel(float *A, float *B, float *C, int m, int k, int n) {
    int row = blockDim.y * blockIdx.y + threadIdx.y;
    int col = blockDim.x * blockIdx.x + threadIdx.x;

    if (row < m && col < n) {
        float sum = 0.0f;
        for (int i = 0; i < k; ++i) {
            sum += A[IDX2C(row, i, m)] * B[IDX2C(i, col, k)];
        }
        C[IDX2C(row, col, m)] = sum;
    }
}
2. 内存分配和初始化
3. 配置 CUDA 核函数的执行参数
这部分是使用GPU并发能力的关键，配置线程块、计算网格大小
// 配置 CUDA 核函数的执行参数
dim3 threadsPerBlock(16, 16);
dim3 blocksPerGrid((N + threadsPerBlock.x - 1) / threadsPerBlock.x,
                   (M + threadsPerBlock.y - 1) / threadsPerBlock.y);

// 调用矩阵乘法核函数
matrixMulKernel<<<blocksPerGrid, threadsPerBlock(devPtrA, devPtrB, devPtrC, M, K, N);
4. 结果传回主机和释放内存

