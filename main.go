package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func main() {
	// 创建一个新的Colly爬虫
	c := colly.NewCollector()

	// 创建一个互斥锁，用于保护文件写入操作
	var mutex sync.Mutex

	// 创建文件
	file, err := os.Create("data.txt")
	if err != nil {
		log.Fatal("创建文件失败:", err)
	}
	defer file.Close()

	// 获取用户输入的目标网站
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入目标网站: ")
	targetSite, _ := reader.ReadString('\n')
	targetSite = strings.TrimSpace(targetSite)

	// 获取用户输入的目录字典
	fmt.Print("请输入目录字典（每行一个目录）: ")
	dictionaryFile, _ := reader.ReadString('\n')

	// 创建等待组
	var wg sync.WaitGroup

	// 设置回调函数，处理每个访问的网页
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasSuffix(link, "/") {
			fmt.Println("目录:", link)

			// 使用等待组增加计数
			wg.Add(1)

			// 使用goroutine并发处理文件写入操作
			go func(link string) {
				// 使用互斥锁保护文件写入操作
				mutex.Lock()
				defer mutex.Unlock()

				// 写入文件
				_, err := file.WriteString(fmt.Sprintf("目录: %s\n", link))
				if err != nil {
					log.Println("写入文件失败:", err)
				}

				// 完成goroutine，减少等待组计数
				wg.Done()
			}(link)
		} else {
			fmt.Println("文件:", link)

			// 使用等待组增加计数
			wg.Add(1)

			// 使用goroutine并发处理文件写入操作
			go func(link string) {
				// 使用互斥锁保护文件写入操作
				mutex.Lock()
				defer mutex.Unlock()

				// 写入文件
				_, err := file.WriteString(fmt.Sprintf("文件: %s\n", link))
				if err != nil {
					log.Println("写入文件失败:", err)
				}

				// 完成goroutine，减少等待组计数
				wg.Done()
