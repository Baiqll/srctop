package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	lib "github.com/baiqll/srctop/pkg"
)

type (
	Tenant struct {
		Name string `json:"name"`
		Tid  int    `json:"tid"`
	}

	TopResult struct {
		Jifen    int    `json:"jifen"`
		Coin     int    `json:"coin"`
		Nickname string `json:"nickname"`
		Team     string `json:"team"`
		Tenant   string `json:"tenant"`
	}
	// 累计数据
	Stats struct {
		TotalJifen int
		TotalCoin  int
		Nickname   string
		Tidnames   []string
	}
)

// httpRequest 函数发起一个 HTTP GET 请求并打印响应体
func AsrcRequest(url string) (body []byte, err error) {

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close() // 确保在函数退出时关闭响应体

	body, err = ioutil.ReadAll(resp.Body)

	return
}

func TopRequest(data *bytes.Buffer, tenant_name string, new_tops chan []TopResult, wg *sync.WaitGroup) {

	defer wg.Done() // 在goroutine完成时，调用Done来通知WaitGroup

	var tops []TopResult

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", "https://asrctenant.security.alibaba.com/profile/", data)
	if err != nil {
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	// 读取响应体
	res_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(tenant_name, "top 英雄榜获取失败", err)
		return
	}

	json.Unmarshal(res_data, &tops)

	for i, _ := range tops {
		tops[i].Tenant = tenant_name
	}

	new_tops <- tops

}

func Total(tops []TopResult) (sorts []*Stats) {

	// 使用映射来追踪统计
	statsMap := make(map[string]*Stats)

	// 遍历玩家，累加积分和金币，拼接tidname
	for _, top := range tops {
		stat, exists := statsMap[top.Nickname]
		if !exists {
			statsMap[top.Nickname] = &Stats{}
			stat = statsMap[top.Nickname]
			stat.Nickname = top.Nickname
		}
		stat.TotalJifen += top.Jifen
		stat.TotalCoin += top.Coin
		if top.Tenant != "" {
			stat.Tidnames = append(stat.Tidnames, top.Tenant)
		}
	}

	sorts = Sort(statsMap)

	return

}

func Sort(statsMap map[string]*Stats) (stats []*Stats) {

	// 打印结果
	for _, stat := range statsMap {

		stats = append(stats, stat)
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalCoin > stats[j].TotalCoin
	})

	return

}

func TopDetail(create_time string) (statsMap []*Stats) {

	var tenants []Tenant
	var all_tops []TopResult

	new_tops := make(chan []TopResult) // 创建缓冲通道

	// 获取业务列表
	res_data, err := AsrcRequest("https://asrctenant.security.alibaba.com/tenant/list")
	if err != nil {
		fmt.Println("ASRC 业务列表获取失败", err)
		return
	}

	json.Unmarshal(res_data, &tenants)

	var wg sync.WaitGroup
	wg.Add(len(tenants)) // 初始化等待组计数器

	for _, item := range tenants {

		data := bytes.NewBufferString(fmt.Sprintf(`{"tid":%d,"create_time":%s}`, item.Tid, create_time))

		go TopRequest(data, item.Name, new_tops, &wg)

	}

	// 让主go程等待其他go程
	go func() {
		wg.Wait()
		close(new_tops)
	}()

	for res := range new_tops {
		all_tops = append(all_tops, res...)
	}

	statsMap = Total(all_tops)

	return

}

func main() {

	// 初始化参数，获取统计时间
	create_time, show_time := lib.CreateTime()

	fmt.Printf(" ASRC Top list: %s\n", show_time)
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%4s %8s %9s %-29s %s\n", "贡献值", "安全币", "", "昵称", "业务")
	fmt.Println(strings.Repeat("-", 100))

	// 获取统计结果
	statsMap := TopDetail(create_time)

	// 打印结果
	for _, stat := range statsMap {
		fmt.Printf("%5d %10d %10s %s %s\n",
			stat.TotalJifen,
			stat.TotalCoin,
			"",
			lib.AlignLeft(stat.Nickname, 30),
			strings.Join(stat.Tidnames, "、"))
	}

}
