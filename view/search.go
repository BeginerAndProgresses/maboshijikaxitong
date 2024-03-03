package view

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"maboshijikaxitong/model"
	util "maboshijikaxitong/utils"
	"reflect"
	"strconv"
	"time"
)

// BindFuncForSearch
// 为entry和button绑定函数
// isStudent == true 绑定学生信息
// isStudent == false 绑定学习记录
func BindFuncForSearch(entry *gtk.Entry, button *gtk.Button, isStudent bool) {
	if entry == nil {
		return
	}
	// 输入回车触发activate信号
	entry.SetActivatesDefault(true)
	entry.Connect("activate", func() {
		Search(entry, isStudent)
	})
	if button == nil {
		return
	}
	button.Connect("clicked", func() {
		Search(entry, isStudent)
	})
}

// Search
// isStudent == true 绑定学生信息
// isStudent == false 绑定学习记录
func Search(entry *gtk.Entry, isStudent bool) *gtk.Window {
	text, _ := entry.GetText()
	if text == "" {
		Tips("请输入姓名或者手机号后四位", 0, nil)
		return nil
	} else {
		infoWinCh := make(chan int)
		var infowin *gtk.Window
		waitWindow := WaitWindow("查询中,请稍后...")
		if util.OneRuneIsNumber(text) {
			infowin = SearchByText(infoWinCh, text, 1, isStudent)
		} else {
			infowin = SearchByText(infoWinCh, text, 0, isStudent)
		}
		go func() {
			select {
			case <-infoWinCh:
				// 如果是1刷新窗口
			}
			if infowin != nil {
				infowin.SetModal(false)
				infowin.Destroy()
			}
		}()
		WaitClose(waitWindow)
		//go func() {
		//	select {
		//	case <-infoWinCh:
		//
		//	}
		//	if infowin != nil {
		//		infowin.SetModal(false)
		//		infowin.Destroy()
		//	}
		//}()
	}
	// 清空搜索框
	ClearEntry("请输入姓名或者手机号后四位", entry)
	return nil
}

// SearchByText
// isStudent == true 绑定学生信息
// isStudent == false 绑定学习记录
func SearchByText(ch chan int, Text string, typ int, isStudent bool) *gtk.Window {
	var win *gtk.Window
	var infos interface{}
	var err error
	if isStudent {
		infos, err = model.GetStudents(Text, typ, false)
	} else {
		infos, err = model.GetLearnTimes(Text, typ)
	}

	a := reflect.ValueOf(infos)

	if err != nil {
		Tips("请确保EXCEL表格未被打开", 1, ch)
		return nil
	}
	if isStudent {
		win = InfoBox(a.Interface(), ch, 1350, 1550)
	} else {
		win = InfoBox(a.Interface(), ch, 700, 1000)
	}
	return win
}

// InfoBox a 是要展示的类型数组
func InfoBox(a interface{}, ch chan int, winWide int, layWide uint) *gtk.Window {
	switch a.(type) {
	case []model.Student:
		students, ok := a.([]model.Student)
		if !ok {
			erlog.Info("断言a为[]model.Student失败")
		}
		if len(students) == 0 {
			Tips("未找到相关信息", 0, ch)
			return nil
		} else {
			builder, err := gtk.BuilderNew()
			if err != nil {
				erlog.Info("BuilderNew 创建builder失败")
			}

			err = builder.AddFromFile(rootPath + "glade/info.glade")
			if err != nil {
				erlog.Info("AddFromFile builder失败")
			}
			window1, err := builder.GetObject("window1")
			win1 := window1.(*gtk.Window)

			win1.SetTitle(util.Conf.Title)
			err = win1.SetIconFromFile(rootPath + "img/icon/maboshi.png")
			if err != nil {
				erlog.Info("SetIconFromFile 读取图标失败")
			}

			win1.Resize(winWide, 600)
			win1.Connect("destroy", func() {
				ch <- 0
			})

			layout1, err := builder.GetObject("layout1")
			if err != nil {
				erlog.Info("layout1 获取layout1失败")
			}

			// 排版
			lay1 := layout1.(*gtk.Layout)

			box1, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
			box1.SetHomogeneous(true)
			for i, student := range students {
				box2, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
				label1, _ := gtk.LabelNew(strconv.Itoa(i+1) + "、 " + student.PrintStudent())
				label1.SetSizeRequest(400, 50)
				label1.SetMarginTop(25)
				button1, _ := gtk.ButtonNewWithLabel("课程减一")
				button1.SetSizeRequest(60, 40)
				// 通过闭包解决数据改变的问题
				button1.Connect("clicked", func() func() {
					//s := student
					//switch s.State {
					//case model.States[2]:
					//	return func() {
					//		Tips("该学生已完成学习", 0, nil)
					//	}
					//case model.States[3]:
					//	return func() {
					//		Tips("处于停卡状态，请先解卡在做操作", 0, nil)
					//	}
					//}
					bh := i
					lab := label1
					row := student.Row
					return func() {
						s := model.GetStudentByRow(row)
						switch s.State {
						case model.States[2]:
							Tips("该学生已完成学习", 0, nil)
							return
						case model.States[3]:
							Tips("处于停卡状态，请先解卡在做操作", 0, nil)
							return
						}

						prompch := make(chan int)
						info := fmt.Sprintf("减少%s的课时数吗？", s.GetName())
						dialog := promptBox(prompch, info, 0)
						go closeDialog(prompch, dialog, 24*time.Hour, func(i int) {
							if i == 1 {
								err2 := model.TimeMinus(s, 1)
								if err2 != nil {
									oplog.Infof("减少%s的课时数失败，err：%s", s.GetName(), err2.Error())
								}
								oplog.Infof("减少%s的课时数成功", s.GetName())
								lab.SetText(strconv.Itoa(bh+1) + "、 " + model.GetStudentByRow(s.Row).PrintStudent())

							}
						})
					}
				}())
				button2, _ := gtk.ButtonNew()
				button2.SetSizeRequest(60, 40)
				switch student.State {
				case model.States[0]:
					// 禁用按钮
					button2.SetLabel("停卡")
					button2.SetSensitive(false)
				case model.States[1]:
					button2.SetLabel("停卡")
				case model.States[2]:
					// 禁用按钮
					button2.SetLabel("停卡")
					button2.SetSensitive(false)
				case model.States[3]:
					// 禁用按钮
					button2.SetLabel("解卡")
				}
				button2.Connect("clicked", func() func() {
					//but := button2
					row := student.Row
					bh := i
					lab := label1
					return func() {
						s := model.GetStudentByRow(row)
						//fmt.Println("state:", s.State)
						if s.State == model.States[1] {
							//fmt.Println("停卡")
							//prompch := make(chan int)
							//info := fmt.Sprintf("确定停%s的卡吗？停卡后将该学生的计时将终止", s.GetName())
							//dialog := promptBox(prompch, info, 0)
							//go closeDialog(prompch, dialog, 24*time.Hour, func(i int) {
							//	// 问题所在
							//	if i == 1 {
							//		err2 := model.SetLockTime(s)
							//		if err2 != nil {
							//			oplog.Infof("对%s停卡失败，err：%s", s.GetName(), err2.Error())
							//		}
							//		oplog.Infof("对%s停卡成功", s.GetName())
							//		lab.SetText(strconv.Itoa(bh+1) + "、 " + model.GetStudentByRow(s.Row).PrintStudent())
							//	} else {
							//		fmt.Println("取消....")
							//	}
							//})
							err2 := model.SetLockTime(s)
							if err2 != nil {
								oplog.Infof("对%s停卡失败，err：%s", s.GetName(), err2.Error())
							}
							oplog.Infof("对%s停卡成功", s.GetName())
							lab.SetText(strconv.Itoa(bh+1) + "、 " + model.GetStudentByRow(s.Row).PrintStudent())
							button2.SetLabel("解卡")
						} else if s.State == model.States[3] {
							//fmt.Println("解卡")
							//prompch := make(chan int)
							//info := fmt.Sprintf("确定对%s解卡吗？解卡后将对该学生继续计时", s.GetName())
							//dialog := promptBox(prompch, info, 0)
							//go closeDialog(prompch, dialog, 24*time.Hour, func(i int) {
							//	// 问题所在
							//	if i == 1 {
							//		err2 := model.SetUnlockTime(s)
							//		if err2 != nil {
							//			oplog.Infof("对%s解卡失败，err：%s", s.GetName(), err2.Error())
							//		}
							//		oplog.Infof("对%s解卡成功", s.GetName())
							//		lab.SetText(strconv.Itoa(bh+1) + "、 " + model.GetStudentByRow(s.Row).PrintStudent())
							//	}
							//})
							err2 := model.SetUnlockTime(s)
							if err2 != nil {
								oplog.Infof("对%s解卡失败，err：%s", s.GetName(), err2.Error())
							}
							oplog.Infof("对%s解卡成功", s.GetName())
							lab.SetText(strconv.Itoa(bh+1) + "、 " + model.GetStudentByRow(s.Row).PrintStudent())
							button2.SetLabel("停卡")
						}
					}

				}())
				box2.Add(label1)
				box2.Add(button1)
				box2.Add(button2)
				box1.Add(box2)
			}

			lay1.Add(box1)
			// 根据元素调整大小
			//lay1.SetSize(1550, uint((len(students)+1)*75))
			lay1.SetSize(layWide, uint((len(students)+1)*75))
			win1.ShowAll()

			win1.SetModal(true)
			return win1
		}
	case []model.LearnTime:
		LearnTimes, ok := a.([]model.LearnTime)
		if !ok {
			erlog.Info("断言a为[]model.LearnTime")
		}
		if len(LearnTimes) == 0 {
			Tips("未找到相关信息", 0, ch)
			return nil
		} else {
			builder, err := gtk.BuilderNew()
			if err != nil {
				erlog.Info("BuilderNew 创建builder失败")
			}

			err = builder.AddFromFile(rootPath + "glade/info.glade")
			if err != nil {
				erlog.Info("AddFromFile builder失败")
			}
			window1, err := builder.GetObject("window1")
			win1 := window1.(*gtk.Window)

			win1.SetTitle(util.Conf.Title)
			err = win1.SetIconFromFile(rootPath + "img/icon/maboshi.png")
			if err != nil {
				erlog.Info("SetIconFromFile 读取图标失败")
			}

			win1.Resize(winWide, 600)

			win1.Connect("destroy", func() {
				ch <- 0
			})

			layout1, err := builder.GetObject("layout1")
			if err != nil {
				erlog.Info("layout1 获取layout1失败")
			}

			// 排版
			lay1 := layout1.(*gtk.Layout)

			box1, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
			box1.SetHomogeneous(true)
			for i, LearnTime := range LearnTimes {
				box2, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
				label1, _ := gtk.LabelNew(strconv.Itoa(i+1) + "、 " + LearnTime.PrintString())
				label1.SetSizeRequest(400, 50)
				label1.SetMarginTop(25)

				box2.Add(label1)
				box1.Add(box2)
			}

			lay1.Add(box1)
			// 根据元素调整大小
			lay1.SetSize(layWide, uint((len(LearnTimes)+1)*75))
			win1.ShowAll()

			win1.SetModal(true)
			return win1
		}
	default:
		return nil
	}
}
