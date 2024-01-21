package main

import (
	"flag"

	"github.com/sirius2001/MediaTool/puller"
	"github.com/sirius2001/MediaTool/pusher"
)

func main() {
	pullUrl := flag.String("i", "", "rtsp input  url")
	pushUrl := flag.String("o", "", "rtmp output  url")
	flag.Parse()

	rtspPullser := puller.NewRTSPPuller(*pullUrl, true)
	go rtspPullser.Start()

	rtmpPusher := pusher.NewRTMPPusher(rtspPullser.RTMPMsg, *pushUrl)
	rtmpPusher.Start()

}
