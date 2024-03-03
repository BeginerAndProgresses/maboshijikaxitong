package view

import (
	"github.com/gotk3/gotk3/gtk"
	"go.uber.org/zap"
	util "maboshijikaxitong/utils"
	"os"
)

var (
	oplog *zap.SugaredLogger
	erlog *zap.SugaredLogger
)

type myHandler func(int)

//var rootPath = "../static/"

var rootPath = "./static/"

//var fun func() func()

// View  开始窗口
func View() {
	rootPath = util.Conf.RootPath
	//fmt.Println(rootPath, "rootPath")
	oplog = util.GetOperateLog()
	erlog = util.GetErrorLog()
	if oplog == nil || erlog == nil {
		//fmt.Println("启动失败")
		return
	}
	gtk.Init(&os.Args)

	builder, err := gtk.BuilderNew()
	if err != nil {
		erlog.Info("BuilderNew 创建builder失败")
	}

	err = builder.AddFromFile(rootPath + "glade/home.glade")
	if err != nil {
		erlog.Info("AddFromFile 创建builder失败")
	}

	window1, err := builder.GetObject("window1")
	win, ok := window1.(*gtk.Window)
	if !ok {
		erlog.Info("GetObject 获取window1失败")
	}
	win.SetTitle(util.Conf.Title)
	win.SetIconFromFile(rootPath + "img/icon/maboshi.png")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	entry1, err := builder.GetObject("entry1")
	ent, ok := entry1.(*gtk.Entry)
	if !ok {
		erlog.Info("GetObject 获取entry1失败")
	}

	button1, err := builder.GetObject("button1")
	searchBut, ok := button1.(*gtk.Button)
	if !ok {
		erlog.Info("GetObject 获取entry1失败")
	}

	BindFuncForSearch(ent, searchBut, true)

	button4, err := builder.GetObject("button4")
	addBut, ok := button4.(*gtk.Button)
	if !ok {
		erlog.Info("GetObject 获取button4失败")
	}
	BindAddBut(addBut)
	button5, err := builder.GetObject("button5")
	otherBut, ok := button5.(*gtk.Button)
	if !ok {
		erlog.Info("GetObject 获取button5失败")
	}

	BindOtherBut(otherBut)

	win.ShowAll()
	gtk.Main()
}
