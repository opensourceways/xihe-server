package server

import (
	"fmt"
	"testing"
	"time"
)

func TestWs(t *testing.T) {
	go sendMessage("唐僧念紧箍咒 到第%d 遍")
	StartWebSocket(9527)
}
func sendMessage(message string) {
	for i := 0; i < 1000; i++ {
		for _, wsChann := range webSocketServer.jobConnections {
			wsmesage := new(WSMessage)

			wsmesage.code = WsCodeText
			if i%5 == 0 {
				//5 的倍数的时候发送  json
				wsmesage.code = WsCodeJson
				wsmesage.payload = &struct {
					Name string `json:"name"`
					Say  string `json:"say"`
					Age  int    `json:"age"`
				}{
					Name: "luonancom",
					Age:  32,
					Say:  fmt.Sprintf(message, i),
				}
			} else {
				// 发送普通文本
				wsmesage.code = WsCodeText
				wsmesage.payload = fmt.Sprintf(message, i)
			}
			if i == 41 {
				// 测试当发送到41的时候 主动断开连接
				wsmesage.code = WsCodeExit
			}
			wsChann <- wsmesage
		}
		time.Sleep(time.Second)
	}
}
