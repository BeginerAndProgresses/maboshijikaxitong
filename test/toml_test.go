package test

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"maboshijikaxitong/utils"
	"os"
	"testing"
	"time"
)

type Config struct {
	Name     string
	Age      int
	Money    float64
	Hobby    []string
	Birthday time.Time
}

type (
	example struct {
		Distros    []distro
		Characters map[string][]struct {
			Name string
			Rank string
		}
	}
	distro struct {
		Name string
		Pkg  string `toml:"Packages"`
	}
)

func (t Config) String() string {
	f := "2006-01-02 15:04:05.999999999"
	if t.Birthday.Hour() == 0 {
		f = "2006-01-02"
	}
	if t.Birthday.Year() == 0 {
		f = "15:04:05.999999999"
	}
	if t.Birthday.Location() == time.UTC {
		f += " UTC"
	} else {
		f += " -0700"
	}
	return t.Birthday.Format(`"` + f + `"`)
}

var f = "../conf/test.toml"

func TestToml1(t *testing.T) {
	if _, err := os.Stat(f); err != nil {
		panic(err)
	}
	var conf Config

	if _, err := toml.DecodeFile(f, &conf); err != nil {
		panic(err)
	}

	fmt.Printf("Name: %v\n", conf.Name)
	fmt.Printf("Age: %v\n", conf.Age)
	fmt.Printf("Money: %v\n\r", conf.Money)
	fmt.Printf("Hobby: %v\n\r", conf.Hobby)
	fmt.Printf("Birthday: %v\n\r", conf.Birthday)
}

// 对象数组
func TestToml2(t *testing.T) {
	if _, err := os.Stat(f); err != nil {
		panic(err)
	}
	var conf example

	if _, err := toml.DecodeFile(f, &conf); err != nil {
		panic(err)
	}
	fmt.Println(conf)
}

// map映射
func TestToml3(t *testing.T) {
	if _, err := os.Stat(f); err != nil {
		panic(err)
	}
	var conf example

	if _, err := toml.DecodeFile(f, &conf); err != nil {
		panic(err)
	}
	fmt.Println(conf)
}

func TestToml(t *testing.T) {
	//utils.ConfInit()
	fmt.Println(utils.Conf)
}
