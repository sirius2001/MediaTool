package main

import (
	"flag"
	"log/slog"

	"github.com/sirius2001/MediaTool/puller"
)

func main() {
	rtmpUrl := flag.String("i", "", "rtmp input  url")
	flag.Parse()

	rtmpPuller := puller.NewRTMPPuller(*rtmpUrl)
	go func() {
		for msg := range rtmpPuller.RTMPMsg {
			slog.Info("get new msg", "pts", msg.Pts(), "dts", msg.Dts())
		}
	}()
	rtmpPuller.Start()
}
