package main

import (
	_ "embed"
	"github.com/wailsapp/wails"
	"math/rand"
	"time"
	"wails-vue3/service"
	"wails-vue3/service/tools"
)

var (
	rd     = rand.New(rand.NewSource(time.Now().UnixNano()))
	msgMap = [...]string{
		"Hello World",
		"你好",
		"안녕하세요",
		"こんにちは",
		"ON LI DAY FaOHE MASHI",
		"hallo! Wie geht es dir?",
	}
)

func basic() string {
	return msgMap[rd.Intn(len(msgMap))]
}

//go:embed frontend/dist/app.js
var js string

//go:embed frontend/dist/app.css
var css string

func main() {
	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "wails-vue3",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(basic)
	app.Bind(service.NewStruct())
	app.Bind(tools.NetWorkStatus) //绑定网络状态
	app.Run()
}
