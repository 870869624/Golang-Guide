package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	// 心跳超时（服务端判断）
	heartbeatTimeout = 5 * time.Second
	// 客户端心跳间隔
	heartbeatInterval = 1 * time.Second
)

// ---------- 服务端 ----------
func server(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	log.Printf("heartbeat server listening on %s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(heartbeatTimeout)) // 首次设置

	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Printf("client %v gone: %v", conn.RemoteAddr(), err)
			return
		}
		// 简单校验
		if msg != "ping\n" {
			log.Printf("invalid heartbeat: %q", msg)
			return
		}
		// 回 pong
		if _, err := conn.Write([]byte("pong\n")); err != nil {
			return
		}
		// 重置超时
		_ = conn.SetDeadline(time.Now().Add(heartbeatTimeout))
		log.Printf("recv heartbeat from %v", conn.RemoteAddr())
	}
}

// ---------- 客户端 ----------
func client(addr string) {
	dialer := &net.Dialer{Timeout: 2 * time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	defer conn.Close()
	log.Printf("connected to %s", addr)

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	r := bufio.NewReader(conn)
	for {
		select {
		case <-ticker.C:
			if _, err := conn.Write([]byte("ping\n")); err != nil {
				log.Printf("write error: %v", err)
				return
			}
			// 等待 pong
			if _, err := r.ReadString('\n'); err != nil {
				log.Printf("read error: %v", err)
				return
			}
			log.Println("heartbeat ok")
		}
	}
}

// ---------- main ----------
func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run main.go <server|client>")
		return
	}

	const addr = "127.0.0.1:6000"

	switch os.Args[1] {
	case "server":
		server(addr)
	case "client":
		client(addr)
	default:
		fmt.Println("unknown command")
	}
}
