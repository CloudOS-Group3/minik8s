# Serverless

### 使用方式

Serverless相关使用：

1. 通过kubectl apply -f <filename>添加一个function，workflow或事件触发器
2. 通过kubectl get function查看集群中的function
3. 通过kubectl get workflow查看集群中的workflow
4. 通过kubectl get job查看function的调用结果/workflow的中间结果
5. 通过kubectl get trigger result查看workflow的最终结果
6. 通过kubectl delete job <name>删除一个job
7. 通过kubectl delete function <name>删除一个function
8. 通过kubectl delete workflow <name>删除一个workflow
9. 通过kubectl delete trigger <function_name>删除一个事件触发器
10. 通过kubectl http触发一个function/workflow
11. 通过给apiserver发送POST请求触发一个function/workflow，路径为/api/v1/namespaces/<function_namespace>/functions(workflow)/<function_name>/run
12. 通过往etcd中存储数据触发一个function，路径/trigger/<function_namespace>/<function_name>



### 架构设计

1. Serverless架构图：

![img](./img/serverless.png)

### Function抽象：

1. 用户可以通过单个python文件定义function，python文件名无要求，但是主函数名必须为`main`；
2. 用户需要指定参数和返回值的数量、类型，也可指定名字（用于workflow之间传递指定参数）；
3. 用户可以指定函数invoke方式：http或event trigger；
4. 用户需在对应文件夹下提供requirement.txt

~~~yaml
yaml文件示例如下：

```yaml
apiVersion: v1
kind: Function
metadata:
  name: BuyTrainTicket
language: python
filePath: /root/minik8s/testdata/workflow1/BuyTrainTicket/
triggerType:
  http: true
params:
  - name: x
    type: int
result:
  - name: x
    type: int
  - name: greeting
    type: string
```
~~~

函数的执行：

1. 通过http、event trigger触发，根据函数配置文件中的参数（数量、类型）定义，组装参数
2. 通过json格式传入容器8080端口，并触发函数运行
3. 结果直接发给`apiserver`

结果的返回：

1. apiserver收到容器返回的结果后，取出对应的Job并将结果存入其中，将Job标记为Ended，代表一次函数调用的结束

### Trigger实现

1. Trigger抽象：无论通过http触发还是事件触发器触发，apiserver都会统一生成一条Trigger消息，Trigger消息包含本次调用的uuid，以及被调用的函数和参数，这条消息通过kafka发给Serverless Controller
2. Job抽象：Serverless Controller接收到Trigger消息后，生成一个Job，Job沿用Trigger消息的uuid，函数名称和参数，但是除此之外，Job还记录了当前调用处于的状态（Created，Running和Ended），调用的结果以及该次调用在哪个Pod上执行
3. Job对象的生成：Serverless Controller接收到Trigger消息后，检查被调用函数对应的Serverless实例（Pod）：
   1. 若当前有空闲的Pod，Serverless Controller直接将该Pod分配给该Job
   2. 若没有空闲的Pod，Serverless Controller创建一个新的Pod分配给该Job，这个Pod会在Job被发回给apiserver的时候被真正创建
   3. Serverless Controller创建的Job处于Created状态
4. Serverless实例的创建：apiserver接收到从Serverless Controller创建的Job后，检查其中的Pod是否已被创建，若Pod仍未创建，表明这是一个新的Pod，按照Pod创建流程创建
5. Job对象的执行：Job Controller监听Job的变化，当一个Job被创建后，检查其中Pod是否已经被创建：
   1. 若Pod没被创建，这个Job会进入等待队列当中，直到Job Controller监听到对应Pod被创建的消息，拿到Pod的IP，才会给Pod IP发http请求，从而调用对应函数
   2. 若Pod已经被创建，直接给Pod发http请求调用函数
   3. 调用函数的http请求会立刻返回，此时Job Controller把Job标记为Running，并发往apiserver





### Workflow：

1. 每个节点的函数用namespace + name指定
2. 规则匹配采用switch-case模式，匹配满足express的第一个，否则执行default
3. expression需要参数名variable、判断关系opt（＞、=等）、值Value、数据类型type，支持多个expression的&&逻辑
4. workflow具体运行触发管理，由workflow_controller控制

workflow 配置文件示例：

```yaml
apiVersion: v1
kind: Workflow
metadata:
  name: my-workflow
  namespace: default
triggerType:
  http: true
graph:
  function:
    name: BuyTrainTicket
    namespace: default
  rule:
    case:
      - expression:
          - variable: status
            opt: "="
            value: "Succeeded"
            type: string
        successor:
          function:
            name: ReserveFlight
            namespace: default
          rule:
            case:
              - expression:
                  - variable: status
                    opt: "="
                    value: "Succeeded"
                    type: string
                successor:
                  function:
                    name: ReserveHotel
                    namespace: default
                  rule:
                    case:
                      - expression:
                          - variable: status
                            opt: "="
                            value: "Failed"
                            type: string
                        successor:
                          function:
                            name: CancelFlight
                            namespace: default
                          rule:
                            default:
                              function:
                                name: CancelTrainTicket
                                namespace: default
                              rule:
                                default:
                                  function:
                                    name: OrderFailed
                                    namespace: default
                    default:
                      function:
                        name: OrderSucceeded
                        namespace: default
            default:
              function:
                name: CancelTrainTicket
                namespace: default
              rule:
                default:
                  function:
                    name: OrderFailed
                    namespace: default
    default:
      function:
        name: OrderFailed
        namespace: default
```

