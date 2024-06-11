# GPU

### GPU任务提交流程：

1. **`kubectl` 通过配置文件创建GPU-job**

   配置文件主要参数有：`job-name`, `ntasks-per-node` 等超算平台需要的参数，`sourcePath` 代码和编译文件存放路径。

2. **`apiServer` 存入GPU-job，并发送创建pod消息**

3. **`kubelet` 创建pod实例**

   创建容器的时候，将配置文件中的**参数**，通过**环境变量**形式传入。

   将所需要的cuda程序代码和编译文件，通过**volumn映射**传入容器，

   创建container所使用的镜像中包含**`gpu_server.py`**，主要作用有：

   - 通过`os.getenv` 获取相关参数（环境变量全大写、下划线）
   - 通过参数编写作业脚本`.slurm` 
   - 连接超算平台，传送对应文件
   - 加载运行环境，编译程序，提交作业
   - 轮询查看是否完成（查看指定路径是否有结果文件输出）
   - 完成后给apiServer发消息（job-name和status），结果文件通过volumn映射传回主机

4. apiServer收到job完成消息，修改job状态



### CUDA程序说明

示例程序位于`minik8s/testdata/Gpu`

程序主要有以下几个部分：

1. CUDA 核函数

   利用当前线程的行和列索引，编写运算逻辑

   ```c
   __global__ void matrixMulKernel(float *A, float *B, float *C, int m, int k, int n) {
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
   ```

2. 内存分配和初始化

3. 配置 CUDA 核函数的执行参数

   这部分是使用GPU并发能力的关键，配置线程块、计算网格大小

   ```c
   // 配置 CUDA 核函数的执行参数
   dim3 threadsPerBlock(16, 16);
   dim3 blocksPerGrid((N + threadsPerBlock.x - 1) / threadsPerBlock.x,
                      (M + threadsPerBlock.y - 1) / threadsPerBlock.y);
   
   // 调用矩阵乘法核函数
   matrixMulKernel<<<blocksPerGrid, threadsPerBlock>>>(devPtrA, devPtrB, devPtrC, M, K, N);
   ```

4. 结果传回主机和释放内存