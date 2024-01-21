package main

import (
	"github.com/sirius2001/MediaTool/pkg"
	"github.com/sirius2001/MediaTool/trans"
)

func main() {
	tsArry, err := pkg.ParaseM3u8("/home/sirius/Desktop/MediaTool/example/hls_ts/playlist.m3u8")
	if err != nil {
		panic(err)
	}
	trans.HlsTsTransFlv(tsArry,"./test/hls_flv.flv")
}
