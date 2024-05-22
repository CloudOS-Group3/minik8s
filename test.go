package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    // defer wg.Done() // 在协程完成时调用 Done() 方法

    fmt.Printf("Worker %d starting\n", id)
    // 在这里执行协程的任务
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    // 启动三个协程
    for i := 1; i <= 3; i++ {
        wg.Add(1) // 告诉 WaitGroup 有一个协程需要等待
        go worker(i, &wg)
    }

    // 等待所有协程完成
    wg.Wait() // 这会阻塞主函数的执行，直到所有协程都完成

    fmt.Println("All workers done")
}
