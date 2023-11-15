package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const BIZ_PID_FILE = "/tmp/qskm-backend-biz.pid"

var done = make(chan int)

func Daemon() {

	//signal.Notify(signalChan, syscall.SIGUSR1)
	//go func() {
	//	sig := <-signalChan
	//	if sig == syscall.SIGUSR1 {
	//		fmt.Printf("i am child, sig=%d", sig)
	//		os.Exit(0)
	//	}
	//}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		sig := <-signalChan
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			fmt.Printf("接收到中断信号-----，结束子进程, sig=%d", sig)
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
	InitRedis()
	if daemon {
		fmt.Println("DAEMON start ....")
		val, _ := GetStatus()
		status := Status(val)
		if status != Init {
			return
		}
		wg := sync.WaitGroup{}
		//childen := make([]*os.Process, 0)
		fmt.Printf("[*] PID: %d PPID: %d ARG: %s\n", os.Getpid(), os.Getppid(), os.Args)
		cmd1 := exec.Command("/Users/zyb/person/multi_process/example_cobra/example_cobra", "ops")
		cmd2 := exec.Command("/Users/zyb/person/multi_process/example_cobra/example_cobra", "serve")
		// 启动子进程
		startChildProcess(cmd1, &wg)
		startChildProcess(cmd2, &wg)

		// 监听操作系统的中断信号，当接收到中断信号时，结束子进程
		//signalChan := make(chan os.Signal, 1)
		//signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		//go func() {
		//	sig := <-signalChan
		//	if sig == syscall.SIGINT || sig == syscall.SIGTERM {
		//		fmt.Printf("接收到中断信号-----，结束子进程, sig=%d", sig)
		//		stopChildProcess(cmd1)
		//		stopChildProcess(cmd2)
		//		//wg.Wait()
		//		fmt.Println("EXIT")
		//		os.Exit(0)
		//	} else if sig == syscall.SIGUSR1 {
		//		fmt.Printf("i am child, sig=%d", sig)
		//		os.Exit(0)
		//	}
		//}()
		//go WatchBIZProcess(cmd2.Args, &wg)
		wg.Wait()
		fmt.Println("子进程结束")
		time.Sleep(10 * time.Second)
		os.Exit(0)

	}
	fmt.Println("...........")
	// biz worker
	//signalChan := make(chan os.Signal, 1)
	//signal.Notify(signalChan, syscall.SIGUSR2)

	// 这快不起作用，收不到信号>>>>

	// 这快不起作用，收不到信号 <<<<<
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World! %v", time.Now())
	})

	err := http.ListenAndServe(":8098", nil)
	if err != nil {
		fmt.Println("启动HTTP服务器失败:", err)
	}
	fmt.Println("start ops server")
	fmt.Printf("[*] BIZ PID: %d PPID: %d ARG: %s\n", os.Getpid(), os.Getppid(), os.Args)

	for {
		fmt.Println("biz process")
		time.Sleep(1 * time.Second)
	}

}

// 启动子进程

func execCMD(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("启动子进程失败:", err)
		os.Exit(1)
	}
	fmt.Printf("启动子进程，PID: %d,PPID: %d, args=%v\n", cmd.Process.Pid, os.Getpid(), cmd.Args)
	// 等待子进程结束
	cmd.Wait()
	fmt.Printf("子进程退出，PID: %d, args=%v\n", cmd.Process.Pid, cmd.Args)
}

func startChildProcess(cmd *exec.Cmd, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		execCMD(cmd)
	}()
}

//func createPidFile(pidfile string) error {
//	f, err := os.Open(pidfile)
//}

func WatchBIZProcess(args []string, wg *sync.WaitGroup) {
	fmt.Println(args)
	cmd := exec.Command(args[0], args[1:]...)
	for {
		val, _ := GetStatus()
		status := Status(val)
		// 如果不是1，5，7 则忽略
		if !(status == Running || status == BizStart || status == Rollback) {
			time.Sleep(1 * time.Second)
			fmt.Printf("Sleep 1 second, PID: %d,\n", os.Getpid())
			continue

		}
		//
		fmt.Printf("status=%d\n", status)
		time.Sleep(5 * time.Second)
		cmd = exec.Command(args[0], args[1:]...)
		startChildProcess(cmd, wg)
		time.Sleep(5 * time.Second)

	}
}

// 结束子进程
func stopChildProcess(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	if err := cmd.Process.Signal(syscall.SIGUSR1); err != nil {
		fmt.Println("结束子进程失败:", err)
	}
}
