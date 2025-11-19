package main

import (
	"fmt"
	"handel-xlsx/date-tool"
	excel_tool "handel-xlsx/excel-tool"
	"handel-xlsx/model"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	listOfContents()
}

func listOfContents() {
	//目录路径
	cataloguePath := []string{
		"D:/workSpace/数据汇总/23年报关数据汇总",
		"D:/workSpace/数据汇总/24年报关数据汇总",
		"D:/workSpace/数据汇总/25年报关数据汇总",
	}
	//var files []string
	//var yearName string
	var excelData []*model.ExcelData
	for _, path := range cataloguePath {
		yearName := date_tool.ExtractYearsFromPaths(path)
		var files []string
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("未找到相关数据： %q: %v\n", path, err)
				return nil
			}
			fmt.Println(path)
			
			if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".xlsx") {
				files = append(files, filepath.Base(path))
			}
			return nil
		})
		if err != nil {
			fmt.Printf("路径错误：%q: %v\n", path, err)
			return
		}
		
		//fmt.Println("找到的xlsx文件:")
		var currentYearData []*model.ExcelData
		for _, file := range files {
			if !date_tool.HasDatePrefix(file) {
				continue
			}
			fmt.Println("file:", file)
			dateName, fileName := date_tool.ExtractFileInfo(file)
			//fmt.Printf("年份：%s\t;日期名：%s\t；文件名：%s\n", yearName, dateName, fileName)
			currentYearData = append(currentYearData, &model.ExcelData{
				Year:     yearName,
				Date:     dateName,
				FileName: fileName,
			})
		}
		
		// 按日期正序排序
		sort.Slice(currentYearData, func(i, j int) bool {
			if currentYearData[i].Year != currentYearData[j].Year {
				return currentYearData[i].Year < currentYearData[j].Year
			}
			
			monthI, dayI := getMonthDay(currentYearData[i].Date)
			monthJ, dayJ := getMonthDay(currentYearData[j].Date)
			
			if monthI != monthJ {
				return monthI < monthJ
			}
			return dayI < dayJ
		})
		
		fmt.Println("\n按日期正序排序后:")
		for _, item := range excelData {
			fmt.Printf("日期:%s %s - 文件名:%s\n", item.Year, item.Date, item.FileName)
		}
		excelData = append(excelData, currentYearData...)
	}
	// 最终对所有数据排序
	sort.Slice(excelData, func(i, j int) bool {
		if excelData[i].Year != excelData[j].Year {
			return excelData[i].Year < excelData[j].Year
		}
		
		monthI, dayI := getMonthDay(excelData[i].Date)
		monthJ, dayJ := getMonthDay(excelData[j].Date)
		
		if monthI != monthJ {
			return monthI < monthJ
		}
		return dayI < dayJ
	})
	excel_tool.CreateExcel(excelData, "23-25")
	
}
func getMonthDay(dateStr string) (month, day int) {
	dotIndex := strings.Index(dateStr, ".")
	if dotIndex == -1 {
		return 0, 0
	}
	month, _ = strconv.Atoi(dateStr[:dotIndex])
	dayPart := dateStr[dotIndex+1:]
	var dayStr string
	for _, char := range dayPart {
		if char >= '0' && char <= '9' {
			dayStr += string(char)
		} else {
			break
		}
	}
	day, _ = strconv.Atoi(dayStr)
	
	return month, day
}
