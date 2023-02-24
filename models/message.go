package models

import (
	// "fmt"
	// "ginchat/utils"
	//"time"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"gopkg.in/fatih/set.v0"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId   int64  //发送者
	TargetId int64  //接收者
	Type     int    //发送类型 群聊 私聊 广播
	Media    int    //消息类型 文字 图片 音频
	Content  string //消息内容
	Pic      string
	Url      string
	Dsc      string
	Amount   int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	//1,获取参数，并检验token等合法性
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	// token := query.Get("token")
	// TargetId := query.Get("TargetId")
	// context := query.Get("context")
	// msgtype := query.Get("type")
	isvalida := true
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//2获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//3用户关系

	//4,userid 跟node 绑定并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	//5,完成发送逻辑
	go sendProc(node)

	//6,完成接受逻辑
	go recvProc(node)

	sendMsg(userId, []byte("welcome chat"))

}

func sendProc(node *Node) {
	for {

		select {
		case data := <-node.DataQueue:
			fmt.Println("sendproc >>>msg>>>", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Print(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Print(err)
			return
		}
		dispatch(data)
		broadMsg(data)
		fmt.Println("[ws]<<<<<", data)
		fmt.Println("recvProc >>>msg>>>", string(data))
	}

}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpsendProc()
	go udpRecvProc()
	fmt.Println("init>>>msg>>>goroutine")
}

// 完成udp数据发送协程
func udpsendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 6),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			fmt.Println("udpSendProc data:", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Print(err)
				return
			}
		}
	}
}

// 完成udp数据接收协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	for {
		var buf [512]byte

		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("udpRecvProc data:", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1:
		fmt.Println("udpRecvProc data:", string(data))
		sendMsg(msg.TargetId, data)
		// case 2:
		// 	sendGroupMsg()
		// case 3:
		// 	sendAllMsg()
		// case 4:
	}
}

func sendMsg(userId int64, msg []byte) {
	fmt.Println("sendMsg userid:", userId, "msg", string(msg))
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
