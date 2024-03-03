package model

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"maboshijikaxitong/utils"
	"strconv"
	"time"
)

// Student 需要重写
type Student struct {
	name             string // 姓名
	phone            string // 电话号码
	memberType       string // 会员类型
	stopDays         uint   // 停学天数
	remainingTime    int    // 剩余可学习周数
	registrationTime string // 注册时间，即录入系统的时间
	startLearnTime   string // 开始学习时间
	updateTime       string // 更新时间，从开始时间开始计算，每过七天，更新学习次数，和剩余学习周数
	numberOfStudies  uint   // 学习次数
	Row              uint   // 所在行数
	State            string // 当前学生状态 “学习中”，“停卡中”，“未开始学习”，“学习已完毕”
	lockTime         string // 停卡时间
	unlockTime       string // 解锁时间
}

var States = []string{"未开始学习", "学习中", "学习已完毕", "停卡中"}
var sheet_name = "Sheet1"

// 将表中的一行数据映射为一个对象
func ref2Student(student *Student, row []string) {
	for i, colCell := range row {
		switch i {
		case 0:
			student.name = colCell
		case 1:
			student.phone = colCell
		case 2:
			student.memberType = colCell
		case 3:
			atio, _ := strconv.Atoi(colCell)
			student.stopDays = uint(atio)
		case 4:
			atio, _ := strconv.Atoi(colCell)
			student.remainingTime = atio
		case 5:
			student.registrationTime = utils.DateFormat(colCell)
		case 6:
			student.startLearnTime = utils.DateFormat(colCell)
		case 7:
			student.updateTime = utils.DateFormat(colCell)
		case 8:
			atio, _ := strconv.Atoi(colCell)
			student.numberOfStudies = uint(atio)
		case 9:
			atio, _ := strconv.Atoi(colCell)
			student.Row = uint(atio)
		case 10:
			student.State = colCell
		case 11:
			student.lockTime = utils.DateFormat(colCell)
		case 12:
			student.unlockTime = utils.DateFormat(colCell)
		}
	}
}

// GetStudents  根据cellText与第col列单元格数据对比，如果相当就返回
// typ == 0 cellText传过来的是姓名
// typ == 1 cellText传过来的是后四位
func GetStudents(cellText string, typ int, excludeZeroTime bool) ([]Student, error) {
	// 打开工作文件
	var students []Student
	f, err := excelize.OpenFile(utils.Conf.DbPath.StudentInfo)
	if err != nil {
		return nil, err
	}

	rows, err := f.Rows("Sheet1")
	if err != nil {
		return nil, err
	}
	// 确保跳过第一行
	rows.Next()
	high := 2
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			erlog.Info(err)
		}
		var student = new(Student)
		student.Row = uint(high)
		high++
		ref2Student(student, row)
		// 是否何以匹配
		flag := (typ == 0 && student.name == cellText) || (typ == 1 && utils.Last4Rune(student.phone) == cellText)
		// flag1 代表学员是否没课
		flag1 := student.remainingTime <= 0
		if excludeZeroTime {
			if flag {
				if !flag1 {
					students = append(students, *student)
				}
			}
		} else {
			if flag {
				students = append(students, *student)
			}
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
	return students, nil
}

func (s *Student) PrintStudent() string {
	info := "哈哈"
	switch s.State {
	case States[1]:
		info = fmt.Sprintf("学生姓名：%s，电话号码：%s，会员等级：%s，注册时间：%s，剩余：%d天，本周剩余次数：%d\n",
			s.name,
			s.phone,
			s.memberType,
			s.registrationTime,
			GetDayByLevel(s.memberType)-LearnedDays(s.startLearnTime, s.stopDays),
			s.numberOfStudies)
	case States[0]:
		info = fmt.Sprintf("学生姓名：%s，电话号码：%s，会员等级：%s，注册时间：%s，未开始学习，本周剩余次数：%d\n",
			s.name,
			s.phone,
			s.memberType,
			s.registrationTime,
			s.numberOfStudies)
	case States[2]:
		info = fmt.Sprintf("学生姓名：%s，电话号码：%s，会员等级：%s，注册时间：%s，学习已完毕\n",
			s.name,
			s.phone,
			s.memberType,
			s.registrationTime)
	case States[3]:
		info = fmt.Sprintf("学生姓名：%s，电话号码：%s，会员等级：%s，注册时间：%s，停卡中\n",
			s.name,
			s.phone,
			s.memberType,
			s.registrationTime)
	}

	return info
}

func (s *Student) GetName() string {
	return s.name
}

func TimeMinus(s *Student, step int) (err error) {
	if s.State == States[2] {
		return fmt.Errorf("%s已无课程", s.name)
	}
	if s.State == States[3] {
		return fmt.Errorf("%s处于停卡状态，请先解除停卡状态在进行操作", s.name)
	}
	file, err := getFileByPath(utils.Conf.DbPath.StudentInfo)
	if err != nil {
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			erlog.Info(err)
		}
	}()
	cell, _ := excelize.CoordinatesToCellName(9, int(s.Row))
	value, _ := file.GetCellValue(sheet_name, cell)
	atoi, _ := strconv.Atoi(value)
	times := atoi
	if times-step < 0 {
		return
	}
	err = file.SetCellValue(sheet_name, cell, times-step)
	changeState(file, s, sheet_name)
	err = ExcelFileSave(file)
	err = SaveLearnTime(s.name, s.phone)
	if err != nil {
		erlog.Info("保存记录失败，err：", err)
		return err
	}
	return
}

// 在课程减1时，改变s状态
func changeState(file *excelize.File, s *Student, sheet_name string) (err error) {
	// 从未开始状态转为开始状态
	if s.updateTime == "" && s.startLearnTime == "" {
		nowDay := time.Now().Format(time.DateOnly)
		// 开始时间
		cell, _ := excelize.CoordinatesToCellName(7, int(s.Row))
		err = file.SetCellValue(sheet_name, cell, nowDay) // 更新时间
		cell, _ = excelize.CoordinatesToCellName(8, int(s.Row))
		err = file.SetCellValue(sheet_name, cell, nowDay) // 更新时间
		cell, _ = excelize.CoordinatesToCellName(11, int(s.Row))
		err = file.SetCellValue(sheet_name, cell, States[1])
		s.updateTime = nowDay
		s.registrationTime = nowDay
		s.State = States[1]
	}
	// 从学习中转为以学完所有课程
	cell, _ := excelize.CoordinatesToCellName(9, int(s.Row))
	value, _ := file.GetCellValue(sheet_name, cell)
	times, _ := strconv.Atoi(value)
	cell, _ = excelize.CoordinatesToCellName(5, int(s.Row))
	value, _ = file.GetCellValue(sheet_name, cell)
	remWeek, _ := strconv.Atoi(value)
	if (times == 0 && remWeek == 1) || remWeek == 0 {
		cell, _ = excelize.CoordinatesToCellName(11, int(s.Row))
		err = file.SetCellValue(sheet_name, cell, States[2])
		if err != nil {
			return err
		}
		// 将剩余周数变为0
		cell, _ = excelize.CoordinatesToCellName(5, int(s.Row))
		err = file.SetCellValue(sheet_name, cell, 0)
		if err != nil {
			return err
		}
		s.State = States[2]
	}
	return
}

// SetLockTime  停卡
func SetLockTime(s *Student) error {
	// 未完成
	file, err := getFileByPath(utils.Conf.DbPath.StudentInfo)
	if err != nil {
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			erlog.Info(err)
		}
	}()
	now := time.Now().Format(time.DateOnly)
	cell, _ := excelize.CoordinatesToCellName(12, int(s.Row))
	err = file.SetCellValue(sheet_name, cell, now)
	if err != nil {
		return err
	}
	cell, _ = excelize.CoordinatesToCellName(11, int(s.Row))
	err = file.SetCellValue(sheet_name, cell, States[3])
	if err != nil {
		return err
	}
	s.State = States[3]
	s.lockTime = now
	err = ExcelFileSave(file)
	return err
}

// SetUnlockTime  解卡
func SetUnlockTime(s *Student) error {
	file, err := getFileByPath(utils.Conf.DbPath.StudentInfo)
	if err != nil {
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			erlog.Info(err)
		}
	}()

	now := time.Now().Format(time.DateOnly)
	cell, _ := excelize.CoordinatesToCellName(13, int(s.Row))
	err = file.SetCellValue(sheet_name, cell, now)
	if err != nil {
		return err
	}

	// 将停卡天数增加
	lockTime := GetLockTime(file, s.Row)
	unlockTime := GetUnlockTime(file, s.Row)
	stopDays, _ := strconv.Atoi(GetStopDays(file, s.Row))
	thisStopDays := StopDays(lockTime, unlockTime)
	newStopDays := stopDays + thisStopDays
	err = SetStopDays(file, s.Row, newStopDays)
	if err != nil {
		return err
	}

	cell, _ = excelize.CoordinatesToCellName(11, int(s.Row))
	err = file.SetCellValue(sheet_name, cell, States[1])
	if err != nil {
		return err
	}
	s.State = States[1]
	s.unlockTime = now
	err = ExcelFileSave(file)
	return err
}

func GetStudentByRow(row uint) *Student {
	s := new(Student)
	s.Row = row
	file, err := getFileByPath(utils.Conf.DbPath.StudentInfo)
	if err != nil {
		erlog.Info("err:", err)
		return nil
	}

	rows, _ := file.GetRows(sheet_name)
	rowValue := rows[row-1]

	ref2Student(s, rowValue)

	defer func() {
		if err := file.Close(); err != nil {
			erlog.Info(err)
		}
	}()
	return s
}

// SaveStudent 保存用户
func SaveStudent(name, phone, member string) error {
	var err error
	file, err := getFileByPath(utils.Conf.DbPath.StudentInfo)
	if err != nil {
		erlog.Info("err:", err)
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			erlog.Info(err)
		}
	}()

	rows, _ := file.GetRows(sheet_name) //获取行内容
	rowslen := len(rows)
	cell, _ := excelize.CoordinatesToCellName(1, rowslen+1)
	file.SetCellValue(sheet_name, cell, name)
	cell, _ = excelize.CoordinatesToCellName(2, rowslen+1)
	file.SetCellValue(sheet_name, cell, phone)
	cell, _ = excelize.CoordinatesToCellName(3, rowslen+1)
	file.SetCellValue(sheet_name, cell, member)
	times, weeks := GetTimeAndWeekByLevel(member)

	cell, _ = excelize.CoordinatesToCellName(4, rowslen+1)
	file.SetCellValue(sheet_name, cell, 0)
	cell, _ = excelize.CoordinatesToCellName(5, rowslen+1)
	file.SetCellValue(sheet_name, cell, weeks)
	cell, _ = excelize.CoordinatesToCellName(6, rowslen+1)
	file.SetCellValue(sheet_name, cell, time.Now().Format(time.DateOnly))
	cell, _ = excelize.CoordinatesToCellName(9, rowslen+1)
	file.SetCellValue(sheet_name, cell, times)
	cell, _ = excelize.CoordinatesToCellName(10, rowslen+1)
	file.SetCellValue(sheet_name, cell, rowslen+1)
	cell, _ = excelize.CoordinatesToCellName(11, rowslen+1)
	file.SetCellValue(sheet_name, cell, States[0])
	if err != nil {
		erlog.Info(name + " 插入失败")
	}

	//保存工作簿
	err = ExcelFileSave(file)
	return err
}

// GetTimeAndWeekByLevel 根据等级获取每轮次数和总轮数
func GetTimeAndWeekByLevel(member string) (times, weeks int) {
	levellen := len(utils.Conf.Level)
	if levellen == 0 {
		return 0, 0
	}
	i := 0
	for i < levellen && member != utils.Conf.Level[i] {
		i++
	}
	if i >= levellen {
		i = 0
		return utils.Conf.TimeOfLoop[0], utils.Conf.Level2Loops[0]
	}
	times = utils.Conf.TimeOfLoop[i]
	weeks = utils.Conf.Level2Loops[i]
	return
}

// GetDayByLevel 根据等级获取天数
func GetDayByLevel(level string) (day int) {
	levellen := len(utils.Conf.Level)
	if levellen == 0 {
		return 0
	}
	i := 0
	for i < levellen && level != utils.Conf.Level[i] {
		i++
	}
	if i >= levellen {
		i = 0
		return utils.Conf.DayOfLoop[0]
	}
	day = utils.Conf.DayOfLoop[i]
	return
}

func ExcelFileSave(file *excelize.File) (err error) {
	if err = file.Save(); err != nil {
		erlog.Info("保存失败，err：", err)
	}
	return err
}

// LearnedDays 已学习天数
func LearnedDays(startTime string, stopDays uint) (days int) {
	stopDay := stopDays
	allDay := utils.DifferenceOfDay(time.Now().Format(time.DateOnly), startTime)
	days = allDay - int(stopDay)
	return
}

// StopDays 获取停卡天数
func StopDays(lockTime, unlockTime string) (days int) {
	days = 0
	if lockTime != "" && unlockTime != "" {
		days = utils.DifferenceOfDay(utils.DateFormat(unlockTime), utils.DateFormat(lockTime))
	}
	return
}
