package main

import (
	"fmt"
	"strings"

	"github.com/baiqll/srctop/pkg/runner"
)

func main() {

	// 初始化参数，获取统计时间
	create_time, show_time := runner.CreateTime()

	fmt.Printf(" ASRC Top list: %s\n", show_time)
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%4s %8s %9s %-29s %s\n", "贡献值", "安全币", "", "昵称", "业务")
	fmt.Println(strings.Repeat("-", 100))

	// 获取统计结果
	statsMap := runner.TopDetail(create_time)

	// 打印结果
	for _, stat := range statsMap {
		fmt.Printf("%5d %10d %10s %s %s\n",
			stat.TotalJifen,
			stat.TotalCoin,
			"",
			runner.AlignLeft(stat.Nickname, 30),
			strings.Join(stat.Tidnames, "、"))
	}

}
