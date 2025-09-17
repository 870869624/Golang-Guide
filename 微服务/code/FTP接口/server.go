package main

import (
	"github.com/dutch/goftp"
	"log"
	"net"
	"os"
	"path/filepath"
)

func main() {
	// 设置FTP服务器的根目录和监听地址
	root := "."                               // 当前目录作为根目录
	listener, err := net.Listen("tcp", ":21") // 监听21端口（FTP标准端口）
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// 创建FTP服务器实例并启动服务
	server := goftp.NewServer(root) // 使用当前目录作为根目录创建服务器实例
	server.ListenAndServe(listener) // 启动服务并监听连接请求
}
