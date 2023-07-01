# Web-Scanner
go语言开发的web目录扫描，使用goroutine来并发处理每个链接的访问和文件读写实现更快扫描速度与处理速度。

## 使用
1.go编译环境安装"github.com/gocolly/colly"
2.调整file.txt内的目录
3.go run main.go

## 注意事项
Parallelism: 2 //设置线程，默认为2
本项目只用作安全测试，在得到授权后对网站进行目录爬取，请勿用于非法攻击。
