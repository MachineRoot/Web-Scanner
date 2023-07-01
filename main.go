package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	// 获取用户输入的网站
	var website string
	fmt.Print("请输入要扫描的网站: ")
	fmt.Scanln(&website)

	// 获取用户输入的目录字典文件路径
	var filePath string
	fmt.Print("请输入目录字典文件的路径: ")
	fmt.Scanln(&filePath)

	// 读取目录字典文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("无法读取目录字典文件: %s\n", err)
		return
	}

	// 将文件内容按行分割成目录列表
	directories := strings.Split(string(content), "\n")

	// 创建一个等待组，用于等待所有扫描任务完成
	var wg sync.WaitGroup

	// 创建一个缓冲通道，用于保存扫描结果
	resultChan := make(chan string, len(directories))

	// 遍历目录列表，为每个目录启动一个扫描任务
	for _, dir := range directories {
		// 去除目录前后的空白字符
		dir = strings.TrimSpace(dir)

		// 跳过空白行
		if dir == "" {
			continue
		}

		// 增加等待组的计数器
		wg.Add(1)

		// 启动一个Go协程来扫描目录
		go func(directory string) {
			defer wg.Done()

			// 构建完整的URL
			url := website + directory

			// 获取HTTP客户端
			client := &http.Client{}

			// 发起HTTP请求
			resp, err := client.Get(url)
			if err != nil {
				fmt.Printf("无法访问目录 %s: %s\n", directory, err)
				return
			}
			defer resp.Body.Close()

			// 读取响应内容
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("无法读取目录 %s 的响应: %s\n", directory, err)
				return
			}

			// 将扫描结果发送到通道
			resultChan <- fmt.Sprintf("目录 %s 的响应内容:\n%s\n", directory, body)
		}(dir)
	}

	// 等待所有扫描任务完成
	wg.Wait()

	// 关闭通道，确保所有结果都已经发送完毕
	close(resultChan)

	// 创建一个文件，用于保存扫描结果
	file, err := os.Create("data.txt")
	if err != nil {
		fmt.Printf("无法创建文件: %s\n", err)
		return
	}
	defer file.Close()

	// 从通道中读取扫描结果并写入文件
	for result := range resultChan {
		file.WriteString(result)
	}
}
