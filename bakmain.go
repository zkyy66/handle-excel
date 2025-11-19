package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	//handleExcelOne()
	handleExcel()
}
func handleExcel() {

	fmt.Println("处理用户邮箱的excel...")
	filesTwo, err := excelize.OpenFile("./test.xlsx")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//读取excel2de第一个sheet
	//filesTwo.GetCellValue()
	valueRow, err := filesTwo.GetRows("SheetJS")
	for _, row := range valueRow {
		if row[0] != "email" && row[1] != "code" {
			fmt.Println(row)
		}

	}
}
func handleExcelOne() {

	fmt.Println("处理excel1...")
	filesTwo, err := excelize.OpenFile("/Users/yaoyuan/Desktop/Yunao/2022-10-10-19-51-57_EXPORT_XLSX_6953719_059_0.xlsx")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//读取excel2de第一个sheet
	var sheetTwoValue1 []int
	valueRow, err := filesTwo.GetRows("Sheet1")
	for _, row := range valueRow {
		for _, colCell := range row {
			value1, _ := strconv.Atoi(colCell)
			if value1 == 0 {
				continue
			}
			sheetTwoValue1 = append(sheetTwoValue1, value1)
		}
	}
	fmt.Println("excel1的UID总数：", len(sheetTwoValue1))
	fmt.Println("excel1 end")

	fmt.Println("处理excel2...start")
	files, err := excelize.OpenFile("/Users/yaoyuan/Desktop/Yunao/sf_pg_trade_log_connect.xlsx")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//读取第一个sheet
	var sheetValue1 []int
	rows, err := files.GetRows("Result 1")
	for _, row := range rows {
		for _, colCell := range row {
			value1, _ := strconv.Atoi(colCell)
			sheetValue1 = append(sheetValue1, value1)
		}
	}
	fmt.Println("Result1下的UID总数：", len(sheetValue1), "\t")

	//读取第二个sheet
	var sheetValue2 []int
	rowsSheetTwo, err := files.GetRows("Result 2")
	for _, sheetTwo := range rowsSheetTwo {
		for _, colcellTwo := range sheetTwo {
			value2, _ := strconv.Atoi(colcellTwo)
			sheetValue2 = append(sheetValue2, value2)
		}
	}
	fmt.Println("Result2下的UID总数：", len(sheetValue2), "\t")

	var totalUIdValue []int
	totalUIdValue = append(sheetValue1, sheetValue2...)
	fmt.Println("excel2的总数：", len(totalUIdValue))
	fmt.Println("处理excel2...end")

	////差集
	//res := diffItem(sheetTwoValue1, totalUIdValue)
	//fmt.Println("差集")
	//createExcel(res, 1)
	//交集
	res := intersectArray(sheetTwoValue1, totalUIdValue)
	fmt.Println("交集")
	screateExcel(res, 2)
}

// 差集
func diffItem(a []int, b []int) []int {
	var diffArray []int
	temp := map[int]struct{}{}
	am := map[int]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			diffArray = append(diffArray, val)
		}
		am[val] = struct{}{}
	}
	fmt.Println("map len:", len(am))
	return diffArray
}

// 交集
func intersectArray(a []int, b []int) []int {
	var inter []int
	mp := make(map[int]bool)

	for _, s := range a {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range b {
		if _, ok := mp[s]; ok {
			inter = append(inter, s)
		}
	}

	return inter
}

func intersectArrayByString(a []string, b []string) []string {
	var inter []string
	mp := make(map[string]bool)

	for _, s := range a {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range b {
		if _, ok := mp[s]; ok {
			inter = append(inter, s)
		}
	}
	return inter
}

// 生成excel
func screateExcel(res []int, mark int) {
	f := excelize.NewFile()
	// 创建一个工作表
	index, _ := f.NewSheet("Sheet1")
	// 设置单元格的值
	err := f.SetCellValue("Sheet1", "A1", "UID")
	if err != nil {
		return
	}
	for k, v := range res {
		line := k + 2
		lineNum := strconv.Itoa(line)
		err := f.SetCellValue("Sheet1", "A"+lineNum, v)
		if err != nil {
			return
		}
	}

	//f.SetCellValue("Sheet1", "B2", 100)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if mark == 1 {
		if err := f.SaveAs("./UID差集-2022-11-08.xlsx"); err != nil {
			println(err.Error())
		}
	} else {
		if err := f.SaveAs("./UID交集-2022-11-08.xlsx"); err != nil {
			println(err.Error())
		}
	}
}

// img转base64
func handelImgToBase64() {
	//handleMail()
	_, name, err := imgToBase64("https://lmg.jj20.com/up/allimg/1114/040221103339/210402103339-8-1200.jpg")
	if err != nil {
		return
	}

	fmt.Println("tupian:", name)
	srcByte, err := ioutil.ReadFile("./" + name)
	if err != nil {
		log.Fatal(err)
	}
	res := base64.StdEncoding.EncodeToString(srcByte)
	fmt.Println(res)
	//调用此方法需要在函数中声明接受变量：func handleMail(res string)
	//handleMail()
}

func imgToBase64(url string) (int64, string, error) {
	path := strings.Split(url, "/")
	var name string
	if len(path) > 1 {
		name = path[len(path)-1]
	}
	out, _ := os.Create(name)
	defer out.Close()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	n, _ := io.Copy(out, bytes.NewReader(pix))
	return n, name, err
}
