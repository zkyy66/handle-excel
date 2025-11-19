package excel_tool

import (
	"fmt"
	"handel-xlsx/model"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func CreateExcel(excelData []*model.ExcelData, yearName string) {
	sheetName := "Sheet1"
	f := excelize.NewFile()
	index, _ := f.NewSheet(sheetName)
	err := f.SetCellValue(sheetName, "A1", "年")
	if err != nil {
		fmt.Printf("A1 set cel value is error: %v\n", err)
		return
	}
	err = f.SetCellValue(sheetName, "B1", "日期")
	if err != nil {
		fmt.Printf("B1 set cel value is error: %v\n", err)
		return
	}
	err = f.SetCellValue(sheetName, "C1", "文件名")
	if err != nil {
		fmt.Printf("C1 set cel value is error: %v\n", err)
		return
	}
	for k, v := range excelData {
		line := k + 2
		lineNum := strconv.Itoa(line)
		err := f.SetCellValue(sheetName, "A"+lineNum, v.Year)
		if err != nil {
			fmt.Printf("A1 set cel value year error: %v\n", err)
			return
		}
		err = f.SetCellValue(sheetName, "B"+lineNum, v.Date)
		if err != nil {
			fmt.Printf("A1 set cel value Date error: %v\n", err)
			return
		}
		err = f.SetCellValue(sheetName, "C"+lineNum, v.FileName)
		if err != nil {
			fmt.Printf("A1 set cel fileName Date error: %v\n", err)
			return
		}
	}
	f.SetActiveSheet(index)
	excelName := yearName + "年.xlsx"
	if err := f.SaveAs("./" + excelName); err != nil {
		println(err.Error())
	}
}
