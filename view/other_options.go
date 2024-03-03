package view

import (
	"github.com/gotk3/gotk3/gtk"
	"maboshijikaxitong/model"
	"maboshijikaxitong/utils"
	"time"
)

// BindOtherBut 绑定其他操作按钮
func BindOtherBut(button *gtk.Button) {
	button.Connect("clicked", func() {
		OpenOtherBut()
	})
}

// BindOtherWinBut 绑定其他操作按钮
func BindOtherWinBut(filebut *gtk.Button) {
	filebut.Connect("clicked", func() {
		importFile()
	})
}

func OpenOtherBut() *gtk.Window {
	builder, err := gtk.BuilderNew()
	if err != nil {
		erlog.Info("BuilderNew 失败")
	}
	err = builder.AddFromFile(rootPath + "glade/other.glade")
	if err != nil {
		erlog.Info("AddFromFile 失败")
	}

	window1, err := builder.GetObject("window1")
	if err != nil {
		erlog.Info("window1 获取失败")
	}
	win1 := window1.(*gtk.Window)
	// 不让改变大小
	win1.SetResizable(false)
	win1.SetTitle(utils.Conf.Title)
	err = win1.SetIconFromFile(rootPath + "img/icon/maboshi.png")
	if err != nil {
		erlog.Info("SetIconFromFile 读取图标失败")
	}

	entry1, err := builder.GetObject("entry1")
	if err != nil {
		erlog.Info("entry2 获取失败")
	}

	button1, err := builder.GetObject("button1")
	if err != nil {
		erlog.Info("button1 获取失败")
	}
	button2, err := builder.GetObject("button2")
	if err != nil {
		erlog.Info("button2 获取失败")
	}

	ent1 := entry1.(*gtk.Entry)
	searchbut := button1.(*gtk.Button)
	filebut := button2.(*gtk.Button)

	BindFuncForSearch(ent1, searchbut, false)
	BindOtherWinBut(filebut)
	win1.ShowAll()
	return win1
}

// 导入文件
func importFile() {
	// 确保filename已被赋值
	interdictCh := make(chan int)
	filenamech := make(chan string)
	var filename string
	go func() {
		select {
		case value := <-filenamech:
			filename = value
			interdictCh <- 1
		}
	}()
	//	打开新的窗口，该窗口关闭
	fileChooserDialog := fileSelect(filenamech)

	<-interdictCh
	fileChooserDialog.Destroy()
	if filename == "" {
		return
	}

	oplog.Infof("开始导入表%s", filename)

	err2 := model.AppendInfo2Sheet(filename)
	proCh := make(chan int)
	var dialog *gtk.Dialog
	if err2 == nil {
		dialog = promptBox(proCh, "写入成功", 0)
		oplog.Infof("---导入完成---")
	} else {
		dialog = promptBox(proCh, "写入失败，请确保EXCEL表格未被打开", 0)
		oplog.Infof("---导入失败---")
	}
	go closeDialog(proCh, dialog, 5*time.Second, func(i int) {

	})
}
