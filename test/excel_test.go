package test

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
	"testing"
	"time"
)

func TestExcel(t *testing.T) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 创建一个工作表
	index, err := f.NewSheet("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置单元格的值
	f.SetCellValue("Sheet2", "A2", "2023-08-17")
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func TestGetValue(t *testing.T) {
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	value, _ := f.GetCellValue("Sheet2", "A2")
	fmt.Println(DataFormat(value))
}

// 判断获取到的日期格式，如果日期格式为08-17-23则转为2023-08-17这样的格式
func DataFormat(data string) string {
	if len(data) <= 8 {
		split := strings.Split(data, "-")
		s := time.Now().String()[0:2]
		return fmt.Sprintf("%s%s-%s-%s", s, split[2], split[0], split[1])
	}
	return data
}
