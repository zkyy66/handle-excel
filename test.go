package main

import (
	"fmt"
	date_tool "handel-xlsx/date-tool"
	"handel-xlsx/model"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func listOfContentss() {
	cataloguePath := []string{
		"D:/workSpace/数据汇总/23年报关数据汇总",
		"D:/workSpace/数据汇总/24年报关数据汇总",
		"D:/workSpace/数据汇总/25年报关数据汇总",
	}
	
	var allExcelData []*model.ExcelData
	
	for _, path := range cataloguePath {
		yearName := date_tool.ExtractYearsFromPaths(path)
		var files []string // 在每个目录内重新定义files
		
		err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("未找到相关数据： %q: %v\n", filePath, err)
				return nil
			}
			
			if !d.IsDir() && strings.HasSuffix(strings.ToLower(filePath), ".xlsx") {
				fileName := filepath.Base(filePath)
				if date_tool.HasDatePrefix(fileName) {
					files = append(files, fileName)
				}
			}
			return nil
		})
		
		if err != nil {
			fmt.Printf("路径错误：%q: %v\n", path, err)
			continue
		}
		
		// 处理当前目录的文件
		var currentYearData []*model.ExcelData
		for _, file := range files {
			fmt.Println("file:", file)
			dateName, fileName := date_tool.ExtractFileInfo(file)
			
			// 添加年份验证逻辑
			if !isValidDateForYear(yearName, dateName) {
				fmt.Printf("警告: 文件 %s 的日期 %s 与年份 %s 不匹配\n", file, dateName, yearName)
				continue
			}
			
			currentYearData = append(currentYearData, &model.ExcelData{
				Year:     yearName,
				Date:     dateName,
				FileName: fileName,
			})
		}
		
		// 对当前年份数据排序
		sort.Slice(currentYearData, func(i, j int) bool {
			monthI, dayI := getMonthDay(currentYearData[i].Date)
			monthJ, dayJ := getMonthDay(currentYearData[j].Date)
			
			if monthI != monthJ {
				return monthI < monthJ
			}
			return dayI < dayJ
		})
		
		fmt.Printf("\n%s 按日期正序排序后:\n", yearName)
		for _, item := range currentYearData {
			fmt.Printf("日期:%s %s - 文件名:%s\n", item.Year, item.Date, item.FileName)
		}
		
		allExcelData = append(allExcelData, currentYearData...)
	}
	
	// 最终对所有数据排序
	sort.Slice(allExcelData, func(i, j int) bool {
		if allExcelData[i].Year != allExcelData[j].Year {
			return allExcelData[i].Year < allExcelData[j].Year
		}
		
		monthI, dayI := getMonthDay(allExcelData[i].Date)
		monthJ, dayJ := getMonthDay(allExcelData[j].Date)
		
		if monthI != monthJ {
			return monthI < monthJ
		}
		return dayI < dayJ
	})
	
	fmt.Println("\n所有数据最终排序结果:")
	for _, item := range allExcelData {
		fmt.Printf("日期:%s %s - 文件名:%s\n", item.Year, item.Date, item.FileName)
	}
}

// 添加日期验证函数
func isValidDateForYear(yearName, dateName string) bool {
	// 提取年份数字
	yearStr := strings.TrimSuffix(yearName, "年")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return true // 如果解析失败，暂时不过滤
	}
	
	month, _ := getMonthDay(dateName)
	
	// 特殊处理2025年，只允许1-3月
	if year == 2025 || year == 25 {
		return month >= 1 && month <= 3
	}
	
	return month >= 1 && month <= 12
}
