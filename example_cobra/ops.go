package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Ops() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {

		sig := <-signalChan
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			fmt.Printf("OPS 接收到中断信号-----，结束子进程, sig=%d， pid=%d, \n", sig, os.Getpid())
			//stopChildProcess(cmd1)
			//stopChildProcess(cmd2)
			//wg.Wait()
			fmt.Println("EXIT")
			os.Exit(0)
		} else if sig == syscall.SIGUSR1 {
			fmt.Printf("i am child, sig=%d", sig)
			os.Exit(0)
		}
	}()
	fmt.Printf("[*] PID: %d PPID: %d ARG: %s\n", os.Getpid(), os.Getppid(), os.Args)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World! %v", time.Now())
	})

	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println("启动HTTP服务器失败:", err)
	}
	fmt.Println("start ops server")
}
