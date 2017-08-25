package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"autocall/eventsocket"
	"autocall/logging"
)

type PhoneList struct {
	Number       string // 待拨的号码
	IsCalled     bool   // 是否已经外呼过
	StartCalling bool   // 是否正在外呼
	Response     string // 外呼结果
	Time         int64  // 外呼时长
	index        int    // 外呼是所在的go线程编号
}

const (
	gMaxChan = 2 // 根据线路的并发量确定线程的数目
)

var (
	gPhoneList []PhoneList
	gChanStart [gMaxChan]chan bool
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, Docker!\n")
}

func main() {

	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		logging.Fatal("ListenAndServe: ", err)
	}

	// end := make(chan int)
	// initPhoneList()

	// logging.SetLogModel(true, true, ":3360")

	// logging.Info("---------------------------------------------")
	// logging.Info("%s", "AUTO_CALL")
	// logging.Info("Server started, version %s", "V0.1")
	// logging.Info("---------------------------------------------")
	// logging.Info("Start-up time: %s", time.Now().Format("2006-01-02 15:04:05"))
	// logging.Info("Logging Port %s", ":3360")

	// // 初始化chan array
	// for i := range gChanStart {
	// 	gChanStart[i] = make(chan bool)
	// }

	// //	根据线路并发量，开通多个通道外呼
	// for i := 0; i < gMaxChan; i++ {
	// 	cmd := "api originate sofia/gateway/4008000/%s &playback(/tmp/info.wav)"
	// 	go call(i, cmd)
	// }

	// // 处理电话挂断释放线路
	// go listenEvent()

	// end <- 1
}

func call(i int, cmd string) {
	logging.Debug("call start")
	for {

		// 获取下一个待外呼的电话
		phone := getNeedCallingNumber(i)
		if phone == nil {
			return
		}

		logging.Debug(" --> outGoing call :%s", phone.Number)
		conn, err := eventsocket.Dial("192.168.5.206:8021", "ClueCon")
		if err != nil {
			logging.Error("conn err :[%v]", err)
			return
		}
		defer conn.Close()

		cmd = fmt.Sprintf(cmd, phone.Number)
		ev, err := conn.Send(cmd)
		if err != nil {
			logging.Error("send err :[%v]", err)
			return
		}
		logging.Debug("call number:%s ,Result:%v", phone.Number, ev.Body)
		time.Sleep(10 * time.Second)

		// 等待通道释放
		// chanVar := <-gChanStart[i]
	}
}

// 电话挂断事件监听
func listenEvent() {
	conn, conerr := eventsocket.Dial("192.168.5.206:8021", "ClueCon")
	for {
		logging.Info("---> send signal, waiting -->")
		if conerr == nil {
			_, err := conn.Send("event plain  CHANNEL_HANGUP_COMPLETE ")
			if err != nil {
				logging.Error("conn send err :[%v]", err)
			}

			readEvent(conn)

		} else {
			logging.Debug("connectioning")
			conn, conerr = eventsocket.Dial("192.168.5.206:8021", "ClueCon")
		}
	}
}

// 过滤挂断事件
func readEvent(conn *eventsocket.Connection) {
	ev, err := conn.ReadEvent()
	if err != nil {
		logging.Error("Conn ReadEvent Err:%v", err)
	}
	switch ev.EventName() {
	case "CHANNEL_HANGUP_COMPLETE":
		go hungUp(ev)
	default:
	}
}

// 释放已挂断的线路，准备下一通呼叫
func hungUp(ev *eventsocket.Event) {

	logging.Debug("EventName:%s", ev.EventName())
	logging.Debug("DestinationNumber:%s", ev.DestinationNumber())

	destnumber := ev.DestinationNumber()
	callInfo := getNumberInfo(destnumber)
	// 该通记录不是自动外呼中的号码，忽略
	if callInfo == nil {
		logging.Debug("ignore destnumber: %s", destnumber)
		return
	}

	callInfo.IsCalled = true
	callInfo.Response = ev.HangupCause()
	callInfo.Time = ev.ChannelHangupTime() - ev.ChannelCreatedTime()

	time.Sleep(5 * time.Second)

	gChanStart[callInfo.index] <- true
}

// --------------------------------------------------------------------------------

// 初始化自动外呼列表
func initPhoneList() {
	// 自动外呼列表
	list := []string{"01385****", "159*****", "159******", "137****"}
	for i := 0; i < len(list); i++ {
		var p PhoneList
		p.Number = list[i]
		p.IsCalled = false
		p.StartCalling = false
		p.index = -1
		gPhoneList = append(gPhoneList, p)
	}
}

// 获取下一通准备外呼的号码
func getNeedCallingNumber(index int) *PhoneList {
	for i := 0; i < len(gPhoneList); i++ {
		if gPhoneList[i].IsCalled == false && gPhoneList[i].StartCalling == false {
			gPhoneList[i].StartCalling = true
			gPhoneList[i].index = index
			return &gPhoneList[i]
		}
	}
	return nil
}

// 获取外呼号码的信息列表
func getNumberInfo(number string) *PhoneList {
	for i := 0; i < len(gPhoneList); i++ {
		if gPhoneList[i].Number == number {
			return &gPhoneList[i]
		}
	}
	return nil
}
