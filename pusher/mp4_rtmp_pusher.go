package pusher

import (
	"log/slog"
	"os"

	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/remux"
	"github.com/yapingcat/gomedia/go-mp4"
)

type Mp4Puhser struct {
	rtmpUrl    string
	rtmpPusher *RTMPPusher
	demuxer    *mp4.MovDemuxer
	remuxer    *remux.AvPacket2RtmpRemuxer
	mp4File    *os.File
	loop       bool
}

func NewMp4Puhser(mp4Path string, rtmpUrl string, loop bool) *Mp4Puhser {
	mp4File, err := os.Open(mp4Path)
	if err != nil {
		slog.Error("NewMp4Puhser", "err", err)
		return nil
	}
	demuxer := mp4.CreateMp4Demuxer(mp4File)
	RTMPMsg := make(chan base.RtmpMsg, 100)

	rtmpPusher := NewRTMPPusher(RTMPMsg, rtmpUrl)
	remuxer := remux.NewAvPacket2RtmpRemuxer().WithOnRtmpMsg(func(msg base.RtmpMsg) {
		rtmpPusher.RTMPMsg <- msg
	})

	tracks, err := demuxer.ReadHead()
	if err != nil {
		slog.Error("NewMp4Puhser", "err", err)
		return nil
	}
	for _, track := range tracks {
		slog.Info("tracks", "trackId", track.TrackId, "codecId", track.Cid)
	}

	return &Mp4Puhser{
		rtmpUrl:    rtmpUrl,
		demuxer:    demuxer,
		remuxer:    remuxer,
		mp4File:    mp4File,
		loop:       loop,
		rtmpPusher: rtmpPusher,
	}
}

func (p *Mp4Puhser) Start() error {

	go p.rtmpPusher.Start()

	for {
		avPacket, err := p.demuxer.ReadPacket()
		if err != nil {
			if p.loop {
				p.mp4File.Seek(0, 0)
				continue
			}
			break
		}
		var rtmpAvpacket base.AvPacket
		if avPacket.Cid == mp4.MP4_CODEC_AAC {
			rtmpAvpacket = base.AvPacket{
				Payload:     avPacket.Data,
				Timestamp:   int64(avPacket.Dts),
				Pts:         int64(avPacket.Pts),
				PayloadType: base.AvPacketPtAac,
			}
		}

		if avPacket.Cid == mp4.MP4_CODEC_H264 {
			rtmpAvpacket = base.AvPacket{
				Payload:     avPacket.Data,
				Timestamp:   int64(avPacket.Dts),
				Pts:         int64(avPacket.Pts),
				PayloadType: base.AvPacketPtAac,
			}
		}
		p.remuxer.FeedAvPacket(rtmpAvpacket)
	}

	return nil

}
