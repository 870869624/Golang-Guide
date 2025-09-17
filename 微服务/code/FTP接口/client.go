package main

import (
	"fmt"
	"log"

	"github.com/jlaffaye/ftp"
)

func main() {
	// 连接到FTP服务器
	c, err := ftp.Dial("ftp.example.com:21", ftp.DialUser("username", "password"))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	// 列出目录内容
	entries, err := c.List("/")
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		fmt.Println(entry.Name)
	}
}
