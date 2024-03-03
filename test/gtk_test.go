package test

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"os"
	"testing"
)

var fun func(int) func()

// 测试绑定函数
func TestGTK_Button(t *testing.T) {
	gtk.Init(&os.Args)
	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.Connect("destory", func() {
		gtk.MainQuit()
	})
	button, _ := gtk.ButtonNewWithLabel("按钮")
	button.Connect("clicked", func() {
		fmt.Println("这是第一个函数")
	})
	fun = func(i int) func() {
		if i == 0 {
			return nil
		}
		return func() {
			fmt.Println("running......")
		}
	}
	button.Connect("clicked", fun(1))
	window.Add(button)
	window.ShowAll()
	gtk.Main()
}

// 测试绑定函数
// 可以在函数中调用按钮
func TestGTK_ButtonChangeButton(t *testing.T) {
	gtk.Init(&os.Args)
	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.Connect("destory", func() {
		gtk.MainQuit()
	})
	button, _ := gtk.ButtonNewWithLabel("按钮")
	button.Connect("clicked", func() {
		fmt.Println("这是第一个函数")
	})
	num := 0
	button.Connect("clicked", func() {
		button.SetLabel(fmt.Sprintf("%s%d", "按钮", num))
		num++
	})
	window.Add(button)
	window.ShowAll()
	gtk.Main()
}

func TestGtk111(t *testing.T) {
	gtk.Init(&os.Args)
	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})
	button, _ := gtk.ButtonNewWithLabel("按钮")
	// 开始绑定的函数也会触发，只要绑定过的函数都会触发
	//button.Connect("clicked", func() {
	//	fmt.Println("这是第一个函数")
	//})
	//var num = 1
	//fun = func() func() {
	//
	//	return func() {
	//		if num == 0 {
	//			num = 1
	//			fmt.Println("stop...........")
	//		} else {
	//			num = 0
	//			fmt.Println("run...........")
	//		}
	//
	//	}
	//}
	//button.Connect("clicked", fun())
	window.Add(button)
	window.ShowAll()
	gtk.Main()
}
