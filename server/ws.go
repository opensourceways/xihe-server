package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 消息是一个纯文本
	WsCodeText int = 1
	// 消息是json
	WsCodeJson int = 2
	// 消息是通知退出
	WsCodeExit int = 9
)

var (
	webSocketServer *WSServer
)

type WSMessage struct {
	code    int
	payload interface{}
}

type WSServer struct {
	addr string
	// 给每个job  一个chan通道输送数据
	jobConnections map[string]chan *WSMessage
}

//connect each websocket
func (wsServer *WSServer) wsQueryJobStatus(w http.ResponseWriter, r *http.Request) {
	tokenQuery := r.URL.Query()["Authorization"]
	if len(tokenQuery) == 0 {
		log.Printf("Unauthorized connnect from remote address :%s", r.RemoteAddr)
		return
	}
	token := tokenQuery[0]
	log.Printf("receive token:%s", token)
	jobidQuery := r.URL.Query()["jobid"]
	if len(jobidQuery) == 0 {
		log.Println("jobid not allowed empty ")
		return
	}
	jobid := jobidQuery[0]
	if len(jobid) == 0 {
		log.Println("jobid not allowed empty ")
		return
	}
	currentChann := wsServer.jobConnections[jobid]
	if currentChann == nil {
		currentChann = make(chan *WSMessage, 10)
	}
	wsServer.jobConnections[jobid] = currentChann
	//checkToken(token)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Subprotocols:    []string{"wamp"},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("websocket upgrade failed, error:%s", err.Error())
		return
	}
	defer func() {
		log.Printf(" jobid %s 走了", jobid)
		ws.Close()
		//当ws关闭后  把chan 也关掉
		close(currentChann)
		delete(wsServer.jobConnections, jobid)
		log.Println(wsServer.jobConnections)
	}()
	log.Printf(" jobid %s 进来了", jobid)
	log.Println(wsServer.jobConnections)
	ws.SetWriteDeadline(time.Now().Local().Add(1200 * time.Second))
	for {
		sendData := <-currentChann
		switch sendData.code {
		case WsCodeExit:
			// 主动断开连接
			return
		case WsCodeText:
			payload := sendData.payload.(string)
			err = ws.WriteMessage(websocket.TextMessage, []byte(payload))
		case WsCodeJson:
			err = ws.WriteJSON(sendData.payload)
		}
		if err != nil {
			log.Printf(" websocket WriteMessage Error :%s", err.Error())
			return
		}
	}
}

// 使用某个端口开启一个websocket服务器
func StartWebSocket(wsPort int) {
	webSocketServer = new(WSServer)
	webSocketServer.addr = fmt.Sprintf(":%d", wsPort)
	webSocketServer.jobConnections = make(map[string]chan *WSMessage)
	http.HandleFunc("/queryJob", webSocketServer.wsQueryJobStatus)
	log.Printf("websocket Listening and serving at %s   \n", webSocketServer.addr)
	if err := http.ListenAndServe(webSocketServer.addr, nil); err != nil {
		log.Fatalf("websocket startup failed ,error: %s", err.Error())
	}

}
