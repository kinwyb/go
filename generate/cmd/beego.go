package main

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/kinwyb/go/generate"

	"io/ioutil"

	"fmt"

	_ "github.com/kinwyb/go/generate/beego"

	"github.com/spf13/cobra"
)

var in, out string

//./generate beego -i /Users/heldiam/Developer/GO/mysrc/zhifangw.cn/src/CloudFactory/API/models/endPoints -o /Users/heldiam/Developer/GO/mysrc/zhifangw.cn/src/CloudFactory/API/application/web/controllers

var cmdBeego = &cobra.Command{
	Use:   "beego",
	Short: "beego 通过接口生成http接口代码",
	Long:  `beego 命令用于通过接口生成http控制器代码`,
	Run: func(cmd *cobra.Command, args []string) {
		filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			filedata, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = generate.ParseFileByLayoutName(bytes.NewReader([]byte(filedata)), "beego", out)
			if err != nil {
				fmt.Println(err.Error())
			}
			return nil
		})
	},
}
