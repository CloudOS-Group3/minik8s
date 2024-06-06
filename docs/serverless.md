# Severless实现

### 1. Function抽象

**使用方式**：

- 用户可以通过单个python文件定义function，python文件名无要求，但是主函数名必须为`main`；
- 用户需要指定参数和返回值的数量、类型，也可指定名字（用于workflow之间传递指定参数）；
- 用户可以指定函数invoke方式：http或event trigger；
- 用户需在对应文件夹下提供requirement.txt

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



实现方式：

#### 1.1 创建function

通过 `apply -f <filepath>` 指令，通过yaml文件创建function：

1. 存入etcd；

2. build函数镜像

   1. 将所有需要的文件复制到工作路径，安装运行环境
   2. 在容器中写入一个server.py，在容器运行时监听8080端口，并且执行用户函数，将结果发送至集群的apiserver

3. 将镜像push到minik8s所维护的镜像仓库

   > 该镜像仓库是运行在master  IP 下的 5050 端口上的一个 Docker registry









1. ⽀持Function抽象。⽤⼾可以通过单个⽂件（zip包或代码⽂件）定义函数内容，通过指令上传给
minik8s，并且通过http trigger和event trigger调⽤函数。 
◦ 函数需要⾄少⽀持Python语⾔ 
◦ 函数的格式，return的格式，update、invoke指令的格式可以⾃定义 
◦ 函数调⽤的⽅式：⽀持⽤⼾通过http请求、和绑定事件触发两种⽅式调⽤函数。函数可以通过
指令指定要绑定的事件源，事件的类型可以是时间计划、⽂件修改或其他⾃定义内容
2. ⽀持Serverless Workflow抽象：⽤⼾可以定义Serverless DAG，包括以下⼏个组成成分： 
◦ 函数调⽤链：在调⽤函数时传参给第⼀个函数，之后依次调⽤各个函数，前⼀个函数的输出作
为后⼀个函数的输⼊，最终输出结果。
◦ 分⽀：根据上⼀个函数的输出，控制⾯决定接下来运⾏哪⼀个分⽀的函数。
◦ Serverless Workflow可以通过配置⽂件来定义，参考AWS StepFunction或Knative的做法。除
此之外，同学们也可以⾃⾏定义编程模型来构建Serverless Workflow，只要workflow能达到
上述要求即可。
3. Serverless的⾃动扩容（Scale-to-0） 
◦ Serverless的实例应当在函数请求⾸次到来时被创建（冷启动），并且在⻓时间没有函数请求再
次到来时被删除（scale-to-0）。同时，Serverless能够监控请求数变化，当请求数量增多时能
够⾃动扩容⾄>1实例。 
◦ Serverless应当能够正确处理⼤量的并发请求(数⼗并发)，演⽰时使⽤wrk2或者Jmeter进⾏压
⼒测试。
4. 找⼀个较为复杂的开源的Serverless应⽤或者⾃⾏实现⼀个较为复杂的Serverless应⽤，该应⽤必
须有现实的应⽤场景(⽐如下图所⽰的 Image Processing)，不能只是简单的加减乘除或者hello 
world。将该应⽤部署在minik8s上，基于这个应⽤来展⽰Serverless相关功能，并在验收报告中结
合该应⽤的需求和特点详细分析该应⽤使⽤Serverless架构的必要性和优势。 Serverless样例可参考Serverlessbench中的Alexa：ServerlessBench/Testcase4-Applicationbreakdown/alexa at master · SJTU-IPADS/ServerlessBench