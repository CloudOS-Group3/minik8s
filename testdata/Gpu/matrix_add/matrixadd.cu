#include <stdio.h>
#include <stdlib.h>
#include <cuda_runtime.h>
#include "cublas_v2.h"

#define M 1000 // 矩阵的行数
#define N 1000 // 矩阵的列数

#define IDX2C(i,j,ld) (((j)*(ld))+(i)) // 用于将二维索引转换为一维索引

// 矩阵加法核函数
__global__ void matrixAddKernel(float *A, float *B, float *C, int rows, int cols) {
    int idx = blockDim.x * blockIdx.x + threadIdx.x; // 程块中的线程数量 * 程序块索引 + 块内线程索引
    int idy = blockDim.y * blockIdx.y + threadIdx.y;

    if (idx < rows && idy < cols) {
        int index = IDX2C(idx, idy, rows);
        C[index] = A[index] + B[index]; // 计算矩阵 C 的元素
    }
}

int main(void) {
    cudaError_t cudaStat;
    cublasStatus_t stat;
    cublasHandle_t handle;
    float *devPtrA, *devPtrB, *devPtrC;
    float *a = 0, *b = 0, *c = 0;
    int i, j;

    // 在主机上分配内存
    a = (float *)malloc(M * N * sizeof(*a));
    b = (float *)malloc(M * N * sizeof(*b));
    c = (float *)malloc(M * N * sizeof(*c));

    if (!a || !b || !c) {
        printf("Host memory allocation failed\n");
        return EXIT_FAILURE;
    }

    // 初始化矩阵 A 和 B
    for (j = 0; j < N; j++) {
        for (i = 0; i < M; i++) {
            a[IDX2C(i,j,M)] = (float)(i * M + j + 1);
            b[IDX2C(i,j,M)] = (float)(i * M + j + 1) * 2;
        }
    }

    // 在设备上分配内存
    cudaStat = cudaMalloc((void**)&devPtrA, M * N * sizeof(*a));
    if (cudaStat != cudaSuccess) {
        printf("Device memory allocation for A failed\n");
        return EXIT_FAILURE;
    }
    cudaStat = cudaMalloc((void**)&devPtrB, M * N * sizeof(*b));
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

    // 创建 cuBLAS 句柄
    stat = cublasCreate(&handle);
    if (stat != CUBLAS_STATUS_SUCCESS) {
        printf("CUBLAS initialization failed\n");
        cudaFree(devPtrA);
        cudaFree(devPtrB);
        cudaFree(devPtrC);
        return EXIT_FAILURE;
    }

    // 将矩阵 A 和 B 复制到设备
    stat = cublasSetMatrix(M, N, sizeof(*a), a, M, devPtrA, M);
    if (stat != CUBLAS_STATUS_SUCCESS) {
        printf("Data download for A failed\n");
        cudaFree(devPtrA);
        cudaFree(devPtrB);
        cudaFree(devPtrC);
        cublasDestroy(handle);
        return EXIT_FAILURE;
    }
    stat = cublasSetMatrix(M, N, sizeof(*b), b, M, devPtrB, M);
    if (stat != CUBLAS_STATUS_SUCCESS) {
        printf("Data download for B failed\n");
        cudaFree(devPtrA);
        cudaFree(devPtrB);
        cudaFree(devPtrC);
        cublasDestroy(handle);
        return EXIT_FAILURE;
    }

    // 配置 CUDA 核函数的执行参数
    dim3 threadsPerBlock(16, 16); // 16x16的线程块, 共256个线程
    dim3 blocksPerGrid((M + threadsPerBlock.x - 1) / threadsPerBlock.x,
                       (N + threadsPerBlock.y - 1) / threadsPerBlock.y);

    // 调用矩阵加法核函数
    matrixAddKernel<<<blocksPerGrid, threadsPerBlock>>>(devPtrA, devPtrB, devPtrC, M, N);
    cudaStat = cudaDeviceSynchronize();
    if (cudaStat != cudaSuccess) {
        printf("Kernel execution failed\n");
        cudaFree(devPtrA);
        cudaFree(devPtrB);
        cudaFree(devPtrC);
        cublasDestroy(handle);
        return EXIT_FAILURE;
    }

    // 将结果矩阵 C 从设备复制回主机
    stat = cublasGetMatrix(M, N, sizeof(*c), devPtrC, M, c, M);
    if (stat != CUBLAS_STATUS_SUCCESS) {
        printf("Data upload for C failed\n");
        cudaFree(devPtrA);
        cudaFree(devPtrB);
        cudaFree(devPtrC);
        cublasDestroy(handle);
        return EXIT_FAILURE;
    }

    // 释放设备内存和 cuBLAS 句柄
    cudaFree(devPtrA);
    cudaFree(devPtrB);
    cudaFree(devPtrC);
    cublasDestroy(handle);

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
