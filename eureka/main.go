package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/hudl/fargo"
	"gopkg.in/ini.v1"
	"os"
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func main() {
	l, _ := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "./readline.tmp",
		AutoComplete:    nil,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("config.ini 读取失败：%v", err)
		os.Exit(1)
	}
	var eurekaHost = make(map[string]string)
	for _, env := range cfg.Section("").Keys() {
		eurekaHost[env.Name()] = env.String()
		println("环境：", env.Name())
	}
	println("输入环境：")
	cmdStr, _ := l.Readline()
	// 连接 Eureka 服务器
	client := fargo.NewConn(eurekaHost[cmdStr])
	for {
		println("输入应用：")
		cmdStr, _ = l.Readline()
		app, _ := client.GetApp(cmdStr)
		for _, targetInstance := range app.Instances {
			fmt.Println(targetInstance.IPAddr)
		}
	}

}
