package puller

import (
	"log/slog"
	"time"

	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/remux"
	"github.com/q191201771/lal/pkg/rtsp"
)

type RTSPPuller struct {
	RTMPMsg     chan base.RtmpMsg
	Remuxer     *remux.Rtmp2AvPacketRemuxer
	Session     *rtsp.PullSession
	ExitChannel chan bool
	RTSPUrl     string
}

func NewRTSPPuller(url string, useTcp bool) *RTSPPuller {
	newPullser := &RTSPPuller{
		RTMPMsg:     make(chan base.RtmpMsg, 100),
		ExitChannel: make(chan bool),
		RTSPUrl:     url,
	}

	remuxer := remux.NewAvPacket2RtmpRemuxer()
	remuxer.WithOnRtmpMsg(func(msg base.RtmpMsg) {
		newPullser.RTMPMsg <- msg
	})
	newPullser.Session = rtsp.NewPullSession(remuxer, func(option *rtsp.PullSessionOption) {
		option.OverTcp = useTcp
		option.PullTimeoutMs = 10000
	})
	return newPullser
}

func (puller *RTSPPuller) Start() error {
	if err := puller.Session.Pull(puller.RTSPUrl); err != nil {
		slog.Error("RTSPPuller", "Url", puller.RTSPUrl, "err", err)
		return err
	}

	go func() {
		for {
			puller.Session.GetStat()
			puller.Session.UpdateStat(uint32(time.Second))
			time.Sleep(1 * time.Second)
		}
	}()

	<-puller.ExitChannel
	return nil
}

func (puller *RTSPPuller) Stop() error {
	if err := puller.Session.Dispose(); err != nil {
		slog.Error("(puller *RTMPPuller) Stop() ", "err", err)
		return err
	}

	puller.ExitChannel <- true
	return nil
}
