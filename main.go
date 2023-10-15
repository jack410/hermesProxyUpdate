package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type progressBarWriter struct {
	bar *pb.ProgressBar
}

func (pw *progressBarWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	pw.bar.Add(n)
	return
}

func main() {
	url := "https://img.kookapp.cn/attachments/2023-10/15/652bac0f68008.zip"
	outputFile := "hermesproxy-v3.8.zip"

	logFileName := time.Now().Format("2006-01-02-15-04") + ".log"
	logFile, err := os.Create(logFileName)
	if err != nil {
		fmt.Println("创建日志文件失败:", err)
		exit()
		return
	}
	defer logFile.Close()

	purple := "\033[0;35m" // ANSI 转义序列，表示紫色
	reset := "\033[0m"     // 重置颜色

	fmt.Println(purple + "该程序为wow1.142客户端的hermesproxy更新程序。" + reset)
	fmt.Println(purple + "请确保将该程序拷贝至游戏安装目录的hermes-proxy目录下再运行。" + reset)
	fmt.Println(purple + "如果有问题去kook everlook亚服频道联系管理员。" + reset)

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前目录失败:", err)
		logFile.WriteString(fmt.Sprintf("获取当前目录失败: %s\n", err))
		exit()
		return
	}

	baseDir := filepath.Base(currentDir)
	parentDir := filepath.Dir(currentDir)

	if baseDir != "hermes-proxy" {
		fmt.Println("程序运行路径有误，请将此程序保存至正确路径再运行。")
		logFile.WriteString("程序运行路径有误，请将此程序保存至正确路径再运行。\n")
		exit()
		return
	}

	launcherPath := filepath.Join(parentDir, "WinterspringLauncher.exe")
	if _, err := os.Stat(launcherPath); os.IsNotExist(err) {
		fmt.Println("程序运行路径有误，请将此程序保存至正确路径再运行。")
		logFile.WriteString("程序运行路径有误，请将此程序保存至正确路径再运行。\n")
		exit()
		return
	}

	fmt.Println("开始下载压缩包...")
	logFile.WriteString("开始下载压缩包...\n")

	if err := downloadFile(url, outputFile, logFile); err != nil {
		fmt.Println("下载压缩包失败:", err)
		logFile.WriteString(fmt.Sprintf("下载压缩包失败: %s\n", err))
		exit()
	}

	absPath, err := filepath.Abs(outputFile) // 获取输出文件的完整路径
	if err != nil {
		fmt.Println("获取输出文件路径失败:", err)
		logFile.WriteString(fmt.Sprintf("获取输出文件路径失败: %s\n", err))
		exit()
	}

	fmt.Println("下载完成！")
	fmt.Println("压缩包下载保存路径:", absPath)
	logFile.WriteString("下载完成！\n")
	logFile.WriteString(fmt.Sprintf("压缩包下载保存路径: %s\n", absPath))

	fmt.Println("开始解压缩文件...")
	logFile.WriteString("开始解压缩文件...\n")

	if err := unzip(outputFile, ".", logFile); err != nil {
		fmt.Println("解压缩文件失败:", err)
		logFile.WriteString(fmt.Sprintf("解压缩文件失败: %s\n", err))
		exit()
	}

	fmt.Println("升级完成！")
	logFile.WriteString("升级完成！\n")

	fmt.Println("日志已保存到:", logFileName)
	exit()
}

func downloadFile(url string, outputFile string, logFile *os.File) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	fileSize := resp.ContentLength
	bar := pb.Full.Start64(fileSize)
	bar.Set(pb.Bytes, true)

	writer := io.MultiWriter(out, &progressBarWriter{bar: bar})

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}

	bar.Finish()

	return nil
}

func unzip(src, dest string, logFile *os.File) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	fileCount := len(r.File)
	bar := pb.StartNew(fileCount)

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		logFile.WriteString(fmt.Sprintf("解压: %s\n", path))

		fileDir := filepath.Dir(path)
		err := os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			return err
		}

		fileWriter, err := os.Create(path)
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(fileWriter, rc)
		if err != nil {
			return err
		}

		fileWriter.Close()
		rc.Close()
		bar.Increment()

		absPath, err := filepath.Abs(path) // 获取解压文件的完整路径
		if err != nil {
			logFile.WriteString(fmt.Sprintf("获取解压文件路径失败: %s\n", err))
		} else {
			logFile.WriteString(fmt.Sprintf("解压完成: %s\n", absPath))
		}
	}

	bar.Finish()

	return nil
}

func exit() {
	fmt.Println("按下任意键退出...")
	fmt.Scanln()
}
