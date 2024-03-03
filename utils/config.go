package utils

import (
	"github.com/BurntSushi/toml"
	"os"
)

var (
	Conf conf
)

type (
	conf struct {
		Title       string
		RootPath    string // 视图根地址
		Level       []string
		Level2Loops []int // 成员等级映射循环数
		TimeOfLoop  []int // 每次循环的次数
		DayOfLoop   []int // 每次循环的天数
		DbPath      dbPath
		LogPath     logPath
	}
	dbPath struct {
		StudentInfo string
		LearnRecord string
	}
	logPath struct {
		OperationPath string
		ErrorPath     string
	}
)

// ConfInit  将文件内容解析到结构体
func ConfInit(confFilePath string) {
	if _, err := os.Stat(confFilePath); err != nil {
		panic(err)
	}

	if _, err := toml.DecodeFile(confFilePath, &Conf); err != nil {
		panic(err)
	}
}
