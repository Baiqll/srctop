package runner

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

func AlignLeft(str string, width int) string {

	str = ShortenString(str, 20)

	re := regexp.MustCompile(`[\p{Han}]+`)

	// 找出所有匹配的中文字符串
	matches := re.FindAllString(str, -1)

	if matches != nil {

		zh_char := utf8.RuneCountInString(matches[0])
		en_char := utf8.RuneCountInString(str) - zh_char

		real_width := int(math.Round(float64(zh_char) / 0.6111))

		return fmt.Sprintf("%s%s", str, strings.Repeat(" ", width-en_char-real_width))

	}

	return fmt.Sprintf("%s%s", str, strings.Repeat(" ", width-len(str)))

}

func ShortenString(s string, maxLen int) string {
	if utf8.RuneCountInString(s) > maxLen {
		if maxLen > 3 {
			// 留出 "..."" 的位置
			return s[:maxLen-3] + "..."
		} else {
			// 如果 maxLen 小于或等于3，只显示 "..."
			return "..."[:maxLen]
		}
	}
	return s
}

func CreateTime() (create_time string, show_time string) {

	args := os.Args[1:] // 获取除程序名之外的参数

	create_time = "null"

	now := time.Now()
	currentYear := now.Year()

	// 当前年的开始日期
	startFormatted := time.Date(currentYear, time.January, 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	// 当前年的结束日期
	endFormatted := time.Date(currentYear, time.December, 31, 23, 59, 59, 0, now.Location()).Format("2006-01-02")

	show_time = fmt.Sprintf("%s To %s", startFormatted, endFormatted)

	if len(args) < 1 {
		return
	}

	param := args[0]

	// 支持的时间布局
	layout := "2006/01"
	// endOfMonthLayout := "2006-01-02 15:04:05"

	switch {
	case strings.Contains(param, "-"): // 参数形式为 2024/03-2024/05

		// start_type := 1 // {1 : 月份，2: 年份}
		var is_month = true // {1 : 月份，2: 年份}

		dates := strings.Split(param, "-")

		if len(dates) != 2 {
			fmt.Println("时间范围格式错误。需要的格式如：2024/03-2024/05。")
			return
		}
		// 解析开始日期
		startTime, err := time.Parse(layout, dates[0])
		if err != nil {
			startTime, err = time.Parse("2006", param)
			if err != nil {
				fmt.Println("开始日期解析失败")
				return
			}
		}

		// 解析结束日期
		endTime, err := time.Parse(layout, dates[1])
		if err != nil {
			endTime, err = time.Parse("2006", param)
			if err != nil {
				fmt.Println("结束日期解析失败")
				return
			}
			is_month = false
		}

		// 转换为所需的格式并计算结束月的最后一天
		startFormatted := startTime.Format("2006-01-02") + " 00:00:00"
		endFormatted := endTime.AddDate(0, 1, -1).Format("2006-01-02") + " 23:59:59"
		show_endFormatted := endTime.AddDate(0, 1, -1).Format("2006-01-02")
		if !is_month {
			endFormatted = endTime.AddDate(1, 0, -1).Format("2006-01-02") + " 23:59:59"
			show_endFormatted = endTime.AddDate(1, 0, -1).Format("2006-01-02")
		}

		// 输出
		create_time = fmt.Sprintf("[\"%s\",\"%s\"]\n", startFormatted, endFormatted)
		show_time = fmt.Sprintf("%s To %s", startTime.Format("2006-01-02"), show_endFormatted)
		return

	case param == "": // 参数为空
		return
	default: // 参数形式为 2024/03
		var is_month = true
		date, err := time.Parse(layout, param)
		if err != nil {
			date, err = time.Parse("2006", param)
			if err != nil {
				fmt.Println("解析时间错误。")
				return
			}
			is_month = false
		}

		// 转换为所需的格式并计算结束月的最后一天
		startFormatted := date.Format("2006-01-02") + " 00:00:00"
		endFormatted := date.AddDate(0, 1, -1).Format("2006-01-02") + " 23:59:59"
		show_endFormatted := date.AddDate(0, 1, -1).Format("2006-01-02")
		if !is_month {
			endFormatted = date.AddDate(1, 0, -1).Format("2006-01-02") + " 23:59:59"
			show_endFormatted = date.AddDate(1, 0, -1).Format("2006-01-02")
		}

		// 输出
		create_time = fmt.Sprintf("[\"%s\",\"%s\"]\n", startFormatted, endFormatted)
		show_time = fmt.Sprintf("%s To %s", date.Format("2006-01-02"), show_endFormatted)
		return
	}

}
