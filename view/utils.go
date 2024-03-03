package view

import (
	"github.com/gotk3/gotk3/gtk"
	util "maboshijikaxitong/utils"
	"time"
)

// warningWindow 信息提示框
// typ == 0 warning提示框
// typ == 1 error提示框
// 返回0 点击取消（销毁）
// 返回1 点击确定
func promptBox(ch chan int, info string, typ uint) *gtk.Dialog {
	builder, err := gtk.BuilderNew()
	if err != nil {
		erlog.Info("BuilderNew 创建builder失败")
	}
	err = builder.AddFromFile(rootPath + "glade/promptbox.glade")
	dialog1, err := builder.GetObject("dialog1")
	if err != nil {
		erlog.Info("GetObject dialog1 获取失败")
	}

	dig1 := dialog1.(*gtk.Dialog)

	dig1.SetModal(true)
	dig1.SetTitle(util.Conf.Title)

	err = dig1.SetIconFromFile(rootPath + "img/icon/maboshi.png")
	if err != nil {
		erlog.Info("dig1 SetIconFromFile 读取图标失败")
	}
	dig1.SetDeletable(false)

	image1, err := builder.GetObject("image1")
	if err != nil {
		erlog.Info("GetObject image1 获取失败")
	}

	label1, err := builder.GetObject("label1")
	if err != nil {
		erlog.Info("GetObject label1 获取失败")
	}

	label2, err := builder.GetObject("label2")
	if err != nil {
		erlog.Info("GetObject label2 获取失败")
	}

	button1, err := builder.GetObject("button1")
	if err != nil {
		erlog.Info("GetObject button1 获取失败")
	}

	button2, err := builder.GetObject("button2")
	if err != nil {
		erlog.Info("GetObject button2 获取失败")
	}

	img1 := image1.(*gtk.Image)
	lab1 := label1.(*gtk.Label)
	lab2 := label2.(*gtk.Label)
	but1 := button1.(*gtk.Button)
	but2 := button2.(*gtk.Button)

	if typ == 0 {
		img1.SetFromFile(rootPath + "warn.png")
		lab1.SetText("提示！")
		lab2.SetText(info)
		but1.Connect("clicked", func() {
			ch <- 1
		})
		but2.Connect("clicked", func() {
			ch <- 0
		})
	} else if typ == 1 {
		img1.SetFromFile(rootPath + "error.png")
		lab1.SetText("错误！！！")
		lab2.SetText(info)
		but1.Connect("clicked", func() {
			ch <- 1
		})
		but2.Connect("clicked", func() {
			ch <- 0
		})
	}

	dig1.Show()

	return dig1
}

// closeDialog 关闭提示框时的操作
// ch 为无缓冲通道，用于接收在widget窗口的操作值
// dialog 为通知窗口
// outTime 为过期时间,如果不想让窗口自动关闭，可将值设置为很长的时间
func closeDialog(ch chan int, dialog *gtk.Dialog, outTime time.Duration, f myHandler) {
	select {
	case value := <-ch:
		f(value)
	case <-time.After(outTime):
	}
	dialog.Destroy()
}

// Tips 简化代码
// info 提示信息
// typ 种类 0提示 1错误
func Tips(info string, typ uint, ch chan int) {
	var chProm = make(chan int)
	dialog := promptBox(chProm, info, typ)

	go closeDialog(chProm, dialog, 5*time.Second, func(i int) {
		if ch != nil {
			ch <- i
		}
	})
}

func WaitWindow(info string) *gtk.Window {
	window1, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		erlog.Info("WindowNew 创建waitWindow失败")

	}
	window1.SetTitle(info)
	window1.SetModal(true)
	window1.SetIconFromFile(rootPath + "img/icon/maboshi.png")
	window1.SetDeletable(false)
	window1.SetPosition(gtk.WIN_POS_CENTER)
	window1.SetKeepAbove(true)
	window1.SetDefaultSize(400, 50)
	window1.Show()
	return window1
}

func WaitClose(win *gtk.Window) {
	win.Destroy()
}

func ChangeLabel(info string, lab *gtk.Label) {
	lab.SetText(info)
}

// ClearEntry 清空文本框，info 为默认显示
func ClearEntry(info string, ent *gtk.Entry) {
	ent.SetText("")
	ent.SetPlaceholderText(info)
}

// 文件选择
func fileSelect(filename chan string) *gtk.FileChooserDialog {
	with1Button, _ := gtk.FileChooserDialogNewWith2Buttons("选择文件", nil, gtk.FILE_CHOOSER_ACTION_OPEN, "确认", gtk.RESPONSE_OK, "取消", gtk.RESPONSE_CANCEL)
	err := with1Button.SetIconFromFile(rootPath + "icon.png")
	if err != nil {
		erlog.Info("SetIconFromFile 读取图标失败")
	}
	with1Button.SetTitle("选择文件")
	with1Button.SetPosition(gtk.WIN_POS_CENTER)
	with1Button.SetIconFromFile(rootPath + "img/icon/maboshi.png")
	with1Button.SetModal(true)
	with1Button.SetDeletable(false)
	responseType := with1Button.Run()
	// 开一个线程监控是否选择了一个文件
	go func() {
		for {
			if responseType == gtk.RESPONSE_OK {
				//fmt.Println("收到RESPONSE_OK")
				getFilename := with1Button.GetFilename()
				transform := util.PathTransform(getFilename)
				filename <- transform
				break
			} else if responseType == gtk.RESPONSE_CANCEL {
				filename <- ""
				break
			}
		}
	}()

	with1Button.Show()
	return with1Button
}
