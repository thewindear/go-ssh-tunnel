package main

import (
	"flag"
	"fmt"
	"go-ssh-tunnel/tunnel"
)

func main() {
	filepath := flag.String("config", "", "ssh隧道配置文件路径")
	flag.Parse()
	if *filepath == "" {
		fmt.Println("请使用 --config 指定ssh隧道配置文件路径")
		return
	}
	f, err := tunnel.BuildWithConfigFile(*filepath)
	if err != nil {
		fmt.Printf("初始化ssh隧道配置文件失败: %s", err)
		return
	}
	f.Run()
}
