package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	DAEMON  = "daemon"
	FOREVER = "forever"
)

var children = make([]*os.Process, 0)

func DoSomething() {
	//fp, _ := os.OpenFile("./dosomething.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//log.SetOutput(fp)
	for {
		log.Printf("biz worker running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
		time.Sleep(time.Second * 2)
	}
}

func StripSlice(slice []string, element string) []string {
	for i := 0; i < len(slice); {
		if slice[i] == element && i != len(slice)-1 {
			slice = append(slice[:i], slice[i+1:]...)
		} else if slice[i] == element && i == len(slice)-1 {
			slice = slice[:i]
		} else {
			i++
		}
	}
	return slice
}

func daemonWorker() {
	for {
		fmt.Println("listen daemon")
		time.Sleep(1 * time.Second)
	}
}

func SubProcess(args []string) *exec.Cmd {
	fmt.Printf("SubProcess args: %v\n", args)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func main() {
	daemon := flag.Bool(DAEMON, false, "run in daemon")
	forever := flag.Bool(FOREVER, false, "run forever")
	flag.Parse()
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("[*] PID: %d PPID: %d ARG: %s\n", os.Getpid(), os.Getppid(), os.Args)
	if *daemon {
		args := StripSlice(os.Args, "-"+DAEMON)
		args = StripSlice(args, "--"+DAEMON)
		fmt.Printf("Daemon running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
		go func() {
			daemonWorker()
		}()
		go func() {
			sig := <-c
			fmt.Println("daemon receive signal")
			for _, child := range children {
				child.Signal(sig)
			}
			os.Exit(0)
		}()
		cmd1 := SubProcess(args)
		cmd1.Start()
		cmd1.Wait()
		//fmt.Println("daemon exit....")
		//os.Exit(0)
	} else if *forever {
		go func() {
			<-c
			fmt.Println("forever receive signal")

			os.Exit(0)
		}()
		for {
			args := StripSlice(os.Args, "-"+FOREVER)
			args = StripSlice(args, "--"+FOREVER)
			cmd := SubProcess(args)
			cmd.Start()
			fmt.Printf("Forever 进程 running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
			cmd.Wait()
		}
		os.Exit(0)
	} else {

		fmt.Printf("biz 进程 running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
	}
	go func() {
		<-c
		fmt.Println("receive signal")

		os.Exit(0)
	}()
	DoSomething()
}
