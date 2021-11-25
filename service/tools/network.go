package tools

import (
	"fmt"
	"os/exec"
	"time"
)

func NetWorkStatus() bool {
	var netStatus = make(chan bool)

	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		// 总共请求1次，超时时间设置为5s
		cmd := exec.Command("ping", "baidu.com", "-c", "1", "-W", "5")
		for range ticker.C {
			fmt.Println("NetWorkStatus Start:", time.Now().Unix())
			err := cmd.Run()
			fmt.Println("NetWorkStatus End :", time.Now().Unix())
			if err == nil {
				netStatus <- true
			} else {
				netStatus <- false
			}
		}
	}()
	return <-netStatus
}
