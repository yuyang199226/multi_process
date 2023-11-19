package main

import (
	"fmt"
	"log"
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
var cmd2 *exec.Cmd
var cacheKey = "qskm-backend:status"

type Daemon struct {
	BizHealthPath string
	OpsHealthPath string
}

func (d *Daemon) start() {
	InitRedis()

	//if daemon {
	fmt.Println("DAEMON start ....")
	val, _ := GetStatus(cacheKey)
	status := Status(val)
	if status != Init {
		log.Fatalf("status not equal 0, status=%d", status)
		return
	}
	wg := sync.WaitGroup{}
	//childen := make([]*os.Process, 0)
	fmt.Printf("[*] PID: %d PPID: %d ARG: %s\n", os.Getpid(), os.Getppid(), os.Args)
	pwd, _ := os.Getwd()
	fmt.Println("PWD: ", pwd)
	cmd1 := exec.Command(pwd+"/example_cobra", "ops")
	cmd2 = exec.Command(pwd+"/example_cobra", "biz")
	// 启动子进程
	startChildProcess(cmd1, &wg)
	startChildProcess(cmd2, &wg)

	// 监听操作系统的中断信号，当接收到中断信号时，结束子进程
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			fmt.Printf("接收到中断信号-----，结束子进程, sig=%d", sig)
			stopChildProcess(cmd1)
			stopChildProcess(cmd2)
			//wg.Wait()
			fmt.Println("EXIT")
			SetStatus(cacheKey, 0)
			os.Exit(0)
		} else if sig == syscall.SIGUSR1 {
			fmt.Printf("i am child, sig=%d", sig)
			os.Exit(0)
		}
	}()
	go d.WatchBIZProcess(cmd2.Args, &wg)
	wg.Wait()
	fmt.Println("子进程结束")
	time.Sleep(10 * time.Second)
	os.Exit(0)

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
}

func startChildProcess(cmd *exec.Cmd, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		execCMD(cmd)
		cmd.Wait()
		fmt.Printf("子进程退出，PID: %d, args=%v\n", cmd.Process.Pid, cmd.Args)
	}()
}

//func createPidFile(pidfile string) error {
//	f, err := os.Open(pidfile)
//}

func (d *Daemon) WatchBIZProcess(args []string, wg *sync.WaitGroup) {
	fmt.Println(args)
	time.Sleep(3 * time.Second)

	for {
		val, _ := GetStatus(cacheKey)
		status := Status(val)
		// 如果不是1，5，7 则忽略
		if !(status == Running || status == BizStart || status == Rollback) {
			time.Sleep(1 * time.Second)
			fmt.Printf("Sleep 1 second, PID: %d,\n", os.Getpid())
			continue

		}
		if d.healthCheck() {
			fmt.Printf("heath check\n")
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("status=%d\n", status)
		time.Sleep(5 * time.Second)
		cmd2 = exec.Command(args[0], args[1:]...)
		startChildProcess(cmd2, wg)
		time.Sleep(5 * time.Second)

	}
}

// 结束子进程
func stopChildProcess(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		fmt.Println("结束子进程失败:", err)
	}
}

func (d *Daemon) healthCheck() bool {
	resp, err := http.Get(d.BizHealthPath)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
