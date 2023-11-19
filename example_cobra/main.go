package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	daemon  bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start qskm backend service.",
	Long:  "start qskm backend service.",
	Run: func(cmd *cobra.Command, args []string) {
		daemon := Daemon{BizHealthPath: "http://127.0.0.1:8098/health"}
		daemon.start()
	},
}

var opsCmd = &cobra.Command{
	Use:   "ops",
	Short: "stop qskm backend service.",
	Long:  "stop qskm backend service.",
	Run: func(cmd *cobra.Command, args []string) {
		Ops()
	},
}

var bizCmd = &cobra.Command{
	Use:   "biz",
	Short: "stop qskm backend service.",
	Long:  "stop qskm backend service.",
	Run: func(cmd *cobra.Command, args []string) {
		Biz()
	},
}

var rootCmd = &cobra.Command{
	Use:   "qskm-backend",
	Short: "qskm-backend is backend service of QSKM.",
	Long:  `qskm-backend is backend service of QSKM.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("execute %s args:%v error:%v\n", cmd.Name(), args, errors.New("unrecognized command"))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	//serveCmd.Flags().BoolVar(&daemon, "daemon", false, "Run in daemon mode")
	rootCmd.AddCommand(opsCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file(default is config.yaml)")
	rootCmd.AddCommand(bizCmd)
}

func main() {
	Execute()
}
