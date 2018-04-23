package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "generate"}
	cmdBeego.Flags().StringVarP(&in, "in", "i", "", "接口文件路径")
	cmdBeego.Flags().StringVarP(&out, "out", "o", "", "输出文件路径")
	rootCmd.AddCommand(cmdBeego)
	rootCmd.Execute()
}
