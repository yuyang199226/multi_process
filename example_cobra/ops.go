package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func Ops() {
	fmt.Printf("[*] PID: %d PPID: %d ARG: %s\n", os.Getpid(), os.Getppid(), os.Args)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World! %v", time.Now())
	})

	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println("启动HTTP服务器失败:", err)
	}
	fmt.Println("start ops server")
}
