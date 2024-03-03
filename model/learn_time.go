package model

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"maboshijikaxitong/utils"
	"time"
)

type LearnTime struct {
	name      string // 学生姓名
	phone     string // 电话号码
	learnTime string // 学习时间
}

func (t *LearnTime) PrintString() string {
	return fmt.Sprintf("姓名：%s,电话：%s,学习时间：%s", t.name, t.phone, t.learnTime)
}

// GetLearnTimes
// typ == 0 cellText传过来的是姓名
// typ == 1 cellText传过来的是后四位
func GetLearnTimes(cellText string, typ int) ([]LearnTime, error) {
	// 打开工作文件
	var learnTimes []LearnTime
	f, err := excelize.OpenFile(utils.Conf.DbPath.LearnRecord)
	if err != nil {
		return nil, err
	}

	rows, err := f.Rows("Sheet1")
	if err != nil {
		return nil, err
	}
	// 确保跳过第一行
	rows.Next()
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			erlog.Info(err)
		}
		var learnTime = new(LearnTime)
		ref2LearnTime(learnTime, row)
		// 是否何以匹配
		flag := (typ == 0 && learnTime.name == cellText) || (typ == 1 && utils.Last4Rune(learnTime.phone) == cellText)

		if flag {
			learnTimes = append(learnTimes, *learnTime)
		}

	}

	defer func() {
		if err = f.Close(); err != nil {
			erlog.Info(err)
		}
		if err = rows.Close(); err != nil {
			erlog.Info(err)
		}
	}()
	return learnTimes, nil
}

// 将表中的一行数据映射为一个对象
func ref2LearnTime(learn *LearnTime, row []string) {
	for i, colCell := range row {
		switch i {
		case 0:
			learn.name = colCell
		case 1:
			learn.phone = colCell
		case 2:
			learn.learnTime = utils.DateFormat(colCell)
		}
	}
}

// SaveLearnTime 保存记录
func SaveLearnTime(name, phone string) error {
	var err error
	file, err := getFileByPath(utils.Conf.DbPath.LearnRecord)
	if err != nil {
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			erlog.Info(err)
		}
	}()
	sheetName := "Sheet1"

	rows, _ := file.GetRows(sheetName) //获取行内容
	rowslen := len(rows)
	cell, _ := excelize.CoordinatesToCellName(1, rowslen+1)
	file.SetCellValue(sheetName, cell, name)
	cell, _ = excelize.CoordinatesToCellName(2, rowslen+1)
	file.SetCellValue(sheetName, cell, phone)
	cell, _ = excelize.CoordinatesToCellName(3, rowslen+1)
	file.SetCellValue(sheetName, cell, time.Now().Format(time.DateOnly))

	//保存工作簿
	if err = file.Save(); err != nil {
		erlog.Info(err)
	}
	return err
}
