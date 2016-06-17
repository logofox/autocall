package eventsocket

import (
	"strconv"
)

func parseInt(str string) int64 {
	result, _ := strconv.ParseInt(str, 10, 64)
	if result > 0 && len(str) > 10 {
		result = result / 1000000
	}
	return result
}

func (self *Event) ParseTime(str string) int64 {
	return parseInt(str)
}

func (self *Event) EventName() string {
	return self.Header["Event-Name"]
}

func (self *Event) EventSubclass() string {
	return self.Header["Event-Subclass"]
}

func (self *Event) ChannelName() string {
	return self.Header["Channel-Name"]
}

func (self *Event) ChannelState() string {
	return self.Header["Channel-State"]
}

// Caller

func (self *Event) CallerAni() string {
	return self.Header["Caller-Ani"]
}

func (self *Event) CallerRdnis() string {
	return self.Header["Caller-Rdnis"]
}

func (self *Event) CalleeIdName() string {
	return self.Header["Caller-Callee-Id-Name"]
}

func (self *Event) CalleeIdNumber() string {
	return self.Header["Caller-Callee-Id-Number"]
}

func (self *Event) CallerIdName() string {
	return self.Header["Caller-Caller-Id-Name"]
}

func (self *Event) CallerIdNumber() string {
	return self.Header["Caller-Caller-Id-Number"]
}

func (self *Event) Context() string {
	return self.Header["Caller-Context"]
}

func (self *Event) DestinationNumber() string {
	return self.Header["Caller-Destination-Number"]
}

func (self *Event) Direction() string {
	return self.Header["Caller-Direction"]
}

func (self *Event) ChannelAnsweredTime() int64 {
	str := self.Header["Caller-Channel-Answered-Time"]
	return parseInt(str)
}

func (self *Event) ChannelBridgedTime() int64 {
	str := self.Header["Caller-Channel-Bridged-Time"]
	return parseInt(str)
}

func (self *Event) ChannelCreatedTime() int64 {
	str := self.Header["Caller-Channel-Created-Time"]
	return parseInt(str)
}

func (self *Event) ChannelHangupTime() int64 {
	str := self.Header["Caller-Channel-Hangup-Time"]
	return parseInt(str)
}

func (self *Event) ChannelResurrectTime() int64 {
	str := self.Header["Caller-Channel-Resurrect-Time"]
	return parseInt(str)
}

func (self *Event) ChannelTransferTime() int64 {
	str := self.Header["Caller-Channel-Transfer-Time"]
	return parseInt(str)
}

func (self *Event) OrigCallerIdName() string {
	return self.Header["Caller-Orig-Caller-Id-Name"]
}

func (self *Event) OrigCallerIdNumber() string {
	return self.Header["Caller-Orig-Caller-Id-Number"]
}

func (self *Event) CallState() string {
	return self.Header["Channel-Call-State"]
}

func (self *Event) NetworkAddr() string {
	return self.Header["Caller-Network-Addr"]
}

func (self *Event) BridgeAUniqueId() string {
	return self.Header["Bridge-A-Unique-Id"]
}

func (self *Event) BridgeBUniqueId() string {
	return self.Header["Bridge-B-Unique-Id"]
}

func (self *Event) DTMF() string {
	return self.Header["Dtmf-Digit"]
}

func (self *Event) HangupCause() string {
	return self.Header["Hangup-Cause"]
}

func (self *Event) UniqueUUID() string {
	return self.Header["Unique-Id"]
}

func (self *Event) OtherLegUniqueId() string {
	return self.Header["Other-Leg-Unique-Id"]
}

func (self *Event) OtherType() string {
	return self.Header["Other-Type"]
}

// Variable

// 摘机状态
func (self *Event) AirgoOffhook() string {
	return self.Header["Variable_airgo_offhook"]
}

// IVR号码
func (self *Event) AirgoIVRNumber() string {
	return self.Header["Variable_airgo_ivr_number"]
}

// 中继号码
func (self *Event) AirgoTrunkNumber() string {
	return self.Header["Variable_airgo_trunk_number"]
}

// 队列号码
func (self *Event) AirgoQueueNumber() string {
	return self.Header["Variable_airgo_queue_number"]
}

// 双向回呼标记
func (self *Event) AirgoBothway() string {
	return self.Header["Variable_airgo_bothway"]
}

// 队列成员UUID
func (self *Event) QueueMemberId() string {
	return self.Header["Variable_cc_member_uuid"]
}

// Call Center

func (self *Event) CCAction() string {
	return self.Header["Cc-Action"]
}

func (self *Event) CCAgent() string {
	return self.Header["Cc-Agent"]
}

func (self *Event) CCAgentState() string {
	return self.Header["Cc-Agent-State"]
}

func (self *Event) CCAgentStatus() string {
	return self.Header["Cc-Agent-Status"]
}

// Custom

func (self *Event) AgentID() string {
	return self.Header["Agent-Id"]
}

func (self *Event) UUID() string {
	return self.Header["Uuid"]
}

func (self *Event) Action() string {
	return self.Header["Action"]
}

func (self *Event) Caller() string {
	return self.Header["Caller"]
}

func (self *Event) Callee() string {
	return self.Header["Callee"]
}

func (self *Event) Destination() string {
	return self.Header["Destination"]
}

func (self *Event) Type() string {
	return self.Header["Type"]
}

func (self *Event) State() string {
	return self.Header["State"]
}

func (self *Event) CreatedTime() int64 {
	str := self.Header["Created-Time"]
	return parseInt(str)
}

func (self *Event) DialStatus() string {
	return self.Header["Dial-Status"]
}

func (self *Event) Digits() string {
	return self.Header["Digits"]
}

func (self *Event) IVRNumber() string {
	return self.Header["Ivr-Number"]
}

func (self *Event) TrunkNumber() string {
	return self.Header["Trunk-Number"]
}

func (self *Event) QueueNumber() string {
	return self.Header["Queue-Number"]
}

// Other

func (self *Event) Username() string {
	return self.Header["Username"]
}

func (self *Event) BridgeAgentStart() string {
	return self.Header["bridge-agent-start"]
}

func (self *Event) BridgeAgentEnd() string {
	return self.Header["bridge-agent-end"]
}

func (self *Event) BridgeAgentFail() string {
	return self.Header["bridge-agent-fail"]
}
