package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kinwyb/go/generate"
	_ "github.com/kinwyb/go/generate/rpcx"
	"github.com/spf13/cobra"
)

//./generate rpcxclient -i /Users/heldiam/Developer/GO/mysrc/zhifangw.cn/src/CloudFactory/API/models/endPoints -o /Users/heldiam/Developer/GO/mysrc/zhifangw.cn/src/CloudFactory/API/application/rpcx/client

var cmdRpcxClient = &cobra.Command{
	Use:   "rpcxclient",
	Short: "rpcxclient 通过接口生成rpcx客户端代码",
	Long:  `rpcxclient 命令用于通过接口生成rpcx客户端代码`,
	Run: func(cmd *cobra.Command, args []string) {
		filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			filedata, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = generate.ParseFileByLayoutName(bytes.NewReader([]byte(filedata)), "rpcxclient", out)
			if err != nil {
				fmt.Println(err.Error())
			}
			return nil
		})
	},
}
