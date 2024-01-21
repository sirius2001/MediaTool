package main

import (
	"flag"

	"github.com/sirius2001/MediaTool/puller"
	"github.com/sirius2001/MediaTool/pusher"
)

func main() {
	pullUrl := flag.String("i", "", "rtmp input  url")
	pushUrl := flag.String("o", "", "rtmp input  url")
	
	flag.Parse()

	rtmpPuller := puller.NewRTMPPuller(*pullUrl)
	go rtmpPuller.Start()

	rtmpPusher := pusher.NewRTMPPusher(rtmpPuller.RTMPMsg, *pushUrl)
	rtmpPusher.Start()
}
