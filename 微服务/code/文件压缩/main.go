package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 把 src 路径（文件或目录）压缩成 dst 指定的 zip 文件
func zipit(src, dst string) error {
	// 创建目标文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	zw := zip.NewWriter(out)
	defer zw.Close()

	// 遍历 src
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 目录本身不写入 zip，只写入文件
		if info.IsDir() {
			return nil
		}

		// 构造 zip 内部相对路径
		rel, err := filepath.Rel(filepath.Dir(src), path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel) // Windows 下路径分隔符统一为 "/"

		// 创建 zip 内部文件头
		w, err := zw.Create(rel)
		if err != nil {
			return err
		}

		// 打开源文件
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// 拷贝内容
		_, err = io.Copy(w, file)
		return err
	})
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("用法: %s <要压缩的路径> <输出.zip>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	src, dst := os.Args[1], os.Args[2]
	// 如果用户忘记加 .zip 后缀，自动补上
	if !strings.HasSuffix(dst, ".zip") {
		dst += ".zip"
	}

	if err := zipit(src, dst); err != nil {
		fmt.Println("压缩失败:", err)
		os.Exit(1)
	}
	fmt.Println("已生成", dst)
}

func unzip(src, dstDir string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dstDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
