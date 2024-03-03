package view

import (
	"github.com/gotk3/gotk3/gtk"
	"maboshijikaxitong/model"
)

// BindAddBut 绑定函数
func BindAddBut(addbut *gtk.Button) {
	addbut.Connect("clicked", func() {
		Save()
	})
}

func Save() {
	builder, err := gtk.BuilderNew()
	if err != nil {
		erlog.Info("BuilderNew 失败")
	}
	err = builder.AddFromFile(rootPath + "glade/add.glade")
	if err != nil {
		erlog.Info("AddFromFile 失败")
	}

	window1, err := builder.GetObject("window1")
	if err != nil {
		erlog.Info("window1 获取失败")
	}
	win1 := window1.(*gtk.Window)
	err = win1.SetIconFromFile(rootPath + "img/icon/maboshi.png")
	if err != nil {
		erlog.Info("SetIconFromFile 读取图标失败")
	}
	win1.SetDeletable(false)
	entry2, err := builder.GetObject("entry2")
	if err != nil {
		erlog.Info("entry2 获取失败")
	}
	entry3, err := builder.GetObject("entry3")
	if err != nil {
		erlog.Info("entry3 获取失败")
	}
	entry4, err := builder.GetObject("entry4")
	if err != nil {
		erlog.Info("entry4 获取失败")
	}
	button1, err := builder.GetObject("button1")
	if err != nil {
		erlog.Info("button1 获取失败")
	}
	button2, err := builder.GetObject("button2")
	if err != nil {
		erlog.Info("button2 获取失败")
	}

	ent2 := entry2.(*gtk.Entry)
	ent3 := entry3.(*gtk.Entry)
	ent4 := entry4.(*gtk.Entry)
	addbut := button1.(*gtk.Button)
	canbut := button2.(*gtk.Button)

	addbut.Connect("clicked", func() {
		clickAddbut(ent2, ent3, ent4, win1)
	})

	canbut.Connect("clicked", func() {
		win1.Destroy()
	})
	win1.ShowAll()
}

func clickAddbut(ent1, ent2, ent3 *gtk.Entry, win *gtk.Window) {
	name, _ := ent1.GetText()
	phone, _ := ent2.GetText()
	level, _ := ent3.GetText()
	if name == "" {
		Tips("姓名不得为空，添加失败", 1, nil)
		return
	}
	if phone == "" {
		Tips("电话号码不得为空，添加失败", 1, nil)
		return
	}
	if level == "" {
		Tips("会员类型不得为空，添加失败", 1, nil)
		return
	}
	err := model.SaveStudent(name, phone, level)
	if err != nil {
		Tips("Excel表格正在被占用", 0, nil)
		return
	}
	Tips("添加成功", 0, nil)

	oplog.Infof("%s添加成功", name)
	win.Destroy()
	return
}
