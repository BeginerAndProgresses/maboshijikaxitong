package model

import (
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"maboshijikaxitong/utils"
	"strconv"
)

//var filePath = "./db/studentsInfo.xlsx"

var (
	studentInfoFilePath string
	learnRecordFilePath string
	erlog               *zap.SugaredLogger
	oplog               *zap.SugaredLogger
)

// DBInit DB初始化
func DBInit() {
	studentInfoFilePath = utils.Conf.DbPath.StudentInfo
	learnRecordFilePath = utils.Conf.DbPath.LearnRecord
	erlog = utils.GetErrorLog()
	oplog = utils.GetOperateLog()
	if !utils.IsExist(studentInfoFilePath) {
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				erlog.Info(err)
			}
		}()
		// 创建一个工作表
		index, err := f.NewSheet("Sheet1")
		if err != nil {
			erlog.Info(err)
			return
		}
		// 设置单元格的值
		f.SetCellValue("Sheet1", "A1", "姓名")
		f.SetCellValue("Sheet1", "B1", "电话")
		f.SetCellValue("Sheet1", "C1", "会员类型")
		f.SetCellValue("Sheet1", "D1", "停学天数")
		f.SetCellValue("Sheet1", "E1", "剩余学习周数")
		f.SetCellValue("Sheet1", "F1", "注册时间")
		f.SetCellValue("Sheet1", "G1", "开始时间")
		f.SetCellValue("Sheet1", "H1", "更新时间")
		f.SetCellValue("Sheet1", "I1", "本周剩余学习次数")
		f.SetCellValue("Sheet1", "J1", "所在行")
		f.SetCellValue("Sheet1", "K1", "状态")
		f.SetCellValue("Sheet1", "L1", "停卡时间")
		f.SetCellValue("Sheet1", "M1", "解卡时间")
		// 设置工作簿的默认工作表
		f.SetActiveSheet(index)
		// 根据指定路径保存文件
		if err := f.SaveAs(studentInfoFilePath); err != nil {
			erlog.Info(err)
		}
	}
	err := checkTime()
	if err != nil {
		erlog.Info("有错误：", err)
		return
	}
	if !utils.IsExist(learnRecordFilePath) {
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				erlog.Info(err)
			}
		}()
		// 创建一个工作表
		index, err := f.NewSheet("Sheet1")
		if err != nil {
			erlog.Info(err)
			return
		}
		// 设置单元格的值
		f.SetCellValue("Sheet1", "A1", "姓名")
		f.SetCellValue("Sheet1", "B1", "电话")
		f.SetCellValue("Sheet1", "C1", "学习时间")
		// 设置工作簿的默认工作表
		f.SetActiveSheet(index)
		// 根据指定路径保存文件
		if err := f.SaveAs(learnRecordFilePath); err != nil {
			erlog.Info(err)
		}
	}
}

func getFileByPath(path string) (*excelize.File, error) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		erlog.Info("excelize OpenFile err:", err)
		return nil, err
	}
	return file, nil
}

// 检查时间是否更新
// 如果到时间更新，并重置次数
func checkTime() error {
	file, err := getFileByPath(utils.Conf.DbPath.StudentInfo)
	if err != nil {
		return err
	}
	sheetName := "Sheet1"
	rows, _ := file.Rows(sheetName)
	rows.Next()
	// 确保跳过第一行
	high := 2
	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return err
		}
		stopDays, err := strconv.Atoi(columns[6])
		startTime := utils.DateFormat(columns[6])
		updateTime := utils.DateFormat(columns[7])
		//lockTime := ""
		//unlockTime := ""
		//if len(columns) >= 12 {
		//	lockTime = utils.DateFormat(columns[11])
		//}
		//if len(columns) >= 13 {
		//	unlockTime = utils.DateFormat(columns[12])
		//}
		times, err := strconv.Atoi(columns[8])
		states := columns[10]
		if err != nil {
			return err
		}
		if states == States[0] || states == States[2] || states == States[3] {
			high++
			continue
		}
		day := GetDayByLevel(columns[2])
		// 现在据开始的时间已学习天数 days
		days := LearnedDays(startTime, uint(stopDays))
		// 计算获得剩余周数
		_, weekday := GetTimeAndWeekByLevel(columns[2])

		newRemWeeks := weekday - days/day
		// 更新时间
		// 计算获得的更新时间
		newUpdate := utils.DateOfNDays(-(days % day))
		// 更新次数，根据会员字段更改
		if updateTime != newUpdate && newRemWeeks > 0 {
			times, _ = GetTimeAndWeekByLevel(columns[2])
		} else if newRemWeeks <= 0 {
			times = 0
		}
		updateInfoByRowNumber(file, high, newRemWeeks, times, updateTime)
		high++
	}
	if err = file.Save(); err != nil {
		return err
	}
	return nil
}

// 更新第row行的数据
func updateInfoByRowNumber(f *excelize.File, row, RemWeeks, timeOfWeek int, updateTime string) {
	sheetName := "Sheet1"
	cell, _ := excelize.CoordinatesToCellName(5, int(row))
	f.SetCellValue(sheetName, cell, RemWeeks)
	if RemWeeks == 0 {
		s := GetStudentByRow(uint(row))
		changeState(f, s, sheetName)
	}
	cell, _ = excelize.CoordinatesToCellName(8, int(row))
	f.SetCellValue(sheetName, cell, updateTime)
	cell, _ = excelize.CoordinatesToCellName(9, int(row))
	f.SetCellValue(sheetName, cell, timeOfWeek)
}

// GetAllInfo 获取文件中的信息
func GetAllInfo(filename string) ([][]string, error) {
	file, err := getFileByPath(filename)
	if err != nil {
		erlog.Info("err:", err)
		return nil, err
	}
	defer file.Close()
	rows, err := file.GetRows("Sheet1")
	if err != nil {
		erlog.Info("err:", err)
		return nil, err
	}
	return rows, nil
}

// AppendInfo2Sheet 文件信息追加到已有文件后
func AppendInfo2Sheet(filename string) error {
	oplog := utils.GetOperateLog()
	data, err := GetAllInfo(filename)
	if err != nil {
		return err
	}
	if len(data) <= 1 {
		return nil
	}
	data = data[1:]
	//打开工作簿
	file, err := getFileByPath(utils.Conf.DbPath.StudentInfo)
	if err != nil {
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			erlog.Info("error closing file: ", err)
		}
	}()
	sheet_name := "Sheet1"
	//获取流式写入器
	streamWriter, _ := file.NewStreamWriter(sheet_name)
	if err != nil {
		erlog.Info(err)
		return err
	}

	rows, _ := file.GetRows(sheet_name) //获取行内容
	cols, _ := file.GetCols(sheet_name) //获取列内容
	collen := len(cols)
	rowlen := len(rows)
	//将源文件内容先写入excel
	for rowid, row_pre := range rows {
		row_p := make([]interface{}, collen)
		for colID_p := 0; colID_p < collen; colID_p++ {
			if row_pre == nil {
				row_p[colID_p] = nil
				continue
			}
			if colID_p == 5 || colID_p == 6 || colID_p == 7 || colID_p == 11 || colID_p == 12 {
				// 跳过一行中可能会有的两个空白格
				if colID_p < len(row_pre) {
					row_p[colID_p] = utils.DateFormat(row_pre[colID_p])
				} else {
					row_p[colID_p] = ""
				}
			} else {
				row_p[colID_p] = row_pre[colID_p]
			}
		}
		cell_pre, _ := excelize.CoordinatesToCellName(1, rowid+1)
		if err := streamWriter.SetRow(cell_pre, row_p); err != nil {
			erlog.Info(err)
		}
	}

	//将新加contents写进流式写入器
	for rowID := 0; rowID < len(data); rowID++ {
		row := make([]interface{}, collen)
		times := 0
		week := 0
		for colID := 0; colID < len(data[0]); colID++ {
			switch colID {
			case 2:
				times, week = GetTimeAndWeekByLevel(data[rowID][colID])
				row[colID] = data[rowID][colID]
			case 3:
				row[5] = utils.DateFormat(data[rowID][colID])
			default:
				row[colID] = data[rowID][colID]
			}
		}
		row[3] = 0
		row[4] = week
		row[8] = times
		row[9] = rowID + rowlen + 1
		row[10] = States[0]
		oplog.Infof("- 写入%s成功", data[rowID][0])
		cell, _ := excelize.CoordinatesToCellName(1, rowID+rowlen+1) //决定写入的位置
		if err := streamWriter.SetRow(cell, row); err != nil {
			erlog.Info(err)
		}
	}

	//结束流式写入过程
	if err := streamWriter.Flush(); err != nil {
		erlog.Info(err)
	}
	//保存工作簿
	if err := file.SaveAs(utils.Conf.DbPath.StudentInfo); err != nil {
		erlog.Info(err)
	}

	return nil
}

// GetLockTime 获取停卡时间
func GetLockTime(f *excelize.File, row uint) string {
	cell, _ := excelize.CoordinatesToCellName(12, int(row))
	value, _ := f.GetCellValue(sheet_name, cell)
	return value
}

// GetUnlockTime 获取解卡时间
func GetUnlockTime(f *excelize.File, row uint) string {
	cell, _ := excelize.CoordinatesToCellName(13, int(row))
	value, _ := f.GetCellValue(sheet_name, cell)
	return value
}

// GetStopDays 获取停卡天数
func GetStopDays(f *excelize.File, row uint) string {
	cell, _ := excelize.CoordinatesToCellName(4, int(row))
	value, _ := f.GetCellValue(sheet_name, cell)
	return value
}

// SetStopDays 设置停卡天数
func SetStopDays(f *excelize.File, row uint, value interface{}) error {
	cell, _ := excelize.CoordinatesToCellName(4, int(row))
	err := f.SetCellValue(sheet_name, cell, value)
	return err
}
