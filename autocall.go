package main

import (
	"autocall/eventsocket"
	"autocall/logging"
	"time"
)

// 预期效果 ： 自动外呼。根据线路的并发量，一通接一通的往外呼电话。
// 实现思路 ： 设置全局chan数组，数组长度为线路的并发量。使用死循环，开启n个线程， 每线程一通，chan 写入值，每挂断一通电话，chan 清空一个，如此不断的进行，直到需要外拨的号码拨完

type PhoneList struct {
	Number       string
	IsCalled     bool
	StartCalling bool
	Response     string
	Time         int64
	index        int
}

const (
	// 根据线路的并发量确定线程的数目
	gMaxChan = 2
)

var (
	gPhoneList []PhoneList
	gChanStart [gMaxChan]chan bool
)

func init() {

}

func main() {
	end := make(chan int)
	initPhoneList()

	logging.SetLogModel(true, true, ":3360")

	logging.Info("---------------------------------------------")
	logging.Info("%s", "AUTO_CALL")
	logging.Info("Server started, version %s", "V0.1")
	logging.Info("---------------------------------------------")
	logging.Info("Start-up time: %s", time.Now().Format("2006-01-02 15:04:05"))
	logging.Info("Logging Port %s", ":3360")

	// 初始化chan array
	for i := range gChanStart {
		gChanStart[i] = make(chan bool)
	}

	//	处理电话自动外呼
	for i := 0; i < gMaxChan; i++ {
		go call(i)
	}

	// 处理电话挂断释放线路
	go listenEvent()

	end <- 1
}

func call(i int) {
	for {
		logging.Debug("call start")

		// 获取下一个待外呼的电话
		phone := getNeedCallingNumber(i)
		if phone == nil {
			return
		}
		logging.Debug(" --> starting outcall :%s", phone.Number)

		conn, err := eventsocket.Dial("192.168.5.206:8021", "ClueCon")
		if err != nil {
			logging.Error("conn err :[%v]", err)
			return
		}

		cmd := "api originate sofia/gateway/4008000/" + phone.Number + " &echo"
		ev, err := conn.Send(cmd)
		if err != nil {
			logging.Error("send err :[%v]", err)
			return
		}
		logging.Debug("call number:%s ,Result:%v", phone.Number, ev.Body)
		conn.Close()

		time.Sleep(10 * time.Second)
		// 等待通道释放
		chanVar := <-gChanStart[i]

		logging.Debug("Free Channel %d", i)
	}
}

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

	logging.Debug("11hungUp %d", callInfo.index)

	gChanStart[callInfo.index] <- true

	logging.Debug("22hungUp %d", callInfo.index)
}

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

func getNumberInfo(number string) *PhoneList {
	for i := 0; i < len(gPhoneList); i++ {
		if gPhoneList[i].Number == number {
			return &gPhoneList[i]
		}
	}
	return nil
}
