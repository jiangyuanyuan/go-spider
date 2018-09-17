package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

//https://www.pengfu.com/index_2.html
func main() {
	fmt.Print("hello,word")
	var start, end int
	fmt.Printf("请输入起始页>1")
	fmt.Scan(&start)
	fmt.Printf("请输入终止页>起始页")
	fmt.Scan(&end)

	Dowrok(start, end)

}
func Dowrok(start int, end int) {
	fmt.Printf("准备爬取第%d到%d页网页内容", start, end)

	for i := start; i <= end; i++ {
		SpiderPage(i)
	}

}
func SpiderPage(i int) {
	url := "https://www.pengfu.com/index_" + strconv.Itoa(i) + ".html"
	//url := "http://sz.zu.fang.com/house/i3" + strconv.Itoa(i) + "/"
	//url := "https://www.pengfu.com/index_2.html"
	fmt.Printf("正在爬取第%d个网页：%s\n", i, url)

	result, err := HttpGet(url)
	if err != nil {
		fmt.Printf("爬取出错")
		return
	}

	//已经拿到了网页源码  开始正则匹配
	//fmt.Printf(result)
	//<p class="title" id="rentid_D09_01_02">
	//                <a href="/chuzu/3_210286731_1.htm" data_channel="1,2" target="_blank" title="大学城地铁口 大型青年社区 品牌家电 大阳台独立厨卫">大学城地铁口 大型青年社区 品牌家电 大阳台独立厨卫</a>
	//
	//            </p>

	re := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))"`)
	if re == nil {
		fmt.Printf("正则匹配失败")
		return
	}
	links := re.FindAllStringSubmatch(result, -1)
	//fmt.Println("链接",links)
	for _, data := range links {
		fmt.Println("url = " + data[1])
		title, content, err := SpiderOne(data[1])
		if err != nil {
			fmt.Println("子网页出错", err)
			continue
		}
		fmt.Println("title = ", title)
		fmt.Println("content = ", content)

	}

}
func SpiderOne(url string) (title string, content string, err error) {
	result, err1 := HttpGet(url)
	if err1 != nil {
		err = err1
		return
	}
	//re :=regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))"`)
	re := regexp.MustCompile(`<h1>"(?s:(.*?))"</h1>`)
	if re == nil {
		//fmt.Printf("子网页正则  title 失败")
		err = fmt.Errorf("%s", "子网页正则  title 失败")
		return
	}
	//取标题
	tempTitle := re.FindAllStringSubmatch(result, 1)
	fmt.Println("tempTitle", tempTitle)
	for _, data := range tempTitle {
		title = data[1]
		break
	}
	//取内容
	//<a id="prev" href="https://www.pengfu.com/content_1849980_1.html"></a>

	//<img oldsrc="https://image7.pengfu.com/thumb/180917/5b9f638be87a8.jpg" width="200" height="239" src="https://image7.pengfu.com/origin/180917/5b9f638be87a8.gif">
	re1 := regexp.MustCompile(`src="(?s:(.*?))"<a id="prev" href=`)
	if re == nil {
		//fmt.Printf("子网页正则  gif 失败")
		err = fmt.Errorf("%s", "子网页正则  gif 失败")
		return
	}
	tempContent := re1.FindAllStringSubmatch(result, -1)
	for _, data := range tempContent {
		fmt.Println("data = ", data)
	}
	return

}
func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024*4)

	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}
	return

}
