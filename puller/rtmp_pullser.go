package puller

import (
	"log/slog"
	"time"

	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/rtmp"
)

type RTMPPuller struct {
	RTMPMsg     chan base.RtmpMsg
	Session     *rtmp.PullSession
	ExitChannel chan bool
	RTMPUrl     string
}

func NewRTMPPuller(rtmpUrl string) *RTMPPuller {
	newPullser := &RTMPPuller{
		RTMPMsg:     make(chan base.RtmpMsg, 1024),
		ExitChannel: make(chan bool),
		RTMPUrl:     rtmpUrl,
	}

	newPullser.Session = rtmp.NewPullSession(func(option *rtmp.PullSessionOption) {

	})

	newPullser.Session.WithOnReadRtmpAvMsg(func(msg base.RtmpMsg) {
		newPullser.RTMPMsg <- msg
	})

	return newPullser
}

func (puller *RTMPPuller) Start() error {

	if err := puller.Session.Pull(puller.RTMPUrl); err != nil {
		slog.Error("NewRTMPPuller", "rtmpUrl", puller.RTMPUrl, "err", err)
		return err
	}

	go func() {
		for {
			puller.Session.UpdateStat(uint32(time.Second))
		}
	}()

	<-puller.ExitChannel

	return nil
}

func (puller *RTMPPuller) Stop() error {
	if err := puller.Session.Dispose(); err != nil {
		slog.Error("(puller *RTMPPuller) Stop() ", "err", err)
		return err
	}

	puller.ExitChannel <- true

	return nil
}
