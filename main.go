package main

import (
	"fmt"
	"time"
	"github.com/robfig/config"
	"os/exec"
	"strings"
	"io/ioutil"
)

func main() {
	// 获取程序的配置
	c, _ := config.ReadDefault("./config.conf")
	sleepTime, _ := c.Int("", "sleep");
	//daemonCheck("frps_ssh.ini", "nohup /usr/local/frp/frps -c /usr/local/frp/frps_ssh.ini > myout.file 2>&1 &")
	for true {
		files, _ := ioutil.ReadDir("./process")
		for _, file := range files {
			if (!strings.HasSuffix(file.Name(), ".conf")) {
				continue
			}

			configSubProcess, _ := config.ReadDefault("./process/" + file.Name())
			keyword, _ := configSubProcess.String("", "keyword")
			execShell, _ := configSubProcess.String("", "exec-shell")
			daemonCheck(keyword, execShell)
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}

/**
keyword 关键字
startShell 启动命令
logFilePath 日志存放目录
 */
func daemonCheck(keyword string, startShell string) {
	fmt.Println("开始执行")
	ret := exec.Command("bash", "-c", fmt.Sprintf("ps -ef | grep %s", keyword))
	//ret := exec.Command("ipconfig")
	retBytes, err := ret.Output()
	if (err != nil) {
		fmt.Println(err, "报错！！！")
		return
	}
	fmt.Println("执行后的输出如下")
	retStr := string(retBytes)
	fmt.Println(retStr)

	if (strings.Count(retStr, keyword) >= 3) {
		fmt.Println("已存在", keyword, "进程")
	} else {
		fmt.Println("不存在", keyword, " 需要重新启动")
		retFrp := exec.Command("bash", "-c", startShell)
		retFrpBytes, errFrp := retFrp.Output()

		if (errFrp != nil) {
			fmt.Println("出现错误", string(retFrpBytes), errFrp.Error())
		}
	}

	fmt.Println("执行完毕")
}
