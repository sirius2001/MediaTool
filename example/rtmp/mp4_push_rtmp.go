package main

import (
	"flag"

	"github.com/sirius2001/MediaTool/pusher"
)

func main() {
	pullUrl := flag.String("i", "", "mp4 input  path")
	pushUrl := flag.String("o", "", "rtmp output  url")
	loop := flag.Bool("loop", false, "loop to push")
	flag.Parse()

	var pusher pusher.Mp4Pusher

	pusher.Start(*pullUrl, *pushUrl, *loop)
}
