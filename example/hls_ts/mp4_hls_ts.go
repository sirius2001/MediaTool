package main

import (
	"flag"

	"github.com/sirius2001/MediaTool/trans"
)

func main() {
	mp4Path := flag.String("i", "", "mp4 filePath")
	hlsTime := flag.Int("t", 4000, "hls ts duation MS")
	prefix := flag.String("p", "./hls/", "hls prefix")
	hlsName := flag.String("o", "hls", "")
	flag.Parse()
	trans.Mp4TransTs(*mp4Path, *prefix, uint64(*hlsTime), *hlsName)
}
