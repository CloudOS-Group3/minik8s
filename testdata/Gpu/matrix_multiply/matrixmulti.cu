#include <stdio.h>
#include <stdlib.h>
#include <cuda_runtime.h>

#define M 500 // 矩阵 A 的行数
#define K 500 // 矩阵 A 的列数和矩阵 B 的行数
#define N 500 // 矩阵 B 的列数

#define IDX2C(i,j,ld) (((j)*(ld))+(i))

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

int main(void) {
    cudaError_t cudaStat;
    float *devPtrA, *devPtrB, *devPtrC;
    float *a = 0, *b = 0, *c = 0;
    int i, j;

    // 在主机上分配内存
    a = (float *)malloc(M * K * sizeof(*a));
    b = (float *)malloc(K * N * sizeof(*b));
    c = (float *)malloc(M * N * sizeof(*c));

    if (!a || !b || !c) {
        printf("Host memory allocation failed\n");
        return EXIT_FAILURE;
    }

    // 初始化矩阵 A 和 B
    for (j = 0; j < K; j++) {
        for (i = 0; i < M; i++) {
            a[IDX2C(i, j, M)] = (float)(rand() % 10); // 用随机值初始化矩阵 A
        }
    }
    for (j = 0; j < N; j++) {
        for (i = 0; i < K; i++) {
            b[IDX2C(i, j, K)] = (float)(rand() % 10); // 用随机值初始化矩阵 B
        }
    }

    // 在设备上分配内存
    cudaStat = cudaMalloc((void**)&devPtrA, M * K * sizeof(*a));
    if (cudaStat != cudaSuccess) {
        printf("Device memory allocation for A failed\n");
        return EXIT_FAILURE;
    }
    cudaStat = cudaMalloc((void**)&devPtrB, K * N * sizeof(*b));
    if (cudaStat != cudaSuccess) {
        printf("Device memory allocation for B failed\n");
        cudaFree(devPtrA);
        return EXIT_FAILURE;
    }
    cudaStat = cudaMalloc((void**)&devPtrC, M * N * sizeof(*c));
    if (cudaStat != cudaSuccess) {
        printf("Device memory allocation for C failed\n");
        cudaFree(devPtrA);
        cudaFree(devPtrB);
        return EXIT_FAILURE;
    }

    // 将矩阵 A 和 B 复制到设备
    cudaStat = cudaMemcpy(devPtrA, a, M * K * sizeof(*a), cudaMemcpyHostToDevice);
    if (cudaStat != cudaSuccess) {
        printf("Data transfer for A failed\n");
        cudaFree(devPtrA);
        cudaFree(devPtrB);
        cudaFree(devPtrC);
        return EXIT_FAILURE;
    }
    cudaStat = cudaMemcpy(devPtrB, b, K * N * sizeof(*b), cudaMemcpyHostToDevice);
    if (cudaStat != cudaSuccess) {
        printf("Data transfer for B failed\n");
        cudaFree(devPtrA);
        cudaFree(devPtrB);
        cudaFree(devPtrC);
        return EXIT_FAILURE;
    }

    // 配置 CUDA 核函数的执行参数
    dim3 threadsPerBlock(16, 16);
    dim3 blocksPerGrid((N + threadsPerBlock.x - 1) / threadsPerBlock.x,
                       (M + threadsPerBlock.y - 1) / threadsPerBlock.y);

    // 调用矩阵乘法核函数
    matrixMulKernel<<<blocksPerGrid, threadsPerBlock>>>(devPtrA, devPtrB, devPtrC, M, K, N);
    cudaStat = cudaDeviceSynchronize();
    if (cudaStat != cudaSuccess) {
        printf("Kernel execution failed\n");
        cudaFree(devPtrA);
        cudaFree(devPtrB);
        cudaFree(devPtrC);
        return EXIT_FAILURE;
    }

    // 将结果矩阵 C 从设备复制回主机
    cudaStat = cudaMemcpy(c, devPtrC, M * N * sizeof(*c), cudaMemcpyDeviceToHost);
    if (cudaStat != cudaSuccess) {
        printf("Data transfer for C failed\n");
        cudaFree(devPtrA);
        cudaFree(devPtrB);
        cudaFree(devPtrC);
        return EXIT_FAILURE;
    }

    // 释放设备内存
    cudaFree(devPtrA);
    cudaFree(devPtrB);
    cudaFree(devPtrC);

    // 打印部分结果矩阵
    printf("Result matrix C (partial):\n");
    for (j = 0; j < (N < 10 ? N : 10); j++) {
        for (i = 0; i < (M < 10 ? M : 10); i++) {
            printf("%7.0f", c[IDX2C(i,j,M)]);
        }
        printf("\n");
    }

    // 释放主机内存
    free(a);
    free(b);
    free(c);

    return EXIT_SUCCESS;
}
