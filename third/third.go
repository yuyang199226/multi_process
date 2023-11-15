package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	// 创建两个子进程
	cmd1 := exec.Command("ls")
	cmd2 := exec.Command("echo", "Hello World")

	// 创建一个子进程列表
	children := []*os.Process{cmd1.Process, cmd2.Process}

	// 启动子进程
	if err := cmd1.Start(); err != nil {
		fmt.Println("启动子进程1失败:", err)
		os.Exit(1)
	}

	if err := cmd2.Start(); err != nil {
		fmt.Println("启动子进程2失败:", err)
		os.Exit(1)
	}

	// 监听操作系统的中断信号，当接收到中断信号时，结束子进程
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// 等待接收中断信号
		sig := <-signalChan
		fmt.Println("接收到中断信号:", sig)

		// 将信号传递给子进程
		for _, child := range children {
			child.Signal(sig)
		}

		// 等待子进程结束
		cmd1.Wait()
		cmd2.Wait()

		fmt.Println("子进程结束")
	}()

	// 等待父进程结束
	select {}
}
