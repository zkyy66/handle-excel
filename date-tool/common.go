package date_tool

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FileInfo struct {
	Date     string
	Filename string
	FullName string
	Ext      string
}

// 判断文件名是否已日期开头
func HasDatePrefix(filename string) bool {
	baseName := filepath.Base(filename)

	// 定义日期模式的正则表达式
	// 支持: 8.1, 8.15, 12.25 等格式
	datePattern := `^\d{1,2}\.\d{1,2}`

	// 编译正则表达式
	re := regexp.MustCompile(datePattern)

	// 检查是否匹配
	return re.MatchString(baseName)
}

// removeExtension 移除文件扩展名
func RemoveExtension(filename string) string {
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		return filename[:idx]
	}
	return filename
}

// extractYearsFromPaths 批量从路径中提取年份
func ExtractYearsFromPaths(paths string) string {
	year, _ := ExtractYearSmart(paths)
	return year
}
func ExtractFileInfo(input string) (dateName, fileName string) {
	info := FileInfo{
		FullName: input,
	}

	// 提取扩展名
	if idx := strings.LastIndex(input, "."); idx != -1 {
		info.Ext = input[idx:]
	}

	// 使用正则表达式匹配日期（支持多种格式）
	datePatterns := []string{
		`(\d{1,2}\.\d{1,2})`,
		`(\d{1,2}-\d{1,2})`,
		`(\d{4}\.\d{1,2}\.\d{1,2})`,
	}

	var date string
	for _, pattern := range datePatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(input); len(matches) > 0 {
			date = matches[0]
			break
		}
	}

	info.Date = date

	// 提取文件名
	if date != "" {
		// 移除日期部分
		filenamePart := strings.Replace(input, date, "", 1)
		// 清理空白字符
		filenamePart = strings.TrimSpace(filenamePart)
		// 移除扩展名
		if info.Ext != "" {
			filenamePart = strings.TrimSuffix(filenamePart, info.Ext)
		}
		info.Filename = filenamePart
	} else {
		// 没有日期的情况
		if info.Ext != "" {
			info.Filename = strings.TrimSuffix(input, info.Ext)
		} else {
			info.Filename = input
		}
	}

	return info.Date, info.Filename
}

// 从字符串中提取日期和文件名
func ExtractDateAndFilename(input string) (dateName, fileName string) {
	//分割字符串并清理
	parts := strings.Fields(input) // 按空白字符分割
	if len(parts) >= 2 {
		// 第一个部分通常是日期
		dateName = parts[0]
		// 剩余部分组合为文件名
		fileName = strings.Join(parts[1:], " ")
	} else if len(parts) == 1 {
		// 如果只有一个部分，可能是没有空格分隔的情况
		fileName = parts[0]
	}
	// 从文件名中移除扩展名
	fileName = RemoveExtension(fileName)

	return dateName, fileName

}

// extractYearSmart 智能提取年份
func ExtractYearSmart(path string) (string, error) {
	// 尝试从完整路径中提取
	re := regexp.MustCompile(`(\d{2,4})年`)
	matches := re.FindStringSubmatch(path)
	if len(matches) < 2 {
		return "", fmt.Errorf("未找到年份")
	}

	yearStr := matches[1]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return "", fmt.Errorf("年份格式错误")
	}

	// 智能处理2位数年份
	if year < 100 {
		currentYear := time.Now().Year() // 可以替换为 time.Now().Year()
		century := currentYear / 100 * 100

		// 如果2位数年份大于当前年份的最后两位，认为是上世纪
		if year > currentYear%100 {
			year = century - 100 + year
		} else {
			year = century + year
		}
	}

	// 验证年份合理性
	if year < 1900 || year > 2100 {
		return "", fmt.Errorf("年份超出合理范围")
	}
	return strconv.Itoa(year), nil
}

//func ParseSpecialDate(dateStr string, year int) (time.Time, error) {
//	if year == 0 {
//		year = time.Now().Year() // 默认年份
//	}
//
//	parts := strings.Split(dateStr, ".")
//	if len(parts) < 2 {
//		return time.Time{}, fmt.Errorf("无效的日期格式: %s", dateStr)
//	}
//
//	monthStr := parts[0]
//	month, err := strconv.Atoi(monthStr)
//	if err != nil {
//		return time.Time{}, fmt.Errorf("月份格式错误: %s", monthStr)
//	}
//
//	dayStr := parts[1]
//	var dayDigits strings.Builder
//	for _, ch := range dayStr {
//		if ch >= '0' && ch <= '9' {
//			dayDigits.WriteRune(ch)
//		} else {
//			break
//		}
//	}
//	day, err := strconv.Atoi(dayDigits.String())
//	if err != nil {
//		return time.Time{}, fmt.Errorf("日期格式错误: %s", dayStr)
//	}
//
//	if month < 1 || month > 12 || day < 1 || day > 31 {
//		return time.Time{}, fmt.Errorf("无效的日期: %d.%d", month, day)
//	}
//
//	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
//}

// parseSpecialForDate 解析特殊日期格式
func ParseSpecialForDate(dateStr, yearName string) (time.Time, error) {
	// 提取数字部分（去除后面的非数字字符） 使用正则表达式或者字符串分割来提取日期部分

	//按点分割，取前两部分
	parts := strings.Split(dateStr, ".")
	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("无效的日期格式: %s", dateStr)
	}

	// 提取月份（第一个部分）
	monthStr := parts[0]
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("月份格式错误: %s", monthStr)
	}

	// 提取日期（第二个部分，只取数字部分）
	dayStr := parts[1]
	// 从日期字符串中提取连续的数字
	var dayDigits strings.Builder
	for _, ch := range dayStr {
		if ch >= '0' && ch <= '9' {
			dayDigits.WriteRune(ch)
		} else {
			break // 遇到非数字字符就停止
		}
	}

	day, err := strconv.Atoi(dayDigits.String())
	if err != nil {
		return time.Time{}, fmt.Errorf("日期格式错误: %s", dayStr)
	}

	// 使用固定年份（可根据需要调整）
	year, _ := strconv.Atoi(yearName)

	// 验证日期是否有效
	if month < 1 || month > 12 || day < 1 || day > 31 {
		return time.Time{}, fmt.Errorf("无效的日期: %d.%d", month, day)
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}
