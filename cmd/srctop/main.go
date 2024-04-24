package main

import (
	"fmt"
	"strings"

	"github.com/baiqll/srctop/pkg/runner"
	"github.com/shenwei356/stable"
)

func main() {

	// 初始化参数，获取统计时间
	create_time, show_time := runner.CreateTime()

	fmt.Printf(" ASRC Top list: %s\n", show_time)

	// 获取统计结果
	statsMap := runner.TopDetail(create_time)

	tbl := stable.New()
	tbl.HumanizeNumbers()
	tbl.Style(stable.StyleThreeLine)
	tbl.Header([]string{
		"贡献值",
		"安全币",
		"昵称",
		"业务",
	})

	// 打印结果
	for _, stat := range statsMap {

		tbl.AddRow([]interface{}{stat.TotalJifen, stat.TotalCoin, stat.Nickname, strings.Join(stat.Tidnames, "、")})

	}

	fmt.Printf("%s\n", tbl.Render(stable.StyleThreeLine))

}
