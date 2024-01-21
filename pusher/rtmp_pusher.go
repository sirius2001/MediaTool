package pusher

import (
	"log/slog"

	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/remux"
	"github.com/q191201771/lal/pkg/rtmp"
)

type RTMPPusher struct {
	Remuxer     *remux.AvPacket2RtmpRemuxer
	Session     *rtmp.PushSession
	ExitChannel chan bool
	RTMPMsg     chan base.RtmpMsg
	PushUrl     string
}

func NewRTMPPusher(RTMPMsg chan base.RtmpMsg, pushUrl string) *RTMPPusher {
	var rtmpPusher RTMPPusher
	rtmpPusher.Session = rtmp.NewPushSession(func(option *rtmp.PushSessionOption) {
		option.WriteAvTimeoutMs = 10000
		option.PushTimeoutMs = 10000
	})

	rtmpPusher.RTMPMsg = RTMPMsg
	rtmpPusher.ExitChannel = make(chan bool)
	rtmpPusher.PushUrl = pushUrl

	return &rtmpPusher
}

func (pusher *RTMPPusher) Start() error {
	if err := pusher.Session.Push(pusher.PushUrl); err != nil {
		slog.Error("(pusher *RTMPPusher) Start()", "Push Url", pusher.PushUrl)
		return err
	}

	go func() {
		for msg := range pusher.RTMPMsg {
			pusher.Session.Write(rtmp.Message2Chunks(msg.Payload, &msg.Header))
		}
	}()

	<-pusher.ExitChannel

	return nil
}
