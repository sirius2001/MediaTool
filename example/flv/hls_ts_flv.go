package main

import (
	"flag"

	"github.com/sirius2001/MediaTool/pkg"
	"github.com/sirius2001/MediaTool/trans"
)

func main() {
	path := flag.String("i", "", "hls m3u8 Path")
	flvPath := flag.String("o", "", "flv output path")
	flag.Parse()
	
	tsArry, err := pkg.ParaseM3u8(*path)
	if err != nil {
		panic(err)
	}
	trans.HlsTsTransFlv(tsArry, *flvPath)
}
