package main

import (
	"fmt"
	"os"

	"github.com/q191201771/lal/pkg/hls"
	"github.com/yapingcat/gomedia/go-mp4"
	"github.com/yapingcat/gomedia/go-mpeg2"
)

func generateM3U8(mp4Path string, outPath string, durationMs int) {

	mp4file, err := os.OpenFile(mp4Path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	demxuer := mp4.CreateMp4Demuxer(mp4file)


	tsMuxer := mpeg2.NewTSMuxer()

	tsMuxer.OnPacket = func(pkg []byte) {
		
	}

	vid := tsMuxer.AddStream(mpeg2.TS_STREAM_H264)
	aid := tsMuxer.AddStream(mpeg2.TS_STREAM_AAC)

	for {
		packet, err := demxuer.ReadPacket()
		if err != nil {
			break
		}

		if packet.Cid == mp4.MP4_CODEC_AAC {
			tsMuxer.Write(aid, packet.Data, packet.Pts, packet.Dts)
		}

		if packet.Cid == mp4.MP4_CODEC_H264 {
			tsMuxer.Write(vid, packet.Data, packet.Pts, packet.Dts)
		}
	}

	

}

func main() {
	generateM3U8("./ocrean.mp4", "./hls/", 4000)

}
