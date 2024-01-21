package main

import (
	"flag"

	"github.com/sirius2001/MediaTool/trans"
)

func main() {
	path := flag.String("i", "", "mp4 Path")
	flvPath := flag.String("o", "", "flv output path")
	flag.Parse()

	trans.Mp4TransFlv(*path, *flvPath)
}
