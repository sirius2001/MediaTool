package main

import (
	"flag"

	"github.com/sirius2001/MediaTool/pusher"
)

func main() {
	pullUrl := flag.String("i", "", "flv input  path")
	pushUrl := flag.String("o", "", "rtmp output  url")
	loop := flag.Bool("loop", false, "loop to push")
	flag.Parse()
	flvPusher := pusher.FlvPuhser{}
	flvPusher.FlvStart(*pullUrl, *pushUrl, *loop)
}
